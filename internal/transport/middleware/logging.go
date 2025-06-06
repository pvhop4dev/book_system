package middleware

import (
	"book_system/internal/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func defaultLogFormatter(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[GIN] |%s | %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		param.Request.Context().Value(utils.TraceIDContextKey),
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

func CustomLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: defaultLogFormatter,
		Output:    gin.DefaultWriter,
		SkipPaths: []string{"/health/live", "/health/ready"},
	})
}
