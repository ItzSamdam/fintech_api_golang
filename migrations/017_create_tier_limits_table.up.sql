CREATE TABLE tier_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tier INTEGER NOT NULL UNIQUE CHECK (tier BETWEEN 0 AND 3),
    daily_limit BIGINT NOT NULL CHECK (daily_limit >= 0),
    weekly_limit BIGINT NOT NULL CHECK (weekly_limit >= 0),
    monthly_limit BIGINT NOT NULL CHECK (monthly_limit >= 0),
    single_tx_limit BIGINT NOT NULL CHECK (single_tx_limit >= 0),
    can_send_money BOOLEAN NOT NULL DEFAULT true,
    can_buy_airtime BOOLEAN NOT NULL DEFAULT true,
    can_buy_data BOOLEAN NOT NULL DEFAULT true,
    can_pay_electricity BOOLEAN NOT NULL DEFAULT true,
    can_fund_betting BOOLEAN NOT NULL DEFAULT true,
    can_save BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tier_limits_tier ON tier_limits(tier);

COMMENT ON COLUMN tier_limits.daily_limit IS 'Daily transaction limit in kobo';
COMMENT ON COLUMN tier_limits.weekly_limit IS 'Weekly transaction limit in kobo';
COMMENT ON COLUMN tier_limits.monthly_limit IS 'Monthly transaction limit in kobo';
COMMENT ON COLUMN tier_limits.single_tx_limit IS 'Single transaction limit in kobo';