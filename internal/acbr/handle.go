package acbr

/*
#cgo LDFLAGS: -L${SRCDIR}/../../lib -lacbrnfe64
#cgo CFLAGS: -I${SRCDIR}
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

// allocCString and libError are defined in errors.go

// Handle wraps a raw ACBrLibNFe handle pointer.
type Handle struct {
	h             C.handle
	mu            sync.Mutex
	ConfiguredFor uuid.UUID // CompanyID currently loaded in this handle
	LastUsed      time.Time
	IniPath       string // Path to the INI file backing this handle
}

// NewHandle initializes a new ACBrLibNFe handle.
// iniPath MUST be a valid file path — ACBr will create it with all default
// values if it doesn't exist. This ensures all config sections ([NFe], [DANFE], etc.)
// are properly registered and accessible via ConfigGravarValor.
func NewHandle(iniPath string) (*Handle, error) {
	var h C.handle

	cConfig, freeConfig := allocCString(iniPath)
	defer freeConfig()

	cCrypt := C.CString("")
	defer C.free(unsafe.Pointer(cCrypt))

	slog.Debug("Calling NFE_Inicializar", "ini_path", iniPath)
	res := C.NFE_Inicializar(&h, cConfig, cCrypt)
	if res != 0 {
		return nil, fmt.Errorf("failed to initialize ACBrLibNFe (code %d)", res)
	}

	slog.Debug("Handle initialized with INI", "path", iniPath)
	return &Handle{
		h:        h,
		LastUsed: time.Now(),
		IniPath:  iniPath,
	}, nil
}

// Destroy cleans up the handle memory in C.
func (hd *Handle) Destroy() error {
	if hd.h == nil {
		return nil
	}

	res := C.NFE_Finalizar(hd.h)
	if res != 0 {
		return fmt.Errorf("failed to finalize ACBrLibNFe (code %d)", res)
	}

	hd.h = nil
	return nil
}

// ConfigGravarValor sets a specific configuration value.
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
// sections is a map of section -> key -> value.
func (hd *Handle) ApplyCompanyConfig(sections map[string]map[string]string) {
	hd.LastUsed = time.Now()

	for section, keys := range sections {
		for key, value := range keys {
			if err := hd.ConfigGravarValor(section, key, value); err != nil {
				slog.Warn("Config key skipped",
					"section", section, "key", key, "error", err)
			}
		}
	}
}

// ConfigLerValor reads a specific configuration value from the handle.
func (hd *Handle) ConfigLerValor(section, key string) (string, error) {
	hd.LastUsed = time.Now()

	cSection, freeSection := allocCString(section)
	defer freeSection()

	cKey, freeKey := allocCString(key)
	defer freeKey()

	var bufferSize C.int = 4096
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	res := C.NFE_ConfigLerValor(hd.h, cSection, cKey, buffer, &bufferSize)
	if res != 0 {
		return "", libError(hd.h, fmt.Sprintf("failed to read config %s/%s", section, key))
	}
	return strings.TrimSpace(C.GoString(buffer)), nil
}
