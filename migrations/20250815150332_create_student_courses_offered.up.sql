CREATE TABLE student_courses_offered(
    student_id UUID NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(course_id) ON DELETE CASCADE,
    PRIMARY KEY(student_id,course_id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);