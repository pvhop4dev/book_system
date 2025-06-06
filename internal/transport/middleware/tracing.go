package middleware

import (
	"book_system/internal/utils"

	"github.com/gin-gonic/gin"
)

func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.SetTraceID(c, utils.RandomTraceID())
		c.Next()
	}
}
