package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "codehub_services_requestDuration",
		Buckets: prometheus.LinearBuckets(0.01, 0.01, 10),
	}, []string{"path", "method"})

	DBQueryDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "codehub_services_SqlDuration",
		Buckets: prometheus.LinearBuckets(0.01, 0.01, 10),
	}, []string{"rep", "method"})

	Hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "codehub_services_hits",
	}, []string{"status", "path"})
)
