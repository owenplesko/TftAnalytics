-- name: GetMatchComps :many
SELECT 
	sqlc.embed(tft_summoner), comp_data
FROM
	tft_comp
	JOIN tft_summoner ON summoner_puuid = puuid
WHERE match_id = $1;

-- name: CreateMatch :exec
INSERT INTO tft_match (
    id,
	data_version,
	game_version,
	queue_id,
	game_type,
	set_name,
	set_number,
	match_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);

-- name: CreateComp :exec
INSERT INTO tft_comp (
    match_id,
	summoner_puuid,
	comp_data,
	match_date,
	rank
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: MatchExists :one
SELECT EXISTS (
    SELECT * FROM tft_match WHERE id = $1
);

-- name: GetSummonerMatchHistory :many
SELECT 
	comp_data, sqlc.embed(tft_match)
FROM
	tft_comp JOIN tft_match 
	ON tft_comp.match_id = tft_match.id
WHERE
	summoner_puuid = $1 AND
	tft_comp.match_date < @before::TIMESTAMP AND
	(sqlc.narg(queue_id)::INT IS NULL OR tft_match.queue_id = sqlc.narg(queue_id)::INT) AND
	set_number = $2
ORDER BY
	tft_comp.match_date DESC
LIMIT $3;
