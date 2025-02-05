-- name: GetSummonerRank :one
SELECT * FROM tft_rank
WHERE summoner_puuid = $1;

-- name: InsertSummonerRank :exec
INSERT INTO tft_rank (
    summoner_puuid,
    tier,
    rank,
    league_points,
    wins,
    losses
) VALUES ( $1, $2, $3, $4, $5, $6);
