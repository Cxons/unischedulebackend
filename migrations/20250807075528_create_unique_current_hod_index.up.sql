CREATE UNIQUE INDEX unique_current_hod
ON current_hod(department_id) 
WHERE end_date IS NULL;