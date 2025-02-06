package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"TFTAnalyticsServer/types"
	"context"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (env Service) GetMatchComps(ctx context.Context, matchId string) ([]db.GetMatchCompsRow, error) {
	comps, err := env.Queries.GetMatchComps(ctx, matchId)

	// prevent returning nil when list is empty
	if comps == nil {
		comps = []db.GetMatchCompsRow{}
	}

	return comps, err
}

func (env Service) GetMatchHistory(ctx context.Context, puuid string, limit int32, after time.Time) ([]db.SummonerMatchHistoryRow, error) {
	matches, err := env.Queries.SummonerMatchHistory(context.Background(), db.SummonerMatchHistoryParams{
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

func (env Service) CollectMatchHistory(ctx context.Context, cluster, puuid string, matchesAfter time.Time) error {
	updatedAt := time.Now()

	res, err := riot.GetMatchHistory(cluster, puuid, matchesAfter)
	if err != nil {
		return err
	}

	env.Queries.SetBackgroundUpdateTimestamp(ctx, db.SetBackgroundUpdateTimestampParams{
		Puuid: puuid,
		BackgroundUpdateTimestamp: pgtype.Timestamp{
			Time:  updatedAt,
			Valid: true,
		},
	})

	log.Printf("Got %v matchIds from summoner %v\n", len(res), puuid)

	for _, matchId := range res {
		if exists, _ := env.Queries.MatchExists(ctx, matchId); exists {
			log.Printf("Skipping match %v...\n", matchId)
			return nil
		}

		res, err := riot.GetMatchDetails(cluster, matchId)
		if err != nil {
			return err
		}

		err = env.storeMatchDetails(ctx, res)

		log.Printf("Stored match %v!\n", matchId)
	}

	return err
}

func extractPuuidsFromMatchDetails(matchDetails *riot.Match) []db.BatchUpsertPuuidsParams {
	region := strings.Split(matchDetails.MetaData.MatchId, "_")[0]
	upsertParams := make([]db.BatchUpsertPuuidsParams, len(matchDetails.MetaData.Participants))

	for i, puuid := range matchDetails.MetaData.Participants {
		upsertParams[i] = db.BatchUpsertPuuidsParams{
			Puuid:  puuid,
			Region: region,
		}
	}

	return upsertParams
}

func (env Service) storeMatchDetails(ctx context.Context, matchDetails *riot.Match) error {
	var err error

	env.batchStoreSummonerPuuid(ctx, extractPuuidsFromMatchDetails(matchDetails))

	tx, err := env.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := env.Queries.WithTx(tx)

	// create match
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

	// create comps
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
