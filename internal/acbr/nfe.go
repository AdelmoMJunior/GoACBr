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
	"time"
	"unsafe"
)

// All methods in this file follow the Handle LOCKING CONTRACT:
// NO internal mu.Lock() — the HandlePool guarantees exclusive access.

// CarregarINI loads an NFe INI string into the handle.
func (hd *Handle) CarregarINI(iniContent string) error {
	hd.LastUsed = time.Now()

	cINI, freeINI := allocCString(iniContent)
	defer freeINI()

	res := C.NFE_CarregarINI(hd.h, cINI)
	if res != 0 {
		return libError(hd.h, "failed to load NFe INI")
	}
	return nil
}

// LimparLista clears the loaded NFes.
func (hd *Handle) LimparLista() error {
	hd.LastUsed = time.Now()

	res := C.NFE_LimparLista(hd.h)
	if res != 0 {
		return libError(hd.h, "failed to clear NFe list")
	}
	return nil
}

// Assinar signs the loaded NFes.
func (hd *Handle) Assinar() error {
	hd.LastUsed = time.Now()

	res := C.NFE_Assinar(hd.h)
	if res != 0 {
		return libError(hd.h, "failed to sign NFe")
	}
	return nil
}

// Validar validates the loaded NFes against schemas.
func (hd *Handle) Validar() error {
	hd.LastUsed = time.Now()

	res := C.NFE_Validar(hd.h)
	if res != 0 {
		return libError(hd.h, "failed to validate NFe schema")
	}
	return nil
}

// Enviar sends the loaded NFes to SEFAZ.
func (hd *Handle) Enviar(lote int, imprimir, sincrono, zipado bool) (string, error) {
	hd.LastUsed = time.Now()

	var bufferSize C.int = 65536 // 64KB for response
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	cImprimir := C.int(0)
	if imprimir {
		cImprimir = 1
	}
	cSincrono := C.int(0)
	if sincrono {
		cSincrono = 1
	}
	cZipado := C.int(0)
	if zipado {
		cZipado = 1
	}

	res := C.NFE_Enviar(hd.h, C.int(lote), cImprimir, cSincrono, cZipado, buffer, &bufferSize)
	if res != 0 {
		return "", libError(hd.h, "failed to send NFe to SEFAZ")
	}

	return readBuffer(buffer), nil
}

// Consultar queries an NFe status by chave.
func (hd *Handle) Consultar(chave string) (string, error) {
	hd.LastUsed = time.Now()

	var bufferSize C.int = 16384
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	cChave, freeChave := allocCString(chave)
	defer freeChave()

	res := C.NFE_Consultar(hd.h, cChave, C.int(0), buffer, &bufferSize)
	if res != 0 {
		return "", libError(hd.h, fmt.Sprintf("failed to query NFe %s", chave))
	}

	return readBuffer(buffer), nil
}

// CarregarEventoINI loads an event INI string.
func (hd *Handle) CarregarEventoINI(iniContent string) error {
	hd.LastUsed = time.Now()

	cINI, freeINI := allocCString(iniContent)
	defer freeINI()

	res := C.NFE_CarregarEventoINI(hd.h, cINI)
	if res != 0 {
		return libError(hd.h, "failed to load event INI")
	}
	return nil
}

// EnviarEvento sends the loaded events.
func (hd *Handle) EnviarEvento(lote int) (string, error) {
	hd.LastUsed = time.Now()

	var bufferSize C.int = 16384
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	res := C.NFE_EnviarEvento(hd.h, C.int(lote), buffer, &bufferSize)
	if res != 0 {
		return "", libError(hd.h, "failed to send events")
	}

	return readBuffer(buffer), nil
}

// DistribuicaoDFePorUltNSU queries DFe by UltNSU.
func (hd *Handle) DistribuicaoDFePorUltNSU(ufAutor int, cnpj, ultNSU string) (string, error) {
	hd.LastUsed = time.Now()

	var bufferSize C.int = 1048576 // 1MB for distribution response
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	cCNPJ, free1 := allocCString(cnpj)
	defer free1()
	cUltNSU, free2 := allocCString(ultNSU)
	defer free2()

	slog.Debug("Calling NFE_DistribuicaoDFePorUltNSU", "cnpj", cnpj, "ultNSU", ultNSU)
	res := C.NFE_DistribuicaoDFePorUltNSU(hd.h, C.int(ufAutor), cCNPJ, cUltNSU, buffer, &bufferSize)
	if res != 0 {
		err := libError(hd.h, "failed to query distribution DFe by UltNSU")
		slog.Error("NFE_DistribuicaoDFePorUltNSU failed", "error", err)
		return "", err
	}
	slog.Debug("NFE_DistribuicaoDFePorUltNSU success")

	return readBuffer(buffer), nil
}

