CREATE TABLE IF NOT EXISTS provider_models (
    id                 BIGSERIAL PRIMARY KEY,
    provider_id        BIGINT NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    model_name         VARCHAR(100) NOT NULL,
    model_id           VARCHAR(100) NOT NULL,
    model_type         VARCHAR(20)  DEFAULT 'chat',
    enabled            BOOLEAN DEFAULT true,
    weight             INTEGER DEFAULT 1,
    priority           INTEGER DEFAULT 0,
    max_context_tokens INTEGER,
    input_price        DECIMAL(10,6),
    output_price       DECIMAL(10,6),
    config             JSONB DEFAULT '{}',
    created_at         TIMESTAMPTZ DEFAULT NOW(),
    updated_at         TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(provider_id, model_name)
);

CREATE INDEX IF NOT EXISTS idx_provider_models_name ON provider_models(model_name);
CREATE INDEX IF NOT EXISTS idx_provider_models_enabled ON provider_models(enabled);
CREATE INDEX IF NOT EXISTS idx_provider_models_priority ON provider_models(model_name, priority, enabled);
