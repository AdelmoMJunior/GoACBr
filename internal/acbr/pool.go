package acbr

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// handleCounter generates unique IDs for handle INI files.
var handleCounter int64

// HandlePool manages a pool of ACBrLibNFe handles.
type HandlePool struct {
	handles    []*Handle
	maxHandles int
	mu         sync.Mutex

	// SchemasPath is the absolute path to the NFe XML schemas directory.
	SchemasPath string
	// LogPath is the absolute path where ACBrLib should write its logs.
	LogPath string
	// IniDir is the directory where per-handle INI files are stored.
	IniDir string
	CryptKey string
	// IniServicosPath is the path to ACBrNFeServicos.ini
	IniServicosPath string
}

// NewHandlePool creates a new pool with the given capacity.
func NewHandlePool(maxHandles int, schemasPath, logPath, cryptKey, iniServicosPath string) (*HandlePool, error) {
	iniDir := "/tmp/acbr_ini"
	os.MkdirAll(iniDir, 0700)

	pool := &HandlePool{
		handles:     make([]*Handle, 0, maxHandles),
		maxHandles:  maxHandles,
		SchemasPath: schemasPath,
		LogPath:     logPath,
		IniDir:      iniDir,
		CryptKey:    cryptKey,
		IniServicosPath: iniServicosPath,
	}

	// Pre-warm the pool with at least one handle
	h, err := pool.createHandle()
	if err != nil {
		return nil, fmt.Errorf("failed to pre-warm pool: %w", err)
	}
	pool.handles = append(pool.handles, h)

	// Start janitor to clean up idle handles
	go pool.janitor()

	return pool, nil
}

// createHandle creates a new ACBr handle with a unique INI file.
func (p *HandlePool) createHandle() (*Handle, error) {
	id := atomic.AddInt64(&handleCounter, 1)
	iniPath := filepath.Join(p.IniDir, fmt.Sprintf("handle_%d.ini", id))
	return NewHandle(iniPath, p.CryptKey)
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
				newH, err := p.createHandle()
				if err == nil {
					p.mu.Lock()
					if len(p.handles) < p.maxHandles {
						p.handles = append(p.handles, newH)
						newH.mu.Lock() // Mark as in-use
						p.mu.Unlock()
						return newH, nil
					}
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
				return h, false // Already configured — best case
			}
			if bestMatch == nil {
				bestMatch = h
			} else {
				h.mu.Unlock()
			}
		}
	}

	if bestMatch != nil {
		bestMatch.ConfiguredFor = uuid.Nil // Will need reconfiguration
		return bestMatch, false
	}

	// 2. Can we create a new one?
	canCreate := len(p.handles) < p.maxHandles
	return nil, canCreate
}

// ReleaseHandle returns a handle to the pool (unlocks the "in-use" flag).
func (p *HandlePool) ReleaseHandle(h *Handle) {
	h.mu.Unlock()
}

// Close destroys all handles in the pool.
func (p *HandlePool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, h := range p.handles {
		h.mu.TryLock()
		_ = h.Destroy()
	}
	p.handles = nil
}

// janitor cleans up handles that have been idle for too long.
func (p *HandlePool) janitor() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.Lock()
		activeHandles := make([]*Handle, 0, p.maxHandles)
		now := time.Now()

		for _, h := range p.handles {
			if h.mu.TryLock() {
				if now.Sub(h.LastUsed) > 30*time.Minute && len(activeHandles) > 0 {
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
