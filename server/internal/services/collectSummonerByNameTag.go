package services

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/owenplesko/TftAnalytics/internal/db"
	"github.com/owenplesko/TftAnalytics/pkg/riot"
)

func (service *Service) CollectSummonerByNameTag(ctx context.Context, cluster, name, tag string) error {
	key := fmt.Sprintf("SUMMONER_BY_NAME_TAG_%s_%s", name, tag)

	_, err, _ := service.group.Do(key, func() (interface{}, error) {
		account, err := service.riot.GetAccountByName(ctx, cluster, name, tag)
		if err != nil {
			return nil, fmt.Errorf("Riot.GetAccountByName failed with err: %w", err)
		}

		summoner, region, err := service.findSummonerAndRegion(ctx, account.Puuid)
		if err != nil {
			return nil, fmt.Errorf("findSummonerRegion failed err: %w", err)
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
			return nil, fmt.Errorf("Queries.UpsertSummoner failed with err: %w", err)
		}

		log.Printf("collected summoner %v#%v on region %v\n", account.Name, account.Tag, region)

		return nil, nil
	})

	return err
}

type regionSuccessRes struct {
	summoner *riot.RiotSummonerRes
	region   string
}

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
