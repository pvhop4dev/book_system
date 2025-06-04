package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Success(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, "success", data)
}

func Created(c *gin.Context, data interface{}) {
	JSON(c, http.StatusCreated, "created", data)
}

func BadRequest(c *gin.Context, message string) {
	JSON(c, http.StatusBadRequest, message, nil)
}

func Unauthorized(c *gin.Context, message string) {
	JSON(c, http.StatusUnauthorized, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	JSON(c, http.StatusForbidden, message, nil)
}

func NotFound(c *gin.Context, message string) {
	JSON(c, http.StatusNotFound, message, nil)
}

func InternalServerError(c *gin.Context, message string) {
	JSON(c, http.StatusInternalServerError, message, nil)
}
