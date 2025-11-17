-- Add UNIQUE constraint to lecturer_id in dean_waiting_list
ALTER TABLE dean_waiting_list
ADD CONSTRAINT dean_waiting_list_lecturer_id_key UNIQUE (lecturer_id);

-- Add UNIQUE constraint to lecturer_id in hod_waiting_list
ALTER TABLE hod_waiting_list
ADD CONSTRAINT hod_waiting_list_lecturer_id_key UNIQUE (lecturer_id);

-- Add UNIQUE constraint to lecturer_id in lecturer_waiting_list
ALTER TABLE lecturer_waiting_list
ADD CONSTRAINT lecturer_waiting_list_lecturer_id_key UNIQUE (lecturer_id);
