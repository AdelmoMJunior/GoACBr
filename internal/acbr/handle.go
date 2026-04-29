package acbr

/*
#cgo LDFLAGS: -L${SRCDIR}/../../lib -lacbrnfe64 -ldl
#include <stdlib.h>
#include "nfe.h"
*/
import "C"
import (
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/google/uuid"
)

// Handle represents a single instance of ACBrLibNFe.
//
// LOCKING CONTRACT: The Handle's mutex (mu) is used EXCLUSIVELY by the
// HandlePool as an "in-use" flag. The pool acquires the lock via TryLock()
// when handing out a handle, and releases it via Unlock() when the handle
// is returned. Individual methods do NOT lock internally — the pool
// guarantees that only one goroutine uses a handle at any given time.
// This avoids deadlocks caused by nested locking.
type Handle struct {
	h             C.handle
	mu            sync.Mutex
	ConfiguredFor uuid.UUID // CompanyID currently loaded in this handle
	LastUsed      time.Time
}

// NewHandle initializes a new ACBrLibNFe handle in memory (no INI file).
func NewHandle() (*Handle, error) {
	var h C.handle

	cConfig := C.CString("")
	defer C.free(unsafe.Pointer(cConfig))

	cCrypt := C.CString("")
	defer C.free(unsafe.Pointer(cCrypt))

	slog.Debug("Calling NFE_Inicializar (memoria)...")
	res := C.NFE_Inicializar(&h, cConfig, cCrypt)
	if res != 0 {
		return nil, fmt.Errorf("failed to initialize ACBrLibNFe (code %d)", res)
	}

	slog.Debug("Handle criado em memória")
	return &Handle{
		h:        h,
		LastUsed: time.Now(),
	}, nil
}

// Destroy cleans up the handle memory in C.
// Must be called with the pool lock held (or standalone after final use).
func (hd *Handle) Destroy() error {
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
	hd.LastUsed = time.Now()

	cSection, freeSection := allocCString(section)
	defer freeSection()

	cKey, freeKey := allocCString(key)
	defer freeKey()

	cValue, freeValue := allocCString(value)
	defer freeValue()

	res := C.NFE_ConfigGravarValor(hd.h, cSection, cKey, cValue)
	if res != 0 {
		err := libError(hd.h, fmt.Sprintf("failed to set config %s/%s", section, key))
		slog.Error("NFE_ConfigGravarValor failed", "section", section, "key", key, "error", err)
		return err
	}
	return nil
}

// ConfigGravar persists the current in-memory config to a file.
func (hd *Handle) ConfigGravar(path string) error {
	hd.LastUsed = time.Now()

	cPath, freePath := allocCString(path)
	defer freePath()

	res := C.NFE_ConfigGravar(hd.h, cPath)
	if res != 0 {
		return libError(hd.h, "failed to save config file")
	}
	return nil
}

// ConfigLer loads configuration from a file into the handle.
func (hd *Handle) ConfigLer(path string) error {
	hd.LastUsed = time.Now()

	cPath, freePath := allocCString(path)
	defer freePath()

	res := C.NFE_ConfigLer(hd.h, cPath)
	if res != 0 {
		var bufferSize C.int = 4096
		buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
		defer C.free(unsafe.Pointer(buffer))
		C.NFE_UltimoRetorno(hd.h, buffer, &bufferSize)
		msg := strings.TrimSpace(C.GoString(buffer))
		return fmt.Errorf("failed to load config file: [acbr] %s", msg)
	}
	return nil
}

// ApplyCompanyConfig writes all company-specific configuration into the handle.
func (hd *Handle) ApplyCompanyConfig(companyID uuid.UUID, configs map[string]map[string]string) error {
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
