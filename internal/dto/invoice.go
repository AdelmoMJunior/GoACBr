package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// InvoiceListResponse is a simplified version of Invoice for lists.
type InvoiceListResponse struct {
	ID        uuid.UUID `json:"id"`
	Chave     string    `json:"chave"`
	Numero    int       `json:"numero"`
	Serie     int       `json:"serie"`
	Modelo    int16     `json:"modelo"`
	DHEmissao time.Time `json:"dh_emissao"`
	Status    string    `json:"status"`
	TotVNF    decimal.Decimal `json:"tot_v_nf"`
	DestNome  string    `json:"dest_nome,omitempty"`
	DestCNPJCPF string  `json:"dest_cnpj_cpf,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// EmittedDocsSummaryResponse provides statistics for emitted documents.
type EmittedDocsSummaryResponse struct {
	TotalAutorizadas int             `json:"total_autorizadas"`
	TotalCanceladas  int             `json:"total_canceladas"`
	TotalInutilizadas int            `json:"total_inutilizadas"`
	TotalDenegadas   int             `json:"total_denegadas"`
	ValorTotalAutorizadas decimal.Decimal `json:"valor_total_autorizadas"`
	ValorTotalCanceladas  decimal.Decimal `json:"valor_total_canceladas"`
}
