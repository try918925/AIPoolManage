CREATE TABLE IF NOT EXISTS user_api_keys (
    id             BIGSERIAL PRIMARY KEY,
    user_id        VARCHAR(100) NOT NULL,
    name           VARCHAR(100),
    key_hash       VARCHAR(64)  NOT NULL UNIQUE,
    key_prefix     VARCHAR(12)  NOT NULL,
    enabled        BOOLEAN DEFAULT true,
    rate_limit     INTEGER DEFAULT 60,
    quota_limit    BIGINT  DEFAULT 0,
    quota_used     BIGINT  DEFAULT 0,
    allowed_models TEXT[],
    expires_at     TIMESTAMPTZ,
    last_used_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ DEFAULT NOW(),
    updated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_api_keys_hash ON user_api_keys(key_hash);
CREATE INDEX IF NOT EXISTS idx_user_api_keys_user ON user_api_keys(user_id);
