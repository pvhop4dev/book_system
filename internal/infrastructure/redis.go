package infrastructure

import (
	"book_system/internal/config"
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
	redisErr    error
)

// InitRedis initializes the Redis client with connection pooling
func InitRedis() (*redis.Client, error) {
	once.Do(func() {
		cfg := config.MustGet()

		// Create Redis options
		options := &redis.Options{
			Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
			Username: cfg.Redis.User,
			Password: cfg.Redis.Password,
			DB:       0, // Use default DB

			// Connection pool settings
			PoolSize:     10, // Maximum number of connections
			MinIdleConns: 5,  // Minimum number of idle connections

			// Timeouts
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolTimeout:  4 * time.Second,

			// Keep connections alive for 5 minutes
			ConnMaxIdleTime: 5 * time.Minute,
		}

		// Create new Redis client
		redisClient = redis.NewClient(options)

		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := redisClient.Ping(ctx).Result()
		if err != nil {
			redisErr = err
			slog.Error("Failed to connect to Redis", "error", err)
			return
		}

		slog.Info("Successfully connected to Redis",
			slog.String("host", cfg.Redis.Host),
			slog.String("port", cfg.Redis.Port))
	})

	return redisClient, redisErr
}

// GetRedis returns the Redis client instance
// Panics if Redis is not initialized
func GetRedis() *redis.Client {
	if redisClient == nil {
		slog.Error("Redis client not initialized, call InitRedis() first")
		panic("redis client not initialized")
	}
	return redisClient
}

// CloseRedis closes the Redis connection
func CloseRedis() {
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			slog.Error("Error closing Redis connection", "error", err)
			return
		}
		slog.Info("Closed Redis connection")
	}
}

// HealthCheck verifies the Redis connection is still alive
func HealthCheck() error {
	if redisClient == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	return err
}
