package acbr

/*
#include <stdlib.h>
#include "nfe.h"
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

// NFE_CarregarINI loads an NFe INI string into the handle.
func (hd *Handle) CarregarINI(iniContent string) error {
	hd.mu.Lock()
	defer hd.mu.Unlock()
	hd.LastUsed = time.Now()

	cINI, freeINI := allocCString(iniContent)
	defer freeINI()

	res := C.NFE_CarregarINI(hd.h, cINI)
	if res != 0 {
		return libError(hd.h, "failed to load NFe INI")
	}
	return nil
}

// NFE_LimparLista clears the loaded NFes.
func (hd *Handle) LimparLista() error {
	hd.mu.Lock()
	defer hd.mu.Unlock()
	hd.LastUsed = time.Now()

	res := C.NFE_LimparLista(hd.h)
	if res != 0 {
		return libError(hd.h, "failed to clear NFe list")
	}
	return nil
}

// NFE_Assinar signs the loaded NFes.
func (hd *Handle) Assinar() error {
	hd.mu.Lock()
	defer hd.mu.Unlock()
	hd.LastUsed = time.Now()

	res := C.NFE_Assinar(hd.h)
	if res != 0 {
		return libError(hd.h, "failed to sign NFe")
	}
	return nil
}

// NFE_Validar validates the loaded NFes against schemas.
func (hd *Handle) Validar() error {
	hd.mu.Lock()
	defer hd.mu.Unlock()
	hd.LastUsed = time.Now()

	res := C.NFE_Validar(hd.h)
	if res != 0 {
		return libError(hd.h, "failed to validate NFe schema")
	}
	return nil
}

// NFE_Enviar sends the loaded NFes to SEFAZ.
func (hd *Handle) Enviar(lote int, imprimir, sincrono, zipado bool) (string, error) {
	hd.mu.Lock()
	defer hd.mu.Unlock()
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

// NFE_Consultar queries an NFe status by chave.
func (hd *Handle) Consultar(chave string) (string, error) {
	hd.mu.Lock()
	defer hd.mu.Unlock()
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

// NFE_Cancelar sends a cancellation event.
func (hd *Handle) Cancelar(chave, justificativa, cnpj string, lote int) (string, error) {
	hd.mu.Lock()
	defer hd.mu.Unlock()
	hd.LastUsed = time.Now()

	var bufferSize C.int = 16384
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	cChave, free := allocCString(chave)
	defer free()
	cJust, free2 := allocCString(justificativa)
	defer free2()
	cCNPJ, free3 := allocCString(cnpj)
	defer free3()

	res := C.NFE_Cancelar(hd.h, cChave, cJust, cCNPJ, C.int(lote), buffer, &bufferSize)
	if res != 0 {
		return "", libError(hd.h, "failed to cancel NFe")
	}

	return readBuffer(buffer), nil
}

// NFE_DistribuicaoDFePorUltNSU queries DFe by UltNSU.
func (hd *Handle) DistribuicaoDFePorUltNSU(ufAutor int, cnpj, ultNSU string) (string, error) {
	hd.mu.Lock()
	defer hd.mu.Unlock()
	hd.LastUsed = time.Now()

	var bufferSize C.int = 1048576 // 1MB for distribution response
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	cCNPJ, free1 := allocCString(cnpj)
	defer free1()
	cUltNSU, free2 := allocCString(ultNSU)
	defer free2()

	res := C.NFE_DistribuicaoDFePorUltNSU(hd.h, C.int(ufAutor), cCNPJ, cUltNSU, buffer, &bufferSize)
	if res != 0 {
		return "", libError(hd.h, "failed to query distribution DFe by UltNSU")
	}

	return readBuffer(buffer), nil
}

// NFE_DistribuicaoDFePorNSU queries a specific NSU.
func (hd *Handle) DistribuicaoDFePorNSU(ufAutor int, cnpj, nsu string) (string, error) {
	hd.mu.Lock()
	defer hd.mu.Unlock()
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
