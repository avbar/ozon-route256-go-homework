package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	RequestsCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "requests_total",
	},
		[]string{"handler"},
	)
	ErrorsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "errors_total",
	})
	HistogramResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "histogram_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"status", "handler"},
	)
	HistogramClientResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "grpc",
		Name:      "histogram_client_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"status", "handler"},
	)
	HistogramDBResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "route256",
		Subsystem: "db",
		Name:      "histogram_db_response_time_seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	},
		[]string{"query"},
	)
)

func New() http.Handler {
	return promhttp.Handler()
}

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	RequestsCounter.WithLabelValues(info.FullMethod).Inc()

	timeStart := time.Now()
	res, err := handler(ctx, req)
	elapsedTime := time.Since(timeStart)

	status := status.Convert(err).Code()
	if status != codes.OK {
		ErrorsCounter.Inc()
	}

	HistogramResponseTime.WithLabelValues(status.String(), info.FullMethod).Observe(elapsedTime.Seconds())

	return res, err
}

func MetricsClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	timeStart := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	elapsedTime := time.Since(timeStart)

	status := status.Convert(err).Code().String()
	HistogramClientResponseTime.WithLabelValues(status, method).Observe(elapsedTime.Seconds())

	return err
}
