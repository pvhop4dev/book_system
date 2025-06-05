package service

import (
	"book_system/internal/model"
	"context"
	"mime/multipart"
)

type IUploadService interface {
	UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, customPath ...string) (string, error)
	DeleteFile(ctx context.Context, objectName string) error
	GetFileURL(ctx context.Context, objectName string) (string, error)
	GetFile(ctx context.Context, objectName string) (*multipart.FileHeader, error)
}

type IUserService interface {
	// User management
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
	UpdateUser(ctx context.Context, id string, req *model.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error

	// Authentication
	Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
	Register(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error)
	RefreshToken(ctx context.Context, token string) (*model.LoginResponse, error)
}

// TokenService defines the interface for token operations
type ITokenService interface {
	// GenerateToken generates a new access and refresh token pair
	GenerateToken(userID, role string) (*model.TokenPair, error)
	// ValidateToken validates a token and returns its claims
	ValidateToken(tokenString string) (*model.TokenClaims, error)
	// RefreshToken generates a new token pair using a refresh token
	RefreshToken(refreshToken string) (*model.TokenPair, error)
}

// IBookService defines the interface for book operations
type IBookService interface {
	// CreateBook creates a new book
	CreateBook(ctx context.Context, req *model.CreateBookRequest) (*model.BookResponse, error)
	// GetBookByID gets a book by ID
	GetBookByID(ctx context.Context, id string) (*model.BookResponse, error)
	// ListBooks gets a paginated list of books
	ListBooks(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*model.BookListResponse, error)
	// UpdateBook updates a book
	UpdateBook(ctx context.Context, id string, req *model.UpdateBookRequest) (*model.BookResponse, error)
	// DeleteBook deletes a book
	DeleteBook(ctx context.Context, id string) error
}
