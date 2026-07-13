package middleware

import (
	"iptv-spider-sh/global"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

func AccessLog(ctx iris.Context) {
	if global.ACCESS_LOG == nil {
		ctx.Next()
		return
	}

	method := ctx.Method()
	path := ctx.Path()
	clientIP := ctx.RemoteAddr()

	ctx.Next()

	statusCode := ctx.GetStatusCode()

	global.ACCESS_LOG.Info(
		"access",
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", statusCode),
		zap.String("client_ip", clientIP),
	)
}
