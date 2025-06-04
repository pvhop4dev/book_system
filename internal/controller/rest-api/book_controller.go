package restapi

import (
	"book_system/internal/model/dto"
	"book_system/internal/service"
	"book_system/pkg/response"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BookController handles book related HTTP requests
type BookController struct {
	bookService service.IBookService
}

// NewBookController creates a new book controller
func NewBookController(bookService service.IBookService) *BookController {
	return &BookController{
		bookService: bookService,
	}
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book with the input payload
// @Tags books
// @Accept  json
// @Produce  json
// @Param input body dto.CreateBookRequest true "Book data"
// @Success 201 {object} response.APIResponse{data=dto.BookResponse} "Successfully created book"
// @Failure 400 {object} response.APIResponse "Invalid input"
// @Failure 409 {object} response.APIResponse "Book with this ISBN already exists"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /api/v1/books [post]
func (c *BookController) CreateBook(ctx *gin.Context) {
	var req dto.CreateBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body")
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	book, err := c.bookService.CreateBook(ctx.Request.Context(), &req)
	if err != nil {
		if err.Error() == "book with this ISBN already exists" {
			response.JSON(ctx, http.StatusConflict, err.Error(), nil)
			return
		}
		slog.Error("Failed to create book", slog.Any("error", err))
		response.InternalServerError(ctx, "Failed to create book")
		return
	}

	response.Created(ctx, book)
}

// GetBookByID godoc
// @Summary Get a book by ID
// @Description Get a book by its ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {object} response.APIResponse{data=dto.BookResponse} "Successfully retrieved book"
// @Failure 400 {object} response.APIResponse "Invalid book ID"
// @Failure 404 {object} response.APIResponse "Book not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /api/v1/books/{id} [get]
func (c *BookController) GetBookByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.BadRequest(ctx, "Book ID is required")
		return
	}

	book, err := c.bookService.GetBookByID(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "book not found" {
			response.NotFound(ctx, "Book not found")
			return
		}
		slog.Error("Failed to get book", slog.Any("error", err))
		response.InternalServerError(ctx, "Failed to get book")
		return
	}

	response.Success(ctx, book)
}

// ListBooks godoc
// @Summary List all books with pagination
// @Description Get a paginated list of books with optional filters
// @Tags books
// @Accept  json
// @Produce  json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10, max: 100)"
// @Param author query string false "Filter by author"
// @Success 200 {object} response.APIResponse{data=dto.BookListResponse} "Successfully retrieved books"
// @Failure 400 {object} response.APIResponse "Invalid query parameters"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /api/v1/books [get]
func (c *BookController) ListBooks(ctx *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	// Build filters
	filters := make(map[string]interface{})
	if author := ctx.Query("author"); author != "" {
		filters["author = ?"] = author
	}

	// Get books from service
	result, err := c.bookService.ListBooks(ctx.Request.Context(), page, pageSize, filters)
	if err != nil {
		slog.Error("Failed to list books", slog.Any("error", err))
		response.InternalServerError(ctx, "Failed to list books")
		return
	}

	response.Success(ctx, result)
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update a book with the input payload
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Param book body dto.UpdateBookRequest true "Update book"
// @Success 200 {object} response.APIResponse{data=dto.BookResponse} "Successfully updated book"
// @Failure 400 {object} response.APIResponse "Invalid input"
// @Failure 404 {object} response.APIResponse "Book not found"
// @Failure 409 {object} response.APIResponse "Book with this ISBN already exists"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /api/v1/books/{id} [put]
func (c *BookController) UpdateBook(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.BadRequest(ctx, "Book ID is required")
		return
	}

	var req dto.UpdateBookRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body")
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	book, err := c.bookService.UpdateBook(ctx.Request.Context(), id, &req)
	if err != nil {
		switch err.Error() {
		case "book not found":
			response.JSON(ctx, http.StatusNotFound, err.Error(), nil)
		case "book with this ISBN already exists":
			response.JSON(ctx, http.StatusConflict, err.Error(), nil)
		default:
			slog.Error("Failed to update book", slog.Any("error", err))
			response.InternalServerError(ctx, "Failed to update book")
		}
		return
	}

	response.Success(ctx, book)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book by its ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {object} response.APIResponse "Successfully deleted book"
// @Failure 400 {object} response.APIResponse "Invalid book ID"
// @Failure 404 {object} response.APIResponse "Book not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /api/v1/books/{id} [delete]
func (c *BookController) DeleteBook(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.BadRequest(ctx, "Book ID is required")
		return
	}

	err := c.bookService.DeleteBook(ctx.Request.Context(), id)
	if err != nil {
		if err.Error() == "book not found" {
			response.NotFound(ctx, "Book not found")
			return
		}
		slog.Error("Failed to delete book", slog.Any("error", err))
		response.InternalServerError(ctx, "Failed to delete book")
		return
	}

	response.Success(ctx, nil)
}
