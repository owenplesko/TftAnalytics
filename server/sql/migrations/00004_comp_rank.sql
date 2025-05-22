-- +goose Up
-- +goose StatementBegin
ALTER TABLE tft_comp ADD COLUMN rank VARCHAR NOT NULL DEFAULT 'unknown';
ALTER TABLE tft_comp ALTER COLUMN rank DROP DEFAULT;
-- +goose StatementEnd

