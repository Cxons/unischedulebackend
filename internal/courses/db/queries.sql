-- name: CreateCourse :one
INSERT INTO courses(
    course_code,
    course_title,
    course_credit_unit,
    department_id,
    university_id,
    lecturer_id,
    sessions_per_week,
    level,
    semester,
    course_duration
)
VALUES(
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
)
RETURNING *;

-- name: RetrieveCoursesForADepartment :many
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
    level,
    semester
FROM courses
WHERE department_id = $1 AND university_id = $2;


-- name: UpdateCourse :one
UPDATE courses
SET
    course_code = $1,
    course_title = $2,
    course_credit_unit = $3,
    course_duration = $4,
    sessions_per_week = $5,
    lecturer_id =$6,
    level = $7,
    semester = $8
WHERE course_id = $9
RETURNING *;


-- name: SetCoursePossibleVenue :exec
INSERT INTO courses_possible_venues(
    course_id,
    venue_id,
    university_id
)VALUES($1,$2,$3);


-- name: DeleteCoursePossibleVenue :exec
DELETE FROM courses_possible_venues
WHERE course_id = $1
AND venue_id = $2;

-- name: FetchCoursePossibleVenues :many
SELECT 
    cpv.venue_id,
    v.venue_name,
    v.capacity
FROM courses_possible_venues cpv
INNER JOIN 
    venues v
ON 
    cpv.venue_id = v.venue_id
WHERE
    course_id = $1;




-- name: RetrieveAllCoursesAndTheirVenueIds :many
SELECT
    c.course_id,
    c.course_code,
    c.course_title,
    c.course_credit_unit,
    c.course_duration,
    c.department_id,
    c.university_id,
    c.lecturer_id,
    c.sessions_per_week,
    c.level,
    c.semester,
    cpv.venue_id
FROM courses c
INNER JOIN 
    courses_possible_venues cpv 
ON 
    cpv.course_id = c.course_id
WHERE c.university_id = $1;



-- name: DeleteCourse :exec
DELETE FROM courses
WHERE course_id = $1;


-- name: SetStudentCourse :one
INSERT INTO student_courses_offered(
    student_id,
    course_id
)VALUES(
    $1,$2
)
RETURNING *;


-- name: FetchStudentCourses :many
SELECT
    c.course_id,
    c.course_code,
    c.course_title,
    c.course_credit_unit,
    c.course_duration,
    c.department_id,
    c.sessions_per_week,
    c.lecturer_id,
    c.semester,
    c.level,
    c.department_id
FROM student_courses_offered sco
JOIN courses c
ON c.course_id = sco.course_id
WHERE student_id = $1;


-- name: RemoveStudentCourse :one
DELETE FROM student_courses_offered
WHERE student_id = $1
AND course_id = $2
RETURNING *;


-- name: SetCourseLecturers :one
INSERT INTO courses_lecturers(
    course_id,lecturer_id
)VALUES(
    $1,$2
)
RETURNING *;


-- name: UpdateCourseLecturers :one
UPDATE courses_lecturers
SET 
    lecturer_id = $1
WHERE course_id = $2 AND lecturer_id = $3
RETURNING *;

-- name: CreateCohortCourse :one
INSERT INTO cohort_courses_offered(
    cohort_id,course_id,university_id
)VALUES(
    $1,$2,$3
)
RETURNING *;


-- -- name: FetchCoursesForACohort :many
-- SELECT
--     cohort_id
-- FROM cohort_courses_offered
-- WHERE cohort_id = $1
-- AND university_id = $2;

-- name: FetchAllCourses :many
SELECT
    course_id,
    course_code,
    course_title,
    course_credit_unit,
    course_duration,
    sessions_per_week,
    level,
    semester
FROM courses
WHERE university_id = $1;
