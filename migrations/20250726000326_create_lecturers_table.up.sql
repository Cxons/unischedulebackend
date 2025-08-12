CREATE TABLE lecturers(
    lecturer_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lecturer_first_name TEXT NOT NULL,
    lecturer_last_name TEXT NOT NULL,
    lecturer_middle_name TEXT DEFAULT NULL,
    lecturer_email TEXT NOT NULL UNIQUE,
    lecturer_password TEXT NOT NULL,
    lecturer_phone_number VARCHAR(15) DEFAULT NULL,
    lecturer_profile_pic TEXT DEFAULT NULL,
    lecturer_staff_id TEXT UNIQUE DEFAULT NULL,
    lecturer_university_id UUID REFERENCES universities(university_id) ON DELETE CASCADE,
    lecturer_faculty_id UUID REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    lecturer_department_id UUID REFERENCES departments(department_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);