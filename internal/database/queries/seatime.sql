-- name: GetShipTypes :many
SELECT * FROM ship_types ORDER BY name;

-- name: GetVoyageTypes :many
SELECT * FROM voyage_types ORDER BY name;

-- name: GetSeatimePeriodTypes :many
SELECT * FROM seatime_period_types ORDER BY name;

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
INSERT INTO ships (id, created_at, updated_at, name, ship_type_id, imo_number, gt, flag, propulsion_power)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

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
