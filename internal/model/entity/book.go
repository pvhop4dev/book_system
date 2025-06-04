package entity

import (
	"book_system/internal/model/dto"
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Title       string    `gorm:"size:255;not null"`
	Author      string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	CoverImage  string    `gorm:"size:512"`
	Price       float64   `gorm:"type:decimal(10,2);not null"`
	Stock       int       `gorm:"not null;default:0"`
	ISBN        string    `gorm:"size:20;uniqueIndex"`
	PublishedAt time.Time `gorm:"type:date"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (Book) TableName() string {
	return "books"
}

// ToDTO converts Book entity to Book DTO
func (b *Book) ToDTO() *dto.BookResponse {
	return &dto.BookResponse{
		ID:          b.ID,
		Title:       b.Title,
		Author:      b.Author,
		Description: b.Description,
		CoverImage:  b.CoverImage,
		Price:       b.Price,
		Stock:       b.Stock,
		ISBN:        b.ISBN,
		PublishedAt: b.PublishedAt,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}
