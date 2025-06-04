package dto

import (
	"book_system/internal/infrastructure"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username" validate:"required,min=3,max=50"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"-"`
	FullName  string    `json:"full_name" validate:"required"`
	Role      string    `json:"role" validate:"oneof=admin user"`
	Avatar    string    `json:"avatar,omitempty"`
	IsActive  bool      `json:"is_active"`
	LastLogin time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
	Role     string `json:"role" validate:"oneof=admin user"`
}

func (r *CreateUserRequest) Validate() error {
	return infrastructure.Validate.Struct(r)
}

type UpdateUserRequest struct {
	FullName *string `json:"full_name,omitempty"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Role     *string `json:"role,omitempty" validate:"omitempty,oneof=admin user"`
	IsActive *bool   `json:"is_active,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

func (r *UpdateUserRequest) Validate() error {
	return infrastructure.Validate.Struct(r)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r *LoginRequest) Validate() error {
	return infrastructure.Validate.Struct(r)
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

func (r *RegisterRequest) Validate() error {
	return infrastructure.Validate.Struct(r)
}

type RegisterResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
