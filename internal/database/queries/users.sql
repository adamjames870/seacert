-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, forename, surname, nationality)
VALUES ($1, $2, $3, $4, $5, $6, $7)
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
    updated_at = NOW()
WHERE id = $1
RETURNING *;