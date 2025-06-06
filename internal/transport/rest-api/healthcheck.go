package restapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheckResponse defines the structure of health check response
type HealthCheckResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Service   string                 `json:"service"`
	Version   string                 `json:"version"`
	Uptime    string                 `json:"uptime,omitempty"`
	Details   map[string]CheckDetail `json:"details,omitempty"`
}

// CheckDetail represents the status of a specific dependency
type CheckDetail struct {
	Status      string        `json:"status"`
	Latency     string        `json:"latency,omitempty"`
	Error       string        `json:"error,omitempty"`
	Version     string        `json:"version,omitempty"`
	Connections *int          `json:"connections,omitempty"`
	PingTime    time.Duration `json:"-"`
}

// HealthConfig holds the configuration for health check
type HealthConfig struct {
	DB          interface{}
	RedisClient interface{}
	MinioClient interface{}
	Version     string
	StartTime   time.Time
}

var healthConfig *HealthConfig

// InitHealthCheck initializes health check configuration
func InitHealthCheck(db, redisClient, minioClient interface{}, version string) {
	healthConfig = &HealthConfig{
		DB:          db,
		RedisClient: redisClient,
		MinioClient: minioClient,
		Version:     version,
		StartTime:   time.Now(),
	}
}

func SetupHealthCheckRoutes(healthGroup *gin.RouterGroup) {
	healthGroup.GET("/health", func(c *gin.Context) {
		healthCheckHandler(c, false, false)
	})
	healthGroup.GET("/health/live", func(c *gin.Context) {
		healthCheckHandler(c, true, false)
	})
	healthGroup.GET("/health/ready", func(c *gin.Context) {
		healthCheckHandler(c, false, true)
	})
}

// healthCheckHandler handles all health check requests
// @Summary      Health Check
// @Description  Check if the service is healthy
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthCheckResponse
// @Router       /health [get]
// @Router       /health/live [get]
// @Router       /health/ready [get]
func healthCheckHandler(c *gin.Context, livenessOnly bool, readinessOnly bool) {
	// Initialize response
	response := HealthCheckResponse{
		Status:    "UP",
		Timestamp: time.Now().UTC(),
		Service:   "book_system",
		Version:   healthConfig.Version,
		Details:   make(map[string]CheckDetail),
	}

	// Calculate uptime
	uptime := time.Since(healthConfig.StartTime)
	response.Uptime = uptime.String()

	// Check database connection if not liveness-only check
	if !livenessOnly && healthConfig.DB != nil {
		if db, ok := healthConfig.DB.(interface{ PingContext(context.Context) error }); ok {
			start := time.Now()
			detail := CheckDetail{Status: "UP"}
			err := db.PingContext(c.Request.Context())
			if err != nil {
				detail.Status = "DOWN"
				detail.Error = err.Error()
				response.Status = "DEGRADED"
			}
			detail.PingTime = time.Since(start)
			detail.Latency = detail.PingTime.String()
			response.Details["database"] = detail
		}
	}

	// Check Redis connection if not liveness-only check
	if !livenessOnly && healthConfig.RedisClient != nil {
		if rc, ok := healthConfig.RedisClient.(interface {
			Ping(context.Context) interface{ Result() (string, error) }
		}); ok {
			start := time.Now()
			detail := CheckDetail{Status: "UP"}
			_, err := rc.Ping(c.Request.Context()).Result()
			if err != nil {
				detail.Status = "DOWN"
				detail.Error = err.Error()
				response.Status = "DEGRADED"
			}
			detail.PingTime = time.Since(start)
			detail.Latency = detail.PingTime.String()
			response.Details["redis"] = detail
		}
	}

	// Check MinIO connection if not liveness-only check
	if !livenessOnly && healthConfig.MinioClient != nil {
		if mc, ok := healthConfig.MinioClient.(interface {
			ListBuckets(context.Context) (interface{}, error)
		}); ok {
			start := time.Now()
			detail := CheckDetail{Status: "UP"}
			_, err := mc.ListBuckets(c.Request.Context())
			if err != nil {
				detail.Status = "DOWN"
				detail.Error = err.Error()
				response.Status = "DEGRADED"
			}
			detail.PingTime = time.Since(start)
			detail.Latency = detail.PingTime.String()
			response.Details["storage"] = detail
		}
	}

	// For readiness check, only check critical dependencies
	if readinessOnly {
		if db, ok := healthConfig.DB.(interface{ PingContext(context.Context) error }); ok && db != nil {
			start := time.Now()
			detail := CheckDetail{Status: "UP"}
			err := db.PingContext(c.Request.Context())
			if err != nil {
				detail.Status = "DOWN"
				detail.Error = err.Error()
				response.Status = "DOWN"
			}
			detail.PingTime = time.Since(start)
			detail.Latency = detail.PingTime.String()
			response.Details["database"] = detail
		} else if healthConfig.DB != nil {
			detail := CheckDetail{
				Status: "DOWN",
				Error:  "Database client has no PingContext method",
			}
			response.Details["database"] = detail
			response.Status = "DOWN"
		}
	}

	// Set status code based on health status
	statusCode := http.StatusOK
	if response.Status == "DOWN" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}
