package services

import (
	"TFTAnalyticsServer/internal/leaderboard"
	"context"
	"fmt"
	"log"
)

func (service *Service) GetSummonerRank(ctx context.Context, region string, summonerId string) (leaderboard.Rank, error) {
	rank, err := service.leaderboard.GetRank(ctx, region, summonerId)
	if err != nil {
		return leaderboard.Rank{}, fmt.Errorf("Leaderboard.GetRank failed with err: %w", err)
	}

	return rank, nil
}

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

func (service *Service) CollectRankEntries(ctx context.Context, region, tier, division string) error {
	page := 0
	for {
		page++
		rankEntries, err := service.riot.GetRankEntries(ctx, region, tier, division, page)
		if err != nil {
			return fmt.Errorf("Riot.GetRankEntries failed with err: %w", err)
		}
		if len(rankEntries) == 0 {
			break
		}

		setRankParams := make([]leaderboard.SetRankParams, len(rankEntries))
		for i, rankEntry := range rankEntries {
			setRankParams[i] = leaderboard.SetRankParams{
				SummonerId: rankEntry.SummonerId,
				RankData: leaderboard.RankData{
					Tier:         rankEntry.Tier,
					Rank:         rankEntry.Rank,
					LeaguePoints: int(rankEntry.LeaguePoints),
				}}
		}

		err = service.leaderboard.SetRank(ctx, region, setRankParams...)
		if err != nil {
			return fmt.Errorf("Leaderboard.SetRank failed with err: %w", err)
		}
	}

	log.Printf("collected rank %v %v entries on region %v\n", tier, division, region)

	return nil
}

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
