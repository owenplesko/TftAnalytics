package riot

import (
	"fmt"
)

type Account struct {
	Puuid string `json:"puuid"`
	Name  string `json:"gameName"`
	Tag   string `json:"tagLine"`
}

func (c *Riot) GetAccountByName(cluster string, name string, tag string) (*Account, error) {
	// QUICK FIX: sea and asia are merged for account route
	if cluster == "sea" {
		cluster = "asia"
	}

	accountRes := new(Account)
	route := fmt.Sprintf("riot/account/v1/accounts/by-riot-id/%v/%v", name, tag)
	err := c.request(cluster, route, accountRes)
	if err != nil {
		return nil, err
	}

	return accountRes, err
}

func (c *Riot) GetAccountByPuuid(cluster string, puuid string) (*Account, error) {
	// QUICK FIX: sea and asia are merged for account route
	if cluster == "sea" {
		cluster = "asia"
	}

	accountRes := new(Account)
	route := fmt.Sprintf("riot/account/v1/accounts/by-puuid/%v", puuid)
	err := c.request(cluster, route, accountRes)
	if err != nil {
		return nil, err
	}

	return accountRes, err
}
