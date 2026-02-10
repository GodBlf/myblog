package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "myblog_http_requests_total",
		Help: "Total number of HTTP requests.",
	}, []string{"method", "path", "status"})

	requestInFlight = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "myblog_http_requests_in_flight",
		Help: "Current number of in-flight HTTP requests.",
	}, []string{"method", "path"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "myblog_http_request_duration_seconds",
		Help:    "HTTP request duration in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path", "status"})
)

func Metric() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.FullPath()
		if path == "" {
			path = "unknown"
		}

		if path == "/metrics" {
			ctx.Next()
			return
		}

		method := ctx.Request.Method
		start := time.Now()

		requestInFlight.WithLabelValues(method, path).Inc()
		defer requestInFlight.WithLabelValues(method, path).Dec()

		ctx.Next()

		status := strconv.Itoa(ctx.Writer.Status())

		requestTotal.WithLabelValues(method, path, status).Inc()
		requestDuration.WithLabelValues(method, path, status).Observe(time.Since(start).Seconds())
	}
}
