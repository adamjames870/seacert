-- +goose Up

CREATE TABLE notifications (
                               id UUID PRIMARY KEY,

                               user_id UUID NULL,

                               notification_type TEXT NOT NULL,

                               notification_key TEXT NOT NULL UNIQUE,

                               status TEXT NOT NULL DEFAULT 'pending'
                                   CHECK (status IN (
                                                     'pending',
                                                     'processing',
                                                     'completed',
                                                     'failed'
                                       )),

                               payload JSONB NOT NULL DEFAULT '{}'::jsonb,

                               scheduled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                               processed_at TIMESTAMPTZ NULL,

                               created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                               updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                               CONSTRAINT fk_notifications_user
                                   FOREIGN KEY (user_id)
                                       REFERENCES users(id)
                                       ON DELETE SET NULL
);

CREATE INDEX idx_notifications_status_scheduled_at
    ON notifications (status, scheduled_at);

CREATE INDEX idx_notifications_user_id
    ON notifications (user_id);

CREATE INDEX idx_notifications_type
    ON notifications (notification_type);


-- +goose Down

DROP TABLE IF EXISTS notifications;