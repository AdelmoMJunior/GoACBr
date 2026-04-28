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

type DistributionHandler struct {
	distService service.DistributionService
	compRepo    repository.CompanyRepository
	certRepo    repository.CertificateRepository
}

func NewDistributionHandler(distService service.DistributionService, compRepo repository.CompanyRepository, certRepo repository.CertificateRepository) *DistributionHandler {
	return &DistributionHandler{
		distService: distService,
		compRepo:    compRepo,
		certRepo:    certRepo,
	}
}

func (h *DistributionHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.CompanyGuard(h.compRepo))
		r.Use(middleware.CertValidator(h.certRepo))

		r.Get("/distribution/control", h.GetControl)
		r.Post("/distribution/query", h.Query)
	})
}

func (h *DistributionHandler) GetControl(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	res, err := h.distService.GetControl(r.Context(), companyID)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}

func (h *DistributionHandler) Query(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	var req dto.DistributionQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid json payload"))
		return
	}

	var res *dto.DistributionQueryResponse
	var err error

	if req.NSU != "" {
		res, err = h.distService.QueryByNSU(r.Context(), companyID, req.NSU)
	} else if req.UltNSU != "" {
		res, err = h.distService.QueryByUltNSU(r.Context(), companyID, req.UltNSU)
	} else {
		httputil.SendError(w, apperror.NewBadRequest("must provide either nsu or ult_nsu"))
		return
	}

	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, res)
}
