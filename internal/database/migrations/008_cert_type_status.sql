-- +goose Up
ALTER TABLE certificate_types
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'approved',
    ADD COLUMN created_by UUID;

ALTER TABLE certificate_types
    ADD CONSTRAINT fk_cert_types_created_by
        FOREIGN KEY (created_by) REFERENCES users(id);

-- +goose Down

ALTER TABLE certificate_types DROP CONSTRAINT IF EXISTS fk_cert_types_created_by;
ALTER TABLE certificate_types DROP COLUMN IF EXISTS created_by;
ALTER TABLE certificate_types DROP COLUMN IF EXISTS status;
