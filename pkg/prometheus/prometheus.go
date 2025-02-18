package prom

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HisMetrics = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_http_request_duration_seconds",
		Help:    "count the total server histogram",
		Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"service", "method"})

	ReqMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_request_total_vec",
		Help: "Total number of gRPC requests processed",
	}, []string{"service", "method", "code"})
)

func init() {
	prometheus.MustRegister(HisMetrics, ReqMetrics)
}