// DistribuicaoDFePorNSU queries a specific NSU.
func (hd *Handle) DistribuicaoDFePorNSU(ufAutor int, cnpj, nsu string) (string, error) {
	hd.LastUsed = time.Now()

	var bufferSize C.int = 1048576
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	cCNPJ, free1 := allocCString(cnpj)
	defer free1()
	cNSU, free2 := allocCString(nsu)
	defer free2()

	res := C.NFE_DistribuicaoDFePorNSU(hd.h, C.int(ufAutor), cCNPJ, cNSU, buffer, &bufferSize)
	if res != 0 {
		return "", libError(hd.h, "failed to query distribution DFe by NSU")
	}

	return readBuffer(buffer), nil
}

// ObterXml returns the XML content of a loaded NFe.
func (hd *Handle) ObterXml(index int) (string, error) {
	hd.LastUsed = time.Now()

	var bufferSize C.int = 1048576 // 1MB
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	res := C.NFE_ObterXml(hd.h, C.int(index), buffer, &bufferSize)
	if res != 0 {
		return "", libError(hd.h, fmt.Sprintf("failed to get XML for index %d", index))
	}

	return readBuffer(buffer), nil
}

// Inutilizar sends an Inutilizacao.
func (hd *Handle) Inutilizar(cnpj, justificativa string, ano, modelo, serie, nInicial, nFinal int) (string, error) {
	hd.LastUsed = time.Now()

	var bufferSize C.int = 16384
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	cCNPJ, free := allocCString(cnpj)
	defer free()
	cJust, free2 := allocCString(justificativa)
	defer free2()

	res := C.NFE_Inutilizar(hd.h, cCNPJ, cJust, C.int(ano), C.int(modelo), C.int(serie), C.int(nInicial), C.int(nFinal), buffer, &bufferSize)
	if res != 0 {
		return "", libError(hd.h, "failed to inutilizar numbers")
	}

	return readBuffer(buffer), nil
}

// ImprimirPDF generates the PDF for the loaded NFes.
func (hd *Handle) ImprimirPDF() error {
	hd.LastUsed = time.Now()

	res := C.NFE_ImprimirPDF(hd.h)
	if res != 0 {
		return libError(hd.h, "failed to generate PDF")
	}
	return nil
}

// StatusServico queries the SEFAZ service status.
func (hd *Handle) StatusServico() (string, error) {
	hd.LastUsed = time.Now()

	var bufferSize C.int = 16384
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

    slog.Debug("Calling NFE_StatusServico")
    res := C.NFE_StatusServico(hd.h, buffer, &bufferSize)
    // Even on success, fetch the last ACBr retorno to aid debugging in case
    // the service status is not fully informative (e.g., when cStat=0).
    if res == 0 {
        errBufSize := C.int(8192)
        errBuf := (*C.char)(C.malloc(C.size_t(errBufSize)))
        if errBuf != nil {
            defer C.free(unsafe.Pointer(errBuf))
            C.NFE_UltimoRetorno(hd.h, errBuf, &errBufSize)
            acbrErr := C.GoString(errBuf)
            slog.Debug("NFE_StatusServico last retorno", "acbr_err", acbrErr)
        }
    }
	if res != 0 {
		// Get detailed error from ACBr
		var errBufSize C.int = 8192
		errBuf := (*C.char)(C.malloc(C.size_t(errBufSize)))
		defer C.free(unsafe.Pointer(errBuf))
		C.NFE_UltimoRetorno(hd.h, errBuf, &errBufSize)
		acbrErr := C.GoString(errBuf)

		slog.Error("NFE_StatusServico failed",
			"res_code", res,
			"acbr_error", acbrErr,
		)
		return "", fmt.Errorf("SEFAZ StatusServico error (code %d): %s", res, acbrErr)
	}
	slog.Debug("NFE_StatusServico success")

	return readBuffer(buffer), nil
}
