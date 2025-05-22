-- +goose Up
-- +goose StatementBegin
ALTER TABLE tft_match 
ALTER COLUMN match_date
TYPE timestamptz
USING match_date AT TIME ZONE 'UTC';

ALTER TABLE tft_comp 
ALTER COLUMN match_date
TYPE timestamptz
USING match_date AT TIME ZONE 'UTC';
-- +goose StatementEnd

