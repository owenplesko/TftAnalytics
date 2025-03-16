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
	matches_before_timestamp TIMESTAMP NOT NULL
);

CREATE TABLE tft_comp (
	match_id VARCHAR NOT NULL REFERENCES tft_match,
	summoner_puuid VARCHAR NOT NULL,
	comp_data JSONB NOT NULL,
	match_date TIMESTAMP NOT NULL,
	PRIMARY KEY(match_id, summoner_puuid)
);

CREATE INDEX idx_name_tag_insensitive ON tft_summoner (REPLACE(LOWER(name), ' ', ''), REPLACE(LOWER(tag), ' ', ''));

CREATE INDEX idx_region_update ON tft_summoner (region, matches_after_timestamp ASC NULLS FIRST);

CREATE INDEX idx_match_history ON tft_comp (summoner_puuid, match_date DESC);
