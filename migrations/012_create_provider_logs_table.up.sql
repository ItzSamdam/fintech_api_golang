CREATE TABLE provider_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
    transaction_id UUID REFERENCES transactions(id) ON DELETE SET NULL,
    endpoint VARCHAR(255) NOT NULL,
    request TEXT,
    response TEXT,
    status_code INTEGER,
    response_time INTEGER,
    is_error BOOLEAN NOT NULL DEFAULT false,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_provider_logs_provider ON provider_logs(provider_id);
CREATE INDEX idx_provider_logs_txn ON provider_logs(transaction_id);
CREATE INDEX idx_provider_logs_created ON provider_logs(created_at);
CREATE INDEX idx_provider_logs_error ON provider_logs(is_error);

COMMENT ON TABLE provider_logs IS 'API call logs to third-party providers';