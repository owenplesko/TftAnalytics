package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/owenplesko/TftAnalytics/pkg/riot"
)

// TODO revamp this to have states and better error handling and stuff
func (service *Service) UpdateSummoner(ctx context.Context, puuid string) error {
	key := fmt.Sprintf("UPDATE_SUMMONER_%s", puuid)

	_, err, _ := service.group.Do(key, func() (interface{}, error) {
		summoner, err := service.GetSummonerByPuuid(ctx, puuid)
		if err != nil {
			return nil, fmt.Errorf("GetSummonerByPuuid failed with err: %w", err)
		}

		err = service.CollectSummonerRank(ctx, summoner.Region, summoner.Puuid)
		// riot.ErrNotFound is an expected error and should not cause UpdateSummonerInfo to fail
		if err != nil && !errors.Is(err, riot.ErrNotFound) {
			return nil, fmt.Errorf("CollectSummonerRank failed with err: %w", err)
		}
		err = service.CollectSummonerByPuuid(ctx, summoner.Region, puuid)
		if err != nil {
			return nil, fmt.Errorf("CollectSummonerByPuuid failed with err: %w", err)
		}
		err = service.CollectMatchHistory(ctx, summoner.Region, puuid, summoner.MatchesBeforeTimestamp.Time)
		if err != nil {
			return nil, fmt.Errorf("CollectMatchHistory failed with err: %w", err)
		}

		return nil, nil
	})

	return err
}
