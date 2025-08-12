CREATE UNIQUE INDEX unique_current_dean
ON current_dean(faculty_id)
WHERE end_date IS NULL;