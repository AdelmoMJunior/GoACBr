package dto

import (
	"time"

	"github.com/google/uuid"
)

// CertificateResponse is the standard representation of a digital certificate.
type CertificateResponse struct {
	ID              uuid.UUID `json:"id"`
	CompanyID       uuid.UUID `json:"company_id"`
	SubjectCN       string    `json:"subject_cn"`
	SerialNumber    string    `json:"serial_number"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidUntil      time.Time `json:"valid_until"`
	DaysUntilExpiry int       `json:"days_until_expiry"`
	IsExpired       bool      `json:"is_expired"`
	CreatedAt       time.Time `json:"created_at"`
}

// CertificateUploadRequest (handled via multipart/form-data in HTTP handler)
// We just need the password here.
type CertificateUploadRequest struct {
	Password string `form:"password" binding:"required"`
}
