CREATE TABLE IF NOT EXISTS invoice_events (
    id UUID PRIMARY KEY,
    invoice_id UUID REFERENCES invoices(id) ON DELETE SET NULL,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    chave_nfe VARCHAR(44) NOT NULL,
    c_orgao SMALLINT NOT NULL,
    tp_evento VARCHAR(6) NOT NULL,
    n_seq_evento INT NOT NULL,
    dh_evento TIMESTAMPTZ NOT NULL,
    protocolo VARCHAR(255),
    x_just TEXT,
    x_correcao TEXT,
    tp_nf SMALLINT,
    dest_cnpj_cpf VARCHAR(14),
    dest_uf VARCHAR(2),
    v_nf DECIMAL(15,2),
    v_icms DECIMAL(15,2),
    v_st DECIMAL(15,2),
    xml_b2_key VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoice_events_company_id ON invoice_events(company_id);
CREATE INDEX idx_invoice_events_chave_nfe ON invoice_events(chave_nfe);
