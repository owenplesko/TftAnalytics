package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/owenplesko/TftAnalytics/pkg/dedupe"
	"github.com/owenplesko/TftAnalytics/pkg/riot"
)

// TODO revamp this to have states and better error handling and stuff
func (service *Service) UpdateSummoner(ctx context.Context, puuid string) error {
	return dedupe.Run(service.deduplicator, updateSummonerTask{
		service: service,
		ctx:     ctx,
		puuid:   puuid,
	}).Await()
}

type updateSummonerTask struct {
	service *Service
	ctx     context.Context
	puuid   string
}

func (task updateSummonerTask) ID() string {
	return fmt.Sprintf("UPDATE_SUMMONER_%s", task.puuid)
}

func (task updateSummonerTask) Run() error {
	summoner, err := task.service.GetSummonerByPuuid(task.ctx, task.puuid)
	if err != nil {
		return fmt.Errorf("GetSummonerByPuuid failed with err: %w", err)
	}

	err = task.service.CollectSummonerRank(task.ctx, summoner.Region, summoner.SummonerID)
	// riot.ErrNotFound is an expected error and should not cause UpdateSummonerInfo to fail
	if err != nil && !errors.Is(err, riot.ErrNotFound) {
		return fmt.Errorf("CollectSummonerRank failed with err: %w", err)
	}
	err = task.service.CollectSummonerByPuuid(task.ctx, summoner.Region, task.puuid)
	if err != nil {
		return fmt.Errorf("CollectSummonerByPuuid failed with err: %w", err)
	}
	err = task.service.CollectMatchHistory(task.ctx, summoner.Region, task.puuid, summoner.MatchesBeforeTimestamp.Time)
	if err != nil {
		return fmt.Errorf("CollectMatchHistory failed with err: %w", err)
	}

	err = task.service.queries.UpdateSummonerStats(task.ctx, task.puuid)
	if err != nil {
		return fmt.Errorf("Queries.UpdateSummonerStats failed with err: %w", err)
	}

	return nil

}
