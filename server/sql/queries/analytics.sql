-- name: GetSummonerStats :one
WITH
	selected_data AS (
		SELECT
			(comp_data ->> 'placement')::INT AS placement,
			(comp_data ->> 'timeEliminated')::FLOAT AS seconds_ingame
		FROM
			tft_comp
			JOIN tft_match ON tft_comp.match_id = tft_match.id
		WHERE
			summoner_puuid = $1
			AND set_number = $2
			AND (
				sqlc.narg('QueueID')::INT IS NULL OR
				sqlc.narg('QueueID')::INT = queue_id
			)
	),
	aggregate_stats AS (
		SELECT
			COUNT(*) AS comp_count,
			COUNT(*) FILTER (
				WHERE
					placement <= 4
			) AS top_4_count,
			COUNT(*) FILTER (
				WHERE
					placement = 1
			) AS top_1_count,
			AVG(placement) AS avg_placement,
			SUM(seconds_ingame) AS total_seconds_ingame
		FROM
			selected_data
	),
	derived_stats AS (
		SELECT
			top_1_count::float / comp_count AS top_1_rate,
			top_4_count::float / comp_count AS top_4_rate
		FROM
			aggregate_stats
	)
SELECT
	*
FROM
	aggregate_stats
	CROSS JOIN derived_stats;
