package model

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name" validate:"required,min=3"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8"`
}
type Session struct {
	SessionID      string
	UserID         uuid.UUID
	ExpirationTime time.Time
}
