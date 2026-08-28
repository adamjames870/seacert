-- name: CreateNotification :one
INSERT INTO notifications (id, user_id, notification_type, notification_key, payload, scheduled_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPendingNotifications :many
SELECT *
FROM notifications
WHERE status = 'pending' AND scheduled_at <= NOW()
ORDER BY scheduled_at ASC, created_at ASC;

-- name: GetNotificationsStatusSummary :many
SELECT status, COUNT(*) AS count
FROM notifications
GROUP BY status;

-- name: MarkNotificationProcessing :one
UPDATE notifications
SET status = 'processing', updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: MarkNotificationCompleted :one
UPDATE notifications
SET status = 'completed', processed_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: MarkNotificationFailed :one
UPDATE notifications
SET status = 'failed', processed_at = NOW(), updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetUsersEligibleForNoCertificates7Day :many
SELECT u.id
FROM users u
WHERE u.created_at <= NOW() - INTERVAL '7 days'
  AND u.created_at >= NOW() - INTERVAL '21 days'
  AND NOT EXISTS (
      SELECT 1
      FROM certificates c
      WHERE c.user_id = u.id
  )
  AND NOT EXISTS (
      SELECT 1
      FROM notifications n
      WHERE n.user_id = u.id
        AND n.notification_type = 'no_certificates_7d'
  );

-- name: GetUsersEligibleForNoCertificates1Month :many
SELECT u.id
FROM users u
WHERE u.created_at <= NOW() - INTERVAL '1 month'
  AND NOT EXISTS (
      SELECT 1
      FROM certificates c
      WHERE c.user_id = u.id
  )
  AND NOT EXISTS (
      SELECT 1
      FROM notifications n
      WHERE n.user_id = u.id
        AND n.notification_type = 'no_certificates_1m'
  );
