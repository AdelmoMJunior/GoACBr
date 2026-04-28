package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

// SessionRepository defines data access for Sessions and Token Blacklist.
type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSessionByID(ctx context.Context, id uuid.UUID) (*domain.Session, error)
	RevokeSession(ctx context.Context, id uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
	
	BlacklistToken(ctx context.Context, jti string, expiresAt time.Time) error
	IsTokenBlacklisted(ctx context.Context, jti string) (bool, error)
	CleanupExpiredBlacklist(ctx context.Context) error
}

type sessionRepository struct {
	db *DBWrapper
}

// NewSessionRepository creates a new session repository.
func NewSessionRepository(db *DBWrapper) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, refresh_token_hash, ip_address, user_agent, is_revoked, expires_at, created_at)
		VALUES (:id, :user_id, :refresh_token_hash, :ip_address, :user_agent, :is_revoked, :expires_at, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return fmtDBError(err, "session")
	}
	return nil
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	var session domain.Session
	query := `SELECT * FROM sessions WHERE id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &session, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("session")
		}
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) RevokeSession(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE sessions SET is_revoked = true WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperror.NewNotFound("session")
	}
	return nil
}

func (r *sessionRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE sessions SET is_revoked = true WHERE user_id = $1 AND is_revoked = false`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *sessionRepository) BlacklistToken(ctx context.Context, jti string, expiresAt time.Time) error {
	query := `INSERT INTO token_blacklist (jti, expires_at) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, jti, expiresAt)
	if err != nil {
		return fmtDBError(err, "token_blacklist")
	}
	return nil
}

func (r *sessionRepository) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	var count int
	query := `SELECT COUNT(1) FROM token_blacklist WHERE jti = $1`
	err := r.db.GetContext(ctx, &count, query, jti)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *sessionRepository) CleanupExpiredBlacklist(ctx context.Context) error {
	query := `DELETE FROM token_blacklist WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
