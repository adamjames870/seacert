-- name: ResetSuccessions :exec
DELETE from successions;

-- name: ResetCerts :exec
DELETE from certificates;

-- name: ResetCertTypes :exec
DELETE from certificate_types;

-- name: ResetIssuers :exec
DELETE from issuers;

-- name: ResetUsers :exec
DELETE from users;