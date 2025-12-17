-- name: CreateCert :one
INSERT INTO certificates (id, created_at, updated_at, user_id, cert_type_id, cert_number, issuer_id, issued_date, alternative_name, remarks)
VALUES (
           $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
       )
RETURNING *;

-- name: GetCerts :many
SELECT * FROM certificates WHERE user_id=$1;

-- name: GetCertFromId :one
SELECT * FROM certificates WHERE id=$1 AND user_id=$2;

-- name: UpdateCertificate :one
UPDATE certificates
SET
    cert_number=COALESCE(sqlc.narg('cert_number'), cert_number),
    issued_date=COALESCE(sqlc.narg('issued_date'), issued_date),
    cert_type_id=COALESCE(sqlc.narg('cert_type_id'), cert_type_id),
    alternative_name=COALESCE(sqlc.narg('alternative_name'), alternative_name),
    remarks=COALESCE(sqlc.narg('remarks'), remarks),
    issuer_id=COALESCE(sqlc.narg('issuer_id'), issuer_id),
    updated_at=NOW()
WHERE id=$1
RETURNING *;