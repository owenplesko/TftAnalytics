-- name: GetSummonerRank :one
SELECT * FROM tft_rank
WHERE summoner_id = $1;

-- name: UpsertRank :exec
INSERT INTO tft_rank (
    summoner_id,
    tier,
    rank,
    league_points,
    wins,
    losses
) VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (summoner_id) 
DO UPDATE SET 
    tier = EXCLUDED.tier,
    rank = EXCLUDED.rank,
    league_points = EXCLUDED.league_points,
    wins = EXCLUDED.wins,
    losses = EXCLUDED.losses;

-- name: BatchUpsertRank :batchexec
INSERT INTO tft_rank (
    summoner_id,
    tier,
    rank,
    league_points,
    wins,
    losses
) VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (summoner_id) 
DO UPDATE SET 
    tier = EXCLUDED.tier,
    rank = EXCLUDED.rank,
    league_points = EXCLUDED.league_points,
    wins = EXCLUDED.wins,
    losses = EXCLUDED.losses;
