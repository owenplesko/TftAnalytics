-- +goose Up
-- +goose StatementBegin

/* Remove tft_summoner_stats and related procedures / columns */
DROP TABLE IF EXISTS tft_summoner_stats;
DROP PROCEDURE IF EXISTS update_tft_summoner_stats (VARCHAR);
ALTER TABLE tft_comp DROP COLUMN IF EXISTS insert_timestamp;
ALTER TABLE tft_summoner DROP COLUMN IF EXISTS stats_update_timestamp;

/* Refactor idx_name_tag_insensitive to use fn_normalize_str to be consistent */
DROP INDEX IF EXISTS idx_name_tag_insensitive;
CREATE INDEX IF NOT EXISTS idx_name_tag_insensitive ON tft_summoner (fn_normalize_str(name), fn_normalize_str(tag));

-- +goose StatementEnd
