package repository

import (
	"context"

	"github.com/AdelmoMJunior/GoACBr/internal/domain"
)

// AuditRepository defines data access for Audit Logs.
type AuditRepository interface {
	Log(ctx context.Context, log *domain.AuditLog) error
}

type auditRepository struct {
	db *DBWrapper
}

// NewAuditRepository creates a new audit repository.
func NewAuditRepository(db *DBWrapper) AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Log(ctx context.Context, log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_log (id, user_id, company_cnpj, action, resource, details, ip_address, created_at)
		VALUES (:id, :user_id, :company_cnpj, :action, :resource, :details, :ip_address, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, log)
	if err != nil {
		return fmtDBError(err, "audit_log")
	}
	return nil
}
