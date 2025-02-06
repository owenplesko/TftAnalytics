-- name: GetSummonerByPuuid :one
SELECT *
FROM tft_summoner 
WHERE puuid = $1;

-- name: SummonerExistsByNameTag :one
SELECT EXISTS (
    SELECT * FROM tft_summoner 
    WHERE REPLACE(LOWER(name), ' ', '') = REPLACE(LOWER(@name::VARCHAR), ' ', '')
    AND REPLACE(LOWER(tag), ' ', '') = REPLACE(LOWER(@tag::VARCHAR), ' ', '')
);

-- name: GetSummonerByNameTag :one
SELECT *
FROM tft_summoner 
WHERE REPLACE(LOWER(name), ' ', '') = REPLACE(LOWER(@name::VARCHAR), ' ', '')
AND REPLACE(LOWER(tag), ' ', '') = REPLACE(LOWER(@tag::VARCHAR), ' ', '');

-- name: BatchUpsertPuuids :batchexec
INSERT INTO tft_summoner (
    puuid,
    region
) VALUES (
    $1, @region::VARCHAR
) ON CONFLICT (puuid) DO NOTHING;

-- name: UpsertAccount :exec
INSERT INTO tft_summoner (
    puuid,
    name,
    tag
) VALUES (
    $1, @name::VARCHAR, @tag::VARCHAR
) ON CONFLICT(puuid) DO UPDATE SET name = @name::VARCHAR, tag = @tag::VARCHAR;

-- name: UpdateRegion :exec
UPDATE tft_summoner
SET region = @region::VARCHAR
WHERE puuid = $1;

-- name: UpdateSummoner :exec
UPDATE tft_summoner
SET summoner_id = @summoner_id::VARCHAR,
    profile_icon_id = @profile_icon_id::INT,
    summoner_level = @summoner_level::INT
WHERE puuid = $1;

-- name: UpdateAccount :exec
UPDATE tft_summoner
SET name = @name::VARCHAR,
    tag = @tag::VARCHAR
WHERE puuid = $1;
