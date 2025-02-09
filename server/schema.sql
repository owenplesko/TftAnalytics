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
	region VARCHAR,
	name VARCHAR,
	tag VARCHAR,
	summoner_id VARCHAR,
	profile_icon_id INT,
	summoner_level INT,
	full_update_timestamp TIMESTAMP,
	background_update_timestamp TIMESTAMP,
	flags VARCHAR[] NOT NULL DEFAULT '{}'
);

CREATE TABLE tft_comp (
	match_id VARCHAR NOT NULL REFERENCES tft_match,
	summoner_puuid VARCHAR NOT NULL REFERENCES tft_summoner,
	comp_data JSONB NOT NULL,
	match_date TIMESTAMP NOT NULL,
	PRIMARY KEY(match_id, summoner_puuid)
);

CREATE TABLE tft_rank (
	summoner_id VARCHAR NOT NULL PRIMARY KEY,
	tier VARCHAR NOT NULL,
	rank VARCHAR NOT NULL,
	league_points INT NOT NULL,
	wins INT NOT NULL,
	losses INT NOT NULL
);

CREATE INDEX idx_name_tag_insensitive ON tft_summoner (REPLACE(LOWER(name), ' ', ''), REPLACE(LOWER(tag), ' ', ''));

CREATE INDEX idx_name_tag_null ON tft_summoner (name, tag) WHERE (name IS NULL OR tag IS NULL) AND NOT 'SKIP_ACCOUNT_DATA' = ANY(flags);

CREATE INDEX idx_region_summoner_id_null ON tft_summoner (region, puuid) WHERE summoner_id IS NULL;

CREATE INDEX idx_region_null ON tft_summoner (puuid) WHERE region IS NULL AND NOT 'SKIP_REGION_MATCH' = ANY(flags);

CREATE INDEX idx_region_background_update ON tft_summoner (region, background_update_timestamp ASC NULLS FIRST);

CREATE INDEX idx_match_history ON tft_comp (summoner_puuid, match_date DESC);
