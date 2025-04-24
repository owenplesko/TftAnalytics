package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/owenplesko/TftAnalytics/internal/db"

	"github.com/jackc/pgx/v5"
)

func (service *Service) GetOrCollectSummonerByNameTag(ctx context.Context, cluster, name, tag string) (db.TftSummoner, error) {
	summoner, err := service.queries.GetSummonerByNameTag(ctx, db.GetSummonerByNameTagParams{
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

	summoner, err = service.queries.GetSummonerByNameTag(ctx, db.GetSummonerByNameTagParams{Name: name, Tag: tag})
	if err != nil {
		return summoner, fmt.Errorf("Queries.GetSummonerByNameTag failed with err: %w", err)
	}

	return summoner, nil
}
