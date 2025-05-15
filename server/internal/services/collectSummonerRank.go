package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/internal/leaderboard"
	"github.com/owenplesko/TftAnalytics/pkg/dedupe"
)

func (service *Service) CollectSummonerRank(ctx context.Context, region, summonerId string) error {
	return dedupe.Run(service.deduplicator, SummonerRankTask{
		service:    service,
		ctx:        ctx,
		region:     region,
		summonerId: summonerId,
	}).Await()
}

type SummonerRankTask struct {
	service    *Service
	ctx        context.Context
	region     string
	summonerId string
}

func (task SummonerRankTask) ID() string {
	return fmt.Sprintf("SUMMONER_RANK_%s_%s", task.region, task.summonerId)
}

func (task SummonerRankTask) Run() error {
	rankEntry, err := task.service.riot.GetRank(task.ctx, task.region, task.summonerId)
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

	err = task.service.leaderboard.SetRank(task.ctx, task.region, setRankParam)
	if err != nil {
		return fmt.Errorf("Leaderboard.SetRank failed with err: %w", err)
	}

	log.Printf("collected rank for summoner with summonerId %v on region %v", task.summonerId, task.region)

	return nil
}
