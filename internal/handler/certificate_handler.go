package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AdelmoMJunior/GoACBr/internal/auth"
	"github.com/AdelmoMJunior/GoACBr/internal/middleware"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
	"github.com/AdelmoMJunior/GoACBr/internal/service"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
	"github.com/AdelmoMJunior/GoACBr/pkg/httputil"
)

type CertificateHandler struct {
	certService service.CertificateService
	compRepo    repository.CompanyRepository
}

func NewCertificateHandler(certService service.CertificateService, compRepo repository.CompanyRepository) *CertificateHandler {
	return &CertificateHandler{
		certService: certService,
		compRepo:    compRepo,
	}
}

func (h *CertificateHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.CompanyGuard(h.compRepo))

		r.Get("/certificates", h.Get)
		r.Post("/certificates", h.Upload)
		r.Delete("/certificates", h.Delete)
	})
}

func (h *CertificateHandler) Get(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	res, err := h.certService.Get(r.Context(), companyID)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *CertificateHandler) Upload(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	// Handle multipart form data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		httputil.SendError(w, apperror.NewBadRequest("failed to parse form data"))
		return
	}

	password := r.FormValue("password")
	if password == "" {
		httputil.SendError(w, apperror.NewBadRequest("password is required"))
		return
	}

	file, _, err := r.FormFile("certificate")
	if err != nil {
		httputil.SendError(w, apperror.NewBadRequest("certificate file is required"))
		return
	}
	defer file.Close()

	res, err := h.certService.Upload(r.Context(), companyID, password, file)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusCreated, res)
}

func (h *CertificateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	if err := h.certService.Delete(r.Context(), companyID); err != nil {
		httputil.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
