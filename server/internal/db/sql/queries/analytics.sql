-- name: UpdateSummonerStats :exec
CALL update_tft_summoner_stats($1);

-- name: GetSummonerStats :many
SELECT
    queue_id,
	comp_count,
	avg_placement,
	top_4_count,
	top_4_rate,
	top_1_count,
	top_1_rate
FROM
	tft_summoner_stats
WHERE
	summoner_puuid = $1;
