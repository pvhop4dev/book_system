package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
// @name Response
// @Description Standard API response format
type Response struct {
	Code    int         `json:"code" swaggertype:"integer"`
	Message string      `json:"message" swaggertype:"string"`
	Data    any `json:"data,omitempty" swaggertype:"object"`
}

func JSON(c *gin.Context, code int, message string, data any) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Success(c *gin.Context, data any) {
	JSON(c, http.StatusOK, "success", data)
}

func Created(c *gin.Context, data any) {
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
