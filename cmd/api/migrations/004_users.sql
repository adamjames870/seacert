-- +goose Up
CREATE TABLE users (
                       id uuid PRIMARY KEY,
                       created_at TIMESTAMP NOT NULL,
                       updated_at TIMESTAMP NOT NULL,
                       forename VARCHAR,
                       surname VARCHAR,
                       email VARCHAR NOT NULL UNIQUE,
                       nationality VARCHAR(2)
);

ALTER TABLE certificates
    ADD COLUMN user_id uuid NOT NULL default '00000000-0000-0000-0000-000000000000';

ALTER TABLE certificates
    ADD CONSTRAINT fk_certificates_users
        FOREIGN KEY (user_id)
            REFERENCES users(id);

-- +goose Down

ALTER TABLE certificates
    DROP COLUMN user_id;

DROP TABLE users;