package entity

import (
	"book_system/internal/model/dto"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Username  string    `gorm:"size:100;not null;uniqueIndex"`
	Email     string    `gorm:"size:100;not null;uniqueIndex"`
	Password  string    `gorm:"size:255;not null"`
	FullName  string    `gorm:"size:100;not null"`
	Role      string    `gorm:"size:20;not null;default:'user'"`
	Avatar    string    `gorm:"size:255"`
	IsActive  bool      `gorm:"not null;default:true"`
	LastLogin time.Time
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// ToDTO converts User entity to User DTO
func (u *User) ToDTO() *dto.User {
	return &dto.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Role:      u.Role,
		Avatar:    u.Avatar,
		IsActive:  u.IsActive,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
