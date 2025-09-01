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

