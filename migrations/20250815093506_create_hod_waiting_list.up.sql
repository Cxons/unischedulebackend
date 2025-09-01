CREATE TABLE hod_waiting_list(
    wait_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lecturer_id UUID NOT NULL REFERENCES lecturers(lecturer_id) ON DELETE CASCADE,
    potential_department TEXT NOT NULL,
    additional_message TEXT DEFAULT NULL,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    faculty_id UUID NOT NULL REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    approved boolean DEFAULT FALSE
);
