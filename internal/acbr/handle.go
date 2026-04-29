package acbr

/*
#include "nfe.h"
*/
import "C"
import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/google/uuid"
)

// Handle represents a single thread-safe instance of ACBrLibNFe.
// ACBrLib is NOT thread-safe for parallel calls on the SAME handle.
// Thus, each handle must be protected by a Mutex.
type Handle struct {
	h             C.handle
	mu            sync.Mutex
	ConfiguredFor uuid.UUID // CompanyID currently loaded in this handle
	LastUsed      time.Time
	ConfigPath    string
}

// NewHandle initializes a new ACBrLibNFe handle.
func NewHandle(libPath, configPath, cryptKey string) (*Handle, error) {
	var h C.handle

	cConfigPath, freeConfigPath := allocCString(configPath)
	defer freeConfigPath()

	cCryptKey, freeCryptKey := allocCString(cryptKey)
	defer freeCryptKey()

	// 0. Clean up any existing config file to avoid permission issues
	_ = os.Remove(configPath)

	// 1. Initialize
	slog.Debug("Calling NFE_Inicializar...", "config_path", configPath)
	res := C.NFE_Inicializar(&h, cConfigPath, cCryptKey)
	if res != 0 {
		slog.Error("NFE_Inicializar failed", "code", res)
		return nil, fmt.Errorf("failed to initialize ACBrLibNFe (code %d)", res)
	}

	slog.Debug("New ACBrLibNFe handle initialized successfully")

	return &Handle{
		h:          h,
		LastUsed:   time.Now(),
		ConfigPath: configPath,
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

	hd.LastUsed = time.Now()

	for section, keys := range configs {
		cSection, freeSection := allocCString(section)
		for key, val := range keys {
			cKey, freeKey := allocCString(key)
			cVal, freeVal := allocCString(val)

			slog.Debug("Setting ACBr config", "section", section, "key", key)
			res := C.NFE_ConfigGravarValor(hd.h, cSection, cKey, cVal)

			freeKey()
			freeVal()

			if res != 0 {
				slog.Error("Failed to set ACBr config", "section", section, "key", key, "res", res)
				freeSection()
				return libError(hd.h, fmt.Sprintf("failed to set config %s/%s", section, key))
			}
			slog.Debug("ACBr config set successfully", "section", section, "key", key)
		}
		freeSection()
	}

	hd.ConfiguredFor = companyID
	return nil
}

func (hd *Handle) ConfigLer(path string) error {
	hd.mu.Lock()
	defer hd.mu.Unlock()
	hd.LastUsed = time.Now()

	cPath, freePath := allocCString(path)
	defer freePath()

	res := C.NFE_ConfigLer(hd.h, cPath)
	if res != 0 {
		// Captura mensagem real da ACBrLib
		var bufferSize C.int = 4096
		buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
		defer C.free(unsafe.Pointer(buffer))
		C.NFE_UltimoRetorno(hd.h, buffer, &bufferSize)
		msg := strings.TrimSpace(C.GoString(buffer))
		return fmt.Errorf("failed to load config file: [acbr] %s", msg)
	}
	return nil
}
