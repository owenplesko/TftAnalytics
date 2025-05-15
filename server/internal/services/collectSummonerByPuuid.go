package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/internal/db"
	"github.com/owenplesko/TftAnalytics/pkg/dedupe"
	"github.com/owenplesko/TftAnalytics/pkg/riot"
)

func (service *Service) CollectSummonerByPuuid(ctx context.Context, region, puuid string) error {
	return dedupe.Run(service.deduplicator, SummonerByPuuidTask{
		service: service,
		ctx:     ctx,
		region:  region,
		puuid:   puuid,
	}).Await()
}

type SummonerByPuuidTask struct {
	service *Service
	ctx     context.Context
	region  string
	puuid   string
}

func (task SummonerByPuuidTask) ID() string {
	return fmt.Sprintf("SUMMONER_BY_PUUID_%s_%s", task.region, task.puuid)
}

func (task SummonerByPuuidTask) Run() error {
	account, err := task.service.riot.GetAccountByPuuid(task.ctx, riot.RegionToCluster[task.region], task.puuid)
	if err != nil {
		return fmt.Errorf("Riot.GetAccountByPuuid failed with err: %w", err)
	}

	summoner, err := task.service.riot.GetSummonerByPuuid(task.ctx, task.region, task.puuid)
	if err != nil {
		return fmt.Errorf("Riot.GetSummonerByPuuid failed with err: %w", err)
	}

	err = task.service.queries.UpsertSummoner(task.ctx, db.UpsertSummonerParams{
		Puuid:         summoner.Puuid,
		Region:        task.region,
		Name:          account.Name,
		Tag:           account.Tag,
		SummonerID:    summoner.SummonerId,
		ProfileIconID: summoner.ProfileIconId,
		SummonerLevel: summoner.SummonerLevel,
	})
	if err != nil {
		return fmt.Errorf("Queries.UpsertSummoner failed with err: %w", err)
	}

	log.Printf("collected summoner %v#%v on region %v\n", account.Name, account.Tag, task.region)

	return nil

}
