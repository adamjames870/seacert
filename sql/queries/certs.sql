-- name: CreateCert :one
INSERT INTO certificates (id, created_at, updated_at, name, cert_number, issuer, issued_date)
VALUES (
           $1, $2, $3, $4, $5, $6, $7
       )
RETURNING *;

-- name: GetCerts :many
SELECT * FROM certificates;

-- name: GetCertFromId :one
SELECT * FROM certificates WHERE id=$1;