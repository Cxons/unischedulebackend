CREATE TABLE venues(
    venue_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    venue_name VARCHAR(50) NOT NULL,
    venue_longitude DOUBLE PRECISION DEFAULT NULL,
    venue_latitude DOUBLE PRECISION DEFAULT NULL,
    location TEXT DEFAULT NULL,
    venue_image TEXT DEFAULT NULL,
    capacity INT NOT NULL,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);