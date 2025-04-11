package models

import (
	"github.com/google/uuid"
	"time"
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

type WrapUser struct {
	User
	ExpiredAt time.Time
}
