package dto

import (
	"book_system/internal/infrastructure"
	"time"

	"github.com/google/uuid"
)

// BookResponse represents the book data sent in responses
type BookResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description,omitempty"`
	CoverImage  string    `json:"cover_image,omitempty"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	ISBN        string    `json:"isbn"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateBookRequest represents the data needed to create a new book
type CreateBookRequest struct {
	Title       string    `json:"title" validate:"required,min=1,max=255"`
	Author      string    `json:"author" validate:"required,min=1,max=255"`
	Description string    `json:"description"`
	CoverImage  string    `json:"cover_image"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Stock       int       `json:"stock" validate:"gte=0"`
	ISBN        string    `json:"isbn" validate:"required,isbn"`
	PublishedAt time.Time `json:"published_at" validate:"required"`
}

// Validate validates the CreateBookRequest
func (r *CreateBookRequest) Validate() error {
	return infrastructure.Validate.Struct(r)
}

// UpdateBookRequest represents the data needed to update a book
type UpdateBookRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Author      *string    `json:"author,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string    `json:"description,omitempty"`
	CoverImage  *string    `json:"cover_image,omitempty"`
	Price       *float64   `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock       *int       `json:"stock,omitempty" validate:"omitempty,gte=0"`
	ISBN        *string    `json:"isbn,omitempty" validate:"omitempty,isbn"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// Validate validates the UpdateBookRequest
func (r *UpdateBookRequest) Validate() error {
	return infrastructure.Validate.Struct(r)
}

// BookListResponse represents a paginated list of books
type BookListResponse struct {
	Data       []*BookResponse `json:"data"`
	Pagination Pagination      `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}
