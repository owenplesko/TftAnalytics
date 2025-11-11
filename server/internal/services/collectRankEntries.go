package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/internal/leaderboard"
)

func (service *Service) CollectRankEntries(ctx context.Context, region, tier, division string) error {
	key := fmt.Sprintf("RANK_ENTRIES_%s_%s_%s", region, tier, division)

	_, err, _ := service.group.Do(key, func() (interface{}, error) {
		page := 0
		for {
			page++
			rankEntries, err := service.riot.GetRankEntries(ctx, region, tier, division, page)
			if err != nil {
				return nil, fmt.Errorf("Riot.GetRankEntries failed with err: %w", err)
			}
			if len(rankEntries) == 0 {
				break
			}

			setRankParams := make([]leaderboard.SetRankParams, len(rankEntries))
			for i, rankEntry := range rankEntries {
				setRankParams[i] = leaderboard.SetRankParams{
					Puuid: rankEntry.Puuid,
					RankData: leaderboard.RankData{
						Tier:         rankEntry.Tier,
						Rank:         rankEntry.Rank,
						LeaguePoints: int(rankEntry.LeaguePoints),
					}}
			}

			err = service.leaderboard.SetRank(ctx, region, setRankParams...)
			if err != nil {
				return nil, fmt.Errorf("Leaderboard.SetRank failed with err: %w", err)
			}
		}

		log.Printf("collected rank %v %v entries on region %v\n", tier, division, region)

		return nil, nil
	})

	return err
}
