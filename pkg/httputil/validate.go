package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// validate is the shared validator instance (thread-safe after init).
var validate = validator.New()

// DecodeAndValidate decodes JSON from the request body into dest and
// validates it using go-playground/validator struct tags.
// Returns nil on success or an AppError on failure.
func DecodeAndValidate(r *http.Request, dest interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("invalid json payload: %w", err)
	}

	if err := validate.Struct(dest); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var msgs []string
			for _, fe := range validationErrors {
				msgs = append(msgs, formatValidationError(fe))
			}
			return fmt.Errorf("validation failed: %s", strings.Join(msgs, "; "))
		}
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

func formatValidationError(fe validator.FieldError) string {
	field := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	case "gtefield":
		return fmt.Sprintf("%s must be >= %s", field, fe.Param())
	default:
		return fmt.Sprintf("%s failed validation (%s=%s)", field, fe.Tag(), fe.Param())
	}
}
