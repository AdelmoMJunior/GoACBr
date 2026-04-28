CREATE TABLE IF NOT EXISTS distribution_documents (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    nsu VARCHAR(50) NOT NULL,
    schema_type VARCHAR(50) NOT NULL,
    chave_nfe VARCHAR(44),
    tp_nf SMALLINT,
    emit_cnpj_cpf VARCHAR(14),
    emit_nome VARCHAR(255),
    emit_ie VARCHAR(50),
    dest_cnpj_cpf VARCHAR(14),
    dh_emissao TIMESTAMPTZ,
    modelo SMALLINT,
    serie INT,
    numero INT,
    nat_op VARCHAR(255),
    c_sit_nfe SMALLINT,
    tot_v_nf DECIMAL(15,2),
    tot_v_icms DECIMAL(15,2),
    tot_v_st DECIMAL(15,2),
    tot_v_pis DECIMAL(15,2),
    tot_v_cofins DECIMAL(15,2),
    tot_v_prod DECIMAL(15,2),
    tot_v_desc DECIMAL(15,2),
    tot_v_frete DECIMAL(15,2),
    tot_v_outro DECIMAL(15,2),
    tp_evento VARCHAR(6),
    desc_evento VARCHAR(255),
    n_seq_evento INT,
    dh_evento TIMESTAMPTZ,
    x_just TEXT,
    x_correcao TEXT,
    protocolo VARCHAR(255),
    dh_recebimento TIMESTAMPTZ,
    xml_b2_key VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE (company_id, nsu)
);

CREATE INDEX idx_dist_docs_company_id_nsu ON distribution_documents(company_id, nsu);
CREATE INDEX idx_dist_docs_company_id_chave ON distribution_documents(company_id, chave_nfe);

CREATE TABLE IF NOT EXISTS distribution_control (
    company_id UUID PRIMARY KEY REFERENCES companies(id) ON DELETE CASCADE,
    last_nsu VARCHAR(50) NOT NULL,
    max_nsu VARCHAR(50) NOT NULL,
    last_query_at TIMESTAMPTZ,
    is_running BOOLEAN NOT NULL DEFAULT false,
    status VARCHAR(50) NOT NULL,
    error_message TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
