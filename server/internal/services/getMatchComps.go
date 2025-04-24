package services

import (
	"TFTAnalyticsServer/internal/db"
	"context"
	"fmt"
)

func (service *Service) GetMatchComps(ctx context.Context, matchId string) ([]db.GetMatchCompsRow, error) {
	comps, err := service.queries.GetMatchComps(ctx, matchId)
	if err != nil {
		return []db.GetMatchCompsRow{}, fmt.Errorf("Queries.GetMatchComps failed with err: %w", err)
	}

	// prevent returning nil when list is empty
	if comps == nil {
		comps = []db.GetMatchCompsRow{}
	}

	return comps, nil
}
