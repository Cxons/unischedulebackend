CREATE TABLE faculties(
    faculty_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    faculty_name VARCHAR(30) NOT NULL,
    faculty_code VARCHAR(10) DEFAULT NULL,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);