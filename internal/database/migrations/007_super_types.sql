-- +goose Up

CREATE DOMAIN certificate_succession AS text
CHECK (VALUE IN ('updated', 'replaced'));

CREATE TABLE certificate_type_successions (
    id uuid PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    replacing_cert_type uuid NOT NULL,
    replaceable_cert_type uuid NOT NULL,
    replace_reason certificate_succession NOT NULL,
    CONSTRAINT fk_replacing_cert_type_id
     FOREIGN KEY (replacing_cert_type)
         REFERENCES certificate_types(id)
         ON DELETE CASCADE,
    CONSTRAINT fk_replaceable_cert_type_id
     FOREIGN KEY (replaceable_cert_type)
         REFERENCES certificate_types(id)
         ON DELETE CASCADE
);

CREATE INDEX idx_successions_replacing_cert_type ON certificate_type_successions (replacing_cert_type);
CREATE INDEX idx_successions_replaceable_cert_type ON certificate_type_successions (replaceable_cert_type);

ALTER TABLE successions
    ALTER COLUMN reason
    TYPE certificate_succession
    USING reason::certificate_succession;

-- +goose Down

ALTER TABLE successions
    ALTER COLUMN reason
    TYPE varchar(12);

DROP TABLE certificate_type_successions;

DROP DOMAIN certificate_succession;