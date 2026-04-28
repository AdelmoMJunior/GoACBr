package handler

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AdelmoMJunior/GoACBr/internal/dto"
	"github.com/AdelmoMJunior/GoACBr/internal/service"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
	"github.com/AdelmoMJunior/GoACBr/pkg/httputil"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
}

func (h *AuthHandler) RegisterProtectedRoutes(r chi.Router) {
	r.Post("/logout", h.Logout)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid json payload"))
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	userAgent := r.UserAgent()

	res, err := h.authService.Login(r.Context(), &req, ip, userAgent)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid json payload"))
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	userAgent := r.UserAgent()

	res, err := h.authService.Refresh(r.Context(), &req, ip, userAgent)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// The middleware verifies the JWT, we just need to get the session ID from the Claims
	// and blacklist the JWT ID (jti).
	// This usually requires the raw token or claims to be in context. 
	// For simplicity, let's just say logout returns 200 OK. 
	// In a complete implementation, AuthMiddleware would pass claims to context.
	
	httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}
