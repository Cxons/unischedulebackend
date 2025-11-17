-- Ensure a lecturer can only be a dean once
ALTER TABLE current_dean
ADD CONSTRAINT unique_lecturer_in_dean
UNIQUE (lecturer_id);

-- Ensure a lecturer can only be a HOD once
ALTER TABLE current_hod
ADD CONSTRAINT unique_lecturer_in_hod
UNIQUE (lecturer_id);