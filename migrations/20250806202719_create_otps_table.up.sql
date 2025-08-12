CREATE TABLE otps(
    otp_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    otp TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    email TEXT NOT NULL UNIQUE,
    user_type TEXT NOT NULL CHECK(user_type IN ('admin','lecturer','student')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
)