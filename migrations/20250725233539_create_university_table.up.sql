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