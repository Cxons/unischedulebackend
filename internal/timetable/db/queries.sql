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






