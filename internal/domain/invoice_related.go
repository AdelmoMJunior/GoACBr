package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// InvoicePayment represents a payment entry in an invoice.
type InvoicePayment struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	InvoiceID uuid.UUID       `json:"invoice_id" db:"invoice_id"`
	NPag      int             `json:"n_pag" db:"n_pag"`
	TpPag     string          `json:"tp_pag" db:"tp_pag"`
	XPag      string          `json:"x_pag,omitempty" db:"x_pag"`
	VPag      decimal.Decimal `json:"v_pag" db:"v_pag"`
	IndPag    *int16          `json:"ind_pag,omitempty" db:"ind_pag"`
	TpIntegra *int16          `json:"tp_integra,omitempty" db:"tp_integra"`
	CNPJPag   string          `json:"cnpj_pag,omitempty" db:"cnpj_pag"`
	TBand     string          `json:"t_band,omitempty" db:"t_band"`
	CAut      string          `json:"c_aut,omitempty" db:"c_aut"`
	VTroco    decimal.Decimal `json:"v_troco" db:"v_troco"`
}

// InvoiceTransport represents the transport data of an invoice.
type InvoiceTransport struct {
	ID             uuid.UUID        `json:"id" db:"id"`
	InvoiceID      uuid.UUID        `json:"invoice_id" db:"invoice_id"`
	ModFrete       int16            `json:"mod_frete" db:"mod_frete"`
	TranspCNPJCPF  string           `json:"transp_cnpj_cpf,omitempty" db:"transp_cnpj_cpf"`
	TranspNome     string           `json:"transp_nome,omitempty" db:"transp_nome"`
	TranspIE       string           `json:"transp_ie,omitempty" db:"transp_ie"`
	TranspEndereco string           `json:"transp_endereco,omitempty" db:"transp_endereco"`
	TranspMunicipio string          `json:"transp_municipio,omitempty" db:"transp_municipio"`
	TranspUF       string           `json:"transp_uf,omitempty" db:"transp_uf"`
	VServ          *decimal.Decimal `json:"v_serv,omitempty" db:"v_serv"`
	VBCRet         *decimal.Decimal `json:"v_bc_ret,omitempty" db:"v_bc_ret"`
	PICMSRet       *decimal.Decimal `json:"p_icms_ret,omitempty" db:"p_icms_ret"`
	VICMSRet       *decimal.Decimal `json:"v_icms_ret,omitempty" db:"v_icms_ret"`
	Placa          string           `json:"placa,omitempty" db:"placa"`
	UFPlaca        string           `json:"uf_placa,omitempty" db:"uf_placa"`
	QVol           *int             `json:"q_vol,omitempty" db:"q_vol"`
	EspVol         string           `json:"esp_vol,omitempty" db:"esp_vol"`
	PesoL          *decimal.Decimal `json:"peso_l,omitempty" db:"peso_l"`
	PesoB          *decimal.Decimal `json:"peso_b,omitempty" db:"peso_b"`
}

// InvoiceEvent represents an event on an invoice (cancellation, correction, etc).
type InvoiceEvent struct {
	ID           uuid.UUID        `json:"id" db:"id"`
	InvoiceID    *uuid.UUID       `json:"invoice_id,omitempty" db:"invoice_id"`
	CompanyID    uuid.UUID        `json:"company_id" db:"company_id"`
	ChaveNFe     string           `json:"chave_nfe" db:"chave_nfe"`
	COrgao       int16            `json:"c_orgao" db:"c_orgao"`
	TpEvento     string           `json:"tp_evento" db:"tp_evento"`
	NSeqEvento   int              `json:"n_seq_evento" db:"n_seq_evento"`
	DHEvento     time.Time        `json:"dh_evento" db:"dh_evento"`
	Protocolo    string           `json:"protocolo,omitempty" db:"protocolo"`
	XJust        string           `json:"x_just,omitempty" db:"x_just"`
	XCorrecao    string           `json:"x_correcao,omitempty" db:"x_correcao"`
	TpNF         *int16           `json:"tp_nf,omitempty" db:"tp_nf"`
	DestCNPJCPF  string           `json:"dest_cnpj_cpf,omitempty" db:"dest_cnpj_cpf"`
	DestUF       string           `json:"dest_uf,omitempty" db:"dest_uf"`
	VNF          *decimal.Decimal `json:"v_nf,omitempty" db:"v_nf"`
	VICMS        *decimal.Decimal `json:"v_icms,omitempty" db:"v_icms"`
	VST          *decimal.Decimal `json:"v_st,omitempty" db:"v_st"`
	XMLB2Key     string           `json:"xml_b2_key,omitempty" db:"xml_b2_key"`
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
}

// InvoiceInutilizacao represents a number range voiding.
type InvoiceInutilizacao struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	CompanyID      uuid.UUID  `json:"company_id" db:"company_id"`
	Ano            int        `json:"ano" db:"ano"`
	Modelo         int16      `json:"modelo" db:"modelo"`
	Serie          int        `json:"serie" db:"serie"`
	NumInicial     int        `json:"num_inicial" db:"num_inicial"`
	NumFinal       int        `json:"num_final" db:"num_final"`
	Justificativa  string     `json:"justificativa" db:"justificativa"`
	Protocolo      string     `json:"protocolo,omitempty" db:"protocolo"`
	DHRecebimento  *time.Time `json:"dh_recebimento,omitempty" db:"dh_recebimento"`
	Status         string     `json:"status" db:"status"`
	XMLB2Key       string     `json:"xml_b2_key,omitempty" db:"xml_b2_key"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}
