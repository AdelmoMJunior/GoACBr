package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DistributionDocument represents a document received via DFe distribution.
type DistributionDocument struct {
	ID             uuid.UUID        `json:"id" db:"id"`
	CompanyID      uuid.UUID        `json:"company_id" db:"company_id"`
	NSU            string           `json:"nsu" db:"nsu"`
	SchemaType     string           `json:"schema_type" db:"schema_type"`
	ChaveNFe       string           `json:"chave_nfe,omitempty" db:"chave_nfe"`
	TpNF           *int16           `json:"tp_nf,omitempty" db:"tp_nf"`
	EmitCNPJCPF    string           `json:"emit_cnpj_cpf,omitempty" db:"emit_cnpj_cpf"`
	EmitNome       string           `json:"emit_nome,omitempty" db:"emit_nome"`
	EmitIE         string           `json:"emit_ie,omitempty" db:"emit_ie"`
	DestCNPJCPF    string           `json:"dest_cnpj_cpf,omitempty" db:"dest_cnpj_cpf"`
	DHEmissao      *time.Time       `json:"dh_emissao,omitempty" db:"dh_emissao"`
	Modelo         *int16           `json:"modelo,omitempty" db:"modelo"`
	Serie          *int             `json:"serie,omitempty" db:"serie"`
	Numero         *int             `json:"numero,omitempty" db:"numero"`
	NatOp          string           `json:"nat_op,omitempty" db:"nat_op"`
	CSitNFe        *int16           `json:"c_sit_nfe,omitempty" db:"c_sit_nfe"`
	TotVNF         *decimal.Decimal `json:"tot_v_nf,omitempty" db:"tot_v_nf"`
	TotVICMS       *decimal.Decimal `json:"tot_v_icms,omitempty" db:"tot_v_icms"`
	TotVST         *decimal.Decimal `json:"tot_v_st,omitempty" db:"tot_v_st"`
	TotVPIS        *decimal.Decimal `json:"tot_v_pis,omitempty" db:"tot_v_pis"`
	TotVCOFINS     *decimal.Decimal `json:"tot_v_cofins,omitempty" db:"tot_v_cofins"`
	TotVProd       *decimal.Decimal `json:"tot_v_prod,omitempty" db:"tot_v_prod"`
	TotVDesc       *decimal.Decimal `json:"tot_v_desc,omitempty" db:"tot_v_desc"`
	TotVFrete      *decimal.Decimal `json:"tot_v_frete,omitempty" db:"tot_v_frete"`
	TotVOutro      *decimal.Decimal `json:"tot_v_outro,omitempty" db:"tot_v_outro"`
	TpEvento       string           `json:"tp_evento,omitempty" db:"tp_evento"`
	DescEvento     string           `json:"desc_evento,omitempty" db:"desc_evento"`
	NSeqEvento     *int             `json:"n_seq_evento,omitempty" db:"n_seq_evento"`
	DHEvento       *time.Time       `json:"dh_evento,omitempty" db:"dh_evento"`
	XJust          string           `json:"x_just,omitempty" db:"x_just"`
	XCorrecao      string           `json:"x_correcao,omitempty" db:"x_correcao"`
	Protocolo      string           `json:"protocolo,omitempty" db:"protocolo"`
	DHRecebimento  *time.Time       `json:"dh_recebimento,omitempty" db:"dh_recebimento"`
	XMLB2Key       string           `json:"xml_b2_key,omitempty" db:"xml_b2_key"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
}

// DistributionControl tracks the DFe distribution state for each company.
type DistributionControl struct {
	CompanyID    uuid.UUID  `json:"company_id" db:"company_id"`
	LastNSU      string     `json:"last_nsu" db:"last_nsu"`
	MaxNSU       string     `json:"max_nsu" db:"max_nsu"`
	LastQueryAt  *time.Time `json:"last_query_at,omitempty" db:"last_query_at"`
	IsRunning    bool       `json:"is_running" db:"is_running"`
	Status       string     `json:"status" db:"status"`
	ErrorMessage string     `json:"error_message,omitempty" db:"error_message"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}
