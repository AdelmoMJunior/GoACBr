package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

// InvoiceRepository interface defines data access for Invoices and related tables.
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *domain.Invoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error)
	GetByChave(ctx context.Context, chave string) (*domain.Invoice, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	
	CreateEvent(ctx context.Context, event *domain.InvoiceEvent) error
	CreateInutilizacao(ctx context.Context, inut *domain.InvoiceInutilizacao) error
}

type invoiceRepository struct {
	db *DBWrapper
}

// NewInvoiceRepository creates a new invoice repository.
func NewInvoiceRepository(db *DBWrapper) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *domain.Invoice) error {
	return r.db.Transaction(ctx, func(tx *sqlx.Tx) error {
		// 1. Insert Invoice Header
		queryHeader := `
			INSERT INTO invoices (
				id, company_id, chave, c_nf, nat_op, modelo, serie, numero, dh_emissao, dh_sai_ent,
				tp_nf, id_dest, c_mun_fg, tp_imp, tp_emis, tp_amb, fin_nfe, ind_final, ind_pres,
				proc_emi, ver_proc, ind_intermed, c_mun_fg_ibs, tp_nf_debito, tp_nf_credito,
				tp_ente_gov, p_redutor, tp_oper_gov, dest_cnpj_cpf, dest_nome, dest_ie,
				dest_ind_ie_dest, dest_email, dest_logradouro, dest_numero, dest_complemento,
				dest_bairro, dest_cod_municipio, dest_municipio, dest_uf, dest_cep, tot_v_bc,
				tot_v_icms, tot_v_icms_deson, tot_v_fcp, tot_v_bc_st, tot_v_st, tot_v_fcp_st,
				tot_v_prod, tot_v_frete, tot_v_seg, tot_v_desc, tot_v_ii, tot_v_ipi, tot_v_pis,
				tot_v_cofins, tot_v_outro, tot_v_nf, tot_v_tot_trib, tot_v_ipi_devol, tot_v_fcp_st_ret,
				tot_q_bc_mono, tot_v_icms_mono, tot_q_bc_mono_reten, tot_v_icms_mono_reten,
				tot_q_bc_mono_ret, tot_v_icms_mono_ret, tot_v_nf_tot, tot_is_v_is, tot_ibs_cbs_v_bc,
				tot_ibs_v_ibs, tot_ibs_v_cred_pres, tot_ibs_v_cred_pres_cond_sus, tot_ibs_uf_v_dif,
				tot_ibs_uf_v_dev_trib, tot_ibs_uf_v_ibs_uf, tot_ibs_mun_v_dif, tot_ibs_mun_v_dev_trib,
				tot_ibs_mun_v_ibs_mun, tot_cbs_v_dif, tot_cbs_v_dev_trib, tot_cbs_v_cbs,
				tot_cbs_v_cred_pres, tot_cbs_v_cred_pres_cond_sus, tot_mono_v_ibs_mono,
				tot_mono_v_cbs_mono, tot_mono_v_ibs_mono_reten, tot_mono_v_cbs_mono_reten,
				tot_mono_v_ibs_mono_ret, tot_mono_v_cbs_mono_ret, protocolo, dh_recebimento,
				status, inf_ad_fisco, inf_cpl, xml_b2_key, pdf_b2_key, created_at, updated_at
			) VALUES (
				:id, :company_id, :chave, :c_nf, :nat_op, :modelo, :serie, :numero, :dh_emissao, :dh_sai_ent,
				:tp_nf, :id_dest, :c_mun_fg, :tp_imp, :tp_emis, :tp_amb, :fin_nfe, :ind_final, :ind_pres,
				:proc_emi, :ver_proc, :ind_intermed, :c_mun_fg_ibs, :tp_nf_debito, :tp_nf_credito,
				:tp_ente_gov, :p_redutor, :tp_oper_gov, :dest_cnpj_cpf, :dest_nome, :dest_ie,
				:dest_ind_ie_dest, :dest_email, :dest_logradouro, :dest_numero, :dest_complemento,
				:dest_bairro, :dest_cod_municipio, :dest_municipio, :dest_uf, :dest_cep, :tot_v_bc,
				:tot_v_icms, :tot_v_icms_deson, :tot_v_fcp, :tot_v_bc_st, :tot_v_st, :tot_v_fcp_st,
				:tot_v_prod, :tot_v_frete, :tot_v_seg, :tot_v_desc, :tot_v_ii, :tot_v_ipi, :tot_v_pis,
				:tot_v_cofins, :tot_v_outro, :tot_v_nf, :tot_v_tot_trib, :tot_v_ipi_devol, :tot_v_fcp_st_ret,
				:tot_q_bc_mono, :tot_v_icms_mono, :tot_q_bc_mono_reten, :tot_v_icms_mono_reten,
				:tot_q_bc_mono_ret, :tot_v_icms_mono_ret, :tot_v_nf_tot, :tot_is_v_is, :tot_ibs_cbs_v_bc,
				:tot_ibs_v_ibs, :tot_ibs_v_cred_pres, :tot_ibs_v_cred_pres_cond_sus, :tot_ibs_uf_v_dif,
				:tot_ibs_uf_v_dev_trib, :tot_ibs_uf_v_ibs_uf, :tot_ibs_mun_v_dif, :tot_ibs_mun_v_dev_trib,
				:tot_ibs_mun_v_ibs_mun, :tot_cbs_v_dif, :tot_cbs_v_dev_trib, :tot_cbs_v_cbs,
				:tot_cbs_v_cred_pres, :tot_cbs_v_cred_pres_cond_sus, :tot_mono_v_ibs_mono,
				:tot_mono_v_cbs_mono, :tot_mono_v_ibs_mono_reten, :tot_mono_v_cbs_mono_reten,
				:tot_mono_v_ibs_mono_ret, :tot_mono_v_cbs_mono_ret, :protocolo, :dh_recebimento,
				:status, :inf_ad_fisco, :inf_cpl, :xml_b2_key, :pdf_b2_key, :created_at, :updated_at
			)
		`
		_, err := tx.NamedExec(queryHeader, invoice)
		if err != nil {
			return fmtDBError(err, "invoice")
		}

		// 2. Insert Items (simplificado para exemplo, usar NamedExec em lote)
		// Em produção seria melhor fazer um "bulk insert", mas para a estrutura de NamedExec:
		queryItem := `
			INSERT INTO invoice_items (
				id, invoice_id, n_item, c_prod, c_ean, x_prod, ncm, cest, cfop, u_com, q_com,
				v_un_com, v_prod, c_ean_trib, u_trib, q_trib, v_un_trib, v_frete, v_seg, v_desc,
				v_outro, ind_tot, x_ped, n_item_ped, v_tot_trib, inf_ad_prod, ind_escala, cnpj_fab,
				c_benef, ind_bem_movel_usado, v_item, icms_orig, icms_cst, icms_csosn, icms_mod_bc,
				icms_p_red_bc, icms_v_bc, icms_p_icms, icms_v_icms, icms_mod_bc_st, icms_p_mva_st,
				icms_p_red_bc_st, icms_v_bc_st, icms_p_icms_st, icms_v_icms_st, icms_uf_st, icms_p_bc_op,
				icms_v_bc_st_ret, icms_v_icms_st_ret, icms_mot_des, icms_p_cred_sn, icms_v_cred_icms_sn,
				icms_v_icms_deson, icms_v_icms_op, icms_p_dif, icms_v_icms_dif, icms_p_st, icms_v_bc_fcp,
				icms_p_fcp, icms_v_fcp, icms_v_bc_fcp_st_ret, icms_p_fcp_st_ret, icms_v_fcp_st_ret,
				icms_p_red_bc_efet, icms_v_bc_efet, icms_p_icms_efet, icms_v_icms_efet, icms_v_icms_substituto,
				icms_q_bc_mono, icms_ad_rem_icms, icms_v_icms_mono, icms_q_bc_mono_reten, icms_ad_rem_icms_reten,
				icms_v_icms_mono_reten, icms_p_red_ad_rem, icms_mot_red_ad_rem, icms_q_bc_mono_ret,
				icms_v_icms_mono_op, icms_v_icms_mono_dif, icms_ad_rem_icms_ret, icms_v_icms_mono_ret,
				icms_uf_dest_v_bc, icms_uf_dest_v_bc_fcp, icms_uf_dest_p_fcp, icms_uf_dest_p_icms,
				icms_uf_dest_p_icms_inter, icms_uf_dest_p_icms_inter_part, icms_uf_dest_v_fcp,
				icms_uf_dest_v_icms, icms_uf_dest_v_icms_remet, pis_cst, pis_v_bc, pis_p_pis,
				pis_q_bc_prod, pis_v_aliq_prod, pis_v_pis, cofins_cst, cofins_v_bc, cofins_p_cofins,
				cofins_q_bc_prod, cofins_v_aliq_prod, cofins_v_cofins, ii_v_bc, ii_v_desp_adu, ii_v_ii,
				ii_v_iof, ipi_cst, ipi_c_enq, ipi_v_bc, ipi_p_ipi, ipi_v_ipi, is_cst, is_c_class_trib,
				is_v_bc, is_p_is, is_p_is_espec, is_v_is, ibs_cbs_cst, ibs_cbs_c_class_trib,
				ibs_cbs_v_bc, ibs_cbs_v_ibs, ibs_uf_p, ibs_uf_v, ibs_uf_p_dif, ibs_uf_v_dif,
				ibs_uf_v_dev_trib, ibs_uf_p_red_aliq, ibs_uf_p_aliq_efet, ibs_mun_p, ibs_mun_v,
				ibs_mun_p_dif, ibs_mun_v_dif, ibs_mun_v_dev_trib, ibs_mun_p_red_aliq, ibs_mun_p_aliq_efet,
				cbs_p, cbs_v, cbs_p_dif, cbs_v_dif, cbs_v_dev_trib, cbs_p_red_aliq, cbs_p_aliq_efet,
				trib_reg_cst, trib_reg_c_class_trib, trib_reg_p_aliq_ibs_uf, trib_reg_v_ibs_uf,
				trib_reg_p_aliq_ibs_mun, trib_reg_v_ibs_mun, trib_reg_p_aliq_cbs, trib_reg_v_cbs,
				ibs_cred_pres_cod, ibs_cred_pres_p, ibs_cred_pres_v, cbs_cred_pres_cod, cbs_cred_pres_p,
				cbs_cred_pres_v, gov_p_aliq_ibs_uf, gov_v_trib_ibs_uf, gov_p_aliq_ibs_mun,
				gov_v_trib_ibs_mun, gov_p_aliq_cbs, gov_v_trib_cbs, mono_v_tot_ibs, mono_v_tot_cbs,
				mono_pad_q_bc, mono_pad_ad_rem_ibs, mono_pad_ad_rem_cbs, mono_pad_v_ibs, mono_pad_v_cbs,
				mono_reten_q_bc, mono_reten_ad_rem_ibs, mono_reten_v_ibs, mono_reten_ad_rem_cbs,
				mono_reten_v_cbs, mono_ret_q_bc, mono_ret_ad_rem_ibs, mono_ret_v_ibs, mono_ret_ad_rem_cbs,
				mono_ret_v_cbs, mono_dif_p_ibs, mono_dif_v_ibs, mono_dif_p_cbs, mono_dif_v_cbs,
				transf_cred_v_ibs, transf_cred_v_cbs, cred_pres_ibs_zfm_tp, cred_pres_ibs_zfm_v
			) VALUES (
				:id, :invoice_id, :n_item, :c_prod, :c_ean, :x_prod, :ncm, :cest, :cfop, :u_com, :q_com,
				:v_un_com, :v_prod, :c_ean_trib, :u_trib, :q_trib, :v_un_trib, :v_frete, :v_seg, :v_desc,
				:v_outro, :ind_tot, :x_ped, :n_item_ped, :v_tot_trib, :inf_ad_prod, :ind_escala, :cnpj_fab,
				:c_benef, :ind_bem_movel_usado, :v_item, :icms_orig, :icms_cst, :icms_csosn, :icms_mod_bc,
				:icms_p_red_bc, :icms_v_bc, :icms_p_icms, :icms_v_icms, :icms_mod_bc_st, :icms_p_mva_st,
				:icms_p_red_bc_st, :icms_v_bc_st, :icms_p_icms_st, :icms_v_icms_st, :icms_uf_st, :icms_p_bc_op,
				:icms_v_bc_st_ret, :icms_v_icms_st_ret, :icms_mot_des, :icms_p_cred_sn, :icms_v_cred_icms_sn,
				:icms_v_icms_deson, :icms_v_icms_op, :icms_p_dif, :icms_v_icms_dif, :icms_p_st, :icms_v_bc_fcp,
				:icms_p_fcp, :icms_v_fcp, :icms_v_bc_fcp_st_ret, :icms_p_fcp_st_ret, :icms_v_fcp_st_ret,
				:icms_p_red_bc_efet, :icms_v_bc_efet, :icms_p_icms_efet, :icms_v_icms_efet, :icms_v_icms_substituto,
				:icms_q_bc_mono, :icms_ad_rem_icms, :icms_v_icms_mono, :icms_q_bc_mono_reten, :icms_ad_rem_icms_reten,
				:icms_v_icms_mono_reten, :icms_p_red_ad_rem, :icms_mot_red_ad_rem, :icms_q_bc_mono_ret,
				:icms_v_icms_mono_op, :icms_v_icms_mono_dif, :icms_ad_rem_icms_ret, :icms_v_icms_mono_ret,
				:icms_uf_dest_v_bc, :icms_uf_dest_v_bc_fcp, :icms_uf_dest_p_fcp, :icms_uf_dest_p_icms,
				:icms_uf_dest_p_icms_inter, :icms_uf_dest_p_icms_inter_part, :icms_uf_dest_v_fcp,
				:icms_uf_dest_v_icms, :icms_uf_dest_v_icms_remet, :pis_cst, :pis_v_bc, :pis_p_pis,
				:pis_q_bc_prod, :pis_v_aliq_prod, :pis_v_pis, :cofins_cst, :cofins_v_bc, :cofins_p_cofins,
				:cofins_q_bc_prod, :cofins_v_aliq_prod, :cofins_v_cofins, :ii_v_bc, :ii_v_desp_adu, :ii_v_ii,
				:ii_v_iof, :ipi_cst, :ipi_c_enq, :ipi_v_bc, :ipi_p_ipi, :ipi_v_ipi, :is_cst, :is_c_class_trib,
				:is_v_bc, :is_p_is, :is_p_is_espec, :is_v_is, :ibs_cbs_cst, :ibs_cbs_c_class_trib,
				:ibs_cbs_v_bc, :ibs_cbs_v_ibs, :ibs_uf_p, :ibs_uf_v, :ibs_uf_p_dif, :ibs_uf_v_dif,
				:ibs_uf_v_dev_trib, :ibs_uf_p_red_aliq, :ibs_uf_p_aliq_efet, :ibs_mun_p, :ibs_mun_v,
				:ibs_mun_p_dif, :ibs_mun_v_dif, :ibs_mun_v_dev_trib, :ibs_mun_p_red_aliq, :ibs_mun_p_aliq_efet,
				:cbs_p, :cbs_v, :cbs_p_dif, :cbs_v_dif, :cbs_v_dev_trib, :cbs_p_red_aliq, :cbs_p_aliq_efet,
				:trib_reg_cst, :trib_reg_c_class_trib, :trib_reg_p_aliq_ibs_uf, :trib_reg_v_ibs_uf,
				:trib_reg_p_aliq_ibs_mun, :trib_reg_v_ibs_mun, :trib_reg_p_aliq_cbs, :trib_reg_v_cbs,
				:ibs_cred_pres_cod, :ibs_cred_pres_p, :ibs_cred_pres_v, :cbs_cred_pres_cod, :cbs_cred_pres_p,
				:cbs_cred_pres_v, :gov_p_aliq_ibs_uf, :gov_v_trib_ibs_uf, :gov_p_aliq_ibs_mun,
				:gov_v_trib_ibs_mun, :gov_p_aliq_cbs, :gov_v_trib_cbs, :mono_v_tot_ibs, :mono_v_tot_cbs,
				:mono_pad_q_bc, :mono_pad_ad_rem_ibs, :mono_pad_ad_rem_cbs, :mono_pad_v_ibs, :mono_pad_v_cbs,
				:mono_reten_q_bc, :mono_reten_ad_rem_ibs, :mono_reten_v_ibs, :mono_reten_ad_rem_cbs,
				:mono_reten_v_cbs, :mono_ret_q_bc, :mono_ret_ad_rem_ibs, :mono_ret_v_ibs, :mono_ret_ad_rem_cbs,
				:mono_ret_v_cbs, :mono_dif_p_ibs, :mono_dif_v_ibs, :mono_dif_p_cbs, :mono_dif_v_cbs,
				:transf_cred_v_ibs, :transf_cred_v_cbs, :cred_pres_ibs_zfm_tp, :cred_pres_ibs_zfm_v
			)
		`
		for i := range invoice.Items {
			invoice.Items[i].InvoiceID = invoice.ID
		}
		if len(invoice.Items) > 0 {
			_, err = tx.NamedExec(queryItem, invoice.Items)
			if err != nil {
				return fmtDBError(err, "invoice_item")
			}
		}

		// 3. Insert Payments
		queryPayment := `
			INSERT INTO invoice_payments (
				id, invoice_id, n_pag, tp_pag, x_pag, v_pag, ind_pag, tp_integra, cnpj_pag, t_band, c_aut, v_troco
			) VALUES (
				:id, :invoice_id, :n_pag, :tp_pag, :x_pag, :v_pag, :ind_pag, :tp_integra, :cnpj_pag, :t_band, :c_aut, :v_troco
			)
		`
		for i := range invoice.Payments {
			invoice.Payments[i].InvoiceID = invoice.ID
		}
		if len(invoice.Payments) > 0 {
			_, err = tx.NamedExec(queryPayment, invoice.Payments)
			if err != nil {
				return fmtDBError(err, "invoice_payment")
			}
		}

		// 4. Insert Transport
		if invoice.Transport != nil {
			invoice.Transport.InvoiceID = invoice.ID
			queryTransport := `
				INSERT INTO invoice_transport (
					id, invoice_id, mod_frete, transp_cnpj_cpf, transp_nome, transp_ie, transp_endereco,
					transp_municipio, transp_uf, v_serv, v_bc_ret, p_icms_ret, v_icms_ret, placa,
					uf_placa, q_vol, esp_vol, peso_l, peso_b
				) VALUES (
					:id, :invoice_id, :mod_frete, :transp_cnpj_cpf, :transp_nome, :transp_ie, :transp_endereco,
					:transp_municipio, :transp_uf, :v_serv, :v_bc_ret, :p_icms_ret, :v_icms_ret, :placa,
					:uf_placa, :q_vol, :esp_vol, :peso_l, :peso_b
				)
			`
			_, err = tx.NamedExec(queryTransport, invoice.Transport)
			if err != nil {
				return fmtDBError(err, "invoice_transport")
			}
		}

		return nil
	})
}

