package main

import (
	"book_system/i18n"
	"book_system/internal/config"
	restapi "book_system/internal/controller/rest-api"
	"book_system/internal/infrastructure"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

var (
	minioClient *minio.Client
)

func init() {
	i18n.InitI18n([]string{"vi", "en"})

	// Create a new logger with the desired options
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	// Create a new logger with attributes
	logger = logger.With(
		slog.String("app", "book_system"),
		slog.String("version", "1.0.0"),
		slog.String("environment", config.MustGet().Environment),
	)

	// Set the default logger
	slog.SetDefault(logger)
	slog.Info("--------------------------------------------STARTING!--------------------------------------------")
}

func initDB() (*gorm.DB, error) {
	db, err := infrastructure.InitDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func initMinIO() error {
	cfg := config.MustGet()
	var err error

	// Initialize MinIO client
	endpoint := fmt.Sprintf("%s:%s", cfg.Minio.Host, cfg.Minio.Port)
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.Secure,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	// Check if the bucket exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := minioClient.BucketExists(ctx, cfg.Minio.DefaultBucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		slog.Info(fmt.Sprintf("Creating bucket: %s", cfg.Minio.DefaultBucket))
		err = minioClient.MakeBucket(ctx, cfg.Minio.DefaultBucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}

func main() {
	// Initialize database
	db, err := initDB()
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	// Initialize MinIO client
	if err := initMinIO(); err != nil {
		slog.Error("Failed to initialize MinIO", "error", err)
		os.Exit(1)
	}

	// Set Gin mode
	if config.MustGet().Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	router := gin.New()

	// Initialize API router
	apiRouter := restapi.NewRouter(
		minioClient,
		config.MustGet().Minio.DefaultBucket,
		config.MustGet().Minio.ReturnURL,
		db,
	)

	// Setup routes
	apiRouter.SetupRoutes(router)

	// Start server
	port := fmt.Sprintf(":%d", config.MustGet().Port)
	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Graceful shutdown
	go func() {
		slog.Info(fmt.Sprintf("Server is running on port %s", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown:", "error", err)
	}

	slog.Info("Server exiting")
}
