CREATE TABLE IF NOT EXISTS cache_fallback (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_cache_fallback_expires_at ON cache_fallback(expires_at);

CREATE TABLE IF NOT EXISTS audit_log (
    id UUID PRIMARY KEY,
    user_id UUID,
    company_cnpj VARCHAR(14),
    action VARCHAR(255) NOT NULL,
    resource VARCHAR(255) NOT NULL,
    details JSONB,
    ip_address INET,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_log_company_cnpj ON audit_log(company_cnpj);
CREATE INDEX idx_audit_log_created_at ON audit_log(created_at);
