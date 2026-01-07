-- name: CreateTypeSuccession :one
INSERT INTO certificate_type_successions
    (id, replacing_cert_type, replaceable_cert_type, replace_reason)
VALUES
    ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteTypeSuccession :exec
DELETE FROM certificate_type_successions
WHERE id = $1;

-- name: GetUpdateSuccessions :many
SELECT
    succession.id, type.id, type.name
FROM
    certificate_types type
INNER JOIN
    certificate_type_successions succession
ON
    succession.replacing_cert_type = type.id
WHERE
    succession.replaceable_cert_type = $1
AND
    succession.replace_reason = 'updated';

-- name: GetReplaceSuccessions :many
SELECT
    succession.id, type.id, type.name
FROM
    certificate_types type
INNER JOIN
    certificate_type_successions succession
ON
    succession.replacing_cert_type = type.id
WHERE
    succession.replaceable_cert_type = $1
AND
    succession.replace_reason = 'replaced';

-- name: GetAllReplaceableByMe :many
SELECT
    succession.id, type.id, type.name, succession.replace_reason
FROM
    certificate_types type
INNER JOIN
    certificate_type_successions succession
ON
    succession.replacing_cert_type = type.id
WHERE
    succession.replacing_cert_type = $1;

-- name: GetAllThatCanReplaceMe :many
SELECT
    succession.id, type.id, type.name, succession.replace_reason
FROM
    certificate_types type
        INNER JOIN
    certificate_type_successions succession
    ON
        succession.replaceable_cert_type = type.id
WHERE
    succession.replaceable_cert_type = $1;
