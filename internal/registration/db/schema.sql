CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


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

CREATE TABLE university_admin(
    admin_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_first_name TEXT NOT NULL,
    admin_last_name TEXT NOT NULL,
    admin_middle_name TEXT DEFAULT NULL,
    admin_email TEXT NOT NULL UNIQUE,
    admin_password TEXT NOT NULL,
    admin_phone_number VARCHAR(15) DEFAULT NULL,
    admin_staff_card TEXT DEFAULT NULL,
    admin_number TEXT DEFAULT NULL,
    university_id UUID REFERENCES universities(university_id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- making sure that there is no duplicate admin_staff_number for one university
CREATE UNIQUE INDEX unique_admin_number_per_uni 
ON university_admin(university_id,admin_number);

-- Remember to create unique index to prevent 2 deans at the same time. This would be done in the queries.sql for this module
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

-- the unique index to prevent 2 deans at the same time in a particular faculty
CREATE UNIQUE INDEX unique_current_dean
ON current_dean(faculty_id)
WHERE end_date IS NULL;


-- same here also create an index to prevent 2 hods at the same time in the same department
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


CREATE UNIQUE INDEX unique_current_hod
ON current_hod(department_id) 
WHERE end_date IS NULL;



CREATE TABLE dean_waiting_list(
    wait_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lecturer_id UUID NOT NULL REFERENCES lecturers(lecturer_id) ON DELETE CASCADE,
    potential_faculty TEXT NOT NULL,
    additional_message TEXT DEFAULT NULL,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    approved boolean DEFAULT FALSE
);


CREATE TABLE hod_waiting_list(
    wait_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lecturer_id UUID NOT NULL REFERENCES lecturers(lecturer_id) ON DELETE CASCADE,
    potential_department TEXT NOT NULL,
    additional_message TEXT DEFAULT NULL,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    faculty_id UUID NOT NULL REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    approved boolean DEFAULT FALSE
);


CREATE TABLE lecturer_waiting_list(
    wait_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lecturer_id UUID NOT NULL REFERENCES lecturers(lecturer_id) ON DELETE CASCADE,
    additional_message TEXT DEFAULT NULL,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    faculty_id UUID NOT NULL REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    department_id UUID NOT NULL REFERENCES departments(department_id) ON DELETE CASCADE,
    approved boolean DEFAULT FALSE
);



-- for days and times lecturers would not be available
CREATE TABLE lecturer_unavailability (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lecturer_id UUID NOT NULL REFERENCES lecturers(lecturer_id) ON DELETE CASCADE,
    day TEXT NOT NULL CHECK (day IN ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday')),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    reason TEXT DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);


CREATE TABLE student_courses_offered(
    student_id UUID NOT NULL REFERENCES students(student_id) ON DELETE CASCADE,
    course_id UUID NOT NULL REFERENCES courses(course_id) ON DELETE CASCADE,
    PRIMARY KEY(student_id,course_id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);



