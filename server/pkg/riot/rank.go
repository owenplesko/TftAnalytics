package riot

import (
	"context"
	"fmt"
	"strings"
)

type RankPage struct {
	Tier      string          `json:"tier"`
	LeagueId  string          `json:"leagueId"`
	QueueType string          `json:"queue"`
	Name      string          `json:"name"`
	Entries   []ApexRankEntry `json:"entries"`
}

type ApexRankEntry struct {
	SummonerID   string `json:"summonerId"`
	Puuid        string `json:"puuid"`
	LeaguePoints int    `json:"leaguePoints"`
	Rank         string `json:"rank"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Veteran      bool   `json:"veteran"`
	Inactive     bool   `json:"inactive"`
	FreshBlood   bool   `json:"freshBlood"`
	HotStreak    bool   `json:"hotStreak"`
}

type RankEntry struct {
	Puuid        string `json:"puuid"`
	LeagueId     string `json:"leagueId"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	SummonerId   string `json:"summonerId"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Veteran      bool   `json:"veteran"`
	Inactive     bool   `json:"inactive"`
	FreshBlood   bool   `json:"freshBlood"`
	HotStreak    bool   `json:"hotStreak"`
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

func (c *Riot) GetRank(ctx context.Context, region string, puuid string) (RankEntry, error) {
	var rankRes []RankEntry
	route := fmt.Sprintf("tft/league/v1/by-puuid/%v", puuid)
	err := c.request(ctx, region, route, &rankRes)
	if err != nil {
		return RankEntry{}, err
	}

	// TODO: return all rank entries
	for _, rank := range rankRes {
		if rank.QueueType == "RANKED_TFT" {
			return rank, nil
		}
	}

	return RankEntry{}, ErrNotFound
}

func (c *Riot) GetRankEntries(ctx context.Context, region string, tier string, division string, page int) ([]RankEntry, error) {
	var rankPageRes []RankEntry

	route := fmt.Sprintf("tft/league/v1/entries/%v/%v?queue=RANKED_TFT&page=%v", tier, division, page)
	err := c.request(ctx, region, route, &rankPageRes)

	return rankPageRes, err
}

func (c *Riot) GetApexRankPage(ctx context.Context, region string, tier string) (RankPage, error) {
	var rankPageRes RankPage

	route := fmt.Sprintf("tft/league/v1/%v?queue=RANKED_TFT", strings.ToLower(tier))
	err := c.request(ctx, region, route, &rankPageRes)

	return rankPageRes, err
}
