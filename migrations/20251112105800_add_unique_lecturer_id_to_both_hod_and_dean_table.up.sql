
-- Ensure a lecturer can only be dean of one faculty and vice versa
ALTER TABLE current_dean
ADD CONSTRAINT unique_lecturer_faculty_combination
UNIQUE (lecturer_id, faculty_id);

-- Ensure a lecturer can only head one department and vice versa
ALTER TABLE current_hod
ADD CONSTRAINT unique_lecturer_department_combination
UNIQUE (lecturer_id, department_id);