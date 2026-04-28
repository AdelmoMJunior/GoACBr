package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors for common cases.
var (
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrValidation        = errors.New("validation error")
	ErrInternal          = errors.New("internal server error")
	ErrCertificateExpired = errors.New("certificate expired")
	ErrCertificateNotFound = errors.New("certificate not found for this company")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrDistributionCooldown = errors.New("distribution query cooldown not elapsed")
	ErrRedisUnavailable  = errors.New("redis unavailable, using fallback")
)

// AppError is the standard error type returned by the API.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Constructor helpers.

func NewNotFound(resource string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Err:     ErrNotFound,
	}
}

func NewAlreadyExists(resource string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: fmt.Sprintf("%s already exists", resource),
		Err:     ErrAlreadyExists,
	}
}

func NewUnauthorized(detail string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
		Detail:  detail,
		Err:     ErrUnauthorized,
	}
}

func NewForbidden(detail string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: "forbidden",
		Detail:  detail,
		Err:     ErrForbidden,
	}
}

func NewValidation(detail string) *AppError {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "validation error",
		Detail:  detail,
		Err:     ErrValidation,
	}
}

func NewInternal(err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
		Err:     fmt.Errorf("%w: %v", ErrInternal, err),
	}
}

func NewCertificateExpired(cnpj string, detail string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("O certificado digital do CNPJ %s está vencido. Por favor, faça o upload de um novo certificado válido.", cnpj),
		Detail:  detail,
		Err:     ErrCertificateExpired,
	}
}

func NewCertificateNotFound(cnpj string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("Nenhum certificado digital encontrado para o CNPJ %s. Faça o upload do certificado PFX antes de realizar esta operação.", cnpj),
		Err:     ErrCertificateNotFound,
	}
}

func NewRateLimitExceeded() *AppError {
	return &AppError{
		Code:    http.StatusTooManyRequests,
		Message: "too many requests, please try again later",
		Err:     ErrRateLimitExceeded,
	}
}

func NewDistributionCooldown(minutesRemaining int) *AppError {
	return &AppError{
		Code:    http.StatusTooManyRequests,
		Message: fmt.Sprintf("A consulta de distribuição DFe só pode ser realizada a cada 1 hora. Tente novamente em %d minutos.", minutesRemaining),
		Err:     ErrDistributionCooldown,
	}
}

func NewBadRequest(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}
