-- name: CreateCertType :one
INSERT INTO certificate_types (id, created_at, updated_at, name, short_name, stcw_reference, normal_validity_months)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetCertTypes :many
SELECT * FROM certificate_types;

-- name: GetCertTypeFromId :one
SELECT * FROM certificate_types WHERE id=$1;

-- name: GetCertTypeFromName :one
SELECT * FROM certificate_types WHERE name=$1;

-- name: UpdateCertType :one
UPDATE certificate_types
SET
    name = COALESCE(sqlc.narg('name'), name),
    short_name = COALESCE(sqlc.narg('short_name'), short_name),
    stcw_reference = COALESCE(sqlc.narg('stcw_reference'), stcw_reference),
    normal_validity_months = COALESCE(sqlc.narg('normal_validity_months'), normal_validity_months)
WHERE id=$1
RETURNING *;
