package worker

import (
	"context"
	"log/slog"
	"time"
)

// Scheduler manages background jobs.
type Scheduler struct {
	distWorker *DistributionWorker
	stopChan   chan struct{}
}

// NewScheduler creates a new background job scheduler.
func NewScheduler(distWorker *DistributionWorker) *Scheduler {
	return &Scheduler{
		distWorker: distWorker,
		stopChan:   make(chan struct{}),
	}
}

// Start begins the background scheduling.
func (s *Scheduler) Start(ctx context.Context) {
	slog.Info("Starting background scheduler")

	// Trigger distribution check every 15 minutes
	// Companies that haven't been checked in 1 hour will be synced
	ticker := time.NewTicker(15 * time.Minute)

	go func() {
		// Run once immediately on startup
		s.distWorker.RunOnce(ctx)

		for {
			select {
			case <-ticker.C:
				s.distWorker.RunOnce(ctx)
			case <-s.stopChan:
				ticker.Stop()
				slog.Info("Background scheduler stopped")
				return
			case <-ctx.Done():
				ticker.Stop()
				slog.Info("Background scheduler context cancelled")
				return
			}
		}
	}()
}

// Stop gracefully stops the scheduler.
func (s *Scheduler) Stop() {
	close(s.stopChan)
}
