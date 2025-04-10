package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Age          int
	Anonymous    bool
	PasswordHash string
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Anonymous bool      `json:"anonymous"`
}

type UserRequest struct {
	Name      string `json:"name" validate:"required,gte=2"`
	Age       int    `json:"age" validate:"required"`
	Anonymous bool   `json:"anonymous"`
}

func (i UserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(i)
}
