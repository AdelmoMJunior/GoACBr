package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

// SendJSON writes a JSON response with a specific status code.
func SendJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// SendError writes an error response, mapping apperror to appropriate HTTP status codes.
func SendError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if e, ok := err.(*apperror.AppError); ok {
		appErr = e
	} else {
		appErr = apperror.NewInternal(err)
	}

	status := appErr.Code
	if status == 0 {
		status = http.StatusInternalServerError
	}

	SendJSON(w, status, map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"message": appErr.Message,
			"code":    status,
			"detail":  appErr.Detail,
		},
	})
}
