-- name: UpdateSummonerAggregates :exec
CALL update_tft_summoner_comp_aggregate($1);

-- name: GetSummonerAggregates :many
SELECT
    queue_id,
	comp_count,
	avg_placement,
	top_4_rate,
	top_1_rate
FROM
	tft_summoner_comp_aggregate
WHERE
	summoner_puuid = $1;
