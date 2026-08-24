-- +goose Up

CREATE TABLE email_deliveries (
                                  id UUID PRIMARY KEY,

                                  notification_id UUID NOT NULL,

                                  recipient TEXT NOT NULL,

                                  provider TEXT NOT NULL,

                                  provider_message_id TEXT NULL,

                                  status TEXT NOT NULL DEFAULT 'pending'
                                      CHECK (status IN (
                                                        'pending',
                                                        'sent',
                                                        'failed'
                                          )),

                                  attempt INTEGER NOT NULL DEFAULT 1
                                      CHECK (attempt > 0),

                                  error_message TEXT NULL,

                                  sent_at TIMESTAMPTZ NULL,

                                  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                  CONSTRAINT fk_email_deliveries_notification
                                      FOREIGN KEY (notification_id)
                                          REFERENCES notifications(id)
                                          ON DELETE CASCADE
);

CREATE INDEX idx_email_deliveries_notification_id
    ON email_deliveries (notification_id);

CREATE INDEX idx_email_deliveries_status
    ON email_deliveries (status);

CREATE INDEX idx_email_deliveries_provider_message_id
    ON email_deliveries (provider_message_id);


-- +goose Down

DROP TABLE IF EXISTS email_deliveries;