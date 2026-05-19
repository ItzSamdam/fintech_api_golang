CREATE TABLE otps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone_number VARCHAR(15) NOT NULL,
    code VARCHAR(6) NOT NULL,
    purpose VARCHAR(50) NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT false,
    attempts INTEGER NOT NULL DEFAULT 0,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_otps_phone ON otps(phone_number);
CREATE INDEX idx_otps_code ON otps(code);
CREATE INDEX idx_otps_expires ON otps(expires_at);
CREATE INDEX idx_otps_purpose ON otps(purpose, is_used);

COMMENT ON TABLE otps IS 'One-time passwords for verification';