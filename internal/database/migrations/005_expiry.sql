-- +goose Up

ALTER TABLE certificates
    ADD COLUMN manual_expiry timestamp;

-- +goose Down

ALTER TABLE certificates
    DROP COLUMN manual_expiry;