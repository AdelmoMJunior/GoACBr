package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AdelmoMJunior/GoACBr/internal/auth"
	"github.com/AdelmoMJunior/GoACBr/internal/dto"
	"github.com/AdelmoMJunior/GoACBr/internal/service"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
	"github.com/AdelmoMJunior/GoACBr/pkg/httputil"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/register", h.Register)
}

func (h *UserHandler) RegisterProtectedRoutes(r chi.Router) {
	r.Get("/me", h.GetMe)
	r.Put("/me", h.UpdateMe)
	r.Put("/me/password", h.ChangePassword)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid json payload"))
		return
	}

	res, err := h.userService.Create(r.Context(), &req)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusCreated, res)
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewUnauthorized("user not authenticated"))
		return
	}

	res, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewUnauthorized("user not authenticated"))
		return
	}

	var req dto.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid json payload"))
		return
	}

	res, err := h.userService.Update(r.Context(), userID, &req)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewUnauthorized("user not authenticated"))
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid json payload"))
		return
	}

	if err := h.userService.ChangePassword(r.Context(), userID, &req); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "password updated successfully"})
}
