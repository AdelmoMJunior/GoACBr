package dto

import "time"

// DistributionQueryRequest represents a request to fetch DFe documents.
type DistributionQueryRequest struct {
	NSU    string `json:"nsu,omitempty"`     // Query by specific NSU
	UltNSU string `json:"ult_nsu,omitempty"` // Query from an NSU onwards
}

// DistributionQueryResponse represents the result of a DFe query.
type DistributionQueryResponse struct {
	CStat    int    `json:"c_stat"`
	XMotivo  string `json:"x_motivo"`
	UltNSU   string `json:"ult_nsu"`
	MaxNSU   string `json:"max_nsu"`
	Documentos []DistributionDoc `json:"documentos"`
}

// DistributionDoc represents a single document from a DFe query.
type DistributionDoc struct {
	NSU       string `json:"nsu"`
	Schema    string `json:"schema"` // resNFe, procNFe, resEvento, procEventoNFe
	Chave     string `json:"chave"`
	XMLBase64 string `json:"xml_base64"`
}

// DistributionControlResponse returns the current sync status for a company.
type DistributionControlResponse struct {
	CompanyID    string    `json:"company_id"`
	LastNSU      string    `json:"last_nsu"`
	MaxNSU       string    `json:"max_nsu"`
	LastQueryAt  *time.Time `json:"last_query_at"`
	IsRunning    bool      `json:"is_running"`
	Status       string    `json:"status"` // "idle", "syncing", "error", "cooldown"
	ErrorMessage string    `json:"error_message,omitempty"`
	NextAllowedQuery *time.Time `json:"next_allowed_query,omitempty"`
}
