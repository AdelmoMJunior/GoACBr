CREATE TABLE IF NOT EXISTS invoice_inutilizacao (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    ano INT NOT NULL,
    modelo SMALLINT NOT NULL,
    serie INT NOT NULL,
    num_inicial INT NOT NULL,
    num_final INT NOT NULL,
    justificativa TEXT NOT NULL,
    protocolo VARCHAR(255),
    dh_recebimento TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL,
    xml_b2_key VARCHAR(500),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_invoice_inutilizacao_company_id ON invoice_inutilizacao(company_id);
