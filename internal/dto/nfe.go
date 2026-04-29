package dto

// NFeEmitRequest represents the payload to emit a new NFe/NFCe.
// For maximum flexibility with ACBr, we accept the raw INI content as a string.
type NFeEmitRequest struct {
	// The full content of the NFe INI file
	INIContent string `json:"ini_content" validate:"required"`

	// Modelo: 55 for NFe, 65 for NFCe
	Modelo int `json:"modelo" validate:"required,oneof=55 65"`

	// Lote number
	Lote int `json:"lote" validate:"required"`

	// Whether to print/generate PDF immediately
	PrintPDF bool `json:"print_pdf"`
}

// NFeResponse represents the result of an NFe emission.
type NFeResponse struct {
	Chave     string `json:"chave"`
	Protocolo string `json:"protocolo,omitempty"`
	Recibo    string `json:"recibo,omitempty"`
	Status    string `json:"status"`
	CStat     int    `json:"c_stat"`
	XMotivo   string `json:"x_motivo"`
	XMLBase64 string `json:"xml_base64,omitempty"` // The emitted XML
	PDFBase64 string `json:"pdf_base64,omitempty"` // The generated PDF (DANFE)
	XMLB2Key  string `json:"xml_b2_key,omitempty"` // Path in B2 storage
	PDFB2Key  string `json:"pdf_b2_key,omitempty"`
}

// NFeStatusRequest represents a request to query NFe status.
type NFeStatusRequest struct {
	Chave string `json:"chave" validate:"required,len=44"`
}
