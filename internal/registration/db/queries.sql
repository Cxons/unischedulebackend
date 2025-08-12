-- name: UpdateAdminInfo :one
UPDATE university_admin
SET admin_middle_name = $1, admin_phone_number = $2, admin_staff_card = $3, admin_number = $4, university_id = $5
WHERE admin_id = $6
RETURNING *;

-- name: UpdateStudentInfo :one
UPDATE students
SET student_middle_name = $1, student_phone_number = $2, student_profile_pic = $3, student_reg_no = $4, student_level = $5, student_university_id = $6, student_faculty_id = $7, student_department_id = $8
WHERE student_id = $9
RETURNING *;

-- name: UpdateLecturerInfo :one
UPDATE lecturers
SET lecturer_middle_name = $1, lecturer_phone_number = $2, lecturer_profile_pic = $3, lecturer_staff_id = $4, lecturer_university_id = $5, lecturer_faculty_id = $6, lecturer_department_id = $7
WHERE lecturer_id = $8
RETURNING *;

-- name: CreateDean :one
INSERT INTO current_dean(
    lecturer_id,faculty_id,university_id,start_date,end_date
)VALUES(
    $1,$2,$3,$4,$5
)
RETURNING *;

-- name: UpdateDean :one
UPDATE current_dean
SET 
    start_date = $1,
    end_date = $2
WHERE dean_id = $3
RETURNING *;

-- name: CreateHod :one
INSERT INTO current_hod(
    lecturer_id,department_id,university_id,start_date,end_date
)VALUES(
    $1,$2,$3,$4,$5
)
RETURNING *;

-- name: UpdateHod :one
UPDATE current_hod
SET 
    start_date = $1,
    end_date = $2
WHERE hod_id = $3
RETURNING *;

-- name: RetrievePendingDeans :many
SELECT * FROM dean_waiting_list WHERE university_id = $1;

-- name: RetrievePendingHods :many
SELECT * FROM hod_waiting_list 
WHERE university_id = $1 AND faculty_id = $2;

-- name: RetrievePendingLecturers :many
SELECT * FROM lecturer_waiting_list
WHERE university_id = $1 AND faculty_id = $2 AND department_id = $3;

-- name: ApproveDean :one
UPDATE dean_waiting_list
SET approved = TRUE
WHERE wait_id = $1
RETURNING *;

-- name: ApproveHod :one
UPDATE hod_waiting_list
SET approved = TRUE
WHERE wait_id = $1
RETURNING *;

-- name: ApproveLecturer :one
UPDATE lecturer_waiting_list
SET approved = TRUE
WHERE wait_id = $1
RETURNING *;


-- name: RequestDeanConfirmation :one
INSERT INTO dean_waiting_list(
    lecturer_id,potential_faculty,additional_message,university_id
)
VALUES(
    $1,$2,$3,$4
)
RETURNING *;

-- name: RequestHodConfirmation :one
INSERT INTO hod_waiting_list(
    lecturer_id,potential_department,additional_message,university_id,faculty_id
)
VALUES(
    $1,$2,$3,$4,$5
)
RETURNING *;

-- name: RequestLecturerConfirmation :one
INSERT INTO lecturer_waiting_list(
    lecturer_id,additional_message,university_id,faculty_id,department_id
)
VALUES(
    $1,$2,$3,$4,$5
)
RETURNING *;









