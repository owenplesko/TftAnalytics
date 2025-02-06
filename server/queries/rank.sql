-- name: GetSummonerRank :one
SELECT * FROM tft_rank
WHERE summoner_puuid = $1;

-- name: UpsertSummonerRank :exec
INSERT INTO tft_rank (
    summoner_puuid,
    tier,
    rank,
    league_points,
    wins,
    losses
) VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (summoner_puuid) 
DO UPDATE SET 
    tier = EXCLUDED.tier,
    rank = EXCLUDED.rank,
    league_points = EXCLUDED.league_points,
    wins = EXCLUDED.wins,
    losses = EXCLUDED.losses;
