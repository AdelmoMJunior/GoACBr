package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AdelmoMJunior/GoACBr/internal/auth"
	"github.com/AdelmoMJunior/GoACBr/internal/dto"
	"github.com/AdelmoMJunior/GoACBr/internal/middleware"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
	"github.com/AdelmoMJunior/GoACBr/internal/service"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
	"github.com/AdelmoMJunior/GoACBr/pkg/httputil"
)

type CompanyHandler struct {
	companyService service.CompanyService
	compRepo       repository.CompanyRepository
}

func NewCompanyHandler(companyService service.CompanyService, compRepo repository.CompanyRepository) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
		compRepo:       compRepo,
	}
}

func (h *CompanyHandler) RegisterRoutes(r chi.Router) {
	// Protected
	r.Group(func(r chi.Router) {
		r.Get("/companies", h.List)
		r.Post("/companies", h.Create)

		// Routes requiring specific company access
		r.Group(func(r chi.Router) {
			r.Use(middleware.CompanyGuard(h.compRepo))

			r.Get("/companies/{company_id}", h.Get)
			r.Put("/companies/{company_id}", h.Update)
			r.Put("/companies/{company_id}/smtp", h.ConfigureSMTP)
		})
	})
}

func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewUnauthorized("user not authenticated"))
		return
	}

	res, err := h.companyService.ListByUser(r.Context(), userID)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewUnauthorized("user not authenticated"))
		return
	}

	var req dto.CompanyCreateRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest(err.Error()))
		return
	}

	res, err := h.companyService.Create(r.Context(), userID, &req)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusCreated, res)
}

func (h *CompanyHandler) Get(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	res, err := h.companyService.GetByID(r.Context(), companyID)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	var req dto.CompanyCreateRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest(err.Error()))
		return
	}

	res, err := h.companyService.Update(r.Context(), companyID, &req)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *CompanyHandler) ConfigureSMTP(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	var req dto.CompanySMTPRequest
	if err := httputil.DecodeAndValidate(r, &req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest(err.Error()))
		return
	}

	if err := h.companyService.ConfigureSMTP(r.Context(), companyID, &req); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "SMTP configured successfully"})
}
