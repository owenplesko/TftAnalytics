package services

import (
	"TFTAnalyticsServer/internal/leaderboard"
	"context"
	"fmt"
	"log"
)

func (service *Service) CollectApexRankEntries(ctx context.Context, region, tier string) error {
	rankPage, err := service.riot.GetApexRankPage(ctx, region, tier)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	setRankParams := make([]leaderboard.SetRankParams, len(rankPage.Entries))
	for i, rankEntry := range rankPage.Entries {
		setRankParams[i] = leaderboard.SetRankParams{
			SummonerId: rankEntry.SummonerId,
			RankData: leaderboard.RankData{
				Tier:         tier,
				Rank:         rankEntry.Rank,
				LeaguePoints: int(rankEntry.LeaguePoints),
			},
		}
	}

	err = service.leaderboard.SetRank(ctx, region, setRankParams...)
	if err != nil {
		return fmt.Errorf("Leaderboard.SetRank failed with err: %w", err)
	}

	log.Printf("collected rank %v entries on region %v\n", tier, region)

	return nil
}
