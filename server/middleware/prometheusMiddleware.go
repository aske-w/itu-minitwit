package middleware

import (
	"aske-w/itu-minitwit/web/utils"
	"regexp"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of requests.",
	},
	[]string{"path", "method"},
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
}, []string{"path", "method", "status_code"})

func PrometheusRequestCountMiddleware(ctx iris.Context) {
	path := ctx.Path()
	// Remove usernames from path since it otherwise generates enough data
	// to cause Prometheus to run out of memory when scraping webserver metrics.
	match, err := regexp.MatchString("/api/msgs/.+", path)
	utils.CheckError(err)
	if match {
		path = "/api/msgs/<username>"
		goto skip
	}

	match, err = regexp.MatchString("/api/fllws/.+", path)
	utils.CheckError(err)
	if match {
		path = "/api/fllws/<username>"
		goto skip
	}

	match, err = regexp.MatchString("/api/users/.*?/follow", path)
	utils.CheckError(err)
	if match {
		path = "/api/users/<username>/follow"
		goto skip
	}

	match, err = regexp.MatchString("/api/users/.*?/isfollowing", path)
	utils.CheckError(err)
	if match {
		path = "/api/users/<username>/isfollowing"
		goto skip
	}

	match, err = regexp.MatchString("/api/users/.*?/tweets", path)
	utils.CheckError(err)
	if match {
		path = "/api/users/<username>/tweets"
		goto skip
	}

	match, err = regexp.MatchString("/api/users/.*", path)
	utils.CheckError(err)
	if match {
		path = "/api/users/<username>"
	}
skip:

	httpMethod := ctx.Method()
	httpStatusCode := strconv.Itoa(ctx.GetStatusCode())

	totalRequests.WithLabelValues(path, httpMethod).Inc()
	timer := prometheus.NewTimer(httpDuration.WithLabelValues(path, httpMethod, httpStatusCode))

	ctx.Next()
	responseStatus.WithLabelValues(httpStatusCode).Inc()
	timer.ObserveDuration()
}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(httpDuration)
	prometheus.Register(responseStatus)
}
