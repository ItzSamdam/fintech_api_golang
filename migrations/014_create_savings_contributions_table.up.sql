CREATE TABLE savings_contributions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    savings_goal_id UUID NOT NULL REFERENCES savings_goals(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL UNIQUE REFERENCES transactions(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL CHECK (amount > 0),
    interest_earned BIGINT NOT NULL DEFAULT 0 CHECK (interest_earned >= 0),
    contribution_date TIMESTAMP NOT NULL DEFAULT NOW(),
    is_auto_debit BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_savings_contributions_goal ON savings_contributions(savings_goal_id);
CREATE INDEX idx_savings_contributions_txn ON savings_contributions(transaction_id);
CREATE INDEX idx_savings_contributions_date ON savings_contributions(contribution_date);

COMMENT ON COLUMN savings_contributions.amount IS 'Contribution amount in kobo';
COMMENT ON COLUMN savings_contributions.interest_earned IS 'Interest earned on contribution in kobo';