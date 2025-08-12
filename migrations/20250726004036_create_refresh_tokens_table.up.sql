CREATE TABLE refresh_tokens(
    token_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    refresh_token TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    user_id UUID NOT NULL, 
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    is_revoked BOOLEAN DEFAULT FALSE
);