package repository

import (
	"book_system/internal/model"
	"context"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data operations
type IUserRepository interface {
	// Create saves a new user
	Create(ctx context.Context, user *model.User) error

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// FindAll returns a paginated list of users
	FindAll(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)

	// Update updates a user
	Update(ctx context.Context, user *model.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// IBookRepository defines the interface for book data operations
type IBookRepository interface {
	// Create saves a new book
	Create(ctx context.Context, book *model.Book) error

	// FindByID finds a book by ID
	FindByID(ctx context.Context, id uuid.UUID) (*model.Book, error)

	// FindAll returns a paginated list of books
	FindAll(ctx context.Context, page, pageSize int, filters map[string]any) ([]*model.Book, int64, error)

	// Update updates a book
	Update(ctx context.Context, book *model.Book) error

	// Delete deletes a book by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// ExistsByISBN checks if a book with the given ISBN exists
	ExistsByISBN(ctx context.Context, isbn string) (bool, error)
}
