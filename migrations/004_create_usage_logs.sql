CREATE TABLE IF NOT EXISTS usage_logs (
    id                BIGSERIAL PRIMARY KEY,
    user_key_id       BIGINT NOT NULL REFERENCES user_api_keys(id),
    user_id           VARCHAR(100) NOT NULL,
    provider_id       BIGINT NOT NULL REFERENCES providers(id),
    model_name        VARCHAR(100) NOT NULL,
    prompt_tokens     INTEGER DEFAULT 0,
    completion_tokens INTEGER DEFAULT 0,
    total_tokens      INTEGER DEFAULT 0,
    latency_ms        INTEGER,
    status            VARCHAR(20) NOT NULL DEFAULT 'success',
    error_message     TEXT,
    request_ip        VARCHAR(45),
    channel_id        BIGINT DEFAULT 0,
    created_at        TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_usage_logs_user ON usage_logs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_usage_logs_model ON usage_logs(model_name, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_usage_logs_created ON usage_logs(created_at DESC);
