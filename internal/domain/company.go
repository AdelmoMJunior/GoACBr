package domain

import (
	"time"

	"github.com/google/uuid"
)

// Company represents a legal entity (CNPJ) in the system.
type Company struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	CNPJ                string    `json:"cnpj" db:"cnpj"`
	RazaoSocial         string    `json:"razao_social" db:"razao_social"`
	NomeFantasia        string    `json:"nome_fantasia,omitempty" db:"nome_fantasia"`
	InscricaoEstadual   string    `json:"inscricao_estadual,omitempty" db:"inscricao_estadual"`
	InscricaoMunicipal  string    `json:"inscricao_municipal,omitempty" db:"inscricao_municipal"`
	CRT                 int16     `json:"crt" db:"crt"`
	Logradouro          string    `json:"logradouro,omitempty" db:"logradouro"`
	Numero              string    `json:"numero,omitempty" db:"numero"`
	Complemento         string    `json:"complemento,omitempty" db:"complemento"`
	Bairro              string    `json:"bairro,omitempty" db:"bairro"`
	CodMunicipio        string    `json:"cod_municipio,omitempty" db:"cod_municipio"`
	Municipio           string    `json:"municipio,omitempty" db:"municipio"`
	UF                  string    `json:"uf" db:"uf"`
	CEP                 string    `json:"cep,omitempty" db:"cep"`
	Telefone            string    `json:"telefone,omitempty" db:"telefone"`
	CNAE                string    `json:"cnae,omitempty" db:"cnae"`
	Ambiente            int16     `json:"ambiente" db:"ambiente"` // 1=Produção 2=Homologação
	SerieNFe            int       `json:"serie_nfe" db:"serie_nfe"`
	SerieNFCe           int       `json:"serie_nfce" db:"serie_nfce"`
	CSCID               string    `json:"csc_id,omitempty" db:"csc_id"`
	CSCToken            string    `json:"csc_token,omitempty" db:"csc_token"`
	SMTPHost            string    `json:"smtp_host,omitempty" db:"smtp_host"`
	SMTPPort            int       `json:"smtp_port,omitempty" db:"smtp_port"`
	SMTPUser            string    `json:"smtp_user,omitempty" db:"smtp_user"`
	SMTPPasswordEnc     string    `json:"-" db:"smtp_password_enc"`
	SMTPFrom            string    `json:"smtp_from,omitempty" db:"smtp_from"`
	SMTPTLS             bool      `json:"smtp_tls" db:"smtp_tls"`
	IsActive            bool      `json:"is_active" db:"is_active"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// UserCompany represents the many-to-many relationship between users and companies.
type UserCompany struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CompanyID uuid.UUID `json:"company_id" db:"company_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
