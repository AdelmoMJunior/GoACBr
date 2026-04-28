package repository

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/AdelmoMJunior/GoACBr/internal/config"
)

var migrationsFS embed.FS

// DBWrapper wraps the sqlx.DB connection.
type DBWrapper struct {
	*sqlx.DB
}

// NewDB initialize the PostgreSQL database connection via sqlx.
func NewDB(cfg config.DatabaseConfig) (*DBWrapper, error) {
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	slog.Info("Connected to PostgreSQL database successfully", "host", cfg.Host, "db", cfg.Name)

	wrapper := &DBWrapper{DB: db}
	if err := wrapper.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return wrapper, nil
}

func (db *DBWrapper) runMigrations() error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create iofs source: %w", err)
	}

	driver, err := postgres.WithInstance(db.DB.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	slog.Info("Database migrations applied successfully")
	return nil
}

// Transaction executes a function within a database transaction.
func (db *DBWrapper) Transaction(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			slog.Error("Failed to rollback transaction", "error", rbErr, "original_error", err)
		}
		return err
	}

	if cmErr := tx.Commit(); cmErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", cmErr)
	}

	return nil
}
