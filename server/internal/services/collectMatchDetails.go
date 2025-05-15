package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/owenplesko/TftAnalytics/internal/db"
	"github.com/owenplesko/TftAnalytics/pkg/dedupe"
	"github.com/owenplesko/TftAnalytics/pkg/riot"

	"github.com/jackc/pgx/v5/pgtype"
)

func (service *Service) CollectMatchDetails(ctx context.Context, region, matchId string) error {
	return dedupe.Run(service.deduplicator, matchDetailsTask{
		service: service,
		ctx:     ctx,
		region:  region,
		matchId: matchId,
	}).Await()
}

type matchDetailsTask struct {
	service *Service
	ctx     context.Context
	region  string
	matchId string
}

func (task matchDetailsTask) ID() string {
	return fmt.Sprintf("MATCH_DETAILS_%s_%s", task.region, task.matchId)
}

func (task matchDetailsTask) Run() error {
	match, err := task.service.riot.GetMatchDetails(task.ctx, riot.RegionToCluster[task.region], task.matchId)
	if err != nil {
		return fmt.Errorf("Riot.GetMatchDetails failed with err: %w", err)
	}

	for _, puuid := range match.MetaData.Participants {
		if exists, _ := task.service.queries.SummonerExistsByPuuid(task.ctx, puuid); !exists {
			err = task.service.CollectSummonerByPuuid(task.ctx, task.region, puuid)
			if err != nil {
				log.Printf("error in CollectMatchDetails collecting summoner with puuid %v from match %v: CollectSummonerByPuuid failed with err: %v", puuid, task.matchId, err)
			}
		}
	}

	err = task.service.storeMatchDetails(task.ctx, match)
	if err != nil {
		return fmt.Errorf("storeMatchDetails failed with err: %w", err)
	}

	log.Printf("collected match %v on region %v\n", task.matchId, task.region)

	return nil
}

func (service *Service) storeMatchDetails(ctx context.Context, matchDetails *riot.Match) error {
	var err error

	// comp and matches inserted in one transaction
	tx, err := service.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to create db transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := service.queries.WithTx(tx)

	// insert match
	matchDate := pgtype.Timestamp{
		Time:  time.UnixMilli(matchDetails.Info.Date),
		Valid: true,
	}

	err = qtx.CreateMatch(ctx, db.CreateMatchParams{
		ID:          matchDetails.MetaData.MatchId,
		DataVersion: matchDetails.MetaData.DataVersion,
		GameVersion: matchDetails.Info.GameVersion,
		QueueID:     matchDetails.Info.QueueId,
		GameType:    matchDetails.Info.GameType,
		SetName:     matchDetails.Info.SetName,
		SetNumber:   matchDetails.Info.SetNumber,
		MatchDate:   matchDate,
	})
	if err != nil {
		return fmt.Errorf("qtx.CreateMatch failed with err: %w", err)
	}

	// insert comps
	for _, compDetails := range matchDetails.Info.Comps {
		if compDetails.Puuid == "BOT" {
			continue
		}

		err = qtx.CreateComp(ctx, db.CreateCompParams{
			MatchID:       matchDetails.MetaData.MatchId,
			SummonerPuuid: compDetails.Puuid,
			CompData:      db.CompData(compDetails),
			MatchDate:     matchDate,
		})
		if err != nil {
			return fmt.Errorf("qtx.CreateComp failed with err: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit db transaction: %w", err)
	}

	return nil
}
