package middleware

import (
	"book_system/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authorizator(authority ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorities := c.Request.Header["Authorities"]
		if len(authorities) == 0 || !utils.AnyContains(strings.Split(authorities[0], ","), authority) {
			err := Forbidden
			c.JSON(http.StatusForbidden, gin.H{
				"code":    err.Code,
				"message": err.GetMesssageI18n(utils.GetCurrentLang(c)),
			})
			c.Abort()
			return
		}
	}
}
