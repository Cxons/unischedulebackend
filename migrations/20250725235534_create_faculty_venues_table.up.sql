CREATE TABLE faculty_venues(
    venue_id UUID REFERENCES venues(venue_id) ON DELETE CASCADE,
    faculty_id UUID REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    PRIMARY KEY(venue_id,faculty_id),
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);