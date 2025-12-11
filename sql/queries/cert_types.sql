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