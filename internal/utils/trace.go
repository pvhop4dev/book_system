package utils

import (
	"context"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type ContextKey string

const TraceIDContextKey ContextKey = "trace_id"

func GetCurrentTraceID(c *gin.Context) string {
	traceID, ok := c.Get(string(TraceIDContextKey))
	if ok {
		return traceID.(string)
	}
	return ""
}
func SetTraceID(c *gin.Context, traceID string) {
	c.Set(string(TraceIDContextKey), traceID)
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), TraceIDContextKey, traceID))
}

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
func RandomTraceID() string {
	return stringWithCharset(16, base62Chars)
}
