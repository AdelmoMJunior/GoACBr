package acbr

/*
#include "nfe.h"
*/
import "C"
import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Handle represents a single thread-safe instance of ACBrLibNFe.
// ACBrLib is NOT thread-safe for parallel calls on the SAME handle.
// Thus, each handle must be protected by a Mutex.
type Handle struct {
	h            C.handle
	mu           sync.Mutex
	ConfiguredFor uuid.UUID // CompanyID currently loaded in this handle
	LastUsed     time.Time
}

// NewHandle initializes a new ACBrLibNFe handle.
func NewHandle(libPath, configPath, cryptKey string) (*Handle, error) {
	var h C.handle

	cConfigPath, freeConfigPath := allocCString(configPath)
	defer freeConfigPath()

	cCryptKey, freeCryptKey := allocCString(cryptKey)
	defer freeCryptKey()

	// 1. Initialize
	slog.Debug("Calling NFE_Inicializar...")
	res := C.NFE_Inicializar(&h, cConfigPath, cCryptKey)
	if res != 0 {
		slog.Error("NFE_Inicializar failed", "code", res)
		return nil, fmt.Errorf("failed to initialize ACBrLibNFe (code %d)", res)
	}

	slog.Debug("New ACBrLibNFe handle initialized successfully")

	return &Handle{
		h:        h,
		LastUsed: time.Now(),
	}, nil
}

// Destroy cleans up the handle memory in C.
func (hd *Handle) Destroy() error {
	hd.mu.Lock()
	defer hd.mu.Unlock()

	if hd.h == nil {
		return nil
	}

	res := C.NFE_Finalizar(hd.h)
	if res != 0 {
		return fmt.Errorf("failed to finalize ACBrLibNFe (code %d)", res)
	}

	hd.h = nil
	slog.Debug("ACBrLibNFe handle destroyed")
	return nil
}

// ConfigGravarValor sets a configuration value in memory.
func (hd *Handle) ConfigGravarValor(section, key, value string) error {
	hd.mu.Lock()
	defer hd.mu.Unlock()
	hd.LastUsed = time.Now()

	cSection, freeSection := allocCString(section)
	defer freeSection()

	cKey, freeKey := allocCString(key)
	defer freeKey()

	cValue, freeValue := allocCString(value)
	defer freeValue()

	slog.Debug("Calling NFE_ConfigGravarValor", "section", section, "key", key)
	res := C.NFE_ConfigGravarValor(hd.h, cSection, cKey, cValue)
	if res != 0 {
		err := libError(hd.h, fmt.Sprintf("failed to set config %s/%s", section, key))
		slog.Error("NFE_ConfigGravarValor failed", "error", err)
		return err
	}
	slog.Debug("NFE_ConfigGravarValor success")
	return nil
}

// ApplyCompanyConfig writes all the company-specific configuration into the handle.
func (hd *Handle) ApplyCompanyConfig(companyID uuid.UUID, configs map[string]map[string]string) error {
	// Don't lock here, we will lock inside ConfigGravarValor or do a batch lock
	// Better to lock around the whole operation to prevent intermediate states

	hd.mu.Lock()
	defer hd.mu.Unlock()
	hd.LastUsed = time.Now()

	for section, keys := range configs {
		cSection, freeSection := allocCString(section)
		for key, val := range keys {
			cKey, freeKey := allocCString(key)
			cVal, freeVal := allocCString(val)
			
			res := C.NFE_ConfigGravarValor(hd.h, cSection, cKey, cVal)
			
			freeKey()
			freeVal()
			
			if res != 0 {
				freeSection()
				return libError(hd.h, fmt.Sprintf("failed to set config %s/%s", section, key))
			}
		}
		freeSection()
	}

	hd.ConfiguredFor = companyID
	return nil
}
