package dto

// NFeEmitRequest represents the payload to emit a new NFe/NFCe.
// This should map closely to the ACBr INI structure, but as JSON.
// Since the INI has hundreds of fields, we typically accept the full JSON
// that will be converted into an INI string, or we accept a pre-formatted INI string.
// For this API, we will accept the raw INI content as a string for maximum flexibility,
// or a structured JSON if preferred. For simplicity and robustness with ACBr,
// accepting the INI content directly is often the most reliable approach for complex cases.
type NFeEmitRequest struct {
	// The full content of the NFe INI file
	INIContent string `json:"ini_content" binding:"required"`
	
	// Modelo: 55 for NFe, 65 for NFCe
	Modelo int `json:"modelo" binding:"required,oneof=55 65"`
	
	// Lote number
	Lote int `json:"lote" binding:"required"`
	
	// Whether to print/generate PDF immediately
	PrintPDF bool `json:"print_pdf"`
}

// NFeResponse represents the result of an NFe emission.
type NFeResponse struct {
	Chave       string `json:"chave"`
	Protocolo   string `json:"protocolo,omitempty"`
	Recibo      string `json:"recibo,omitempty"`
	Status      string `json:"status"`
	CStat       int    `json:"c_stat"`
	XMotivo     string `json:"x_motivo"`
	XMLBase64   string `json:"xml_base64,omitempty"` // The emitted XML
	PDFBase64   string `json:"pdf_base64,omitempty"` // The generated PDF (DANFE)
	XMLB2Key    string `json:"xml_b2_key,omitempty"` // Path in B2 storage
	PDFB2Key    string `json:"pdf_b2_key,omitempty"`
}

// NFeStatusRequest represents a request to query NFe status.
type NFeStatusRequest struct {
	Chave string `json:"chave" binding:"required,len=44"`
}
