CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone_number VARCHAR(15) NOT NULL UNIQUE,
    email VARCHAR(255) UNIQUE,
    bvn VARCHAR(11),
    nin VARCHAR(11),
    tier INTEGER NOT NULL DEFAULT 0 CHECK (tier BETWEEN 0 AND 3),
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_suspended BOOLEAN NOT NULL DEFAULT false,
    suspended_at TIMESTAMP,
    suspension_reason VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    face_photo_url VARCHAR(500),
    face_embedding TEXT,
    device_id VARCHAR(255),
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_phone ON users(phone_number);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_bvn ON users(bvn);
CREATE INDEX idx_users_nin ON users(nin);
CREATE INDEX idx_users_tier ON users(tier);
CREATE INDEX idx_users_status ON users(is_active, is_suspended);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

COMMENT ON COLUMN users.bvn IS 'Encrypted using pgcrypto';
COMMENT ON COLUMN users.nin IS 'Encrypted using pgcrypto';
COMMENT ON TABLE users IS 'User accounts with tier-based access control';