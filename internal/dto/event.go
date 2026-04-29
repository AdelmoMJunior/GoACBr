package dto

import "time"

// EventRequest is the base request for an NFe event.
type EventRequest struct {
	Chave   string `json:"chave" validate:"required,len=44"`
	CNPJCPF string `json:"cnpj_cpf" validate:"required"`
	Lote    int    `json:"lote" validate:"required"`
}

// CancelRequest is the payload to cancel an NFe.
type CancelRequest struct {
	EventRequest
	Protocolo     string `json:"protocolo" validate:"required"`
	Justificativa string `json:"justificativa" validate:"required,min=15,max=255"`
}

// CCeRequest is the payload to issue a Carta de Correção (CCe).
type CCeRequest struct {
	EventRequest
	NSeqEvento int    `json:"n_seq_evento" validate:"required,min=1"`
	Correcao   string `json:"correcao" validate:"required,min=15,max=1000"`
}

// EventResponse represents the result of an event operation.
type EventResponse struct {
	Chave      string    `json:"chave"`
	TpEvento   string    `json:"tp_evento"`
	NSeqEvento int       `json:"n_seq_evento"`
	CStat      int       `json:"c_stat"`
	XMotivo    string    `json:"x_motivo"`
	Protocolo  string    `json:"protocolo,omitempty"`
	DHEvento   time.Time `json:"dh_evento,omitempty"`
	XMLBase64  string    `json:"xml_base64,omitempty"`
	PDFBase64  string    `json:"pdf_base64,omitempty"`
}

// InutilizacaoRequest represents a request to void a range of NFe numbers.
type InutilizacaoRequest struct {
	Ano           int    `json:"ano" validate:"required,min=0,max=99"`
	Modelo        int    `json:"modelo" validate:"required,oneof=55 65"`
	Serie         int    `json:"serie" validate:"required,min=1"`
	NumInicial    int    `json:"num_inicial" validate:"required,min=1"`
	NumFinal      int    `json:"num_final" validate:"required,min=1,gtefield=NumInicial"`
	Justificativa string `json:"justificativa" validate:"required,min=15,max=255"`
}

// InutilizacaoResponse represents the result of a voiding operation.
type InutilizacaoResponse struct {
	CStat     int    `json:"c_stat"`
	XMotivo   string `json:"x_motivo"`
	Protocolo string `json:"protocolo,omitempty"`
	XMLBase64 string `json:"xml_base64,omitempty"`
}
