CREATE TABLE auto_roundups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    savings_goal_id UUID NOT NULL REFERENCES savings_goals(id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    multiplier INTEGER NOT NULL DEFAULT 1 CHECK (multiplier >= 1),
    max_daily_amount BIGINT NOT NULL DEFAULT 100000 CHECK (max_daily_amount > 0),
    total_roundup BIGINT NOT NULL DEFAULT 0 CHECK (total_roundup >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auto_roundups_user ON auto_roundups(user_id);
CREATE INDEX idx_auto_roundups_active ON auto_roundups(is_active);

COMMENT ON COLUMN auto_roundups.multiplier IS 'Round up to nearest X naira';
COMMENT ON COLUMN auto_roundups.max_daily_amount IS 'Maximum daily roundup amount in kobo';
COMMENT ON COLUMN auto_roundups.total_roundup IS 'Total roundup amount saved in kobo';