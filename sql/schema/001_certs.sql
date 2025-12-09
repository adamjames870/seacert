-- +goose Up
CREATE TABLE certificates (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR NOT NULL,
    cert_number VARCHAR NOT NULL,
    issuer VARCHAR NOT NULL,
    issued_date DATE NOT NULL
);

-- +goose Down
DROP TABLE certificates;