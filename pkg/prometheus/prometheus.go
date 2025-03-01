package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type Metrics struct {
	Histogram *prometheus.HistogramVec
	Counter   *prometheus.CounterVec
}

var (
	once sync.Once
	m    *Metrics
)

var (
	his = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "count the total server histogram",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"service", "method"})

	counter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_request_total_count",
		Help: "Total number of gRPC requests processed",
	}, []string{"service", "method", "code"})
)

func InitPrometheus() *Metrics {
	once.Do(func() {
		m = &Metrics{
			Histogram: his,
			Counter:   counter,
		}
		prometheus.MustRegister(m.Histogram, m.Counter)
	})
	return m
}
