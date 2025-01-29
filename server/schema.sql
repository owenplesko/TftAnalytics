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
	PRIMARY KEY(match_id, summoner_puuid)
);

CREATE INDEX idx_name_tag_insensitive ON tft_summoner (REPLACE(LOWER(name), ' ', ''), REPLACE(LOWER(tag), ' ', ''));
