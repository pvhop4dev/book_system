package restapi

import (
	"book_system/internal/config"
	"book_system/internal/infrastructure"
	"book_system/internal/repository"
	book_service "book_system/internal/service/book_service"
	token_service "book_system/internal/service/token_service"
	upload_service "book_system/internal/service/upload_service"
	user_service "book_system/internal/service/user_service"
	"book_system/internal/transport/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	db *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	return &Router{
		db: db,
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	// Initialize health check with database connection
	db, err := r.db.DB()
	if err != nil {
		// Log error but continue, health check will show database as down
		gin.DefaultWriter.Write([]byte("Failed to get database connection for health check: " + err.Error() + "\n"))
	}

	// Get Redis client if available
	var redisClient interface{} = nil
	if r := infrastructure.GetRedis(); r != nil {
		redisClient = r
	}

	// Get MinIO client if available
	var minioClient interface{} = nil
	if m := infrastructure.GetMinioClient(); m != nil {
		minioClient = m
	}

	// Initialize health check with all dependencies
	healthConfig = &HealthConfig{
		DB:          db,
		RedisClient: redisClient,
		MinioClient: minioClient,
		Version:     "1.0.0",
		StartTime:   time.Now(),
	}

	// Global middleware
	router.Use(middleware.GlobalRecover)

	// Health check endpoints
	healthGroup := router.Group("")
	{
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

	// Swagger
	r.setupSwagger(router)

	// Initialize repositories
	userRepo := repository.NewUserRepository(r.db)
	bookRepo := repository.NewBookRepository(r.db)

	// Initialize services
	tokenSvc := token_service.NewTokenService(
		config.MustGet().JWT.AccessSecret,
		config.MustGet().JWT.RefreshSecret,
		time.Duration(config.MustGet().JWT.AccessExpiry)*time.Minute,
		time.Duration(config.MustGet().JWT.RefreshExpiry)*time.Minute,
	)

	userService := user_service.NewUserService(userRepo, tokenSvc)
	bookService := book_service.NewBookService(bookRepo)
	uploadService := upload_service.NewUploadService()

	// Initialize transports
	userController := NewUserController(userService)
	bookController := NewBookController(bookService)
	uploadController := NewUploadController(uploadService)

	// Public routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		authGroup := v1.Group("/auth")
		authGroup.Use(middleware.AuthMiddleware(tokenSvc))
		userController.SetupAuthRoutes(authGroup)

		// File upload routes
		filesGroup := v1.Group("/files")
		filesGroup.Use(middleware.AuthMiddleware(tokenSvc))
		uploadController.SetupUploadRoutes(filesGroup)

		// User routes (protected)
		usersGroup := v1.Group("/users")
		usersGroup.Use(middleware.AuthMiddleware(tokenSvc))
		userController.SetupUsersRoutes(usersGroup)

		// Book routes (protected)
		booksGroup := v1.Group("/books")
		booksGroup.Use(middleware.AuthMiddleware(tokenSvc))
		bookController.SetupBooksRoutes(booksGroup)
	}
}
