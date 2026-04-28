package dto

import "time"

// LoginRequest is the request payload for login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the response payload for login.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}

// RefreshRequest is the request payload to refresh tokens.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest is the payload to change user password.
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// TokenData is the payload returned when decoding a JWT.
type TokenData struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	CompanyID string    `json:"company_id,omitempty"`
	ExpiresAt time.Time `json:"expires_at"`
}
