package services

import (
	"TFTAnalyticsServer/internal/db"
	"TFTAnalyticsServer/pkg/riot"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"time"
)

func (service *Service) CollectMatchHistory(ctx context.Context, region, puuid string, matchesAfter time.Time) error {
	count := riot.MATCH_HISTORY_MAX_COUNT

	matchesBefore := time.Now()
	if matchesAfter.Before(service.matchesAfterCutoff) {
		matchesAfter = service.matchesAfterCutoff
	}

	matchIds := make([]string, 0, count)

	for {
		startIndex := len(matchIds)
		res, err := service.riot.GetMatchHistoryInTimeRange(ctx, riot.RegionToCluster[region], puuid, count, matchesAfter, matchesBefore, startIndex)
		if err != nil {
			return fmt.Errorf("Riot.GetMatchHistory failed with err: %w", err)
		}

		if len(res) == 0 {
			break
		}

		matchIds = append(matchIds, res...)
	}

	for _, matchId := range matchIds {
		exists, err := service.queries.MatchExists(ctx, matchId)
		if err != nil {
			log.Printf("error getting match exists: %v", err)
		}
		if exists {
			continue
		}

		err = service.CollectMatchDetails(ctx, region, matchId)
		if err != nil {
			// TODO: explore returning CollectMatchDetails err
			log.Printf("error in CollectMatchHistory collecting match %v for summoner with puuid %v: CollectMatchDetails failed with err: %v", matchId, puuid, err)
		}
	}

	err := service.queries.SetMatchesBeforeTimestamp(ctx, db.SetMatchesBeforeTimestampParams{
		Puuid: puuid,
		MatchesBeforeTimestamp: pgtype.Timestamp{
			Time:  matchesBefore,
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("Queries.SetMatchesAfterTimestamp failed with err: %w", err)
	}

	return nil
}
