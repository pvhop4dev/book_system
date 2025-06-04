package repository

import (
	"book_system/internal/model/entity"
	"context"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type IUserRepository interface {
	// Create saves a new user
	Create(ctx context.Context, user *entity.User) error

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindAll returns a paginated list of users
	FindAll(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error)

	// Update updates a user
	Update(ctx context.Context, user *entity.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
