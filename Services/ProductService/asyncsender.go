package senders

import (
	"context"
	"sync"
	"time"
)

type Sender[T any] interface {
	Send(ctx context.Context, items []*T) error
}

type AsyncSender[W Sender[T], T any] struct {
	opts   *Options
	sender Sender[T]

	mu    sync.Mutex
	batch []*T

	sending          bool
	sendingQueue     chan []*T
	sendingBuffer    *CircularBuffer[[]*T]
	bufferProcessing bool

	onSendingFailure func()
}

func NewAsyncSender[W Sender[T], T any](sender Sender[T], opts ...Option) *AsyncSender[W, T] {
	options := newDefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	return &AsyncSender[W, T]{
		opts:             options,
		sender:           sender,
		batch:            make([]*T, 0, options.batchCapacity),
		sendingQueue:     make(chan []*T, options.sendingQueueCapacity),
		sendingBuffer:    NewCircularBuffer[[]*T](options.sendingBufferCapacity),
		onSendingFailure: options.onSendingFailure,
	}
}

func (as *AsyncSender[W, T]) Send(item *T) {
	if !as.sending {
		return
	}

	as.mu.Lock()
	defer as.mu.Unlock()
	as.batch = append(as.batch, item)
	if len(as.batch) >= as.opts.batchCapacity {
		as.sendingQueue <- as.batch
		as.batch = make([]*T, 0, as.opts.batchCapacity)
	}
}

func (as *AsyncSender[W, T]) Run(ctx context.Context) {
	as.sending = true
	ticker := time.NewTicker(as.opts.sendingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			as.sending = false
			return
		case batch := <-as.sendingQueue:
			as.send(ctx, batch)
			ticker.Reset(as.opts.sendingInterval)
		case <-ticker.C:
			as.mu.Lock()
			if len(as.batch) == 0 {
				as.mu.Unlock()
				continue
			}
			batch := as.batch
			as.batch = make([]*T, 0, as.opts.batchCapacity)
			as.mu.Unlock()

			as.send(ctx, batch)
		}
	}
}

func (as *AsyncSender[W, T]) send(ctx context.Context, batch []*T) {
	if len(batch) == 0 {
		return
	}

	err := as.sender.Send(ctx, batch)
	if err != nil {
		as.sendingBuffer.Push(batch)
		if as.onSendingFailure != nil {
			as.onSendingFailure()
		}
		return
	}

	if as.sendingBuffer.Len() > 0 && !as.bufferProcessing {
		as.bufferProcessing = true
		go func(ctx context.Context) {
			defer func() {
				as.bufferProcessing = false
			}()

			for batch := as.sendingBuffer.Pop(); batch != nil; {
				err := as.sender.Send(ctx, batch)
				if err != nil {
					if as.onSendingFailure != nil {
						as.onSendingFailure()
					}
					break
				}
			}
		}(ctx)
	}
}

const (
	defaultBatchCapacity         = 10000
	defaultSendingQueueCapacity  = 0
	defaultSendingBufferCapacity = 100
	defaultSendingInterval       = 1 * time.Minute
)

type Option func(*Options)

type Options struct {
	batchCapacity         int
	sendingQueueCapacity  int
	sendingBufferCapacity int
	sendingInterval       time.Duration
	onSendingFailure      func()
}

func newDefaultOptions() *Options {
	return &Options{
		batchCapacity:         defaultBatchCapacity,
		sendingQueueCapacity:  defaultSendingQueueCapacity,
		sendingBufferCapacity: defaultSendingBufferCapacity,
		sendingInterval:       defaultSendingInterval,
	}
}

func WithBatchCapacity(batchCapacity int) Option {
	return func(opts *Options) {
		opts.batchCapacity = batchCapacity
	}
}

func WithSendingQueueCapacity(sendingQueueCapacity int) Option {
	return func(opts *Options) {
		opts.sendingQueueCapacity = sendingQueueCapacity
	}
}

func WithSendingBufferCapacity(sendingBufferCapacity int) Option {
	return func(opts *Options) {
		opts.sendingBufferCapacity = sendingBufferCapacity
	}
}

func WithSendingInterval(sendingInterval time.Duration) Option {
	return func(opts *Options) {
		opts.sendingInterval = sendingInterval
	}
}

func OnSendingFailure(f func()) Option {
	return func(opts *Options) {
		opts.onSendingFailure = f
	}
}
