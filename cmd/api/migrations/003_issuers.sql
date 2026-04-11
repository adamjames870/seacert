-- +goose Up

CREATE TABLE issuers (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR(255) NOT NULL UNIQUE,
    country VARCHAR(2),
    website VARCHAR(255)
);

ALTER TABLE certificates
    ADD COLUMN issuer_id uuid NOT NULL default '00000000-0000-0000-0000-000000000000';

ALTER TABLE certificates
    ADD CONSTRAINT fk_cert_issuer_id
    FOREIGN KEY (issuer_id)
    REFERENCES issuers(id);

ALTER TABLE certificates
    DROP COLUMN issuer;

-- +goose Down

ALTER TABLE certificates
    ADD COLUMN issuer VARCHAR(255);

ALTER TABLE certificates
    DROP CONSTRAINT fk_cert_issuer_id;

ALTER TABLE certificates
    DROP COLUMN issuer_id;

DROP TABLE issuers;
