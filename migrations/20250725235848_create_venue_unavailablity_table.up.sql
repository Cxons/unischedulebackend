CREATE TABLE venue_unavailability (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    venue_id UUID REFERENCES venues(venue_id),
    reason TEXT,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    day TEXT CHECK (day IN ('Monday','Tuesday','Wednesday','Thursday','Friday','Saturday','Sunday')),
    start_time TIME,
    end_time TIME,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);