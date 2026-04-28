DROP INDEX IF EXISTS idx_audit_log_created_at;
DROP INDEX IF EXISTS idx_audit_log_company_cnpj;
DROP TABLE IF EXISTS audit_log;

DROP INDEX IF EXISTS idx_cache_fallback_expires_at;
DROP TABLE IF EXISTS cache_fallback;
