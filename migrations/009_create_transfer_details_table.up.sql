CREATE TABLE transfer_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL UNIQUE REFERENCES transactions(id) ON DELETE CASCADE,
    recipient_type VARCHAR(20) NOT NULL CHECK (recipient_type IN ('bank', 'wallet')),
    recipient_id VARCHAR(100) NOT NULL,
    recipient_name VARCHAR(255),
    recipient_bank_code VARCHAR(10),
    recipient_bank_name VARCHAR(100),
    nip_session_id VARCHAR(100),
    narration VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transfer_details_txn ON transfer_details(transaction_id);
CREATE INDEX idx_transfer_details_recipient ON transfer_details(recipient_type, recipient_id);

COMMENT ON TABLE transfer_details IS 'Bank/wallet transfer specific details';