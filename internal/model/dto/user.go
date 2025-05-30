package dto

import "book_system/internal/infrastructure"

type LoginRequest struct {
}

func (r *LoginRequest) Validate() error {
	return infrastructure.Validate.Struct(r)
}

type LoginResponse struct {
}

type RegisterRequest struct {
}

func (r *RegisterRequest) Validate() error {
	return infrastructure.Validate.Struct(r)
}

type RegisterResponse struct {
}
