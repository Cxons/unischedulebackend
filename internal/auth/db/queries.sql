-- name: RegisterStudent :one
INSERT INTO students(
    student_first_name,student_last_name,student_email,student_password
) VALUES (
    $1,$2,$3,$4
)
RETURNING *;

-- name: RegisterLecturer :one
INSERT INTO lecturers(
    lecturer_first_name,lecturer_last_name,lecturer_email,lecturer_password
) VALUES(
    $1,$2,$3,$4
)
RETURNING *;

-- name: RegisterUniversityAdmin :one
INSERT INTO university_admin(
    admin_first_name,admin_last_name,admin_email,admin_password
) VALUES (
    $1,$2,$3,$4
)
RETURNING *;

-- name: RetrieveStudentEmail :one
SELECT student_id,student_email,student_password FROM students WHERE student_email = $1;

-- name: RetrieveLecturerEmail :one
SELECT lecturer_id,lecturer_email,lecturer_password FROM lecturers WHERE lecturer_email = $1;

-- name: RetrieveAdminEmail :one
SELECT admin_id,admin_email,admin_password FROM university_admin WHERE admin_email = $1;

-- name: AddRefreshToken :one
INSERT INTO refresh_tokens(
    refresh_token, expires_at, user_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens 
SET refresh_token = $1, expires_at = $2, updated_at = NOW()
WHERE user_id = $3
RETURNING *;

-- name: CheckAndReturnToken :one
SELECT refresh_token,expires_at,user_id,is_revoked
FROM refresh_tokens
WHERE user_id = $1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET is_revoked = TRUE
WHERE user_id = $1
RETURNING *;

-- name: DeleteRefreshToken :one
DELETE FROM refresh_tokens
WHERE user_id = $1
RETURNING *;

-- name: InsertOtp :one
INSERT INTO otps(
    otp,expires_at,email,user_type
)VALUES(
    $1,$2,$3,$4
)
RETURNING *;

-- name: UpdateOtp :one
UPDATE otps
SET otp = $1, expires_at = $2
WHERE email = $3
RETURNING *;

-- name: RetrieveOtp :one
SELECT otp,expires_at FROM otps WHERE email = $1;
