package prom

import (
	"fmt"
	"github.com/peterouob/golang_template/tools"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		tools.Log(fmt.Sprintf("receive a signal %s", sig.String()))
		os.Exit(0)
	}()
	return m
}
