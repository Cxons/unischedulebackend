-- name: UpdateAdminInfo :one
UPDATE university_admin
SET admin_middle_name = $1, admin_phone_number = $2, admin_staff_card = $3, admin_number = $4, university_id = $5
WHERE admin_id = $6
RETURNING *;


-- name: RetrieveAdmin :one
SELECT 
    admin_first_name,
    admin_middle_name,
    admin_email,
    admin_phone_number,
    admin_staff_card,
    admin_number,
    university_id
FROM university_admin
WHERE admin_id = $1;

-- name: RetrieveDean :one
SELECT
    lecturer_id,
    faculty_id,
    university_id,
    start_date,
    end_date
FROM current_dean
WHERE dean_id = $1;


-- name: RetrieveHod :one
SELECT 
    lecturer_id,
    department_id,
    university_id,
    start_date,
    end_date
FROM current_hod
WHERE hod_id = $1;


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
SELECT * FROM dean_waiting_list 
WHERE university_id = $1 
AND approved = FALSE;

-- name: RetrievePendingHods :many
SELECT * FROM hod_waiting_list 
WHERE university_id = $1 
AND faculty_id = $2 
AND approved = FALSE;

-- name: RetrievePendingLecturers :many
SELECT * FROM lecturer_waiting_list
WHERE university_id = $1 
AND faculty_id = $2 
AND department_id = $3
AND approved = FALSE;

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

-- name: CheckDeanConfirmation :one
SELECT 
    l.lecturer_id,
    l.lecturer_first_name,
    l.lecturer_last_name,
    l.lecturer_profile_pic,
    dwl.potential_faculty,
    dwl.additional_message,
    dwl.approved
FROM 
    dean_waiting_list dwl
INNER JOIN 
    lecturers l ON dwl.lecturer_id = l.lecturer_id
WHERE 
    dwl.wait_id = $1;

-- name: RequestHodConfirmation :one
INSERT INTO hod_waiting_list(
    lecturer_id,potential_department,additional_message,university_id,faculty_id
)
VALUES(
    $1,$2,$3,$4,$5
)
RETURNING *;

-- name: CheckHodConfirmation :one
SELECT 
    l.lecturer_id,
    l.lecturer_first_name,
    l.lecturer_last_name,
    l.lecturer_profile_pic,
    hwl.potential_department,
    hwl.additional_message,
    hwl.approved
FROM 
    hod_waiting_list hwl
INNER JOIN 
    lecturers l ON hwl.lecturer_id = l.lecturer_id
WHERE 
    hwl.wait_id = $1;

-- name: RequestLecturerConfirmation :one
INSERT INTO lecturer_waiting_list(
    lecturer_id,additional_message,university_id,faculty_id,department_id
)
VALUES(
    $1,$2,$3,$4,$5
)
RETURNING *;

-- name: CheckLecturerConfirmation :one
SELECT 
    l.lecturer_id,
    l.lecturer_first_name,
    l.lecturer_last_name,
    l.lecturer_profile_pic,
    lwl.additional_message,
    lwl.approved 
FROM lecturer_waiting_list lwl
INNER JOIN 
    lecturers l ON lwl.lecturer_id = l.lecturer_id
WHERE 
    lwl.wait_id = $1;









