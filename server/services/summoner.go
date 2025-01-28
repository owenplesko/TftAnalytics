package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
	"errors"
	"log"
)

func (env ServiceEnv) CollectSummonerDetails(ctx context.Context, region, puuid string) error {
	res, err := riot.GetSummonerByPuuid(region, puuid)
	if err != nil {
		return err
	}

	err = env.Queries.UpdateSummoner(ctx, db.UpdateSummonerParams{
		Puuid:         puuid,
		SummonerID:    res.SummonerId,
		ProfileIconID: res.ProfileIconId,
		SummonerLevel: res.SummonerLevel,
	})

	log.Printf("Summoner details collected for %v\n", puuid)

	return err
}

func (env ServiceEnv) CollectAccountByPuuid(ctx context.Context, cluster, puuid string) error {
	res, err := riot.GetAccountByPuuid(cluster, puuid)
	if errors.Is(err, riot.NotFoundError) {
		err = env.Queries.AddSummonerFlag(ctx, db.AddSummonerFlagParams{
			Puuid: puuid,
			Flag:  "SKIP_ACCOUNT_DATA",
		})
		return err
	}
	if err != nil {
		return err
	}

	err = env.Queries.UpdateAccount(ctx, db.UpdateAccountParams{
		Puuid: puuid,
		Name:  res.Name,
		Tag:   res.Tag,
	})

	log.Printf("Account details collected for %v\n", puuid)

	return err
}

func (env ServiceEnv) CollectAccountByNameTag(ctx context.Context, cluster, name, tag string) error {
	res, err := riot.GetAccountByName(cluster, name, tag)
	if err != nil {
		return err
	}

	err = env.Queries.UpsertAccount(ctx, db.UpsertAccountParams{
		Puuid: res.Puuid,
		Name:  res.Name,
		Tag:   res.Tag,
	})

	log.Printf("Account collected with name %v#%v\n", name, tag)

	return err
}

func (env ServiceEnv) GetSummonerByPuuid(ctx context.Context, puuid string) (db.GetSummonerByPuuidRow, error) {
	return env.Queries.GetSummonerByPuuid(ctx, puuid)
}

func (env ServiceEnv) GetOrCollectSummonerByNameTag(ctx context.Context, cluster, name, tag string) (db.GetSummonerByNameTagRow, error) {
	exists, err := env.Queries.SummonerExistsByNameTag(ctx, db.SummonerExistsByNameTagParams{
		Name: name,
		Tag:  tag,
	})
	if err != nil {
		return db.GetSummonerByNameTagRow{}, err
	}

	if !exists {
		err = env.CollectAccountByNameTag(ctx, cluster, name, tag)
		if err != nil {
			return db.GetSummonerByNameTagRow{}, err
		}
	}

	return env.Queries.GetSummonerByNameTag(ctx, db.GetSummonerByNameTagParams{Name: name, Tag: tag})
}

func (env ServiceEnv) CollectSummonerRegion(ctx context.Context, puuid string) error {
	var regionMatch string

	for region, _ := range riot.RegionToCluster {
		_, err := riot.GetSummonerByPuuid(region, puuid)

		if errors.Is(err, riot.NotFoundError) {
			continue
		}

		if err != nil {
			return err
		}

		regionMatch = region
		break
	}

	if regionMatch == "" {
		env.Queries.AddSummonerFlag(ctx, db.AddSummonerFlagParams{
			Puuid: puuid,
			Flag:  "SKIP_REGION_MATCH",
		})
		return errors.New("no region match found")
	}

	err := env.Queries.UpdateRegion(ctx, db.UpdateRegionParams{
		Puuid:  puuid,
		Region: regionMatch,
	})

	return err
}
