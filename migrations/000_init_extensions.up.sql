-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable pgcrypto for encryption (for BVN/NIN)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create ENUM types for better data integrity
CREATE TYPE user_tier AS ENUM ('0', '1', '2', '3');
CREATE TYPE transaction_status AS ENUM ('pending', 'processing', 'success', 'failed', 'reversed');
CREATE TYPE transaction_type AS ENUM ('credit', 'debit');
CREATE TYPE transaction_category AS ENUM ('transfer', 'airtime', 'data', 'electricity', 'betting', 'savings', 'fee', 'refund');
CREATE TYPE bill_type AS ENUM ('airtime', 'data', 'electricity', 'betting');
CREATE TYPE kyc_status AS ENUM ('pending', 'approved', 'rejected', 'verification_needed');
CREATE TYPE ticket_status AS ENUM ('open', 'in_progress', 'resolved', 'closed');
CREATE TYPE ticket_priority AS ENUM ('low', 'medium', 'high', 'urgent');
CREATE TYPE provider_status AS ENUM ('healthy', 'degraded', 'down', 'unknown');
CREATE TYPE savings_goal_status AS ENUM ('active', 'completed', 'withdrawn', 'cancelled');