package restapi

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheckResponse represents the health check response structure
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// HealthCheckHandler handles health check requests
// @Summary      Health Check
// @Description  Check if the service is healthy
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthCheckResponse
// @Router       /health [get]
func HealthCheckHandler(c *gin.Context) {
	response := HealthCheckResponse{
		Status:    "UP",
		Timestamp: time.Now().UTC(),
		Service:   "book_system",
	}
	c.JSON(http.StatusOK, response)
}
