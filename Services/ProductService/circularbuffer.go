package senders

import "sync"

type CircularBuffer[T any] struct {
	size   int
	mu     sync.Mutex
	buffer chan T
}

func NewCircularBuffer[T any](size int) *CircularBuffer[T] {
	return &CircularBuffer[T]{
		buffer: make(chan T, size),
		size:   size,
	}
}

func (c *CircularBuffer[T]) Push(data T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.buffer) == c.size {
		<-c.buffer
	}

	c.buffer <- data
}

func (c *CircularBuffer[T]) Pop() T {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.buffer) == 0 {
		return c.getDefault()
	}

	return <-c.buffer
}

func (c *CircularBuffer[T]) Len() int {
	return len(c.buffer)
}

func (c *CircularBuffer[T]) getDefault() T {
	var result T
	return result
}
