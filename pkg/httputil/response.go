package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

// Response is the standard API response envelope.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorBody represents the error details in the response.
type ErrorBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Meta contains pagination metadata.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		resp := Response{
			Success: status >= 200 && status < 300,
			Data:    data,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// JSONWithMeta writes a JSON response with pagination metadata.
func JSONWithMeta(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	resp := Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// Error writes a JSON error response.
func Error(w http.ResponseWriter, err *apperror.AppError) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(err.Code)
	resp := Response{
		Success: false,
		Error: &ErrorBody{
			Code:    err.Code,
			Message: err.Message,
			Detail:  err.Detail,
		},
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// ErrorFromStatus writes a simple error response from an HTTP status code and message.
func ErrorFromStatus(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	resp := Response{
		Success: false,
		Error: &ErrorBody{
			Code:    status,
			Message: message,
		},
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// NoContent writes a 204 No Content response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Created writes a 201 Created response with data.
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

// DecodeJSON decodes a JSON request body into the given destination.
func DecodeJSON(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}
