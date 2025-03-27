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

CREATE TABLE tft_summoner_stats (
	summoner_puuid VARCHAR NOT NULL,
	queue_id INT,
	set_number INT NOT NULL,
	comp_count INT NOT NULL,
	top_4_count INT NOT NULL,
	top_1_count INT NOT NULL,
	placement_sum INT NOT NULL,
	avg_placement FLOAT NOT NULL GENERATED ALWAYS AS (placement_sum::FLOAT / comp_count) STORED,
	top_4_rate FLOAT NOT NULL GENERATED ALWAYS AS (top_4_count::FLOAT / comp_count * 100) STORED,
	top_1_rate FLOAT NOT NULL GENERATED ALWAYS AS (top_1_count::FLOAT / comp_count * 100) STORED,
	UNIQUE(summoner_puuid, set_number, queue_id)
);

CREATE OR REPLACE PROCEDURE update_tft_summoner_stats(IN in_summoner_puuid VARCHAR)
LANGUAGE plpgsql
AS $$
DECLARE
    update_old TIMESTAMP;
    update_new TIMESTAMP := CURRENT_TIMESTAMP;
BEGIN
    -- Get the last update timestamp for the summoner
    SELECT update_timestamp INTO update_old
    FROM tft_summoner
    WHERE puuid = in_summoner_puuid;

    -- Calculate stats grouped by queue_id
    WITH queue_id_stats AS (
        SELECT
            queue_id,
            set_number,
            COUNT(*) AS comp_count,
            COUNT(*) FILTER (WHERE (comp_data ->> 'placement')::INT <= 4) AS top_4_count,
            COUNT(*) FILTER (WHERE (comp_data ->> 'placement')::INT = 1) AS top_1_count,
            SUM((comp_data ->> 'placement')::INT) AS placement_sum
        FROM tft_comp
        JOIN tft_match ON tft_comp.match_id = tft_match.id
        WHERE tft_comp.summoner_puuid = in_summoner_puuid
            AND (update_old IS NULL OR tft_comp.insert_timestamp >= update_old)
            AND tft_comp.insert_timestamp < update_new
        GROUP BY set_number, queue_id
    ),
    -- Calculate stats
    stats AS (
        SELECT
            set_number,
            SUM(comp_count) AS comp_count,
            SUM(top_4_count) AS top_4_count,
            SUM(top_1_count) AS top_1_count,
            SUM(placement_sum) AS placement_sum
        FROM queue_id_stats
        GROUP BY set_number
    ),
	-- Combine both stats
	combined_stats AS (
    	SELECT queue_id, set_number, comp_count, top_4_count, top_1_count, placement_sum FROM queue_id_stats
    	UNION ALL
    	SELECT NULL, set_number, comp_count, top_4_count, top_1_count, placement_sum FROM stats
	)

    -- Insert new stats or update existing ones
    INSERT INTO tft_summoner_stats (summoner_puuid, queue_id, set_number, comp_count, top_4_count, top_1_count, placement_sum)
    SELECT 
        in_summoner_puuid, queue_id, set_number, comp_count, top_4_count, top_1_count, placement_sum
    FROM combined_stats
    ON CONFLICT (summoner_puuid, queue_id, set_number) DO UPDATE
    SET
        comp_count = tft_summoner_stats.comp_count + EXCLUDED.comp_count,
        top_4_count = tft_summoner_stats.top_4_count + EXCLUDED.top_4_count,
        top_1_count = tft_summoner_stats.top_1_count + EXCLUDED.top_1_count,
        placement_sum = tft_summoner_stats.placement_sum + EXCLUDED.placement_sum;

    -- Update summoner's last update timestamp
    UPDATE tft_summoner
    SET update_timestamp = update_new
    WHERE puuid = in_summoner_puuid;
END;
$$;
