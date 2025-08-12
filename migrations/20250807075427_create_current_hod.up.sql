CREATE TABLE current_hod(
    hod_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lecturer_id UUID REFERENCES lecturers(lecturer_id) ON DELETE CASCADE,
    department_id UUID REFERENCES departments(department_id) ON DELETE CASCADE,
    university_id UUID REFERENCES universities(university_id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);