func (r *invoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	var invoice domain.Invoice
	query := `SELECT * FROM invoices WHERE id = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &invoice, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("invoice")
		}
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) GetByChave(ctx context.Context, chave string) (*domain.Invoice, error) {
	var invoice domain.Invoice
	query := `SELECT * FROM invoices WHERE chave = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &invoice, query, chave)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("invoice")
		}
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE invoices SET status = $1, updated_at = NOW() WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperror.NewNotFound("invoice")
	}
	return nil
}

func (r *invoiceRepository) CreateEvent(ctx context.Context, event *domain.InvoiceEvent) error {
	query := `
		INSERT INTO invoice_events (
			id, invoice_id, company_id, chave_nfe, c_orgao, tp_evento, n_seq_evento, dh_evento,
			protocolo, x_just, x_correcao, tp_nf, dest_cnpj_cpf, dest_uf, v_nf, v_icms, v_st, xml_b2_key, created_at
		) VALUES (
			:id, :invoice_id, :company_id, :chave_nfe, :c_orgao, :tp_evento, :n_seq_evento, :dh_evento,
			:protocolo, :x_just, :x_correcao, :tp_nf, :dest_cnpj_cpf, :dest_uf, :v_nf, :v_icms, :v_st, :xml_b2_key, :created_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmtDBError(err, "invoice_event")
	}
	return nil
}

func (r *invoiceRepository) CreateInutilizacao(ctx context.Context, inut *domain.InvoiceInutilizacao) error {
	query := `
		INSERT INTO invoice_inutilizacao (
			id, company_id, ano, modelo, serie, num_inicial, num_final, justificativa, protocolo,
			dh_recebimento, status, xml_b2_key, created_at
		) VALUES (
			:id, :company_id, :ano, :modelo, :serie, :num_inicial, :num_final, :justificativa, :protocolo,
			:dh_recebimento, :status, :xml_b2_key, :created_at
		)
	`
	_, err := r.db.NamedExecContext(ctx, query, inut)
	if err != nil {
		return fmtDBError(err, "invoice_inutilizacao")
	}
	return nil
}
