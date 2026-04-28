package domain

import (
	"time"

	"github.com/google/uuid"
)

// Session represents an active user session.
type Session struct {
	ID               uuid.UUID `json:"id" db:"id"`
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	RefreshTokenHash string    `json:"-" db:"refresh_token_hash"`
	IPAddress        string    `json:"ip_address" db:"ip_address"`
	UserAgent        string    `json:"user_agent" db:"user_agent"`
	IsRevoked        bool      `json:"is_revoked" db:"is_revoked"`
	ExpiresAt        time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// TokenBlacklist represents a revoked JWT token.
type TokenBlacklist struct {
	JTI       string    `json:"jti" db:"jti"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// AuditLog represents an audit log entry.
type AuditLog struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	UserID      uuid.UUID              `json:"user_id" db:"user_id"`
	CompanyCNPJ string                 `json:"company_cnpj,omitempty" db:"company_cnpj"`
	Action      string                 `json:"action" db:"action"`
	Resource    string                 `json:"resource" db:"resource"`
	Details     map[string]interface{} `json:"details,omitempty" db:"details"`
	IPAddress   string                 `json:"ip_address" db:"ip_address"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}
