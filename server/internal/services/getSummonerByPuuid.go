package services

import (
	"context"
	"fmt"

	"github.com/owenplesko/TftAnalytics/internal/db"
)

func (service *Service) GetSummonerByPuuid(ctx context.Context, puuid string) (db.TftSummoner, error) {
	summoner, err := service.queries.GetSummonerByPuuid(ctx, puuid)
	if err != nil {
		return summoner, fmt.Errorf("Queries.GetSummonerByPuuid failed with err: %w", err)
	}

	return summoner, nil
}
