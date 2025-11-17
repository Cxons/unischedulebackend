CREATE TABLE cohorts(
    cohort_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cohort_name TEXT NOT NULL,
    cohort_level INT NOT NULL,
    cohort_department_id UUID NOT NULL REFERENCES departments(department_id) ON DELETE CASCADE,
    cohort_faculty_id UUID NOT NULL REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    cohort_university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()  
);