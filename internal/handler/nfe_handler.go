package handler

import (
	"encoding/json"
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

type NFeHandler struct {
	nfeService   service.NFeService
	eventService service.EventService
	compRepo     repository.CompanyRepository
	certRepo     repository.CertificateRepository
}

func NewNFeHandler(nfeService service.NFeService, eventService service.EventService, compRepo repository.CompanyRepository, certRepo repository.CertificateRepository) *NFeHandler {
	return &NFeHandler{
		nfeService:   nfeService,
		eventService: eventService,
		compRepo:     compRepo,
		certRepo:     certRepo,
	}
}

func (h *NFeHandler) RegisterRoutes(r chi.Router) {
	// Need Auth + CompanyGuard + CertValidator
	r.Group(func(r chi.Router) {
		r.Use(middleware.CompanyGuard(h.compRepo))
		r.Use(middleware.CertValidator(h.certRepo))

		r.Post("/nfe/emit", h.Emit)
		r.Post("/nfe/status", h.Status)
		
		r.Post("/nfe/cancel", h.Cancel)
		r.Post("/nfe/cce", h.CCe)
		r.Post("/nfe/inutilizacao", h.Inutilizacao)
	})
}

func (h *NFeHandler) Emit(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	var req dto.NFeEmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid json payload"))
		return
	}

	res, err := h.nfeService.Emit(r.Context(), companyID, &req)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *NFeHandler) Status(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	var req dto.NFeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid json payload"))
		return
	}

	res, err := h.nfeService.QueryStatus(r.Context(), companyID, &req)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *NFeHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	var req dto.CancelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid json payload"))
		return
	}

	res, err := h.eventService.Cancel(r.Context(), companyID, &req)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *NFeHandler) CCe(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	var req dto.CCeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid json payload"))
		return
	}

	res, err := h.eventService.CCe(r.Context(), companyID, &req)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *NFeHandler) Inutilizacao(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	var req dto.InutilizacaoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid json payload"))
		return
	}

	res, err := h.eventService.Inutilizacao(r.Context(), companyID, &req)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}
