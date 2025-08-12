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