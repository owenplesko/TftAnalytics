-- +goose Up
-- +goose StatementBegin
UPDATE tft_match
SET game_version = regexp_replace(
    game_version,
    '.*<Releases/(\d+\.\d+)>.*',
    '\1'
)
WHERE game_version ~ '<Releases/\d+\.\d+>';
-- +goose StatementEnd

