-- +goose Up
CREATE TABLE certificate_types (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR(255) NOT NULL UNIQUE,
    short_name VARCHAR(24) NOT NULL UNIQUE,
    stcw_reference VARCHAR(255),
    normal_validity_months INTEGER
);

ALTER TABLE certificates
    ADD COLUMN cert_type_id UUID NOT NULL default '00000000-0000-0000-0000-000000000000';

ALTER TABLE certificates
    ADD CONSTRAINT fk_cert_cert_type_id
        FOREIGN KEY (cert_type_id)
        REFERENCES certificate_types(id);

ALTER TABLE certificates
    ADD COLUMN alternative_name VARCHAR(255);

ALTER TABLE certificates
    ADD COLUMN remarks TEXT;

ALTER TABLE certificates
DROP COLUMN name;

-- +goose Down

ALTER TABLE certificates ADD COLUMN name VARCHAR(255);
ALTER TABLE certificates DROP COLUMN remarks;
ALTER TABLE certificates DROP COLUMN alternative_name;
ALTER TABLE certificates DROP COLUMN cert_type_id;
DROP TABLE certificate_types;

