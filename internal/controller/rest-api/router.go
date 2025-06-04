package restapi

import (
	"book_system/internal/config"
	"book_system/internal/repository"
	token_service "book_system/internal/service/token_service"
	upload_service "book_system/internal/service/upload_service"
	user_service "book_system/internal/service/user_service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type Router struct {
	minioClient   *minio.Client
	defaultBucket string
	returnURL     string
	db            *gorm.DB
}

func NewRouter(minioClient *minio.Client, defaultBucket, returnURL string, db *gorm.DB) *Router {
	return &Router{
		minioClient:   minioClient,
		defaultBucket: defaultBucket,
		returnURL:     returnURL,
		db:            db,
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	// Health check
	router.GET("/health", HealthCheckHandler)

	// Swagger
	r.setupSwagger(router)

	// Initialize repositories
	userRepo := repository.NewUserRepository(r.db)

	// Initialize services
	tokenSvc := token_service.NewTokenService(
		config.MustGet().JWT.AccessSecret,
		config.MustGet().JWT.RefreshSecret,
		time.Duration(config.MustGet().JWT.AccessExpiry)*time.Minute,
		time.Duration(config.MustGet().JWT.RefreshExpiry)*time.Minute,
	)

	userService := user_service.NewUserService(userRepo, tokenSvc)
	uploadService := upload_service.NewUploadService(r.minioClient, r.defaultBucket, r.returnURL)

	// Initialize controllers
	userController := NewUserController(userService)
	uploadController := NewUploadController(uploadService)

	// Public routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", userController.Register)
			authGroup.POST("/login", userController.Login)
			authGroup.POST("/refresh", userController.RefreshToken)
		}

		// File upload routes
		filesGroup := v1.Group("/files")
		filesGroup.Use(AuthMiddleware(tokenSvc))
		{
			filesGroup.POST("/upload", uploadController.UploadFile)
			filesGroup.POST("/upload/multiple", uploadController.UploadMultipleFiles)
			filesGroup.GET("/:filename", uploadController.GetFile)
			filesGroup.DELETE("/:filename", uploadController.DeleteFile)
			filesGroup.GET("/:filename/url", uploadController.GetFileURL)
		}

		// User routes (protected)
		usersGroup := v1.Group("/users")
		usersGroup.Use(AuthMiddleware(tokenSvc))
		{
			usersGroup.GET("/me", userController.GetUserProfile)
			usersGroup.PUT("/me", userController.UpdateUserProfile)
			usersGroup.GET("", userController.ListUsers)
		}
	}
}
