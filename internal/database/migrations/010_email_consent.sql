-- +goose Up
ALTER TABLE users
    ADD COLUMN email_consent BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN email_consent_timestamp TIMESTAMP,
    ADD COLUMN email_consent_version VARCHAR,
    ADD COLUMN email_consent_source VARCHAR;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS email_consent;
ALTER TABLE users DROP COLUMN IF EXISTS email_consent_timestamp;
ALTER TABLE users DROP COLUMN IF EXISTS email_consent_version;
ALTER TABLE users DROP COLUMN IF EXISTS email_consent_source;