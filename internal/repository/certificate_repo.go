package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

// CertificateRepository interface defines data access for Certificates.
type CertificateRepository interface {
	Save(ctx context.Context, cert *domain.Certificate) error
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*domain.Certificate, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type certificateRepository struct {
	db *DBWrapper
}

// NewCertificateRepository creates a new certificate repository.
func NewCertificateRepository(db *DBWrapper) CertificateRepository {
	return &certificateRepository{db: db}
}

func (r *certificateRepository) Save(ctx context.Context, cert *domain.Certificate) error {
	query := `
		INSERT INTO certificates (
			id, company_id, pfx_data, pfx_password_enc, subject_cn, serial_number, valid_from, valid_until, created_at, updated_at
		) VALUES (
			:id, :company_id, :pfx_data, :pfx_password_enc, :subject_cn, :serial_number, :valid_from, :valid_until, :created_at, :updated_at
		)
		ON CONFLICT (company_id) DO UPDATE SET
			pfx_data = EXCLUDED.pfx_data,
			pfx_password_enc = EXCLUDED.pfx_password_enc,
			subject_cn = EXCLUDED.subject_cn,
			serial_number = EXCLUDED.serial_number,
			valid_from = EXCLUDED.valid_from,
			valid_until = EXCLUDED.valid_until,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.NamedExecContext(ctx, query, cert)
	if err != nil {
		return fmtDBError(err, "certificate")
	}
	return nil
}

func (r *certificateRepository) GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*domain.Certificate, error) {
	var cert domain.Certificate
	query := `SELECT * FROM certificates WHERE company_id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &cert, query, companyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("certificate")
		}
		return nil, err
	}
	return &cert, nil
}

func (r *certificateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM certificates WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
