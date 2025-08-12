CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE refresh_tokens(
    token_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    user_id UUID NOT NULL, 
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    is_revoked BOOLEAN DEFAULT FALSE
);

CREATE TABLE otps(
    otp_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    otp TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    email TEXT NOT NULL UNIQUE,
    user_type TEXT NOT NULL CHECK(user_type IN ('admin','lecturer','student')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
)

-- CREATE TABLE refresh_tokens(
--     token_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     user_id UUID NOT NULL,
--     user_type TEXT NOT NULL CHECK( user_type IN ('admin','lecturer','student')),
--     refresh_token TEXT NOT NULL UNIQUE,
--     ip_address INET DEFAULT NULL,
--     user_agent TEXT DEFAULT NULL,
--     is_revoked BOOLEAN DEFAULT FALSE,
--     created_at TIMESTAMPTZ DEFAULT NOW(),
--     expires_at TIMESTAMPTZ NOT NULL,
--     updated_at TIMESTAMPTZ DEFAULT NOW()
-- );