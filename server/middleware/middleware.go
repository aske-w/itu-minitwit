package middleware

import "github.com/kataras/iris/v12"

func InitMiddleware(ctx iris.Context) {
	PrometheusRequestCountMiddleware(ctx)
}
