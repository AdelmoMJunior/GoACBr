package worker

import (
	"context"
	"log/slog"
	"sync"
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
	maxParallel int // Maximum number of companies to sync in parallel
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
		maxParallel: 5, // Semaphore: at most 5 companies syncing in parallel
	}
}

// RunOnce executes a single pass of the distribution sync for all eligible companies.
// Each company runs in its own goroutine (limited by maxParallel).
func (w *DistributionWorker) RunOnce(ctx context.Context) {
	slog.Info("Starting Distribution Worker pass")

	companies, err := w.compRepo.GetCompaniesEligibleForSync(ctx)
	if err != nil {
		slog.Error("Failed to fetch eligible companies for distribution sync", "error", err)
		return
	}

	if len(companies) == 0 {
		slog.Debug("No companies eligible for distribution sync")
		return
	}

	slog.Info("Companies eligible for sync", "count", len(companies))

	var wg sync.WaitGroup
	sem := make(chan struct{}, w.maxParallel) // Semaphore for concurrency limit

	for _, comp := range companies {
		wg.Add(1)
		go func(companyID uuid.UUID) {
			defer wg.Done()

			// Acquire semaphore slot
			sem <- struct{}{}
			defer func() { <-sem }()

			if err := w.syncCompany(ctx, companyID); err != nil {
				slog.Error("Failed to sync company", "company_id", companyID, "error", err)

				// Save error status to DB
				now := time.Now()
				errMsg := err.Error()
				_ = w.distRepo.UpsertControl(ctx, &domain.DistributionControl{
					CompanyID:    companyID,
					Status:       "error",
					ErrorMessage: errMsg,
					LastQueryAt:  &now,
					UpdatedAt:    now,
				})
			}
		}(comp.ID)
	}

	wg.Wait()
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

    // (Cooldown removed) Always start distribution for this company

	slog.Info("Syncing DFe for company", "company_id", companyID, "last_nsu", ctrl.LastNSU)

	// Mark as running
	now := time.Now()
	ctrl.IsRunning = true
	ctrl.Status = "syncing"
	ctrl.UpdatedAt = now
	_ = w.distRepo.UpsertControl(ctx, ctrl)

	// Loop to fetch all NSUs until LastNSU == MaxNSU (max 50 batches per pass)
	for i := 0; i < 50; i++ {
		res, err := w.distService.QueryByUltNSU(ctx, companyID, ctrl.LastNSU)
		if err != nil {
			// Persist error and stop
			ctrl.IsRunning = false
			ctrl.Status = "error"
			ctrl.ErrorMessage = err.Error()
			ctrl.UpdatedAt = time.Now()
			_ = w.distRepo.UpsertControl(ctx, ctrl)
			return err
		}

		if res.CStat != 138 {
			slog.Warn("DFe sync returned non-success status", "cstat", res.CStat, "motivo", res.XMotivo)
			break
		}

		slog.Info("Batch synced", "company_id", companyID, "docs_count", len(res.Documentos), "new_last_nsu", res.UltNSU, "max_nsu", res.MaxNSU)

		// Update control after each batch
		ctrl.LastNSU = res.UltNSU
		ctrl.MaxNSU = res.MaxNSU
		ctrl.UpdatedAt = time.Now()
		_ = w.distRepo.UpsertControl(ctx, ctrl)

		if ctrl.LastNSU == ctrl.MaxNSU {
			// Fully synced
			break
		}
	}

	// Mark as idle after completion
	ctrl.IsRunning = false
	ctrl.Status = "idle"
	ctrl.ErrorMessage = ""
	queryTime := time.Now()
	ctrl.LastQueryAt = &queryTime
	ctrl.UpdatedAt = queryTime
	_ = w.distRepo.UpsertControl(ctx, ctrl)

	return nil
}
