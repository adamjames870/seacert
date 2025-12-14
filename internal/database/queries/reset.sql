-- name: ResetCerts :exec
DELETE from certificates;

-- name: ResetCertTypes :exec
DELETE from certificate_types;

-- name: ResetIssuers :exec
DELETE from issuers;

