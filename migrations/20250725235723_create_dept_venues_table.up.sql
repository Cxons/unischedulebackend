CREATE TABLE dept_venues(
    venue_id UUID REFERENCES venues(venue_id) ON DELETE CASCADE,
    department_id UUID REFERENCES departments(department_id) ON DELETE CASCADE,
    PRIMARY KEY(venue_id,department_id),
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);