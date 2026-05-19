CREATE TABLE kyc_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    bvn_verified BOOLEAN NOT NULL DEFAULT false,
    bvn_verified_at TIMESTAMP,
    nin_verified BOOLEAN NOT NULL DEFAULT false,
    nin_verified_at TIMESTAMP,
    face_verified BOOLEAN NOT NULL DEFAULT false,
    face_verified_at TIMESTAMP,
    liveness_score DECIMAL(5,2) CHECK (liveness_score >= 0 AND liveness_score <= 100),
    id_card_front VARCHAR(500),
    id_card_back VARCHAR(500),
    passport_photo VARCHAR(500),
    utility_bill VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    rejection_reason VARCHAR(255),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_kyc_user ON kyc_records(user_id);
CREATE INDEX idx_kyc_status ON kyc_records(status);
CREATE INDEX idx_kyc_verified ON kyc_records(bvn_verified, nin_verified, face_verified);

COMMENT ON TABLE kyc_records IS 'KYC verification records for user tier upgrades';