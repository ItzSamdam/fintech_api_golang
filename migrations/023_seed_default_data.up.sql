-- Seed default tier limits
INSERT INTO tier_limits (tier, daily_limit, weekly_limit, monthly_limit, single_tx_limit, can_send_money, can_buy_airtime, can_buy_data, can_pay_electricity, can_fund_betting, can_save) VALUES
(0, 0, 0, 0, 0, false, false, false, false, false, false),
(1, 5000000, 20000000, 50000000, 2500000, true, true, true, true, true, true),
(2, 20000000, 100000000, 300000000, 10000000, true, true, true, true, true, true),
(3, 500000000, 2000000000, 5000000000, 500000000, true, true, true, true, true, true);

-- Seed default fee configs
INSERT INTO fee_configs (bill_type, fee_type, fee_value, cap_amount, min_amount, vat_rate, is_active) VALUES
('transfer', 'percentage', 0.5, 500000, 10000, 7.5, true),
('airtime', 'percentage', 1.0, 100000, 1000, 7.5, true),
('data', 'percentage', 1.0, 100000, 1000, 7.5, true),
('electricity', 'fixed', 10000, 0, 10000, 7.5, true),
('betting', 'percentage', 1.5, 200000, 10000, 7.5, true);

-- Seed default providers
INSERT INTO providers (name, code, type, category, priority, base_url, timeout_seconds, retry_count, margin_percent, is_active) VALUES
('MTN Nigeria', 'MTN', 'airtime', 'mtn', 1, 'https://api.mtn.ng/v1', 30, 2, 0.5, true),
('Glo Nigeria', 'GLO', 'airtime', 'glo', 2, 'https://api.glo.com/v1', 30, 2, 0.5, true),
('Airtel Nigeria', 'AIRTEL', 'airtime', 'airtel', 3, 'https://api.airtel.ng/v1', 30, 2, 0.5, true),
('9mobile', '9MOBILE', 'airtime', '9mobile', 4, 'https://api.9mobile.com/v1', 30, 2, 0.5, true),
('Ikeja Electric', 'IE', 'electricity', 'ikeja', 1, 'https://api.ikejaelectric.com/v1', 30, 2, 1.0, true),
('Eko Electric', 'EKEDC', 'electricity', 'eko', 2, 'https://api.ekoelectric.com/v1', 30, 2, 1.0, true),
('Abuja DISCO', 'AEDC', 'electricity', 'abuja', 3, 'https://api.abujaelectric.com/v1', 30, 2, 1.0, true),
('Bet9ja', 'BET9JA', 'betting', 'bet9ja', 1, 'https://api.bet9ja.com/v1', 30, 2, 2.0, true),
('SportyBet', 'SPORTY', 'betting', 'sportybet', 2, 'https://api.sportybet.com/v1', 30, 2, 2.0, true);

-- Seed default admin roles
INSERT INTO roles (name, permissions, description) VALUES
('super_admin', '["*"]', 'Full system access'),
('admin', '["users.read", "users.update", "transactions.read", "transactions.reverse", "providers.toggle", "reports.view"]', 'Administrative access'),
('viewer', '["users.read", "transactions.read", "reports.view"]', 'Read-only access'),
('support', '["users.read", "transactions.read", "tickets.manage"]', 'Customer support access');

-- Seed default admin user (password: Admin@123 - CHANGE IN PRODUCTION!)
INSERT INTO admin_users (email, password_hash, full_name, role) VALUES 
('admin@fintech.com', '$2a$10$YourHashedPasswordHere', 'System Administrator', 'super_admin');

COMMENT ON TABLE providers IS 'Default providers seeded with initial configuration';