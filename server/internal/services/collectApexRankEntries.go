package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/internal/leaderboard"
	"github.com/owenplesko/TftAnalytics/pkg/dedupe"
)

func (service *Service) CollectApexRankEntries(ctx context.Context, region, tier string) error {
	return dedupe.Run(service.deduplicator, apexRankEntriesTask{
		service: service,
		ctx:     ctx,
		region:  region,
		tier:    tier,
	}).Await()
}

type apexRankEntriesTask struct {
	service *Service
	ctx     context.Context
	region  string
	tier    string
}

func (task apexRankEntriesTask) ID() string {
	return fmt.Sprintf("APEX_RANK_ENTRIES_%s_%s", task.region, task.tier)
}

func (task apexRankEntriesTask) Run() error {
	rankPage, err := task.service.riot.GetApexRankPage(task.ctx, task.region, task.tier)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	setRankParams := make([]leaderboard.SetRankParams, len(rankPage.Entries))
	for i, rankEntry := range rankPage.Entries {
		setRankParams[i] = leaderboard.SetRankParams{
			SummonerId: rankEntry.SummonerId,
			RankData: leaderboard.RankData{
				Tier:         task.tier,
				Rank:         rankEntry.Rank,
				LeaguePoints: int(rankEntry.LeaguePoints),
			},
		}
	}

	err = task.service.leaderboard.SetRank(task.ctx, task.region, setRankParams...)
	if err != nil {
		return fmt.Errorf("Leaderboard.SetRank failed with err: %w", err)
	}

	log.Printf("collected rank %v entries on region %v\n", task.tier, task.region)

	return nil
}
