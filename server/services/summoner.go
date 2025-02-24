package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
	"errors"
	"log"
	"sync"
)

func (service Service) CollectSummonerDetails(ctx context.Context, region, puuid string) error {
	res, err := service.Riot.GetSummonerByPuuid(region, puuid)
	if err != nil {
		return err
	}

	err = service.Queries.UpdateSummoner(ctx, db.UpdateSummonerParams{
		Puuid:         puuid,
		SummonerID:    res.SummonerId,
		ProfileIconID: res.ProfileIconId,
		SummonerLevel: res.SummonerLevel,
	})

	log.Printf("Summoner details collected for %v\n", puuid)

	return err
}

func (service Service) CollectAccountByPuuid(ctx context.Context, cluster, puuid string) error {
	res, err := service.Riot.GetAccountByPuuid(cluster, puuid)
	if errors.Is(err, riot.NotFoundError) {
		err = service.Queries.AddSummonerFlag(ctx, db.AddSummonerFlagParams{
			Puuid: puuid,
			Flag:  "SKIP_ACCOUNT_DATA",
		})
		return err
	}
	if err != nil {
		return err
	}

	err = service.Queries.UpdateAccount(ctx, db.UpdateAccountParams{
		Puuid: puuid,
		Name:  res.Name,
		Tag:   res.Tag,
	})

	log.Printf("Account details collected for %v\n", puuid)

	return err
}

func (service Service) CollectAccountByNameTag(ctx context.Context, cluster, name, tag string) error {
	res, err := service.Riot.GetAccountByName(cluster, name, tag)
	if err != nil {
		return err
	}

	err = service.Queries.UpsertAccount(ctx, db.UpsertAccountParams{
		Puuid: res.Puuid,
		Name:  res.Name,
		Tag:   res.Tag,
	})

	log.Printf("Account collected with name %v#%v\n", name, tag)

	return err
}

func (env Service) GetSummonerByPuuid(ctx context.Context, puuid string) (db.TftSummoner, error) {
	return env.Queries.GetSummonerByPuuid(ctx, puuid)
}

func (env Service) GetOrCollectSummonerByNameTag(ctx context.Context, cluster, name, tag string) (db.TftSummoner, error) {
	exists, err := env.Queries.SummonerExistsByNameTag(ctx, db.SummonerExistsByNameTagParams{
		Name: name,
		Tag:  tag,
	})
	if err != nil {
		return db.TftSummoner{}, err
	}

	if !exists {
		err = env.CollectAccountByNameTag(ctx, cluster, name, tag)
		if err != nil {
			return db.TftSummoner{}, err
		}
	}

	return env.Queries.GetSummonerByNameTag(ctx, db.GetSummonerByNameTagParams{Name: name, Tag: tag})
}

func (service Service) CollectSummonerRegion(ctx context.Context, puuid string) error {
	wg := sync.WaitGroup{}
	successChan := make(chan string)
	for region, _ := range riot.RegionToCluster {
		wg.Add(1)

		go func(region string) {
			defer wg.Done()
			if _, err := service.Riot.GetSummonerByPuuid(region, puuid); err == nil {
				successChan <- region
			}
		}(region)
	}
	go func() {
		wg.Wait()
		close(successChan)
	}()

	region := <-successChan

	if region == "" {
		service.Queries.AddSummonerFlag(ctx, db.AddSummonerFlagParams{
			Puuid: puuid,
			Flag:  "SKIP_REGION_MATCH",
		})
		return errors.New("no region match found")
	}

	err := service.Queries.UpdateRegion(ctx, db.UpdateRegionParams{
		Puuid:  puuid,
		Region: region,
	})

	return err
}
