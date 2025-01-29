-- name: GetOldestMatchHistories :many
SELECT
    puuid,
    background_update_timestamp
FROM tft_summoner
WHERE region = ANY(@regions::VARCHAR[])
ORDER BY background_update_timestamp ASC NULLS FIRST
LIMIT $1;

-- name: GetPuuidsWithNullSummonerData :many
SELECT
    puuid
FROM tft_summoner
WHERE summoner_id IS NULL
AND region = @region::VARCHAR
LIMIT $1;

-- name: GetPuuidsWithNullAccountData :many
SELECT
    puuid
FROM tft_summoner
WHERE (name IS NULL OR tag IS NULL) AND NOT 'SKIP_ACCOUNT_DATA' = ANY(flags)
LIMIT $1;

-- name: AddSummonerFlag :exec
UPDATE tft_summoner
    SET flags = array_append(flags, @flag::VARCHAR)
WHERE puuid = $1 AND NOT @flag::VARCHAR = ANY(flags);

-- name: RemoveSummonerFlag :exec
UPDATE tft_summoner
    SET flags = array_remove(flags, @flag::VARCHAR)
WHERE puuid = $1;


-- name: SetBackgroundUpdateTimestamp :exec
UPDATE tft_summoner
SET background_update_timestamp = @background_update_timestamp::TIMESTAMP
WHERE puuid = $1;

-- name: GetPuuidsWithNullRegion :many
SELECT
    puuid
FROM tft_summoner
    WHERE region IS NULL AND NOT 'SKIP_REGION_MATCH' = ANY(flags)
LIMIT $1;
