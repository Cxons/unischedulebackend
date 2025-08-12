CREATE TABLE students(
    student_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_first_name TEXT NOT NULL,
    student_last_name TEXT NOT NULL,
    student_middle_name TEXT DEFAULT NULL,
    student_email TEXT UNIQUE NOT NULL,
    student_password TEXT NOT NULL,
    student_phone_number TEXT DEFAULT NULL,
    student_profile_pic TEXT DEFAULT NULL,
    student_reg_no TEXT DEFAULT NULL,
    student_level INT DEFAULT NULL, 
    student_university_id UUID REFERENCES universities(university_id) ON DELETE CASCADE,
    student_faculty_id UUID REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    student_department_id UUID REFERENCES departments(department_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);