-- Up migration: modify column types to TEXT
ALTER TABLE universities
    ALTER COLUMN university_name TYPE TEXT,
    ALTER COLUMN university_logo TYPE TEXT,
    ALTER COLUMN university_abbr TYPE TEXT,
    ALTER COLUMN email TYPE TEXT,
    ALTER COLUMN university_addr TYPE TEXT;
