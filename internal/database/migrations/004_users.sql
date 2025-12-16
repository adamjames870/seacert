-- +goose Up
CREATE TABLE users (
                       id uuid PRIMARY KEY,
                       created_at TIMESTAMP NOT NULL,
                       updated_at TIMESTAMP NOT NULL,
                       forename VARCHAR NOT NULL,
                       surname VARCHAR NOT NULL,
                       email VARCHAR NOT NULL UNIQUE,
                       nationality VARCHAR(2) NOT NULL
);

ALTER TABLE certificates
    ADD COLUMN user_id uuid NOT NULL default '00000000-0000-0000-0000-000000000000';

ALTER TABLE certificates
    ADD CONSTRAINT fk_certificates_users
        FOREIGN KEY (user_id)
            REFERENCES users(id);

CREATE TABLE refresh_tokens(
                               token VARCHAR PRIMARY KEY,
                               created_at TIMESTAMP NOT NULL,
                               updated_at TIMESTAMP NOT NULL,
                               user_id UUID NOT NULL,
                               expires_at TIMESTAMP NOT NULL,
                               revoked_at TIMESTAMP,
                               CONSTRAINT fk_user_id
                                   FOREIGN KEY (user_id)
                                       REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down

ALTER TABLE certificates
    DROP COLUMN user_id;

DROP TABLE refresh_tokens;

DROP TABLE users;