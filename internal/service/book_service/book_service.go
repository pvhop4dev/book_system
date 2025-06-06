package book_service

import (
	"book_system/internal/model"
	"book_system/internal/repository"
	"book_system/internal/service"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type bookService struct {
	repo repository.IBookRepository
}

// NewBookService creates a new book service
func NewBookService(repo repository.IBookRepository) service.IBookService {
	return &bookService{
		repo: repo,
	}
}

// CreateBook creates a new book
func (s *bookService) CreateBook(ctx context.Context, req *model.CreateBookRequest) (*model.BookResponse, error) {
	// Check if book with same ISBN already exists
	exists, err := s.repo.ExistsByISBN(ctx, req.ISBN)
	if err != nil {
		return nil, fmt.Errorf("failed to check book existence: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("book with ISBN %s already exists", req.ISBN)
	}

	// Create new book entity
	book := &model.Book{
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Price:       req.Price,
		Stock:       req.Stock,
		ISBN:        req.ISBN,
		PublishedAt: req.PublishedAt,
	}

	// Save to database
	if err := s.repo.Create(ctx, book); err != nil {
		return nil, fmt.Errorf("failed to create book: %v", err)
	}

	return book.ToDTO(), nil
}

// GetBookByID gets a book by ID
func (s *bookService) GetBookByID(ctx context.Context, id string) (*model.BookResponse, error) {
	bookID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID format: %v", err)
	}

	book, err := s.repo.FindByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("book not found: %v", err)
	}

	return book.ToDTO(), nil
}

// ListBooks gets a paginated list of books
func (s *bookService) ListBooks(ctx context.Context, page, pageSize int, filters map[string]any) (*model.BookListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	books, total, err := s.repo.FindAll(ctx, page, pageSize, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list books: %v", err)
	}

	// Convert to DTOs
	bookDTOs := make([]*model.BookResponse, len(books))
	for i, book := range books {
		bookDTOs[i] = book.ToDTO()
	}

	totalPage := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &model.BookListResponse{
		Data: bookDTOs,
		Pagination: model.Pagination{
			Page:      page,
			PageSize:  pageSize,
			Total:     total,
			TotalPage: totalPage,
		},
	}, nil
}

// UpdateBook updates a book
func (s *bookService) UpdateBook(ctx context.Context, id string, req *model.UpdateBookRequest) (*model.BookResponse, error) {
	bookID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid book ID format: %v", err)
	}

	// Get existing book
	book, err := s.repo.FindByID(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("book not found: %v", err)
	}

	// Update fields if provided
	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Author != nil {
		book.Author = *req.Author
	}
	if req.Description != nil {
		book.Description = *req.Description
	}
	if req.CoverImage != nil {
		book.CoverImage = *req.CoverImage
	}
	if req.Price != nil {
		book.Price = *req.Price
	}
	if req.Stock != nil {
		book.Stock = *req.Stock
	}
	if req.ISBN != nil && *req.ISBN != book.ISBN {
		// Check if new ISBN already exists
		exists, err := s.repo.ExistsByISBN(ctx, *req.ISBN)
		if err != nil {
			return nil, fmt.Errorf("failed to check ISBN existence: %v", err)
		}
		if exists {
			return nil, fmt.Errorf("book with ISBN %s already exists", *req.ISBN)
		}
		book.ISBN = *req.ISBN
	}
	if req.PublishedAt != nil {
		book.PublishedAt = *req.PublishedAt
	}

	// Save updates
	if err := s.repo.Update(ctx, book); err != nil {
		return nil, fmt.Errorf("failed to update book: %v", err)
	}

	return book.ToDTO(), nil
}

// DeleteBook deletes a book
func (s *bookService) DeleteBook(ctx context.Context, id string) error {
	bookID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid book ID format: %v", err)
	}

	// Check if book exists
	_, err = s.repo.FindByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("book not found: %v", err)
	}

	// Delete book
	if err := s.repo.Delete(ctx, bookID); err != nil {
		return fmt.Errorf("failed to delete book: %v", err)
	}

	return nil
}
