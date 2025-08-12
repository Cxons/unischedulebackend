CREATE TABLE departments(
    department_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_name VARCHAR(30) NOT NULL,
    department_code VARCHAR(10) DEFAULT NULL,
    faculty_id UUID NOT NULL REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()  
);