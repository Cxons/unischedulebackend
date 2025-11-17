CREATE TABLE cohort_courses_offered (
    cohort_id UUID REFERENCES cohorts(cohort_id) ON DELETE CASCADE,
    course_id UUID REFERENCES courses(course_id) ON DELETE CASCADE,
    PRIMARY KEY(cohort_id,course_id),
    university_id UUID NOT NULL REFERENCES universities(university_id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);