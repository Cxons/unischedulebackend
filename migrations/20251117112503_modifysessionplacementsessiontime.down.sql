-- +migrate Down
ALTER TABLE session_placements 
ALTER COLUMN session_time TYPE TIME USING session_time::TIME;