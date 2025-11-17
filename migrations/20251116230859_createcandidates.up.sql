CREATE TABLE candidates(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    fitness DOUBLE PRECISION NOT NULL,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    candidate_status TEXT NOT NULL CHECK(candidate_status IN ('CURRENT','DEPRECATED')),
    start_of_day TIMESTAMPTZ NOT NULL,
    end_of_day TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);