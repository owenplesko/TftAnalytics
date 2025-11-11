package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/owenplesko/TftAnalytics/internal/db"
	"github.com/owenplesko/TftAnalytics/pkg/riot"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
)

func (service *Service) CollectMatchDetails(ctx context.Context, region, matchId string) error {
	key := fmt.Sprintf("MATCH_DETAILS_%s_%s", region, matchId)
	_, err, _ := service.group.Do(key, func() (interface{}, error) {
		exists, _ := service.queries.MatchExists(ctx, matchId)
		if exists {
			return nil, nil
		}

		match, err := service.riot.GetMatchDetails(ctx, riot.RegionToCluster[region], matchId)
		if err != nil {
			return nil, fmt.Errorf("Riot.GetMatchDetails failed with err: %w", err)
		}

		for _, puuid := range match.MetaData.Participants {
			if exists, _ := service.queries.SummonerExistsByPuuid(ctx, puuid); !exists {
				err = service.CollectSummonerByPuuid(ctx, region, puuid)
				if err != nil {
					log.Printf("error in CollectMatchDetails collecting summoner with puuid %v from match %v: CollectSummonerByPuuid failed with err: %v", puuid, matchId, err)
				}
			}
		}

		err = service.storeMatchDetails(ctx, region, match)
		if err != nil {
			return nil, fmt.Errorf("storeMatchDetails failed with err: %w", err)
		}

		log.Printf("collected match %v on region %v\n", matchId, region)

		return nil, nil

	})
	return err
}

func extractPatchNumber(input string) (string, error) {
	re := regexp.MustCompile(`<Releases/(\d+\.\d+)>`)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return "", fmt.Errorf("patch number not found")
	}

	return matches[1], nil
}

func (service *Service) storeMatchDetails(ctx context.Context, region string, matchDetails *riot.Match) error {
	// get relevant data before starting transaction
	patchNumber, err := extractPatchNumber(matchDetails.Info.GameVersion)
	if err != nil {
		return err
	}

	ranks := make([]string, len(matchDetails.MetaData.Participants))
	for i, puuid := range matchDetails.MetaData.Participants {
		rankEntry, err := service.leaderboard.GetRank(ctx, region, puuid)

		tier := rankEntry.Data.Tier
		if errors.Is(err, redis.Nil) {
			tier = "unknown"
		} else if err != nil {
			return err
		}
		ranks[i] = tier
	}

	// comp and matches inserted in one transaction
	tx, err := service.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to create db transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// insert match
	qtx := service.queries.WithTx(tx)

	matchDate := pgtype.Timestamptz{
		Time:  time.UnixMilli(matchDetails.Info.Date),
		Valid: true,
	}

	err = qtx.CreateMatch(ctx, db.CreateMatchParams{
		ID:          matchDetails.MetaData.MatchId,
		DataVersion: matchDetails.MetaData.DataVersion,
		GameVersion: patchNumber,
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
	for i, compDetails := range matchDetails.Info.Comps {
		if compDetails.Puuid == "BOT" {
			continue
		}

		err = qtx.CreateComp(ctx, db.CreateCompParams{
			MatchID:       matchDetails.MetaData.MatchId,
			SummonerPuuid: compDetails.Puuid,
			CompData:      db.CompData(compDetails),
			MatchDate:     matchDate,
			Rank:          ranks[i],
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
