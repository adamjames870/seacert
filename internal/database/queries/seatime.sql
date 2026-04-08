-- name: GetShipTypes :many
SELECT * FROM ship_types ORDER BY name;

-- name: GetShipTypeById :one
SELECT * FROM ship_types WHERE id = $1;

-- name: CreateShipType :one
INSERT INTO ship_types (id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateShipType :one
UPDATE ship_types
SET name = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteShipType :exec
DELETE FROM ship_types WHERE id = $1;

-- name: GetVoyageTypes :many
SELECT * FROM voyage_types ORDER BY name;

-- name: GetVoyageTypeById :one
SELECT * FROM voyage_types WHERE id = $1;

-- name: CreateVoyageType :one
INSERT INTO voyage_types (id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateVoyageType :one
UPDATE voyage_types
SET name = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteVoyageType :exec
DELETE FROM voyage_types WHERE id = $1;

-- name: GetSeatimePeriodTypes :many
SELECT * FROM seatime_period_types ORDER BY name;

-- name: GetSeatimePeriodTypeById :one
SELECT * FROM seatime_period_types WHERE id = $1;

-- name: CreateSeatimePeriodType :one
INSERT INTO seatime_period_types (id, name, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateSeatimePeriodType :one
UPDATE seatime_period_types
SET name = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteSeatimePeriodType :exec
DELETE FROM seatime_period_types WHERE id = $1;

-- name: GetShipByImo :one
SELECT s.*, st.name as ship_type_name
FROM ships s
JOIN ship_types st ON s.ship_type_id = st.id
WHERE s.imo_number = $1;

-- name: GetShipById :one
SELECT s.*, st.name as ship_type_name
FROM ships s
JOIN ship_types st ON s.ship_type_id = st.id
WHERE s.id = $1;

-- name: CreateShip :one
INSERT INTO ships (id, created_at, updated_at, name, ship_type_id, imo_number, gt, flag, propulsion_power, status, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetShips :many
SELECT s.*, st.name as ship_type_name
FROM ships s
JOIN ship_types st ON s.ship_type_id = st.id;

-- name: GetShipsForUser :many
SELECT s.*, st.name as ship_type_name
FROM ships s
JOIN ship_types st ON s.ship_type_id = st.id
WHERE s.status = 'approved' OR s.created_by = $1;

-- name: UpdateShip :one
UPDATE ships
SET name = $2, ship_type_id = $3, imo_number = $4, gt = $5, flag = $6, propulsion_power = $7, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateShipStatus :one
UPDATE ships
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateShipReferences :exec
UPDATE seatime SET ship_id = $2 WHERE ship_id = $1;

-- name: DeleteShip :exec
DELETE FROM ships WHERE id = $1;

-- name: CreateSeatime :one
INSERT INTO seatime (id, user_id, ship_id, voyage_type_id, created_at, updated_at, start_date, start_location, end_date, end_location, total_days, company, capacity, is_watchkeeping)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: CreateSeatimePeriod :one
INSERT INTO seatime_periods (id, seatime_id, period_type_id, start_date, end_date, days, remarks)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSeatimeByUserId :many
SELECT 
    s.*, 
    sh.name as ship_name, 
    sh.imo_number as ship_imo,
    sh.gt as ship_gt,
    sh.flag as ship_flag,
    sh.propulsion_power as ship_propulsion_power,
    sh.ship_type_id as ship_type_id,
    sh.created_at as ship_created_at,
    sh.updated_at as ship_updated_at,
    sh.status as ship_status,
    sh.created_by as ship_created_by,
    st.name as ship_type_name,
    vt.name as voyage_type_name
FROM seatime s
JOIN ships sh ON s.ship_id = sh.id
JOIN ship_types st ON sh.ship_type_id = st.id
JOIN voyage_types vt ON s.voyage_type_id = vt.id
WHERE s.user_id = $1
ORDER BY s.start_date DESC;

-- name: GetSeatimePeriods :many
SELECT sp.*, spt.name as period_type_name
FROM seatime_periods sp
JOIN seatime_period_types spt ON sp.period_type_id = spt.id
WHERE sp.seatime_id = $1;

-- name: DeleteSeatime :exec
DELETE FROM seatime WHERE id = $1 AND user_id = $2;

-- name: UpdateSeatime :one
UPDATE seatime
SET ship_id = $3, voyage_type_id = $4, updated_at = $5, start_date = $6, start_location = $7, end_date = $8, end_location = $9, total_days = $10, company = $11, capacity = $12, is_watchkeeping = $13
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteSeatimePeriods :exec
DELETE FROM seatime_periods WHERE seatime_id = $1;

-- name: GetOverlappingSeatime :many
SELECT * FROM seatime
WHERE user_id = $1
AND (
    (start_date <= sqlc.arg(new_start_date) AND end_date >= sqlc.arg(new_start_date)) -- existing covers new start
    OR (start_date <= sqlc.arg(new_end_date) AND end_date >= sqlc.arg(new_end_date)) -- existing covers new end
    OR (start_date >= sqlc.arg(new_start_date) AND end_date <= sqlc.arg(new_end_date)) -- new covers existing
)
AND (sqlc.narg(current_id)::UUID IS NULL OR id != sqlc.narg(current_id));
