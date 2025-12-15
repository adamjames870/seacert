-- name: CountCertTypes :one

SELECT COUNT (*) AS count_cert_types FROM certificate_types;

-- name: CountCertificates :one

SELECT COUNT (*) AS count_certificates FROM certificates;

-- name: CountIssuers :one

SELECT COUNT (*) AS count_issuers FROM issuers;