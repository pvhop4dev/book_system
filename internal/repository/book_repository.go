package repository

import (
	"book_system/internal/model/entity"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

// NewBookRepository creates a new book repository
func NewBookRepository(db *gorm.DB) IBookRepository {
	return &bookRepository{
		db: db,
	}
}

// Create saves a new book
func (r *bookRepository) Create(ctx context.Context, book *entity.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

// FindByID finds a book by ID
func (r *bookRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Book, error) {
	var book entity.Book
	err := r.db.WithContext(ctx).First(&book, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// FindAll returns a paginated list of books
func (r *bookRepository) FindAll(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*entity.Book, int64, error) {
	var books []*entity.Book
	var count int64

	offset := (page - 1) * pageSize

	// Start building the query
	query := r.db.WithContext(ctx).Model(&entity.Book{})

	// Apply filters if any
	for key, value := range filters {
		query = query.Where(key, value)
	}

	// Get total count
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated books
	if err := query.Offset(offset).
		Limit(pageSize).
		Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

// Update updates a book
func (r *bookRepository) Update(ctx context.Context, book *entity.Book) error {
	return r.db.WithContext(ctx).Save(book).Error
}

// Delete deletes a book by ID
func (r *bookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Book{}, "id = ?", id).Error
}

// ExistsByISBN checks if a book with the given ISBN exists
func (r *bookRepository) ExistsByISBN(ctx context.Context, isbn string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Book{}).
		Where("isbn = ?", isbn).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
