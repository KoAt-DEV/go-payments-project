package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email" validate:"required,email"`
	PasswordHash string    `json:"password"`
	Role         string    `json:"role"`
	CoffeeCount  int       `json:"coffee_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
