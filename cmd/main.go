package main

import (
	"book_system/i18n"
	"book_system/internal/config"
	"book_system/internal/infrastructure"
	restapi "book_system/internal/transport/rest-api"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// var (
// 	minioClient *minio.Client
// )

func init() {
	i18n.InitI18n([]string{"vi", "en"})

	// Create a new JSON logger with custom source handler
	handler := &infrastructure.CustomSourceHandler{
		Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: false, // Tắt source mặc định
		}),
	}

	// Add attributes to the logger
	logger := slog.New(handler).With(
		slog.String("app", "book_system"),
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

func main() {
	// Initialize database
	db, err := initDB()
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
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
