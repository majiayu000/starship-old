package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/majiayu000/cc-starship/pkg/logger"
)

// Logging creates a middleware that logs request information
func Logging(logManager *logger.LogManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 创建访问日志
		accessLog := &logger.AccessLog{
			BaseLog: logger.BaseLog{
				LogID:     logger.GenerateLogID(),
				Timestamp: time.Now(),
				Type:      logger.LogTypeAccess,
				Level:     logger.LogLevelInfo,
				Service:   "cc-starship",
				TraceID:   c.GetHeader("X-Trace-ID"),
				SpanID:    logger.GenerateSpanID(),
			},
			Request: logger.RequestInfo{
				Method:    c.Request.Method,
				Path:      path,
				Query:     raw,
				Headers:   logger.GetHeaders(c.Request.Header),
				ClientIP:  c.ClientIP(),
				UserAgent: c.Request.UserAgent(),
				Body:      logger.GetBody(c.Request),
			},
		}

		// 处理请求
		c.Next()

		// 更新响应信息
		accessLog.Response = logger.ResponseInfo{
			Status:  c.Writer.Status(),
			Headers: logger.GetHeaders(c.Writer.Header()),
			Size:    int64(c.Writer.Size()),
			TimeMS:  float64(time.Since(start).Milliseconds()),
		}

		// 发送日志
		logManager.Send(accessLog)
	}
}
