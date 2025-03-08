package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (service *Service) GetSummonerByPuuid(ctx context.Context, puuid string) (db.TftSummoner, error) {
	summoner, err := service.Queries.GetSummonerByPuuid(ctx, puuid)
	if err != nil {
		return summoner, fmt.Errorf("Queries.GetSummonerByPuuid failed with err: %w", err)
	}

	return summoner, nil
}

func (service *Service) GetOrCollectSummonerByNameTag(ctx context.Context, cluster, name, tag string) (db.TftSummoner, error) {
	summoner, err := service.Queries.GetSummonerByNameTag(ctx, db.GetSummonerByNameTagParams{
		Name: name,
		Tag:  tag,
	})
	if err == nil {
		// summoner found! no need to collect
		return summoner, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return summoner, fmt.Errorf("Queries.GetSummonerByNameTag failed with err: %w", err)
	}

	err = service.CollectSummonerByNameTag(ctx, cluster, name, tag)
	if err != nil {
		return summoner, fmt.Errorf("CollectSummonerByNameTag failed with err: %w", err)
	}

	summoner, err = service.Queries.GetSummonerByNameTag(ctx, db.GetSummonerByNameTagParams{Name: name, Tag: tag})
	if err != nil {
		return summoner, fmt.Errorf("Queries.GetSummonerByNameTag failed with err: %w", err)
	}

	return summoner, nil
}

func (service *Service) CollectSummonerByPuuid(ctx context.Context, region, puuid string) error {
	account, err := service.Riot.GetAccountByPuuid(riot.RegionToCluster[region], puuid)
	if err != nil {
		return fmt.Errorf("Riot.GetAccountByPuuid failed with err: %w", err)
	}

	summoner, err := service.Riot.GetSummonerByPuuid(region, puuid)
	if err != nil {
		return fmt.Errorf("Riot.GetSummonerByPuuid failed with err: %w", err)
	}

	err = service.Queries.UpsertSummoner(ctx, db.UpsertSummonerParams{
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

func (service *Service) CollectSummonerByNameTag(ctx context.Context, cluster, name, tag string) error {
	account, err := service.Riot.GetAccountByName(cluster, name, tag)
	if err != nil {
		return fmt.Errorf("Riot.GetAccountByName failed with err: %w", err)
	}

	summoner, region, err := service.findSummonerAndRegion(account.Puuid)
	if err != nil {
		return fmt.Errorf("findSummonerRegion failed err: %w", err)
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
func (service *Service) findSummonerAndRegion(puuid string) (*riot.RiotSummonerRes, string, error) {
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
		err = fmt.Errorf("failed to find region for summoner with puuid %v", puuid)
	}

	return success.summoner, success.region, err
}

// TODO revamp this to have states and better error handling and stuff
func (service *Service) UpdateSummonerInfo(ctx context.Context, puuid string) error {
	summoner, err := service.GetSummonerByPuuid(ctx, puuid)
	if err != nil {
		return fmt.Errorf("GetSummonerByPuuid failed with err: %w", err)
	}

	var after time.Time
	if summoner.MatchesAfterTimestamp.Valid {
		after = summoner.MatchesAfterTimestamp.Time
	} else {
		after = time.Now().UTC().Add(-time.Hour * 24 * 30)
	}

	err = service.CollectSummonerRank(ctx, summoner.Region, summoner.SummonerID)
	if err != nil {
		return fmt.Errorf("CollectSummonerRank failed with err: %w", err)
	}
	err = service.CollectSummonerByPuuid(ctx, summoner.Region, puuid)
	if err != nil {
		return fmt.Errorf("CollectSummonerByPuuid failed with err: %w", err)
	}
	err = service.CollectMatchHistory(ctx, summoner.Region, puuid, after)
	if err != nil {
		return fmt.Errorf("CollectMatchHistory failed with err: %w", err)
	}

	err = service.Queries.SetUpdateTimestamp(ctx, db.SetUpdateTimestampParams{
		Puuid: puuid,
		UpdateTimestamp: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("Queries.SetUpdateTimestamp failed with err: %w", err)
	}

	return nil
}
