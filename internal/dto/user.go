package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserCreateRequest is the request payload to create a new user.
type UserCreateRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
	Phone    string `json:"phone,omitempty"`
}

// UserUpdateRequest is the request payload to update an existing user.
type UserUpdateRequest struct {
	FullName string `json:"full_name,omitempty"`
	Phone    string `json:"phone,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// UserResponse is the standard representation of a user in API responses.
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
