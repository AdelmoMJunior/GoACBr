package handler

import (
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AdelmoMJunior/GoACBr/internal/auth"
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
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest(err.Error()))
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
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest(err.Error()))
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
	jti, jtiOk := auth.GetJTI(r.Context())
	sessionID, sidOk := auth.GetSessionID(r.Context())
	expiresAt, expOk := auth.GetExpiresAt(r.Context())

	if !jtiOk || !sidOk || !expOk {
		httputil.SendError(w, apperror.NewUnauthorized("missing token claims for logout"))
		return
	}

	if err := h.authService.Logout(r.Context(), sessionID, jti, expiresAt); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}
