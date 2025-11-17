-- migrations/002_alter_session_time_to_timestamptz.sql

-- +migrate Up
ALTER TABLE session_placements 
ALTER COLUMN session_time TYPE TIMESTAMPTZ 
USING (CURRENT_DATE + session_time) AT TIME ZONE 'UTC';

-- -- +migrate Down  
-- ALTER TABLE session_placements 
-- ALTER COLUMN session_time TYPE TIME 
-- USING (session_time::time);