package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/internal/db"
	"github.com/owenplesko/TftAnalytics/pkg/riot"
)

func (service *Service) CollectSummonerByPuuid(ctx context.Context, region, puuid string) error {
	key := fmt.Sprintf("SUMMONER_BY_PUUID_%s_%s", region, puuid)

	_, err, _ := service.group.Do(key, func() (interface{}, error) {
		account, err := service.riot.GetAccountByPuuid(ctx, riot.RegionToCluster[region], puuid)
		if err != nil {
			return nil, fmt.Errorf("Riot.GetAccountByPuuid failed with err: %w", err)
		}

		summoner, err := service.riot.GetSummonerByPuuid(ctx, region, puuid)
		if err != nil {
			return nil, fmt.Errorf("Riot.GetSummonerByPuuid failed with err: %w", err)
		}

		err = service.queries.UpsertSummoner(ctx, db.UpsertSummonerParams{
			Puuid:         summoner.Puuid,
			Region:        region,
			Name:          account.Name,
			Tag:           account.Tag,
			SummonerID:    summoner.SummonerId,
			ProfileIconID: summoner.ProfileIconId,
			SummonerLevel: summoner.SummonerLevel,
		})
		if err != nil {
			return nil, fmt.Errorf("Queries.UpsertSummoner failed with err: %w", err)
		}

		log.Printf("collected summoner %v#%v on region %v\n", account.Name, account.Tag, region)

		return nil, nil
	})

	return err
}
