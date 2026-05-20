CREATE TABLE savings_goals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    target_amount BIGINT NOT NULL CHECK (target_amount > 0),
    current_amount BIGINT NOT NULL DEFAULT 0 CHECK (current_amount >= 0),
    interest_rate DECIMAL(5,2) NOT NULL DEFAULT 0,
    duration_days INTEGER NOT NULL CHECK (duration_days > 0),
    start_date TIMESTAMP NOT NULL DEFAULT NOW(),
    target_date TIMESTAMP NOT NULL,
    is_auto_debit BOOLEAN NOT NULL DEFAULT false,
    auto_debit_amount BIGINT NOT NULL DEFAULT 0 CHECK (auto_debit_amount >= 0),
    auto_debit_day INTEGER NOT NULL DEFAULT 1 CHECK (auto_debit_day BETWEEN 1 AND 31),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    withdrawn_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CHECK (target_date > start_date),
    CHECK (current_amount <= target_amount)
);

CREATE INDEX idx_savings_goals_user ON savings_goals(user_id);
CREATE INDEX idx_savings_goals_status ON savings_goals(status);
CREATE INDEX idx_savings_goals_target ON savings_goals(target_date);
CREATE INDEX idx_savings_goals_auto ON savings_goals(is_auto_debit);

COMMENT ON COLUMN savings_goals.target_amount IS 'Target amount in kobo';
COMMENT ON COLUMN savings_goals.current_amount IS 'Current saved amount in kobo';
COMMENT ON TABLE savings_goals IS 'User savings goals and targets';