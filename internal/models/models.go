package models

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Anonymous bool      `json:"anonymous"`
}

type userInput struct {
	Name      string `json:"name" validate:"required,gte=2"`
	Age       int    `json:"age" validate:"required"`
	Anonymous bool   `json:"anonymous"`
}
