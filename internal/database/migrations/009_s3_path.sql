-- +goose Up
ALTER TABLE certificates ADD COLUMN document_path VARCHAR(512);

-- +goose Down
ALTER TABLE certificates DROP COLUMN document_path;