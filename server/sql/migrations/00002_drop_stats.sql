-- +goose Up
-- +goose StatementBegin

/* Remove tft_summoner_stats and related procedures / columns */
DROP TABLE IF EXISTS tft_summoner_stats;
DROP PROCEDURE IF EXISTS update_tft_summoner_stats (VARCHAR);
ALTER TABLE tft_comp DROP COLUMN IF EXISTS insert_timestamp;
ALTER TABLE tft_summoner DROP COLUMN IF EXISTS stats_update_timestamp;

/* create fn_normalize_str */
CREATE OR REPLACE FUNCTION fn_normalize_str(input TEXT)
RETURNS TEXT AS $$
BEGIN
  RETURN LOWER(REPLACE(input, ' ', ''));
END;
$$ LANGUAGE plpgsql IMMUTABLE;

/* Refactor idx_name_tag_insensitive to use fn_normalize_str to be consistent */
DROP INDEX IF EXISTS idx_name_tag_insensitive;
CREATE INDEX IF NOT EXISTS idx_name_tag_insensitive ON tft_summoner (fn_normalize_str(name), fn_normalize_str(tag));

-- +goose StatementEnd
