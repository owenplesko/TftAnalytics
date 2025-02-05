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

func GetRank(region string, summonerId string) (RankEntry, error) {
	var rankRes []RankEntry
	route := fmt.Sprintf("tft/league/v1/entries/by-summoner/%v", summonerId)
	err := getJson(region, route, &rankRes)
	if err != nil {
		return RankEntry{}, err
	}

	for _, rank := range rankRes {
		if rank.QueueType == "RANKED_TFT" {
			return rank, nil
		}
	}

	return RankEntry{}, NotFoundError
}

var ApexTiers = []string{
	"CHALLENGER",
	"GRANDMASTER",
	"MASTER",
}

var Tiers = []string{
	"DIAMOND",
	"EMERALD",
	"PLATINUM",
	"GOLD",
	"SILVER",
	"BRONZE",
	"IRON",
}

var Divisions = []string{
	"I",
	"II",
	"III",
	"IV",
}

func GetRankPage(region string, tier string, division string) ([]RankEntry, error) {
	var rankPageRes []RankEntry
	route := fmt.Sprintf("tft/league/v1/entries/%v/%v?queue=RANKED_TFT", tier)
	err := getJson(region, route, &rankPageRes)

	return rankPageRes, err
}

func GetApexRankPage(region string, tier string) ([]RankEntry, error) {
	var rankPageRes []RankEntry
	route := fmt.Sprintf("tft/league/v1/%v?queue=RANKED_TFT", tier)
	err := getJson(region, route, &rankPageRes)

	return rankPageRes, err
}
