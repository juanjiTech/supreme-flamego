package middleware

import (
	"github.com/flamego/flamego"
	"supreme-flamego/core/logx"
	"supreme-flamego/pkg/utils/gen/snowflake"
	"time"
)

// RequestLog 顺便插入 trace log
func RequestLog() flamego.Handler {
	return func(c flamego.Context) {
		traceId := snowflake.Node.Generate().String()
		log := logx.NameSpace("access")
		c.Map(*log)
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 日志格式
		log.Infof("| %s | %3d | %13v | %15s | %s | %s |",
			traceId,
			c.ResponseWriter().Status(),
			time.Now().Sub(startTime),
			c.RemoteAddr(),
			c.Request().Method,
			c.Request().RequestURI,
		)
	}
}
