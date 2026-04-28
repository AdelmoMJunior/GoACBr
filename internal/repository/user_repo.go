package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

// UserRepository interface defines data access for Users.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

type userRepository struct {
	db *DBWrapper
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *DBWrapper) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, full_name, phone, is_active, created_at, updated_at)
		VALUES (:id, :email, :password_hash, :full_name, :phone, :is_active, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmtDBError(err, "user")
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("user")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE email = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("user")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET email = :email, password_hash = :password_hash, full_name = :full_name, 
		    phone = :phone, is_active = :is_active, updated_at = :updated_at
		WHERE id = :id
	`
	res, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmtDBError(err, "user")
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperror.NewNotFound("user")
	}

	return nil
}

// fmtDBError wraps database errors into standard application errors.
func fmtDBError(err error, entity string) error {
	// Simple check for duplicate keys (PostgreSQL error 23505)
	if err != nil && err.Error() != "" && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return apperror.NewAlreadyExists(entity)
	}
	return err
}
