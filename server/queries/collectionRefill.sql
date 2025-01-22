-- name: GetOldestMatchHistories :many
SELECT
    puuid,
    background_update_timestamp
FROM tft_summoner
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
WHERE (name IS NULL OR tag IS NULL) AND NOT skip_account
LIMIT $1;

-- name: SetSkipAccountFlag :exec
UPDATE tft_summoner
SET skip_account = $2
WHERE puuid = $1;

-- name: SetBackgroundUpdateTimestamp :exec
UPDATE tft_summoner
SET background_update_timestamp = @background_update_timestamp::TIMESTAMP
WHERE puuid = $1;

-- name: GetPuuidsWithNullRegion :many
SELECT
    puuid
FROM tft_summoner
WHERE region IS NULL
LIMIT $1;
