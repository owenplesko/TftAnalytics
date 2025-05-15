package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/owenplesko/TftAnalytics/internal/db"
	"github.com/owenplesko/TftAnalytics/pkg/dedupe"
	"github.com/owenplesko/TftAnalytics/pkg/riot"
)

func (service *Service) CollectMatchHistory(ctx context.Context, region, puuid string, matchesAfter time.Time) error {
	return dedupe.Run(service.deduplicator, matchHistoryTask{
		service:      service,
		ctx:          ctx,
		region:       region,
		puuid:        puuid,
		matchesAfter: matchesAfter,
	}).Await()
}

type matchHistoryTask struct {
	service      *Service
	ctx          context.Context
	region       string
	puuid        string
	matchesAfter time.Time
}

func (task matchHistoryTask) ID() string {
	return fmt.Sprintf("MATCH_HISTORY_%s_%s", task.region, task.puuid)
}

func (task matchHistoryTask) Run() error {
	count := riot.MATCH_HISTORY_MAX_COUNT

	matchesBefore := time.Now()
	if task.matchesAfter.Before(task.service.matchesAfterCutoff) {
		task.matchesAfter = task.service.matchesAfterCutoff
	}

	matchIds := make([]string, 0, count)

	for {
		startIndex := len(matchIds)
		res, err := task.service.riot.GetMatchHistoryInTimeRange(task.ctx, riot.RegionToCluster[task.region], task.puuid, count, task.matchesAfter, matchesBefore, startIndex)
		if err != nil {
			return fmt.Errorf("Riot.GetMatchHistory failed with err: %w", err)
		}

		if len(res) == 0 {
			break
		}

		matchIds = append(matchIds, res...)
	}

	for _, matchId := range matchIds {
		exists, err := task.service.queries.MatchExists(task.ctx, matchId)
		if err != nil {
			log.Printf("error getting match exists: %v", err)
		}
		if exists {
			continue
		}

		err = task.service.CollectMatchDetails(task.ctx, task.region, matchId)
		if err != nil {
			// TODO: explore returning CollectMatchDetails err
			log.Printf("error in CollectMatchHistory collecting match %v for summoner with puuid %v: CollectMatchDetails failed with err: %v", matchId, task.puuid, err)
		}
	}

	err := task.service.queries.SetMatchesBeforeTimestamp(task.ctx, db.SetMatchesBeforeTimestampParams{
		Puuid: task.puuid,
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
