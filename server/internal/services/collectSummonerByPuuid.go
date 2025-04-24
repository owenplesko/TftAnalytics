package services

import (
	"TFTAnalyticsServer/internal/db"
	"TFTAnalyticsServer/pkg/riot"
	"context"
	"fmt"
	"log"
)

func (service *Service) CollectSummonerByPuuid(ctx context.Context, region, puuid string) error {
	account, err := service.riot.GetAccountByPuuid(ctx, riot.RegionToCluster[region], puuid)
	if err != nil {
		return fmt.Errorf("Riot.GetAccountByPuuid failed with err: %w", err)
	}

	summoner, err := service.riot.GetSummonerByPuuid(ctx, region, puuid)
	if err != nil {
		return fmt.Errorf("Riot.GetSummonerByPuuid failed with err: %w", err)
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
		return fmt.Errorf("Queries.UpsertSummoner failed with err: %w", err)
	}

	log.Printf("collected summoner %v#%v on region %v\n", account.Name, account.Tag, region)

	return nil
}
