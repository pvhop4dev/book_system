package restapi

import (
	"book_system/internal/config"
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
	// Health check
	router.Use(middleware.GlobalRecover)
	router.GET("/health", HealthCheckHandler)

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
	userTransport := NewUserTransport(userService)
	bookTransport := NewBookTransport(bookService)
	uploadTransport := NewUploadTransport(uploadService)

	// Public routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", userTransport.Register)
			authGroup.POST("/login", userTransport.Login)
			authGroup.POST("/refresh", userTransport.RefreshToken)
		}

		// File upload routes
		filesGroup := v1.Group("/files")
		filesGroup.Use(middleware.AuthMiddleware(tokenSvc))
		{
			filesGroup.POST("/upload", uploadTransport.UploadFile)
			filesGroup.POST("/upload/multiple", uploadTransport.UploadMultipleFiles)
			filesGroup.GET("/:filename", uploadTransport.GetFile)
			filesGroup.DELETE("/:filename", uploadTransport.DeleteFile)
			filesGroup.GET("/:filename/url", uploadTransport.GetFileURL)
		}

		// User routes (protected)
		usersGroup := v1.Group("/users")
		usersGroup.Use(middleware.AuthMiddleware(tokenSvc))
		{
			usersGroup.GET("/me", userTransport.GetUserProfile)
			usersGroup.PUT("/me", userTransport.UpdateUserProfile)
			usersGroup.GET("", userTransport.ListUsers)
		}

		// Book routes (protected)
		booksGroup := v1.Group("/books")
		booksGroup.Use(middleware.AuthMiddleware(tokenSvc))
		{
			booksGroup.POST("", bookTransport.CreateBook)
			booksGroup.GET("", bookTransport.ListBooks)
			booksGroup.GET("/:id", bookTransport.GetBookByID)
			booksGroup.PUT("/:id", bookTransport.UpdateBook)
			booksGroup.DELETE("/:id", bookTransport.DeleteBook)
		}
	}
}
