CREATE TABLE two_fa (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    secret VARCHAR(255) NOT NULL,
    backup_codes TEXT,
    is_enabled BOOLEAN NOT NULL DEFAULT false,
    verified_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_2fa_user ON two_fa(user_id);
CREATE INDEX idx_2fa_enabled ON two_fa(is_enabled);

COMMENT ON TABLE two_fa IS 'Two-factor authentication settings';