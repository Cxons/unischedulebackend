-- name: CountTotalLecturers :one
SELECT COUNT(*) FROM lecturers
WHERE lecturer_university_id = $1;

-- name: CountFacultyLecturers :one
SELECT COUNT(*) FROM lecturers 
WHERE lecturer_university_id = $1 AND lecturer_faculty_id = $2;

-- name: CountDepartmentLecturers :one
SELECT COUNT(*) FROM lecturers
WHERE lecturer_university_id = $1 AND lecturer_department_id = $2;


-- name: CountTotalCohorts :one
SELECT COUNT(*) FROM cohorts
WHERE cohort_university_id = $1;


-- name: CountTotalCourses :one
SELECT COUNT(*) FROM courses
WHERE university_id = $1;


-- name: RetrieveTotalLecturers :many
SELECT 
    lecturer_id,
    lecturer_first_name,
    lecturer_last_name,
    lecturer_middle_name,
    lecturer_email,
    lecturer_profile_pic
FROM lecturers
WHERE lecturer_university_id = $1;


-- name: RetrieveLecturersForFaculty :many
SELECT 
    lecturer_id,
    lecturer_first_name,
    lecturer_last_name,
    lecturer_middle_name,
    lecturer_email,
    lecturer_profile_pic
FROM lecturers
WHERE lecturer_university_id = $1 AND lecturer_faculty_id = $2;


-- name: RetrieveLecturersForDepartment :many
SELECT 
    lecturer_id,
    lecturer_first_name,
    lecturer_last_name,
    lecturer_middle_name,
    lecturer_email,
    lecturer_profile_pic
FROM lecturers
WHERE 
    lecturer_university_id = $1 
    AND 
    lecturer_faculty_id = $2 
    AND
    lecturer_department_id = $3;


-- name: RetrieveTotalLecturerUnavailability :many
SELECT 
    l.lecturer_id,
    lu.day,
    lu.start_time,
    lu.end_time,
    lu.reason
FROM lecturers l
INNER JOIN lecturer_unavailability lu
ON l.lecturer_id = lu.lecturer_id
WHERE l.lecturer_university_id = $1;


-- name: RetrieveTotalVenueUnavailability :many
SELECT 
    v.venue_id,
    vu.reason,
    vu.day,
    vu.start_time,
    vu.end_time
FROM venues v
INNER JOIN venue_unavailability vu
ON v.venue_id = vu.venue_id
WHERE v.university_id = $1;

-- name: RetrieveAllVenues :many
SELECT 
    venue_id,
    venue_name,
    venue_longitude,
    venue_latitude,
    location,
    venue_image,
    is_active
FROM venues
WHERE university_id = $1;

-- name: RetrieveAllCourses :many
SELECT
    course_id,
    course_code,
    course_title,
    course_credit_unit,
    course_duration,
    department_id,
    university_id,
    lecturer_id,
    sessions_per_week,
    semester
FROM courses
WHERE university_id = $1;


-- name: RetrieveCohortsForAllCourses :many
SELECT 
    cohort_id,
    course_id,
    university_id
FROM cohort_courses_offered
WHERE university_id = $1;


-- name: CreateCandidate :one
INSERT INTO candidates(
    fitness,university_id,candidate_status,start_of_day,end_of_day
)VALUES($1,$2,$3,$4,$5)
RETURNING *;


-- name: CreateSessionPlacements :one
INSERT INTO session_placements(
    candidate_id,session_idx,course_id,venue_id,day,session_time,university_id
)VALUES($1,$2,$3,$4,$5,$6,$7)
RETURNING *;

-- name: DeprecateLatestCandidate :exec
UPDATE candidates AS c
SET candidate_status = 'DEPRECATED',
    updated_at = NOW()
WHERE c.id = (
  SELECT id
  FROM candidates
  WHERE candidates.university_id = $1
  ORDER BY created_at DESC
  LIMIT 1
);

-- name: RestoreCurrentCandidate :exec
UPDATE candidates AS c
SET candidate_status = 'CURRENT',
    updated_at = NOW()
WHERE c.id = (
  SELECT id
  FROM candidates
  WHERE candidates.university_id = $1
  ORDER BY created_at DESC
  LIMIT 1
);


-- name: GetCohortSessionsInCurrentTimetable :many
SELECT 
    sp.id AS session_id,
    sp.session_idx,
    sp.course_id,
    sp.venue_id,
    sp.day,
    sp.session_time,
    sp.university_id,
    c.fitness,
    c.candidate_status,
    c.start_of_day,
    c.end_of_day
FROM session_placements sp
JOIN candidates c 
    ON sp.candidate_id = c.id
JOIN cohort_courses_offered cco
    ON sp.course_id = cco.course_id
WHERE 
    cco.cohort_id = $1
    AND cco.university_id = $2
    AND c.university_id = $2
    AND c.candidate_status = 'CURRENT';


-- -- name: UpdateOtherCandidateStatus :one
-- UPDATE 






