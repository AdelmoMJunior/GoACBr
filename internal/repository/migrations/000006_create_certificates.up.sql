CREATE TABLE IF NOT EXISTS certificates (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL UNIQUE REFERENCES companies(id) ON DELETE CASCADE,
    pfx_data BYTEA NOT NULL,
    pfx_password_enc VARCHAR(255) NOT NULL,
    subject_cn VARCHAR(255),
    serial_number VARCHAR(100),
    valid_from TIMESTAMPTZ,
    valid_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
