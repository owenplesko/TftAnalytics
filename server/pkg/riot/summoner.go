package riot

import (
	"context"
	"fmt"
)

type RiotSummonerRes struct {
	Puuid         string `json:"puuid"`
	SummonerId    string `json:"id"`
	AccountId     string `json:"accountId"`
	ProfileIconId int32  `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int32  `json:"summonerLevel"`
}

func (c *Riot) GetSummonerByPuuid(ctx context.Context, region string, puuid string) (*RiotSummonerRes, error) {
	summonerRes := new(RiotSummonerRes)

	route := fmt.Sprintf("tft/summoner/v1/summoners/by-puuid/%v", puuid)
	err := c.request(ctx, region, route, summonerRes)

	return summonerRes, err
}
