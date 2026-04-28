package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

// CompanyRepository interface defines data access for Companies.
type CompanyRepository interface {
	Create(ctx context.Context, company *domain.Company) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error)
	GetByCNPJ(ctx context.Context, cnpj string) (*domain.Company, error)
	Update(ctx context.Context, company *domain.Company) error
	LinkUser(ctx context.Context, userID, companyID uuid.UUID) error
	GetUsersByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.User, error)
	GetCompaniesByUser(ctx context.Context, userID uuid.UUID) ([]domain.Company, error)
	GetCompaniesEligibleForSync(ctx context.Context) ([]domain.Company, error)
}

type companyRepository struct {
	db *DBWrapper
}

// NewCompanyRepository creates a new company repository.
func NewCompanyRepository(db *DBWrapper) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Create(ctx context.Context, company *domain.Company) error {
	query := `
		INSERT INTO companies (
			id, cnpj, razao_social, nome_fantasia, inscricao_estadual, inscricao_municipal,
			crt, logradouro, numero, complemento, bairro, cod_municipio, municipio, uf, cep,
			telefone, cnae, ambiente, serie_nfe, serie_nfce, csc_id, csc_token, smtp_host,
			smtp_port, smtp_user, smtp_password_enc, smtp_from, smtp_tls, is_active, created_at, updated_at
		) VALUES (
			:id, :cnpj, :razao_social, :nome_fantasia, :inscricao_estadual, :inscricao_municipal,
			:crt, :logradouro, :numero, :complemento, :bairro, :cod_municipio, :municipio, :uf, :cep,
			:telefone, :cnae, :ambiente, :serie_nfe, :serie_nfce, :csc_id, :csc_token, :smtp_host,
			:smtp_port, :smtp_user, :smtp_password_enc, :smtp_from, :smtp_tls, :is_active, :created_at, :updated_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, company)
	if err != nil {
		return fmtDBError(err, "company")
	}
	return nil
}

func (r *companyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	var company domain.Company
	query := `SELECT * FROM companies WHERE id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &company, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("company")
		}
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) GetByCNPJ(ctx context.Context, cnpj string) (*domain.Company, error) {
	var company domain.Company
	query := `SELECT * FROM companies WHERE cnpj = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &company, query, cnpj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("company")
		}
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) Update(ctx context.Context, company *domain.Company) error {
	query := `
		UPDATE companies SET
			razao_social = :razao_social, nome_fantasia = :nome_fantasia,
			inscricao_estadual = :inscricao_estadual, inscricao_municipal = :inscricao_municipal,
			crt = :crt, logradouro = :logradouro, numero = :numero, complemento = :complemento,
			bairro = :bairro, cod_municipio = :cod_municipio, municipio = :municipio, uf = :uf, cep = :cep,
			telefone = :telefone, cnae = :cnae, ambiente = :ambiente, serie_nfe = :serie_nfe,
			serie_nfce = :serie_nfce, csc_id = :csc_id, csc_token = :csc_token, smtp_host = :smtp_host,
			smtp_port = :smtp_port, smtp_user = :smtp_user, smtp_password_enc = :smtp_password_enc,
			smtp_from = :smtp_from, smtp_tls = :smtp_tls, is_active = :is_active, updated_at = :updated_at
		WHERE id = :id
	`
	res, err := r.db.NamedExecContext(ctx, query, company)
	if err != nil {
		return fmtDBError(err, "company")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperror.NewNotFound("company")
	}

	return nil
}

func (r *companyRepository) LinkUser(ctx context.Context, userID, companyID uuid.UUID) error {
	query := `INSERT INTO user_companies (user_id, company_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, companyID)
	return err
}

func (r *companyRepository) GetUsersByCompany(ctx context.Context, companyID uuid.UUID) ([]domain.User, error) {
	var users []domain.User
	query := `
		SELECT u.* FROM users u
		JOIN user_companies uc ON u.id = uc.user_id
		WHERE uc.company_id = $1
	`
	err := r.db.SelectContext(ctx, &users, query, companyID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *companyRepository) GetCompaniesByUser(ctx context.Context, userID uuid.UUID) ([]domain.Company, error) {
	var companies []domain.Company
	query := `
		SELECT c.* FROM companies c
		JOIN user_companies uc ON c.id = uc.company_id
		WHERE uc.user_id = $1
	`
	err := r.db.SelectContext(ctx, &companies, query, userID)
	if err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *companyRepository) GetCompaniesEligibleForSync(ctx context.Context) ([]domain.Company, error) {
	var companies []domain.Company
	query := `
		SELECT DISTINCT c.*
		FROM companies c
		JOIN certificates cert ON c.id = cert.company_id
		LEFT JOIN distribution_control dc ON c.id = dc.company_id
		WHERE c.is_active = true
		  AND cert.valid_until > NOW()
		  AND (
			  dc.last_query_at IS NULL 
			  OR dc.last_query_at < NOW() - INTERVAL '1 hour'
			  OR dc.last_nsu != dc.max_nsu
		  )
	`
	err := r.db.SelectContext(ctx, &companies, query)
	if err != nil {
		return nil, err
	}
	return companies, nil
}
