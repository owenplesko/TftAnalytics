-- name: RefreshUnitPlacement :exec
REFRESH MATERIALIZED VIEW CONCURRENTLY unit_placement;

-- name: GetSummonerStats :one
SELECT
  COUNT(*) AS comp_count,
  COUNT(*) FILTER (WHERE (comp_data ->> 'placement')::INT <= 4) AS top_4_count,
  COUNT(*) FILTER (WHERE (comp_data ->> 'placement')::INT = 1) AS top_1_count,
  AVG((comp_data ->> 'placement')::INT) AS avg_placement,
  COALESCE(SUM((comp_data -> 'timeEliminated')::FLOAT), 0) AS total_seconds_ingame
FROM
  tft_comp
  JOIN tft_match ON tft_comp.match_id = tft_match.id
WHERE
  summoner_puuid = $1
  AND set_number = $2
  AND (
    sqlc.narg('QueueID')::INT IS NULL OR
    sqlc.narg('QueueID')::INT = queue_id
);

