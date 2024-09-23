package main

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "store"

var (
	initMetricsOnce     sync.Once
	PromMetrics         *Metrics
	initRestMetricsOnce sync.Once
	restMetrics         *RestMetrics
)

type Metrics struct {
	Rest *RestMetrics
}

func InitMetrics() *Metrics {
	initMetricsOnce.Do(func() {
		PromMetrics = newMetrics()
	})

	return PromMetrics
}

func newMetrics() *Metrics {
	return &Metrics{
		Rest: InitRestMetrics(),
	}
}

// RestMetrics представляет набор метрик формируемых компонентом rest
type RestMetrics struct {
	ReqCnt                        *prometheus.CounterVec
	ReqDur                        *prometheus.HistogramVec
	ReqSz                         prometheus.Summary
	ResSz                         prometheus.Summary
	GetRtm70eBarcodesDistribution prometheus.Histogram
}

func InitRestMetrics() *RestMetrics {
	initRestMetricsOnce.Do(func() {
		restMetrics = newRestMetrics()
	})

	return restMetrics
}

func newRestMetrics() *RestMetrics {
	return &RestMetrics{
		ReqCnt: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "requests_total",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		}, []string{"code", "method", "host", "url"}),
		ReqDur: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		}, []string{"code", "method", "url"}),
		ReqSz: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "request_size_bytes",
			Help:      "The HTTP request sizes in bytes.",
		}),
		ResSz: promauto.NewSummary(prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "response_size_bytes",
			Help:      "The HTTP response sizes in bytes.",
		}),
	}
}
