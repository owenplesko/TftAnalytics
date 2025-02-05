package riot

import (
	"fmt"
)

type RankEntry struct {
	Puuid        string `json:"puuid"`
	LeagueId     string `json:"leagueId"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	SummonerId   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	LeaguePoints int32  `json:"leaguePoints"`
	Wins         int32  `json:"wins"`
	Losses       int32  `json:"losses"`
	Veteran      bool   `json:"veteran"`
	Inactive     bool   `json:"inactive"`
	FreshBlood   bool   `json:"freshBlood"`
	HotStreak    bool   `json:"hotStreak"`
}

func GetRank(region string, summonerId string) (*RankEntry, error) {
	var rankRes []RankEntry
	route := fmt.Sprintf("tft/league/v1/entries/by-summoner/%v", summonerId)
	err := getJson(region, route, &rankRes)
	if err != nil {
		return nil, err
	}

	for _, rank := range rankRes {
		if rank.QueueType == "RANKED_TFT" {
			return &rank, nil
		}
	}
	return nil, nil
}
