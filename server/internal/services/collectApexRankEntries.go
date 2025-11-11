package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/internal/leaderboard"
)

func (service *Service) CollectApexRankEntries(ctx context.Context, region, tier string) error {
	key := fmt.Sprintf("APEX_RANK_ENTRIES_%s_%s", region, tier)

	_, err, _ := service.group.Do(key, func() (interface{}, error) {
		rankPage, err := service.riot.GetApexRankPage(ctx, region, tier)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		setRankParams := make([]leaderboard.SetRankParams, len(rankPage.Entries))
		for i, rankEntry := range rankPage.Entries {
			setRankParams[i] = leaderboard.SetRankParams{
				Puuid: rankEntry.Puuid,
				RankData: leaderboard.RankData{
					Tier:         tier,
					Rank:         rankEntry.Rank,
					LeaguePoints: int(rankEntry.LeaguePoints),
				},
			}
		}

		err = service.leaderboard.SetRank(ctx, region, setRankParams...)
		if err != nil {
			return nil, fmt.Errorf("Leaderboard.SetRank failed with err: %w", err)
		}

		log.Printf("collected rank %v entries on region %v\n", tier, region)

		return nil, nil

	})

	return err
}
