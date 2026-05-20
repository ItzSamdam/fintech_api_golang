CREATE TABLE providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('airtime', 'data', 'electricity', 'betting', 'bank')),
    category VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT true,
    priority INTEGER NOT NULL DEFAULT 100,
    base_url VARCHAR(255),
    api_key VARCHAR(255),
    api_secret VARCHAR(255),
    timeout_seconds INTEGER NOT NULL DEFAULT 30,
    retry_count INTEGER NOT NULL DEFAULT 2,
    margin_percent DECIMAL(5,2) NOT NULL DEFAULT 0,
    metadata JSONB,
    last_health_check TIMESTAMP,
    health_status VARCHAR(20) NOT NULL DEFAULT 'unknown',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_providers_code ON providers(code);
CREATE INDEX idx_providers_type ON providers(type);
CREATE INDEX idx_providers_active ON providers(is_active);
CREATE INDEX idx_providers_health ON providers(health_status);
CREATE INDEX idx_providers_priority ON providers(priority);

COMMENT ON TABLE providers IS 'Third-party service providers configuration';