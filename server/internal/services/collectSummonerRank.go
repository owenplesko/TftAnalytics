package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/internal/leaderboard"
)

func (service *Service) CollectSummonerRank(ctx context.Context, region, summonerId string) error {
	rankEntry, err := service.riot.GetRank(ctx, region, summonerId)
	if err != nil {
		return fmt.Errorf("Riot.GetRank failed with err: %w", err)
	}

	setRankParam := leaderboard.SetRankParams{
		SummonerId: rankEntry.SummonerId,
		RankData: leaderboard.RankData{
			Tier:         rankEntry.Tier,
			Rank:         rankEntry.Rank,
			LeaguePoints: int(rankEntry.LeaguePoints),
		}}

	err = service.leaderboard.SetRank(ctx, region, setRankParam)
	if err != nil {
		return fmt.Errorf("Leaderboard.SetRank failed with err: %w", err)
	}

	log.Printf("collected rank for summoner with summonerId %v on region %v", summonerId, region)

	return nil
}
