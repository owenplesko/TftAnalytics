package services

import (
	"TFTAnalyticsServer/internal/db"
	"TFTAnalyticsServer/pkg/riot"
	"context"
	"fmt"
	"log"
	"sync"
)

func (service *Service) CollectSummonerByNameTag(ctx context.Context, cluster, name, tag string) error {
	account, err := service.riot.GetAccountByName(ctx, cluster, name, tag)
	if err != nil {
		return fmt.Errorf("Riot.GetAccountByName failed with err: %w", err)
	}

	summoner, region, err := service.findSummonerAndRegion(ctx, account.Puuid)
	if err != nil {
		return fmt.Errorf("findSummonerRegion failed err: %w", err)
	}

	err = service.queries.UpsertSummoner(ctx, db.UpsertSummonerParams{
		Puuid:         account.Puuid,
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

type regionSuccessRes struct {
	summoner *riot.RiotSummonerRes
	region   string
}

// TODO: explore passing riot errors on no region found
// should implement request retrying for riot requests first tho..
func (service *Service) findSummonerAndRegion(ctx context.Context, puuid string) (*riot.RiotSummonerRes, string, error) {
	wg := sync.WaitGroup{}
	successChan := make(chan regionSuccessRes)
	for region := range riot.RegionToCluster {
		wg.Add(1)

		go func(region string) {
			defer wg.Done()
			if summoner, err := service.riot.GetSummonerByPuuid(ctx, region, puuid); err == nil {
				successChan <- regionSuccessRes{summoner: summoner, region: region}
			}
		}(region)
	}
	go func() {
		wg.Wait()
		close(successChan)
	}()

	var err error
	success := <-successChan
	if success.region == "" {
		err = fmt.Errorf("failed to find region for summoner with puuid %v", puuid)
	}

	return success.summoner, success.region, err
}
