CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE courses(
    course_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_code TEXT NOT NULL,
    course_title TEXT NOT NULL,
    course_credit_unit INT NOT NULL,
    course_duration INT NOT NULL,
    department_id UUID NOT NULL REFERENCES departments(department_id) ON DELETE CASCADE,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    UNIQUE(course_code,university_id),
    lecturer_id UUID REFERENCES lecturers(lecturer_id) ON DELETE SET NULL,
    sessions_per_week INT NOT NULL,
    level INT NOT NULL,
    semester TEXT NOT NULL CHECK (semester IN ('First','Second')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE courses_possible_venues(
    course_id UUID REFERENCES courses(course_id) ON DELETE CASCADE,
    venue_id UUID REFERENCES venues(venue_id) ON DELETE CASCADE,
    university_id UUID REFERENCES universities(university_id) ON DELETE CASCADE,
    PRIMARY KEY(course_id,venue_id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);  


CREATE TABLE courses_lecturers(
    course_id UUID REFERENCES courses(course_id) ON DELETE CASCADE,
    lecturer_id UUID REFERENCES lecturers(lecturer_id) ON DELETE CASCADE,
    PRIMARY KEY(course_id,lecturer_id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);