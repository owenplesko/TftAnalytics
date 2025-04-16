-- name: GetSummonerByPuuid :one
SELECT *
FROM tft_summoner 
WHERE puuid = $1;

-- name: GetSummonerBySummonerId :one
SELECT *
FROM tft_summoner
WHERE summoner_id = $1;

-- name: GetSummonerByNameTag :one
SELECT *
FROM tft_summoner 
WHERE fn_normalize_str(name) = fn_normalize_str(@name::VARCHAR)
AND fn_normalize_str(tag) = fn_normalize_str(@tag::VARCHAR);

-- name: SummonerExistsByPuuid :one
SELECT EXISTS (
    SELECT * FROM tft_summoner WHERE puuid = $1
);

-- name: UpsertSummoner :exec
INSERT INTO tft_summoner (
    puuid,
    region,
    name,
    tag,
    summoner_id,
    profile_icon_id,
    summoner_level
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) ON CONFLICT (puuid) DO UPDATE
SET name = EXCLUDED.name,
    tag = EXCLUDED.tag,
    summoner_id = EXCLUDED.summoner_id,
    profile_icon_id = EXCLUDED.profile_icon_id,
    summoner_level = EXCLUDED.summoner_level;

-- name: UpdateRegion :exec
UPDATE tft_summoner
SET region = @region::VARCHAR
WHERE puuid = $1;

-- name: GetOldestMatchesAfter :many
SELECT
    puuid,
    matches_before_timestamp
FROM tft_summoner
WHERE region = @region::VARCHAR
ORDER BY matches_before_timestamp ASC NULLS FIRST
LIMIT $1;

-- name: SetMatchesBeforeTimestamp :exec
UPDATE tft_summoner
    SET matches_before_timestamp = @matches_before_timestamp::TIMESTAMP
WHERE puuid = $1;
