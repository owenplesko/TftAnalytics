package services

import (
	"context"
	"fmt"

	"github.com/owenplesko/TftAnalytics/internal/leaderboard"
)

func (service *Service) GetSummonerRank(ctx context.Context, region string, summonerId string) (leaderboard.Rank, error) {
	rank, err := service.leaderboard.GetRank(ctx, region, summonerId)
	if err != nil {
		return leaderboard.Rank{}, fmt.Errorf("Leaderboard.GetRank failed with err: %w", err)
	}

	return rank, nil
}
