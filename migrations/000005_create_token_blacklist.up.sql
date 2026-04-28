CREATE TABLE IF NOT EXISTS token_blacklist (
    jti VARCHAR(255) PRIMARY KEY,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_token_blacklist_expires_at ON token_blacklist(expires_at);
