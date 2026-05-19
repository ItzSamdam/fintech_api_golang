CREATE TABLE wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    ledger_balance BIGINT NOT NULL DEFAULT 0 CHECK (ledger_balance >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'NGN',
    is_locked BOOLEAN NOT NULL DEFAULT false,
    locked_at TIMESTAMP,
    lock_reason VARCHAR(255),
    daily_spent BIGINT NOT NULL DEFAULT 0,
    weekly_spent BIGINT NOT NULL DEFAULT 0,
    monthly_spent BIGINT NOT NULL DEFAULT 0,
    last_daily_reset TIMESTAMP NOT NULL DEFAULT NOW(),
    last_weekly_reset TIMESTAMP NOT NULL DEFAULT NOW(),
    last_monthly_reset TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_wallets_user ON wallets(user_id);
CREATE INDEX idx_wallets_locked ON wallets(is_locked);

COMMENT ON COLUMN wallets.balance IS 'Current balance in kobo (1 NGN = 100 kobo)';
COMMENT ON COLUMN wallets.ledger_balance IS 'Ledger balance for reconciliation in kobo';
COMMENT ON COLUMN wallets.daily_spent IS 'Daily spent amount in kobo for tier limits';
COMMENT ON TABLE wallets IS 'User wallets with balance in kobo';