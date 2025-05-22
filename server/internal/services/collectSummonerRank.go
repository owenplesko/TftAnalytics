package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/internal/leaderboard"
	"github.com/owenplesko/TftAnalytics/pkg/dedupe"
)

func (service *Service) CollectSummonerRank(ctx context.Context, region, puuid string) error {
	return dedupe.Run(service.deduplicator, summonerRankTask{
		service: service,
		ctx:     ctx,
		region:  region,
		puuid:   puuid,
	}).Await()
}

type summonerRankTask struct {
	service *Service
	ctx     context.Context
	region  string
	puuid   string
}

func (task summonerRankTask) ID() string {
	return fmt.Sprintf("SUMMONER_RANK_%s_%s", task.region, task.puuid)
}

func (task summonerRankTask) Run() error {
	rankEntry, err := task.service.riot.GetRank(task.ctx, task.region, task.puuid)
	if err != nil {
		return fmt.Errorf("Riot.GetRank failed with err: %w", err)
	}

	setRankParam := leaderboard.SetRankParams{
		Puuid: rankEntry.Puuid,
		RankData: leaderboard.RankData{
			Tier:         rankEntry.Tier,
			Rank:         rankEntry.Rank,
			LeaguePoints: int(rankEntry.LeaguePoints),
		}}

	err = task.service.leaderboard.SetRank(task.ctx, task.region, setRankParam)
	if err != nil {
		return fmt.Errorf("Leaderboard.SetRank failed with err: %w", err)
	}

	log.Printf("collected rank for summoner with puuid %v on region %v", task.puuid, task.region)

	return nil
}
