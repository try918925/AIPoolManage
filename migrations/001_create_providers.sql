CREATE TABLE IF NOT EXISTS providers (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL UNIQUE,
    type          VARCHAR(50)  NOT NULL,
    base_url      VARCHAR(500) NOT NULL,
    api_key_encrypted TEXT NOT NULL,
    org_id        VARCHAR(200),
    enabled       BOOLEAN DEFAULT true,
    config        JSONB DEFAULT '{}',
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_providers_type ON providers(type);
CREATE INDEX IF NOT EXISTS idx_providers_enabled ON providers(enabled);
