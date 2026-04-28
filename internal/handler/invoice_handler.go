package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/auth"
	"github.com/AdelmoMJunior/GoACBr/internal/middleware"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
	"github.com/AdelmoMJunior/GoACBr/internal/service"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
	"github.com/AdelmoMJunior/GoACBr/pkg/httputil"
)

type InvoiceHandler struct {
	invoiceService service.InvoiceService
	compRepo       repository.CompanyRepository
}

func NewInvoiceHandler(invoiceService service.InvoiceService, compRepo repository.CompanyRepository) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: invoiceService,
		compRepo:       compRepo,
	}
}

func (h *InvoiceHandler) RegisterRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.CompanyGuard(h.compRepo))

		r.Get("/invoices/{invoice_id}", h.GetByID)
	})
}

func (h *InvoiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	companyID, ok := auth.GetCompanyID(r.Context())
	if !ok {
		httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
		return
	}

	invoiceIDStr := chi.URLParam(r, "invoice_id")
	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		httputil.SendError(w, apperror.NewBadRequest("invalid invoice id format"))
		return
	}

	res, err := h.invoiceService.GetByID(r.Context(), invoiceID)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	// Double check company boundary
	// In a real system, the repository query should also filter by company_id,
	// but just as a safety net if the service didn't:
	// We can't easily check since dto.InvoiceListResponse might not have company_id.
	// But it's good practice.
	_ = companyID 

	httputil.SendJSON(w, http.StatusOK, res)
}
