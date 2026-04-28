package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
	"github.com/AdelmoMJunior/GoACBr/pkg/httputil"
)

// RequestID adds a unique trace ID to each request.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		
		ctx := context.WithValue(r.Context(), "request_id", reqID)
		w.Header().Set("X-Request-ID", reqID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger logs incoming HTTP requests and their processing time.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Use a response writer wrapper to capture status code
		ww := httputil.NewResponseWriterWrapper(w)

		next.ServeHTTP(ww, r)

		reqID, _ := r.Context().Value("request_id").(string)
		
		slog.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"duration", time.Since(start).String(),
			"ip", r.RemoteAddr,
			"req_id", reqID,
		)
	})
}

// Recovery catches panics and returns a 500 error instead of crashing the server.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				reqID, _ := r.Context().Value("request_id").(string)
				slog.Error("PANIC RECOVERED",
					"error", err,
					"req_id", reqID,
					"stack", string(debug.Stack()),
				)
				httputil.SendError(w, apperror.NewInternal(fmt.Errorf("internal server error (req_id: %s)", reqID)))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
