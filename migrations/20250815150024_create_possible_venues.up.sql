CREATE TABLE courses_possible_venues(
    course_id UUID REFERENCES courses(course_id) ON DELETE CASCADE,
    venue_id UUID REFERENCES venues(venue_id) ON DELETE CASCADE,
    university_id UUID REFERENCES universities(university_id) ON DELETE CASCADE,
    PRIMARY KEY(course_id,venue_id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);  