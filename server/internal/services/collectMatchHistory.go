package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/owenplesko/TftAnalytics/internal/db"
	"github.com/owenplesko/TftAnalytics/pkg/riot"
)

func (service *Service) CollectMatchHistory(ctx context.Context, region, puuid string, matchesAfter time.Time) error {
	key := fmt.Sprintf("MATCH_HISTORY_%s_%s", region, puuid)

	_, err, _ := service.group.Do(key, func() (interface{}, error) {
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
				return nil, fmt.Errorf("Riot.GetMatchHistory failed with err: %w", err)
			}

			if len(res) == 0 {
				break
			}

			matchIds = append(matchIds, res...)
		}

		for _, matchId := range matchIds {
			err := service.CollectMatchDetails(ctx, region, matchId)
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
			return nil, fmt.Errorf("Queries.SetMatchesAfterTimestamp failed with err: %w", err)
		}

		return nil, nil

	})
	return err
}

