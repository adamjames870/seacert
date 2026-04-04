-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, forename, surname, nationality, email_consent, email_consent_timestamp, email_consent_version, email_consent_source)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET
    forename = COALESCE(sqlc.narg('forename'), forename),
    surname = COALESCE(sqlc.narg('surname'), surname),
    nationality = COALESCE(sqlc.narg('nationality'), nationality),
    email_consent = COALESCE(sqlc.narg('email_consent'), email_consent),
    email_consent_timestamp = COALESCE(sqlc.narg('email_consent_timestamp'), email_consent_timestamp),
    email_consent_version = COALESCE(sqlc.narg('email_consent_version'), email_consent_version),
    email_consent_source = COALESCE(sqlc.narg('email_consent_source'), email_consent_source),
    updated_at = NOW()
WHERE id = $1
RETURNING *;