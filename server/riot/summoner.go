package riot

import (
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

func (c *Riot) GetSummonerByPuuid(region string, puuid string) (*RiotSummonerRes, error) {
	summonerRes := new(RiotSummonerRes)
	route := fmt.Sprintf("tft/summoner/v1/summoners/by-puuid/%v", puuid)
	err := c.request(region, route, summonerRes)
	if err != nil {
		return nil, err
	}

	return summonerRes, err
}
