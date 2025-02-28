package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"TFTAnalyticsServer/types"
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (service Service) GetMatchComps(ctx context.Context, matchId string) ([]db.GetMatchCompsRow, error) {
	comps, err := service.Queries.GetMatchComps(ctx, matchId)

	// prevent returning nil when list is empty
	if comps == nil {
		comps = []db.GetMatchCompsRow{}
	}

	return comps, err
}

func (service Service) GetMatchHistory(ctx context.Context, puuid string, limit int32, after time.Time) ([]db.SummonerMatchHistoryRow, error) {
	matches, err := service.Queries.SummonerMatchHistory(context.Background(), db.SummonerMatchHistoryParams{
		SummonerPuuid: puuid,
		Limit:         limit,
		After: pgtype.Timestamp{
			Time:  after,
			Valid: true,
		},
	})

	// prevent returning nil when list is empty
	if matches == nil {
		matches = []db.SummonerMatchHistoryRow{}
	}

	return matches, err
}

func (service Service) CollectMatchHistory(ctx context.Context, region, puuid string, matchesAfter time.Time) error {
	updatedAt := time.Now()

	matchIds, err := service.Riot.GetMatchHistory(region, puuid, matchesAfter)
	if err != nil {
		return err
	}

	for _, matchId := range matchIds {
		if exists, _ := service.Queries.MatchExists(ctx, matchId); !exists {
			_ = service.CollectMatchDetails(ctx, region, matchId)
		}
	}

	service.Queries.SetMatchesAfterTimestamp(ctx, db.SetMatchesAfterTimestampParams{
		Puuid: puuid,
		MatchesAfterTimestamp: pgtype.Timestamp{
			Time:  updatedAt,
			Valid: true,
		},
	})

	return nil
}

func (service Service) CollectMatchDetails(ctx context.Context, region, matchId string) error {
	match, err := service.Riot.GetMatchDetails(region, matchId)
	if err != nil {
		return err
	}

	for _, puuid := range match.MetaData.Participants {
		if exists, _ := service.Queries.SummonerExistsByPuuid(ctx, puuid); !exists {
			_ = service.CollectSummonerByPuuid(ctx, region, puuid)
		}
	}

	err = service.storeMatchDetails(ctx, match)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("Stored match %v!\n", matchId)

	return nil
}

func (service Service) storeMatchDetails(ctx context.Context, matchDetails *riot.Match) error {
	var err error

	// comp and matches inserted in one transaction
	tx, err := service.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := service.Queries.WithTx(tx)

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
		return err
	}

	// insert comps
	for _, compDetails := range matchDetails.Info.Comps {
		err = qtx.CreateComp(ctx, db.CreateCompParams{
			MatchID:       matchDetails.MetaData.MatchId,
			SummonerPuuid: compDetails.Puuid,
			CompData:      types.CompData(compDetails),
			MatchDate:     matchDate,
		})
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)

	return err
}
