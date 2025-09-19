-- name: CreateUniversity :one
INSERT INTO universities(
    university_name,university_logo,university_abbr,email,website,phone_number,university_addr,current_session
)VALUES(
    $1,$2,$3,$4,$5,$6,$7,$8
) RETURNING *;

-- name: UpdateUniversity :one
UPDATE universities
SET 
    university_name = $1, 
    university_logo = $2, 
    university_abbr = $3, 
    email = $4, 
    website = $5, 
    phone_number = $6, 
    university_addr = $7, 
    current_session = $8
WHERE university_id = $9
RETURNING *;

-- name: RetrieveUniversitiesWithLimit :many
SELECT * FROM universities LIMIT $1;

-- name: RetrieveAllUniversities :many
SELECT * FROM universities;

-- name: CreateFaculty :one
INSERT INTO faculties(
    faculty_name,faculty_code,university_id
)VALUES(
    $1,$2,$3
)
RETURNING *;

-- name: UpdateFaculty :one
UPDATE faculties
SET 
    faculty_name = $1,
    faculty_code = $2
WHERE faculty_id = $3
RETURNING *;

-- name: RetrieveFacultiesForAUni :many
SELECT * FROM faculties
WHERE university_id = $1;

-- name: CreateDepartment :one
INSERT INTO departments(
    department_name,department_code,faculty_id,university_id,number_of_levels
)VALUES(
    $1,$2,$3,$4,$5
)
RETURNING *;

-- name: UpdateDepartment :one
UPDATE departments
SET 
    department_name = $1,
    department_code = $2,
    number_of_levels = $3
WHERE department_id = $3
RETURNING *;

-- name: CreateCohort :one
INSERT INTO cohorts(cohort_name,
    cohort_level,
    cohort_department_id,
    cohort_faculty_id,
    cohort_university_id
)VALUES(
    $1,$2,$3,$4,$5
)
RETURNING *;

-- name: UpdateCohort :one
UPDATE cohorts
SET
    cohort_level = $1,
    cohort_department_id = $2,
    cohort_faculty_id = $3,
    cohort_university_id = $4
WHERE cohort_id = $5
RETURNING *;

-- name: RetrieveAllCohorts :many
SELECT * FROM cohorts
WHERE cohort_university_id = $1;


-- name: RetrieveDeptsForAFaculty :many
SELECT * FROM departments
WHERE faculty_id = $1 AND university_id = $2;


-- name: CreateVenue :one
INSERT INTO venues(
    venue_name,
    venue_longitude,
    venue_latitude,
    location,
    venue_image,
    capacity,
    university_id
)VALUES(
    $1,$2,$3,$4,$5,$6,$7
)
RETURNING *;

-- name: SetFacultyVenue :exec
INSERT INTO faculty_venues(
    venue_id,
    faculty_id,
    university_id
)VALUES(
    $1,$2,$3
);

-- name: SetDepartmentVenue :exec
INSERT INTO dept_venues(
    venue_id,
    department_id,
    university_id
)
VALUES(
    $1,$2,$3
);

