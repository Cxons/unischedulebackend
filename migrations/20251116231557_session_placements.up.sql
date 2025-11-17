CREATE TABLE session_placements(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    candidate_id UUID NOT NULL REFERENCES candidates(id) ON DELETE CASCADE,
    session_idx INT NOT NULL,
    course_id UUID NOT NULL REFERENCES courses(course_id) ON DELETE CASCADE,
    venue_id UUID NOT NULL REFERENCES venues(venue_id) ON DELETE CASCADE,
    day TEXT NOT NULL CHECK (day IN('Monday','Tuesday','Wednesday','Thursday','Friday')),
    session_time TIME NOT NULL,
    university_id UUID NOT NULL REFERENCES  universities(university_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);