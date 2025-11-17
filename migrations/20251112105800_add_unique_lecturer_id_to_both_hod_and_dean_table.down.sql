-- Remove the unique constraints if rollback is needed
ALTER TABLE current_dean
DROP CONSTRAINT IF EXISTS unique_lecturer_faculty_combination;

ALTER TABLE current_hod
DROP CONSTRAINT IF EXISTS unique_lecturer_department_combination;