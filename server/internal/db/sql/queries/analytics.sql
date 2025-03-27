-- name: UpdateSummonerStats :exec
CALL update_tft_summoner_stats($1);

-- name: GetSummonerStats :one
SELECT
    comp_count,
    avg_placement,
    top_4_count,
    top_4_rate,
    top_1_count,
    top_1_rate
FROM
    tft_summoner_stats
WHERE
    summoner_puuid = $1 
    AND set_number = $2
    AND (queue_id = $3 OR ($3 IS NULL AND queue_id IS NULL));