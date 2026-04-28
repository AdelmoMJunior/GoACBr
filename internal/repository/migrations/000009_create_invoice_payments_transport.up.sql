CREATE TABLE IF NOT EXISTS invoice_payments (
    id UUID PRIMARY KEY,
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    n_pag INT NOT NULL,
    tp_pag VARCHAR(2) NOT NULL,
    x_pag VARCHAR(255),
    v_pag DECIMAL(15,2) NOT NULL,
    ind_pag SMALLINT,
    tp_integra SMALLINT,
    cnpj_pag VARCHAR(14),
    t_band VARCHAR(2),
    c_aut VARCHAR(255),
    v_troco DECIMAL(15,2) NOT NULL DEFAULT 0
);

CREATE INDEX idx_invoice_payments_invoice_id ON invoice_payments(invoice_id);

CREATE TABLE IF NOT EXISTS invoice_transport (
    id UUID PRIMARY KEY,
    invoice_id UUID NOT NULL UNIQUE REFERENCES invoices(id) ON DELETE CASCADE,
    mod_frete SMALLINT NOT NULL,
    transp_cnpj_cpf VARCHAR(14),
    transp_nome VARCHAR(255),
    transp_ie VARCHAR(50),
    transp_endereco VARCHAR(255),
    transp_municipio VARCHAR(255),
    transp_uf VARCHAR(2),
    v_serv DECIMAL(15,2),
    v_bc_ret DECIMAL(15,2),
    p_icms_ret DECIMAL(7,4),
    v_icms_ret DECIMAL(15,2),
    placa VARCHAR(50),
    uf_placa VARCHAR(2),
    q_vol INT,
    esp_vol VARCHAR(255),
    peso_l DECIMAL(15,3),
    peso_b DECIMAL(15,3)
);

CREATE INDEX idx_invoice_transport_invoice_id ON invoice_transport(invoice_id);
