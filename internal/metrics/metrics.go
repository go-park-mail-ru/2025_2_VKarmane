package metrics
import (
    "github.com/prometheus/client_golang/prometheus"
)


var (
    HttpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    HttpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duration of HTTP requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path", "status"},
    )
)

func init() {
    prometheus.MustRegister(HttpRequestsTotal)
    prometheus.MustRegister(HttpRequestDuration)
}