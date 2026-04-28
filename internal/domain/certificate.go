package domain

import (
	"time"

	"github.com/google/uuid"
)

// Certificate represents a digital certificate (PFX) associated with a company.
type Certificate struct {
	ID             uuid.UUID `json:"id" db:"id"`
	CompanyID      uuid.UUID `json:"company_id" db:"company_id"`
	PFXData        []byte    `json:"-" db:"pfx_data"`             // AES-256-GCM encrypted
	PFXPasswordEnc string    `json:"-" db:"pfx_password_enc"`     // AES-256-GCM encrypted
	SubjectCN      string    `json:"subject_cn" db:"subject_cn"`
	SerialNumber   string    `json:"serial_number" db:"serial_number"`
	ValidFrom      time.Time `json:"valid_from" db:"valid_from"`
	ValidUntil     time.Time `json:"valid_until" db:"valid_until"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// IsExpired checks if the certificate has expired.
func (c *Certificate) IsExpired() bool {
	return time.Now().After(c.ValidUntil)
}

// DaysUntilExpiry returns the number of days until the certificate expires.
func (c *Certificate) DaysUntilExpiry() int {
	return int(time.Until(c.ValidUntil).Hours() / 24)
}
