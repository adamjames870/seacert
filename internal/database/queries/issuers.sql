-- name: CreateIssuer :one

INSERT INTO issuers (id, created_at, updated_at, name, country, website)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetIssuers :many
SELECT * FROM issuers;

-- name: GetIssuerById :one
SELECT * FROM issuers WHERE id = $1;

-- name: GetIssuerByName :one
SELECT * FROM issuers WHERE name = $1;

-- name: UpdateIssuer :one
UPDATE issuers
SET
    name = COALESCE(sqlc.narg('name'), name),
    country = COALESCE(sqlc.narg('country'), country),
    website = COALESCE(sqlc.narg('website'), website)
WHERE id = $1
RETURNING *;
