CREATE TABLE bill_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL UNIQUE REFERENCES transactions(id) ON DELETE CASCADE,
    bill_type VARCHAR(50) NOT NULL CHECK (bill_type IN ('airtime', 'data', 'electricity', 'betting')),
    provider_id UUID NOT NULL REFERENCES providers(id),
    provider_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(15),
    meter_number VARCHAR(50),
    meter_type VARCHAR(20) CHECK (meter_type IN ('prepaid', 'postpaid')),
    customer_name VARCHAR(255),
    customer_address VARCHAR(500),
    betting_account VARCHAR(100),
    betting_operator VARCHAR(50),
    data_plan_id VARCHAR(50),
    data_plan_name VARCHAR(100),
    data_volume VARCHAR(20),
    data_validity VARCHAR(50),
    electricity_token VARCHAR(100),
    electricity_units INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bill_details_txn ON bill_details(transaction_id);
CREATE INDEX idx_bill_details_type ON bill_details(bill_type);
CREATE INDEX idx_bill_details_provider ON bill_details(provider_id);
CREATE INDEX idx_bill_details_customer ON bill_details(phone_number, meter_number, betting_account);

COMMENT ON TABLE bill_details IS 'Bill payment specific details (airtime, data, electricity, betting)';