package riot

import (
	"fmt"
	"strings"
	"time"
)

type Match struct {
	MetaData struct {
		DataVersion  string   `json:"data_version"`
		MatchId      string   `json:"match_id"`
		Participants []string `json:"participants"`
	} `json:"metadata"`
	Info struct {
		Date        int64   `json:"game_datetime"`
		Length      float64 `json:"game_length"`
		GameVersion string  `json:"game_version"`
		QueueId     int32   `json:"queue_id"`
		GameType    string  `json:"tft_game_type"`
		SetName     string  `json:"tft_set_core_name"`
		SetNumber   int32   `json:"tft_set_number"`
		Comps       []Comp  `json:"participants"`
	} `json:"info"`
}

type Comp struct {
	Companion struct {
		ContentId string `json:"content_ID"`
		ItemId    int32  `json:"item_ID"`
		SkinId    int32  `json:"skin_ID"`
		Species   string `json:"species"`
	} `json:"companion"`
	RemainingGold     int32   `json:"gold_left"`
	LastRound         int32   `json:"last_round"`
	Level             int32   `json:"level"`
	Placement         int32   `json:"placement"`
	PlayersEliminated int32   `json:"players_eliminated"`
	Puuid             string  `json:"puuid"`
	TimeEliminated    float64 `json:"time_eliminated"`
	DamageToPlayers   int32   `json:"total_damage_to_players"`
	Traits            []struct {
		Name       string `json:"name"`
		NumUnits   int32  `json:"num_units"`
		Style      int32  `json:"style"`
		TierActive int32  `json:"tier_current"`
		TierMax    int32  `json:"tier_total"`
	} `json:"traits"`
	Units []struct {
		CharacterId string   `json:"character_id"`
		ItemNames   []string `json:"itemNames"`
		Rarity      int32    `json:"rarity"`
		Tier        int32    `json:"tier"`
	} `json:"units"`
}

func getTftMatchCluster(region string) string {
	// return SEA cluster if applicable
	if cluster, ok := seaCluster[region]; ok {
		return cluster
	}

	return RegionToCluster[region]
}

func (c *Riot) GetMatchRegion(matchId string) string {
	return strings.Split(matchId, "_")[0]
}

func (c *Riot) GetMatchDetails(region, matchId string) (*Match, error) {
	matchRes := new(Match)
	route := fmt.Sprintf("tft/match/v1/matches/%v", matchId)
	err := c.request(getTftMatchCluster(region), RegionToCluster[region], route, matchRes)
	if err != nil {
		return nil, err
	}

	return matchRes, err
}

func (c *Riot) GetMatchHistory(region string, puuid string, matchesAfter time.Time) ([]string, error) {
	var history []string
	count := 200

	route := fmt.Sprintf("tft/match/v1/matches/by-puuid/%v/ids?count=%v&startTime=%v", puuid, count, matchesAfter.Unix())
	err := c.request(getTftMatchCluster(region), RegionToCluster[region], route, &history)
	return history, err
}
