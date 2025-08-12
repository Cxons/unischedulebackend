CREATE TABLE current_dean(
    dean_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lecturer_id UUID REFERENCES lecturers(lecturer_id) ON DELETE CASCADE,
    faculty_id UUID REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    university_id UUID REFERENCES universities(university_id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);