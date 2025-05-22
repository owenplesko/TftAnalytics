package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/internal/leaderboard"
	"github.com/owenplesko/TftAnalytics/pkg/dedupe"
)

func (service *Service) CollectRankEntries(ctx context.Context, region, tier, division string) error {
	return dedupe.Run(service.deduplicator, rankEntriesTask{
		service:  service,
		ctx:      ctx,
		region:   region,
		tier:     tier,
		division: division,
	}).Await()
}

type rankEntriesTask struct {
	service  *Service
	ctx      context.Context
	region   string
	tier     string
	division string
}

func (task rankEntriesTask) ID() string {
	return fmt.Sprintf("RANK_ENTRIES_%s_%s_%s", task.region, task.tier, task.division)
}

func (task rankEntriesTask) Run() error {
	page := 0
	for {
		page++
		rankEntries, err := task.service.riot.GetRankEntries(task.ctx, task.region, task.tier, task.division, page)
		if err != nil {
			return fmt.Errorf("Riot.GetRankEntries failed with err: %w", err)
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

		err = task.service.leaderboard.SetRank(task.ctx, task.region, setRankParams...)
		if err != nil {
			return fmt.Errorf("Leaderboard.SetRank failed with err: %w", err)
		}
	}

	log.Printf("collected rank %v %v entries on region %v\n", task.tier, task.division, task.region)

	return nil

}
