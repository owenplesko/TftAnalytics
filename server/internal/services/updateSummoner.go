package services

import (
	"TFTAnalyticsServer/pkg/riot"
	"context"
	"errors"
	"fmt"
)

// TODO revamp this to have states and better error handling and stuff
func (service *Service) UpdateSummoner(ctx context.Context, puuid string) error {
	summoner, err := service.GetSummonerByPuuid(ctx, puuid)
	if err != nil {
		return fmt.Errorf("GetSummonerByPuuid failed with err: %w", err)
	}

	err = service.CollectSummonerRank(ctx, summoner.Region, summoner.SummonerID)
	// riot.ErrNotFound is an expected error and should not cause UpdateSummonerInfo to fail
	if err != nil && !errors.Is(err, riot.ErrNotFound) {
		return fmt.Errorf("CollectSummonerRank failed with err: %w", err)
	}
	err = service.CollectSummonerByPuuid(ctx, summoner.Region, puuid)
	if err != nil {
		return fmt.Errorf("CollectSummonerByPuuid failed with err: %w", err)
	}
	err = service.CollectMatchHistory(ctx, summoner.Region, puuid, summoner.MatchesBeforeTimestamp.Time)
	if err != nil {
		return fmt.Errorf("CollectMatchHistory failed with err: %w", err)
	}

	err = service.queries.UpdateSummonerStats(ctx, puuid)
	if err != nil {
		return fmt.Errorf("Queries.UpdateSummonerStats failed with err: %w", err)
	}

	return nil
}
