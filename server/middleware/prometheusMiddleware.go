package middleware

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func PrometheusRequestCountMiddleware(ctx iris.Context) {
	path := ctx.Path()
	totalRequests.WithLabelValues(path).Inc()
	timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))

	ctx.Next()
	responseStatus.WithLabelValues(strconv.Itoa(ctx.GetStatusCode())).Inc()
	timer.ObserveDuration()
}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(httpDuration)
	prometheus.Register(responseStatus)
}
