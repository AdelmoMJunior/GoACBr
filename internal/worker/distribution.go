package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
	"github.com/AdelmoMJunior/GoACBr/internal/service"
)

// DistributionWorker handles asynchronous synchronization of DFes from SEFAZ.
type DistributionWorker struct {
	compRepo    repository.CompanyRepository
	distRepo    repository.DistributionRepository
	distService service.DistributionService
}

// NewDistributionWorker creates a new distribution worker.
func NewDistributionWorker(
	compRepo repository.CompanyRepository,
	distRepo repository.DistributionRepository,
	distService service.DistributionService,
) *DistributionWorker {
	return &DistributionWorker{
		compRepo:    compRepo,
		distRepo:    distRepo,
		distService: distService,
	}
}

// RunOnce executes a single pass of the distribution sync for all eligible companies.
func (w *DistributionWorker) RunOnce(ctx context.Context) {
	slog.Info("Starting Distribution Worker pass")

	companies, err := w.compRepo.GetCompaniesEligibleForSync(ctx)
	if err != nil {
		slog.Error("Failed to fetch eligible companies for distribution sync", "error", err)
		return
	} 
	
	for _, comp := range companies {
		if err := w.syncCompany(ctx, comp.ID); err != nil {
			slog.Error("Failed to sync company", "company_id", comp.ID, "error", err)
		}
	}

	slog.Info("Distribution Worker pass completed")
}

func (w *DistributionWorker) syncCompany(ctx context.Context, companyID uuid.UUID) error {
	ctrl, err := w.distRepo.GetControl(ctx, companyID)
	if err != nil {
		// Initialize control if not exists
		ctrl = &domain.DistributionControl{
			CompanyID: companyID,
			LastNSU:   "0",
			MaxNSU:    "0",
		}
	}

	// 1 hr cooldown check
	if ctrl.LastQueryAt != nil && time.Since(*ctrl.LastQueryAt) < 1*time.Hour && ctrl.LastNSU == ctrl.MaxNSU {
		slog.Debug("Skipping company sync due to SEFAZ cooldown", "company_id", companyID)
		return nil
	}

	slog.Info("Syncing DFe for company", "company_id", companyID, "last_nsu", ctrl.LastNSU)

	// In a real loop, we would call QueryByUltNSU until LastNSU == MaxNSU
	// We wrap in a short loop to prevent infinite loops in case of ACBr error
	for i := 0; i < 50; i++ {
		res, err := w.distService.QueryByUltNSU(ctx, companyID, ctrl.LastNSU)
		if err != nil {
			return err
		}

		if res.CStat != 138 {
			slog.Warn("DFe sync returned non-success status", "cstat", res.CStat, "motivo", res.XMotivo)
			break
		}

		// Update control
		ctrl.LastNSU = res.UltNSU
		ctrl.MaxNSU = res.MaxNSU
		
		if ctrl.LastNSU == ctrl.MaxNSU {
			// Fully synced
			break
		}
	}

	return nil
}
