package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

func PrometheusRequestCountMiddleware(ctx iris.Context) {
	totalRequests.WithLabelValues(ctx.Path()).Inc()
	ctx.Next()
}

func init() {
	prometheus.Register(totalRequests)
}
