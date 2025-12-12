-- name: CreateCert :one
INSERT INTO certificates (id, created_at, updated_at, cert_type_id, cert_number, issuer, issued_date, alternative_name, remarks)
VALUES (
           $1, $2, $3, $4, $5, $6, $7, $8, $9
       )
RETURNING *;

-- name: GetCerts :many
SELECT * FROM certificates;

-- name: GetCertFromId :one
SELECT * FROM certificates WHERE id=$1;