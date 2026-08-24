-- name: CreateEmailDelivery :one
INSERT INTO email_deliveries (id, notification_id, recipient, provider, attempt)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: MarkEmailDeliverySent :one
UPDATE email_deliveries
SET status = 'sent', provider_message_id = $2, sent_at = NOW()
WHERE id = $1
RETURNING *;

-- name: MarkEmailDeliveryFailed :one
UPDATE email_deliveries
SET status = 'failed', error_message = $2
WHERE id = $1
RETURNING *;

-- name: GetEmailDeliveriesForNotification :many
SELECT *
FROM email_deliveries
WHERE notification_id = $1
ORDER BY attempt ASC, created_at ASC;
