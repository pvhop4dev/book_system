package service

import (
	"book_system/internal/model/dto"
	"book_system/internal/model/entity"
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
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.User, error)
	GetUserByID(ctx context.Context, id string) (*dto.User, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.User, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]*dto.User, int64, error)
	UpdateUser(ctx context.Context, id string, req *dto.UpdateUserRequest) (*dto.User, error)
	DeleteUser(ctx context.Context, id string) error

	// Authentication
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	RefreshToken(ctx context.Context, token string) (*dto.LoginResponse, error)
}

// TokenService defines the interface for token operations
type ITokenService interface {
	// GenerateToken generates a new access and refresh token pair
	GenerateToken(userID, role string) (*entity.TokenPair, error)
	// ValidateToken validates a token and returns its claims
	ValidateToken(tokenString string) (*entity.TokenClaims, error)
	// RefreshToken generates a new token pair using a refresh token
	RefreshToken(refreshToken string) (*entity.TokenPair, error)
}
