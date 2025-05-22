-- +goose Up
-- +goose StatementBegin
DROP MATERIALIZED VIEW IF EXISTS unit_placement;
CREATE MATERIALIZED VIEW unit_placement AS SELECT
	tft_match.set_number,
	tft_match.game_version,
	tft_comp.rank,
	AVG((tft_comp.comp_data -> 'placement')::INTEGER) AS avg_placement,
	units ->> 'characterId' as unit,
	COUNT(*) AS frequency
FROM
	tft_comp JOIN tft_match ON tft_comp.match_id = tft_match.id,
	LATERAL jsonb_array_elements(comp_data -> 'units') AS units
WHERE tft_comp.match_date >= date_trunc('day', CURRENT_TIMESTAMP - INTERVAL '7 days')
GROUP BY
	set_number, game_version, rank, unit;
-- +goose StatementEnd

