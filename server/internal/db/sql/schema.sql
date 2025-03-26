CREATE TABLE tft_match (
	id VARCHAR PRIMARY KEY NOT NULL,
	data_version VARCHAR NOT NULL,
	game_version VARCHAR NOT NULL,
	queue_id INT NOT NULL,
	game_type VARCHAR NOT NULL,
	set_name VARCHAR NOT NULL,
	set_number INT NOT NULL,
	match_date TIMESTAMP NOT NULL
);

CREATE TABLE tft_summoner (
	puuid VARCHAR PRIMARY KEY NOT NULL,
	region VARCHAR NOT NULL,
	name VARCHAR NOT NULL,
	tag VARCHAR NOT NULL,
	summoner_id VARCHAR NOT NULL,
	profile_icon_id INT NOT NULL,
	summoner_level INT NOT NULL,
	update_timestamp TIMESTAMP,
	matches_before_timestamp TIMESTAMP
);

CREATE TABLE tft_comp (
	match_id VARCHAR NOT NULL REFERENCES tft_match,
	summoner_puuid VARCHAR NOT NULL,
	comp_data JSONB NOT NULL,
	match_date TIMESTAMP NOT NULL,
	insert_timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(match_id, summoner_puuid)
);

CREATE TABLE tft_summoner_comp_aggregate (
	summoner_puuid VARCHAR NOT NULL,
	queue_id INT NOT NULL,
	comp_count INT NOT NULL,
	top_4_count INT NOT NULL,
	top_1_count INT NOT NULL,
	placement_sum INT NOT NULL,
	avg_placement FLOAT GENERATED ALWAYS AS (placement_sum::FLOAT / comp_count) STORED,
	top_4_rate FLOAT GENERATED ALWAYS AS (top_4_count::FLOAT / comp_count * 100) STORED,
	top_1_rate FLOAT GENERATED ALWAYS AS (top_1_count::FLOAT / comp_count * 100) STORED,
	PRIMARY KEY(summoner_puuid, queue_id)
);

CREATE
OR REPLACE FUNCTION fn_normalize_str (VARCHAR) RETURNS VARCHAR AS $$
BEGIN
	RETURN REPLACE(LOWER($1), ' ', '');
END;
$$ LANGUAGE plpgsql IMMUTABLE;

CREATE INDEX idx_name_tag_normalized ON tft_summoner (fn_normalize_str(name), fn_normalize_str(tag));

CREATE INDEX idx_region_update ON tft_summoner (region, matches_after_timestamp ASC NULLS FIRST);

CREATE INDEX idx_match_history ON tft_comp (summoner_puuid, match_date DESC);

CREATE OR REPLACE PROCEDURE update_tft_summoner_comp_aggregate(IN in_summoner_puuid VARCHAR)
LANGUAGE plpgsql
AS $$
DECLARE
    update_old TIMESTAMP;
    update_new TIMESTAMP := CURRENT_TIMESTAMP;
BEGIN
    -- Fetch the old update timestamp
    SELECT update_timestamp INTO update_old
    FROM tft_summoner
    WHERE puuid = in_summoner_puuid;

    -- Insert or increment aggregated data into tft_summoner_comp_aggregate
    INSERT INTO tft_summoner_comp_aggregate (summoner_puuid, queue_id, comp_count, top_4_count, top_1_count, placement_sum)
    SELECT
        in_summoner_puuid,
        queue_id,
        COUNT(*) AS comp_count,
        COUNT(*) FILTER (WHERE (comp_data ->> 'placement')::INT <= 4) AS top_4_count,
        COUNT(*) FILTER (WHERE (comp_data ->> 'placement')::INT = 1) AS top_1_count,
        SUM((comp_data ->> 'placement')::INT) AS placement_sum
    FROM tft_comp
    JOIN tft_match ON tft_comp.match_id = tft_match.id
    WHERE tft_comp.summoner_puuid = in_summoner_puuid
        AND (
            (update_old IS NULL AND tft_comp.insert_timestamp <= update_new)
            OR (tft_comp.insert_timestamp BETWEEN update_old AND update_new)
        )
    GROUP BY queue_id
    ON CONFLICT (summoner_puuid, queue_id) DO UPDATE
    SET
        comp_count = tft_summoner_comp_aggregate.comp_count + EXCLUDED.comp_count,
        top_4_count = tft_summoner_comp_aggregate.top_4_count + EXCLUDED.top_4_count,
        top_1_count = tft_summoner_comp_aggregate.top_1_count + EXCLUDED.top_1_count,
        placement_sum = tft_summoner_comp_aggregate.placement_sum + EXCLUDED.placement_sum;

    -- Update the summoner's update timestamp
    UPDATE tft_summoner
    SET update_timestamp = update_new
    WHERE puuid = in_summoner_puuid;
END;
$$;