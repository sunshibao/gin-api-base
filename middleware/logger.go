package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logger 请求日志中间件（基于 Go 标准库 slog）
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			path = path + "?" + raw
		}

		// 生成或获取 traceId
		traceID := c.GetHeader("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Header("X-Trace-Id", traceID)
		c.Set("trace_id", traceID)

		// 请求日志
		slog.Info("请求开始",
			"trace_id", traceID,
			"method", c.Request.Method,
			"path", path,
			"client_ip", c.ClientIP(),
		)

		c.Next()

		// 响应日志
		slog.Info("请求完成",
			"trace_id", traceID,
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency", time.Since(start).String(),
			"size", c.Writer.Size(),
		)
	}
}
