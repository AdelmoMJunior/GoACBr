package acbr

/*
#cgo LDFLAGS: -L../../lib -lacbrnfe64 -ldl
#include <stdlib.h>
#include "nfe.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

// libError extracts the last error from ACBrLib using NFE_UltimoRetorno.
func libError(h C.handle, defaultMsg string) error {
	var bufferSize C.int = 4096
	buffer := (*C.char)(C.malloc(C.size_t(bufferSize)))
	defer C.free(unsafe.Pointer(buffer))

	res := C.NFE_UltimoRetorno(h, buffer, &bufferSize)
	if res == 0 {
		msg := C.GoString(buffer)
		// Clean up string
		msg = strings.TrimSpace(msg)
		return fmt.Errorf("%s: %s", defaultMsg, msg)
	}

	return errors.New(defaultMsg)
}

// Result struct to standardize returning ACBrLib string responses.
type Result struct {
	Response string
}

// allocCString is a helper to allocate a C string and return it along with a free function.
func allocCString(s string) (*C.char, func()) {
	cstr := C.CString(s)
	return cstr, func() {
		C.free(unsafe.Pointer(cstr))
	}
}

// readBuffer is a helper to read the buffer populated by ACBrLib functions.
func readBuffer(buffer *C.char) string {
	return C.GoString(buffer)
}
