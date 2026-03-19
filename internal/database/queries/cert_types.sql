-- name: CreateCertType :one
INSERT INTO certificate_types (id, created_at, updated_at, name, short_name, stcw_reference, normal_validity_months, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetCertTypes :many
SELECT * FROM certificate_types;

-- name: GetCertTypesForUser :many
SELECT * FROM certificate_types
WHERE status = 'approved' OR created_by = $1;

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
    normal_validity_months = COALESCE(sqlc.narg('normal_validity_months'), normal_validity_months),
    status = COALESCE(sqlc.narg('status'), status),
    updated_at = NOW()
WHERE id=$1
RETURNING *;

-- name: UpdateCertTypeReferences :exec
UPDATE certificates SET cert_type_id = $2 WHERE cert_type_id = $1;

-- name: DeleteCertType :exec
DELETE FROM certificate_types WHERE id = $1;
