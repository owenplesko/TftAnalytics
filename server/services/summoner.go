package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (service Service) GetSummonerByPuuid(ctx context.Context, puuid string) (db.TftSummoner, error) {
	return service.Queries.GetSummonerByPuuid(ctx, puuid)
}

func (service Service) GetOrCollectSummonerByNameTag(ctx context.Context, cluster, name, tag string) (db.TftSummoner, error) {
	summoner, err := service.Queries.GetSummonerByNameTag(ctx, db.GetSummonerByNameTagParams{
		Name: name,
		Tag:  tag,
	})
	if err != pgx.ErrNoRows {
		return summoner, err
	}

	err = service.CollectSummonerByNameTag(ctx, cluster, name, tag)
	if err != nil {
		return summoner, err
	}

	summoner, err = service.Queries.GetSummonerByNameTag(ctx, db.GetSummonerByNameTagParams{Name: name, Tag: tag})

	return summoner, err
}

func (service Service) CollectSummonerByPuuid(ctx context.Context, region, puuid string) error {
	accountDetails, err := service.Riot.GetAccountByPuuid(riot.RegionToCluster[region], puuid)
	if err != nil {
		return err
	}

	summonerDetails, err := service.Riot.GetSummonerByPuuid(region, puuid)
	if err != nil {
		return err
	}

	err = service.Queries.UpsertSummoner(ctx, db.UpsertSummonerParams{
		Puuid:         summonerDetails.Puuid,
		Region:        region,
		Name:          accountDetails.Name,
		Tag:           accountDetails.Tag,
		SummonerID:    summonerDetails.SummonerId,
		ProfileIconID: summonerDetails.ProfileIconId,
		SummonerLevel: summonerDetails.SummonerLevel,
	})
	if err != nil {
		return err
	}

	log.Printf("Summoner details collected for %v\n", puuid)

	return nil
}

func (service Service) CollectSummonerByNameTag(ctx context.Context, cluster, name, tag string) error {
	account, err := service.Riot.GetAccountByName(cluster, name, tag)
	if err != nil {
		return err
	}

	summoner, region, err := service.findSummonerRegion(account.Puuid)
	if err != nil {
		return err
	}

	err = service.Queries.UpsertSummoner(ctx, db.UpsertSummonerParams{
		Puuid:         account.Puuid,
		Region:        region,
		Name:          account.Name,
		Tag:           account.Tag,
		SummonerID:    summoner.SummonerId,
		ProfileIconID: summoner.ProfileIconId,
		SummonerLevel: summoner.SummonerLevel,
	})

	log.Printf("Account collected with name %v#%v\n", name, tag)

	return err
}

type regionSuccessRes struct {
	summoner *riot.RiotSummonerRes
	region   string
}

func (service Service) findSummonerRegion(puuid string) (*riot.RiotSummonerRes, string, error) {
	wg := sync.WaitGroup{}
	successChan := make(chan regionSuccessRes)
	for region := range riot.RegionToCluster {
		wg.Add(1)

		go func(region string) {
			defer wg.Done()
			if summoner, err := service.Riot.GetSummonerByPuuid(region, puuid); err == nil {
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
		err = fmt.Errorf("no region match found")
	}

	return success.summoner, success.region, err
}

func (service Service) UpdateSummonerInfo(ctx context.Context, puuid string) error {
	summoner, err := service.GetSummonerByPuuid(ctx, puuid)
	if err != nil {
		return err
	}

	go service.CollectSummonerRank(ctx, summoner.Region, summoner.SummonerID)
	go service.CollectSummonerByPuuid(ctx, summoner.Region, puuid)
	go service.CollectMatchHistory(ctx, riot.RegionToCluster[summoner.Region], puuid, summoner.MatchesAfterTimestamp.Time)

	service.Queries.SetUpdateTimestamp(ctx, db.SetUpdateTimestampParams{
		Puuid: puuid,
		UpdateTimestamp: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	})

	return nil
}
