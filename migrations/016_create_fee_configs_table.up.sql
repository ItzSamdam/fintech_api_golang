CREATE TABLE fee_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bill_type VARCHAR(50) NOT NULL CHECK (bill_type IN ('transfer', 'airtime', 'data', 'electricity', 'betting')),
    fee_type VARCHAR(20) NOT NULL CHECK (fee_type IN ('percentage', 'fixed')),
    fee_value DECIMAL(10,2) NOT NULL CHECK (fee_value >= 0),
    cap_amount BIGINT NOT NULL DEFAULT 0 CHECK (cap_amount >= 0),
    min_amount BIGINT NOT NULL DEFAULT 0 CHECK (min_amount >= 0),
    max_amount BIGINT NOT NULL DEFAULT 0 CHECK (max_amount >= 0),
    vat_rate DECIMAL(5,2) NOT NULL DEFAULT 7.5 CHECK (vat_rate >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    effective_from TIMESTAMP NOT NULL DEFAULT NOW(),
    effective_to TIMESTAMP,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CHECK (effective_to IS NULL OR effective_to > effective_from)
);

CREATE INDEX idx_fee_configs_type ON fee_configs(bill_type);
CREATE INDEX idx_fee_configs_active ON fee_configs(is_active);
CREATE INDEX idx_fee_configs_effective ON fee_configs(effective_from, effective_to);

COMMENT ON COLUMN fee_configs.fee_value IS 'Fee value (if percentage: 1.5 = 1.5%)';
COMMENT ON COLUMN fee_configs.cap_amount IS 'Maximum fee amount in kobo (0 = no cap)';
COMMENT ON COLUMN fee_configs.min_amount IS 'Minimum transaction amount in kobo';
COMMENT ON COLUMN fee_configs.max_amount IS 'Maximum transaction amount in kobo';