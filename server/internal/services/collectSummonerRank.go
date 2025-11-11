package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/internal/leaderboard"
)

func (service *Service) CollectSummonerRank(ctx context.Context, region, puuid string) error {
	key := fmt.Sprintf("SUMMONER_RANK_%s_%s", region, puuid)

	_, err, _ := service.group.Do(key, func() (interface{}, error) {
		rankEntry, err := service.riot.GetRank(ctx, region, puuid)
		if err != nil {
			return nil, fmt.Errorf("Riot.GetRank failed with err: %w", err)
		}

		setRankParam := leaderboard.SetRankParams{
			Puuid: rankEntry.Puuid,
			RankData: leaderboard.RankData{
				Tier:         rankEntry.Tier,
				Rank:         rankEntry.Rank,
				LeaguePoints: int(rankEntry.LeaguePoints),
			}}

		err = service.leaderboard.SetRank(ctx, region, setRankParam)
		if err != nil {
			return nil, fmt.Errorf("Leaderboard.SetRank failed with err: %w", err)
		}

		log.Printf("collected rank for summoner with puuid %v on region %v", puuid, region)

		return nil, nil
	})

	return err
}
