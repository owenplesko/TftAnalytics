package services

import (
	"TFTAnalyticsServer/internal/db"
	"TFTAnalyticsServer/pkg/riot"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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

func (service *Service) CollectMatchDetails(ctx context.Context, region, matchId string) error {
	match, err := service.riot.GetMatchDetails(ctx, riot.RegionToCluster[region], matchId)
	if err != nil {
		return fmt.Errorf("Riot.GetMatchDetails failed with err: %w", err)
	}

	for _, puuid := range match.MetaData.Participants {
		if exists, _ := service.queries.SummonerExistsByPuuid(ctx, puuid); !exists {
			err = service.CollectSummonerByPuuid(ctx, region, puuid)
			if err != nil {
				log.Printf("error in CollectMatchDetails collecting summoner with puuid %v from match %v: CollectSummonerByPuuid failed with err: %v", puuid, matchId, err)
			}
		}
	}

	err = service.storeMatchDetails(ctx, match)
	if err != nil {
		return fmt.Errorf("storeMatchDetails failed with err: %w", err)
	}

	log.Printf("collected match %v on region %v\n", matchId, region)

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
