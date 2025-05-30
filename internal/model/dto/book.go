package dto

import (
	"book_system/internal/infrastructure"
)

type CreateBookRequest struct {
}

func (r *CreateBookRequest) Validate() error {
	return infrastructure.Validate.Struct(r)
}

type CreateBookResponse struct {
}
