CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference VARCHAR(100) NOT NULL UNIQUE,
    wallet_id UUID NOT NULL REFERENCES wallets(id),
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL CHECK (type IN ('credit', 'debit')),
    category VARCHAR(50) NOT NULL CHECK (category IN ('transfer', 'airtime', 'data', 'electricity', 'betting', 'savings', 'fee', 'refund')),
    sub_category VARCHAR(50),
    amount BIGINT NOT NULL CHECK (amount > 0),
    fee BIGINT NOT NULL DEFAULT 0 CHECK (fee >= 0),
    vat BIGINT NOT NULL DEFAULT 0 CHECK (vat >= 0),
    total_amount BIGINT NOT NULL CHECK (total_amount > 0),
    balance_before BIGINT NOT NULL CHECK (balance_before >= 0),
    balance_after BIGINT NOT NULL CHECK (balance_after >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    description VARCHAR(255),
    metadata JSONB,
    provider_reference VARCHAR(100),
    provider_response TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    is_reversed BOOLEAN NOT NULL DEFAULT false,
    reversed_txn_id UUID,
    ip_address VARCHAR(45),
    device_id VARCHAR(255),
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (reversed_txn_id) REFERENCES transactions(id)
);

CREATE INDEX idx_transactions_reference ON transactions(reference);
CREATE INDEX idx_transactions_wallet ON transactions(wallet_id);
CREATE INDEX idx_transactions_user ON transactions(user_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_category ON transactions(category);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_created ON transactions(created_at);
CREATE INDEX idx_transactions_completed ON transactions(completed_at);
CREATE INDEX idx_transactions_metadata ON transactions USING gin(metadata);

COMMENT ON COLUMN transactions.amount IS 'Transaction amount in kobo';
COMMENT ON COLUMN transactions.fee IS 'Fee amount in kobo';
COMMENT ON COLUMN transactions.vat IS 'VAT amount in kobo (7.5% of fee)';
COMMENT ON COLUMN transactions.total_amount IS 'Total amount debited/credited in kobo (amount + fee + vat)';
COMMENT ON COLUMN transactions.metadata IS 'JSON metadata for flexible data storage';
COMMENT ON TABLE transactions IS 'All financial transactions with balances in kobo';