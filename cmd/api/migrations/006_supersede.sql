-- +goose Up

CREATE TABLE successions (
    id uuid PRIMARY KEY,
    new_cert uuid NOT NULL,
    old_cert uuid NOT NULL,
    reason varchar(12) NOT NULL,
    CONSTRAINT fk_new_cert_id
        FOREIGN KEY (new_cert)
        REFERENCES certificates(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_old_cert_id
        FOREIGN KEY (old_cert)
        REFERENCES certificates(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_successions_new_cert_id ON successions (new_cert);
CREATE INDEX idx_successions_old_cert_id ON successions (old_cert);

ALTER TABLE certificates
    ADD COLUMN deleted boolean NOT NULL DEFAULT false;

-- +goose Down

ALTER TABLE certificates
    DROP COLUMN deleted;

DROP TABLE successions;