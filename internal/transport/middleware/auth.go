package middleware

import (
	"book_system/internal/config"
	"book_system/internal/service"
	"book_system/internal/utils"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
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

var e *casbin.Enforcer

func init() {
	adapter, err := xormadapter.NewAdapter("mysql", config.MustGet().Casbin.DSN)
	if err != nil {
		panic(err)
	}
	e, err = casbin.NewEnforcer("casbin/model.conf", adapter)
	if err != nil {
		panic(err)
	}
	err = e.LoadPolicy()
	if err != nil {
		panic(err)
	}
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		sub := c.GetString("user")
		obj := c.Request.URL.Path
		act := c.Request.Method

		ok, err := e.Enforce(sub, obj, act)
		if err != nil {
			errType := InternalServerError
			c.JSON(http.StatusInternalServerError, gin.H{"code": errType.Code, "message": errType.GetMesssageI18n(utils.GetCurrentLang(c))})
			c.Abort()
			return
		}

		if !ok {
			errType := Forbidden
			c.JSON(http.StatusForbidden, gin.H{"code": errType.Code, "message": errType.GetMesssageI18n(utils.GetCurrentLang(c))})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthMiddleware creates a Gin middleware for JWT authentication
func AuthMiddleware(tokenService service.ITokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract the token from the header (format: "Bearer <token>")
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate the token
		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Add user ID to context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}
