package dto

import (
	"time"

	"github.com/google/uuid"
)

// CompanyCreateRequest is the payload to create a new company.
type CompanyCreateRequest struct {
	CNPJ               string `json:"cnpj" validate:"required"`
	RazaoSocial        string `json:"razao_social" validate:"required"`
	NomeFantasia       string `json:"nome_fantasia,omitempty"`
	InscricaoEstadual  string `json:"inscricao_estadual,omitempty"`
	InscricaoMunicipal string `json:"inscricao_municipal,omitempty"`
	CRT                int16  `json:"crt" validate:"required"`
	Logradouro         string `json:"logradouro,omitempty"`
	Numero             string `json:"numero,omitempty"`
	Complemento        string `json:"complemento,omitempty"`
	Bairro             string `json:"bairro,omitempty"`
	CodMunicipio       string `json:"cod_municipio,omitempty"`
	Municipio          string `json:"municipio,omitempty"`
	UF                 string `json:"uf" validate:"required"`
	CEP                string `json:"cep,omitempty"`
	Telefone           string `json:"telefone,omitempty"`
	CNAE               string `json:"cnae,omitempty"`
	Ambiente           int16  `json:"ambiente"` // defaults to 2 (homologation)
	SerieNFe           int    `json:"serie_nfe"`
	SerieNFCe          int    `json:"serie_nfce"`
	CSCID              string `json:"csc_id,omitempty"`
	CSCToken           string `json:"csc_token,omitempty"`
}

// CompanySMTPRequest configures the email settings for a company.
type CompanySMTPRequest struct {
	SMTPHost     string `json:"smtp_host" validate:"required"`
	SMTPPort     int    `json:"smtp_port" validate:"required"`
	SMTPUser     string `json:"smtp_user" validate:"required"`
	SMTPPassword string `json:"smtp_password" validate:"required"`
	SMTPFrom     string `json:"smtp_from" validate:"required,email"`
	SMTPTLS      *bool  `json:"smtp_tls" validate:"required"`
}

// CompanyResponse is the standard representation of a company.
type CompanyResponse struct {
	ID                 uuid.UUID `json:"id"`
	CNPJ               string    `json:"cnpj"`
	RazaoSocial        string    `json:"razao_social"`
	NomeFantasia       string    `json:"nome_fantasia,omitempty"`
	InscricaoEstadual  string    `json:"inscricao_estadual,omitempty"`
	InscricaoMunicipal string    `json:"inscricao_municipal,omitempty"`
	CRT                int16     `json:"crt"`
	Logradouro         string    `json:"logradouro,omitempty"`
	Numero             string    `json:"numero,omitempty"`
	Complemento        string    `json:"complemento,omitempty"`
	Bairro             string    `json:"bairro,omitempty"`
	CodMunicipio       string    `json:"cod_municipio,omitempty"`
	Municipio          string    `json:"municipio,omitempty"`
	UF                 string    `json:"uf"`
	CEP                string    `json:"cep,omitempty"`
	Telefone           string    `json:"telefone,omitempty"`
	CNAE               string    `json:"cnae,omitempty"`
	Ambiente           int16     `json:"ambiente"`
	SerieNFe           int       `json:"serie_nfe"`
	SerieNFCe          int       `json:"serie_nfce"`
	CSCID              string    `json:"csc_id,omitempty"`
	SMTPConfigured     bool      `json:"smtp_configured"`
	HasCertificate     bool      `json:"has_certificate"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
