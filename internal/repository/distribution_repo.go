package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

// DistributionRepository interface defines data access for DFe Distribution.
type DistributionRepository interface {
	SaveDocument(ctx context.Context, doc *domain.DistributionDocument) error
	GetControl(ctx context.Context, companyID uuid.UUID) (*domain.DistributionControl, error)
	UpsertControl(ctx context.Context, control *domain.DistributionControl) error
}

type distributionRepository struct {
	db *DBWrapper
}

// NewDistributionRepository creates a new distribution repository.
func NewDistributionRepository(db *DBWrapper) DistributionRepository {
	return &distributionRepository{db: db}
}

func (r *distributionRepository) SaveDocument(ctx context.Context, doc *domain.DistributionDocument) error {
	query := `
		INSERT INTO distribution_documents (
			id, company_id, nsu, schema_type, chave_nfe, tp_nf, emit_cnpj_cpf, emit_nome, emit_ie,
			dest_cnpj_cpf, dh_emissao, modelo, serie, numero, nat_op, c_sit_nfe, tot_v_nf, tot_v_icms,
			tot_v_st, tot_v_pis, tot_v_cofins, tot_v_prod, tot_v_desc, tot_v_frete, tot_v_outro,
			tp_evento, desc_evento, n_seq_evento, dh_evento, x_just, x_correcao, protocolo,
			dh_recebimento, xml_b2_key, created_at
		) VALUES (
			:id, :company_id, :nsu, :schema_type, :chave_nfe, :tp_nf, :emit_cnpj_cpf, :emit_nome, :emit_ie,
			:dest_cnpj_cpf, :dh_emissao, :modelo, :serie, :numero, :nat_op, :c_sit_nfe, :tot_v_nf, :tot_v_icms,
			:tot_v_st, :tot_v_pis, :tot_v_cofins, :tot_v_prod, :tot_v_desc, :tot_v_frete, :tot_v_outro,
			:tp_evento, :desc_evento, :n_seq_evento, :dh_evento, :x_just, :x_correcao, :protocolo,
			:dh_recebimento, :xml_b2_key, :created_at
		) ON CONFLICT (company_id, nsu) DO NOTHING
	`
	_, err := r.db.NamedExecContext(ctx, query, doc)
	if err != nil {
		return fmtDBError(err, "distribution_document")
	}
	return nil
}

func (r *distributionRepository) GetControl(ctx context.Context, companyID uuid.UUID) (*domain.DistributionControl, error) {
	var control domain.DistributionControl
	query := `SELECT * FROM distribution_control WHERE company_id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &control, query, companyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("distribution_control")
		}
		return nil, err
	}
	return &control, nil
}

func (r *distributionRepository) UpsertControl(ctx context.Context, control *domain.DistributionControl) error {
	query := `
		INSERT INTO distribution_control (
			company_id, last_nsu, max_nsu, last_query_at, is_running, status, error_message, updated_at
		) VALUES (
			:company_id, :last_nsu, :max_nsu, :last_query_at, :is_running, :status, :error_message, :updated_at
		)
		ON CONFLICT (company_id) DO UPDATE SET
			last_nsu = EXCLUDED.last_nsu,
			max_nsu = EXCLUDED.max_nsu,
			last_query_at = EXCLUDED.last_query_at,
			is_running = EXCLUDED.is_running,
			status = EXCLUDED.status,
			error_message = EXCLUDED.error_message,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.NamedExecContext(ctx, query, control)
	if err != nil {
		return fmtDBError(err, "distribution_control")
	}
	return nil
}
