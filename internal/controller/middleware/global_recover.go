package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func globalRecover(c *gin.Context) {
	defer func(c *gin.Context) {
		if rec := recover(); rec != nil {
			// that recovery also handle XHR's
			// you need handle it
			if XHR(c) {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": rec,
				})
			} else {
				c.HTML(http.StatusOK, "500", gin.H{})
			}
		}
	}(c)
	c.Next()
}
func XHR(c *gin.Context) bool {
	return strings.ToLower(c.Request.Header.Get("X-Requested-With")) == "xmlhttprequest"
}
