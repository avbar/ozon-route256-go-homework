package cache

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HitsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "hits_total",
	})
	ErrorsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "errors_total",
	})
	HistogramResponseTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "cache",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	})
)
