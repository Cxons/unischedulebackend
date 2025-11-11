-- Down migration: revert column types to previous type
ALTER TABLE universities
    ALTER COLUMN university_name TYPE VARCHAR(255),
    ALTER COLUMN university_logo TYPE VARCHAR(255),
    ALTER COLUMN university_abbr TYPE VARCHAR(20),
    ALTER COLUMN email TYPE VARCHAR(100),
    ALTER COLUMN university_addr TYPE VARCHAR(255);
