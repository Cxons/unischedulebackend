-- Migration: Convert session_time from TIME to TIMESTAMPTZ
ALTER TABLE session_placements 
ALTER COLUMN session_time TYPE TIMESTAMPTZ 
USING (CURRENT_DATE + session_time) AT TIME ZONE 'UTC';