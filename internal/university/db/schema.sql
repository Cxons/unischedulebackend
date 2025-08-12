CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE universities (
    university_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    university_name VARCHAR(50) NOT NULL UNIQUE,
    university_logo VARCHAR(50) DEFAULT NULL,
    university_abbr VARCHAR(20) DEFAULT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    website VARCHAR(30) DEFAULT NULL,
    phone_number VARCHAR(15) NOT NULL UNIQUE,
    university_addr VARCHAR(100) DEFAULT NULL,
    current_session TEXT DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE faculties(
    faculty_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    faculty_name VARCHAR(30) NOT NULL,
    faculty_code VARCHAR(10) DEFAULT NULL,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE departments(
    department_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    department_name VARCHAR(30) NOT NULL,
    department_code VARCHAR(10) DEFAULT NULL,
    faculty_id UUID NOT NULL REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()  
);

CREATE TABLE venues(
    venue_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    venue_name VARCHAR(50) NOT NULL,
    venue_longitude DOUBLE PRECISION DEFAULT NULL,
    venue_latitude DOUBLE PRECISION DEFAULT NULL,
    location TEXT DEFAULT NULL,
    venue_image TEXT DEFAULT NULL,
    capacity INT NOT NULL,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE faculty_venues(
    venue_id UUID REFERENCES venues(venue_id) ON DELETE CASCADE,
    faculty_id UUID REFERENCES faculties(faculty_id) ON DELETE CASCADE,
    PRIMARY KEY(venue_id,faculty_id),
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE dept_venues(
    venue_id UUID REFERENCES venues(venue_id) ON DELETE CASCADE,
    department_id UUID REFERENCES departments(department_id) ON DELETE CASCADE,
    PRIMARY KEY(venue_id,department_id),
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE venue_unavailability (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    venue_id UUID REFERENCES venues(venue_id),
    reason TEXT,
    university_id UUID NOT NULL REFERENCES universities(university_id) ON DELETE CASCADE,
    day TEXT CHECK (day IN ('Monday','Tuesday','Wednesday','Thursday','Friday','Saturday','Sunday')),
    start_time TIME,
    end_time TIME,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);


