-- name: CreateCert :one
INSERT INTO certificates (id, created_at, updated_at, user_id, cert_type_id, cert_number, issuer_id, issued_date, alternative_name, remarks, manual_expiry)
VALUES (
           $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
       )
RETURNING *;

-- name: CreateSuccession :one
INSERT INTO successions (id, old_cert, new_cert, reason) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCerts :many
SELECT
    certificates.id AS id,
    certificates.created_at AS created_at,
    certificates.updated_at AS updated_at,
    cert_number,
    issued_date,
    alternative_name,
    remarks,
    manual_expiry,
    deleted,
    EXISTS (
        SELECT true
        FROM successions s
        WHERE s.new_cert = certificates.id
    ) as has_predecessors,
    EXISTS (
        SELECT true
        FROM successions s
        WHERE s.old_cert = certificates.id
    ) as has_successor,
    certificate_types.id AS cert_type_id,
    certificate_types.created_at AS cert_type_created_at,
    certificate_types.updated_at AS cert_type_updated_at,
    certificate_types.name AS cert_type_name,
    certificate_types.short_name AS cert_type_short_name,
    certificate_types.stcw_reference AS cert_type_stcw_reference,
    certificate_types.normal_validity_months AS normal_validity_months,
    issuers.id AS issuer_id,
    issuers.created_at AS issuer_created_at,
    issuers.updated_at AS issuer_updated_at,
    issuers.name AS issuer_name,
    issuers.country AS issuer_country,
    issuers.website AS issuer_website
FROM certificates
INNER JOIN certificate_types ON certificate_types.id=certificates.cert_type_id
INNER JOIN issuers ON issuers.id=certificates.issuer_id
WHERE user_id=$1;

-- name: GetCertFromId :one
SELECT
    certificates.id AS id,
    certificates.created_at AS created_at,
    certificates.updated_at AS updated_at,
    cert_number,
    issued_date,
    alternative_name,
    remarks,
    manual_expiry,
    deleted,
    EXISTS (
        SELECT true
        FROM successions s
        WHERE s.new_cert = certificates.id
    ) as has_predecessors,
    EXISTS (
        SELECT true
        FROM successions s
        WHERE s.old_cert = certificates.id
    ) as has_successor,
    certificate_types.id AS cert_type_id,
    certificate_types.created_at AS cert_type_created_at,
    certificate_types.updated_at AS cert_type_updated_at,
    certificate_types.name AS cert_type_name,
    certificate_types.short_name AS cert_type_short_name,
    certificate_types.stcw_reference AS cert_type_stcw_reference,
    certificate_types.normal_validity_months AS normal_validity_months,
    issuers.id AS issuer_id,
    issuers.created_at AS issuer_created_at,
    issuers.updated_at AS issuer_updated_at,
    issuers.name AS issuer_name,
    issuers.country AS issuer_country,
    issuers.website AS issuer_website
FROM certificates
INNER JOIN certificate_types ON certificate_types.id=certificates.cert_type_id
INNER JOIN issuers ON issuers.id=certificates.issuer_id
WHERE certificates.id=$1 AND certificates.user_id=$2;

-- name: UpdateCertificate :one
UPDATE certificates
SET
    cert_number=COALESCE(sqlc.narg('cert_number'), cert_number),
    issued_date=COALESCE(sqlc.narg('issued_date'), issued_date),
    cert_type_id=COALESCE(sqlc.narg('cert_type_id'), cert_type_id),
    alternative_name=COALESCE(sqlc.narg('alternative_name'), alternative_name),
    remarks=COALESCE(sqlc.narg('remarks'), remarks),
    issuer_id=COALESCE(sqlc.narg('issuer_id'), issuer_id),
    manual_expiry=COALESCE(sqlc.narg('manual_expiry'), manual_expiry),
    deleted=COALESCE(sqlc.narg('deleted'), deleted),
    updated_at=NOW()
WHERE id=$1
RETURNING *;

-- name: GetPredecessors :many
SELECT old_cert FROM successions WHERE new_cert=$1;

-- name: GetSuccessors :many
SELECT new_cert FROM successions WHERE old_cert=$1;

-- name: DeleteCert :exec
DELETE FROM certificates WHERE id=$1 AND user_id=$2;