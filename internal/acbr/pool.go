package acbr

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

// HandlePool manages a pool of ACBrLibNFe handles.
// Since handles are expensive to create and we want to optimize configuration
// time, we maintain a pool of them. If a handle is already configured for a
// specific CNPJ, we prefer to reuse it to avoid rewriting all config sections.
type HandlePool struct {
	handles    []*Handle
	maxHandles int
	mu         sync.Mutex
	libPath    string
	configPath string
	cryptKey   string
}

// NewHandlePool creates a new pool with the given capacity.
func NewHandlePool(maxHandles int, libPath, configPath, cryptKey string) (*HandlePool, error) {
	pool := &HandlePool{
		handles:    make([]*Handle, 0, maxHandles),
		maxHandles: maxHandles,
		libPath:    libPath,
		configPath: configPath,
		cryptKey:   cryptKey,
	}

	// Pre-warm the pool with at least one handle
	h, err := NewHandle(libPath, configPath, cryptKey)
	if err != nil {
		return nil, fmt.Errorf("failed to pre-warm pool: %w", err)
	}
	pool.handles = append(pool.handles, h)

	// Start janitor to clean up idle handles
	go pool.janitor()

	return pool, nil
}

// GetHandle retrieves a handle from the pool, preferably one already configured for companyID.
func (p *HandlePool) GetHandle(ctx context.Context, companyID uuid.UUID) (*Handle, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// 1. Attempt to find a free handle
			h, canCreate := p.tryAcquire(companyID)
			if h != nil {
				return h, nil
			}

			// 2. If we can create a new one, do it outside the lock
			if canCreate {
				newH, err := NewHandle(p.libPath, p.configPath, p.cryptKey)
				if err == nil {
					p.mu.Lock()
					// Double check if someone else filled the pool while we were creating
					if len(p.handles) < p.maxHandles {
						p.handles = append(p.handles, newH)
						newH.mu.Lock()
						p.mu.Unlock()
						return newH, nil
					}
					// Pool filled up, destroy this one and continue loop
					p.mu.Unlock()
					_ = newH.Destroy()
				} else {
					slog.Error("Failed to create new ACBr handle", "error", err)
				}
			}

			// 3. Backoff and retry
			slog.Debug("All ACBr handles busy, waiting...", "company_id", companyID)
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (p *HandlePool) tryAcquire(companyID uuid.UUID) (*Handle, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 1. Try to find a free handle already configured for this company
	var bestMatch *Handle
	for _, h := range p.handles {
		if h.mu.TryLock() {
			if h.ConfiguredFor == companyID {
				return h, false
			}
			if bestMatch == nil {
				bestMatch = h
			} else {
				h.mu.Unlock()
			}
		}
	}

	if bestMatch != nil {
		bestMatch.ConfiguredFor = uuid.Nil
		return bestMatch, false
	}

	// 2. Can we create a new one?
	canCreate := len(p.handles) < p.maxHandles
	return nil, canCreate
}

// ReleaseHandle puts a handle back into the available pool (by simply unlocking it).
func (p *HandlePool) ReleaseHandle(h *Handle) {
	// If a handle got corrupted, we could destroy it and remove it here,
	// but normally we just unlock it.
	h.mu.Unlock()
}

// Close destroys all handles in the pool.
func (p *HandlePool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, h := range p.handles {
		_ = h.Destroy()
	}
	p.handles = nil
}

// janitor cleans up handles that have been idle for too long to save memory.
func (p *HandlePool) janitor() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.Lock()
		activeHandles := make([]*Handle, 0, p.maxHandles)
		now := time.Now()

		for _, h := range p.handles {
			// If handle is free (TryLock succeeds) and hasn't been used in 30 mins
			if h.mu.TryLock() {
				if now.Sub(h.LastUsed) > 30*time.Minute && len(activeHandles) > 0 {
					// Destroy and don't append to active
					_ = h.Destroy()
					h.mu.Unlock()
					continue
				}
				h.mu.Unlock()
			}
			activeHandles = append(activeHandles, h)
		}
		p.handles = activeHandles
		p.mu.Unlock()
	}
}
