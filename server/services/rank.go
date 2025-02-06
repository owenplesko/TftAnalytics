package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
	"log"
)

func (service Service) GetSummonerRank(ctx context.Context, puuid string) (db.TftRank, error) {
	return service.Queries.GetSummonerRank(ctx, puuid)
}

func (service Service) CollectSummonerRank(ctx context.Context, region, summonerId string) error {
	rankEntry, err := riot.GetRank(region, summonerId)
	if err != nil {
		return err
	}

	return service.storeSummonerRank(ctx, rankEntry)
}

func (service Service) storeSummonerRank(ctx context.Context, rankEntry riot.RankEntry) error {
	return service.Queries.UpsertSummonerRank(ctx, db.UpsertSummonerRankParams{
		SummonerPuuid: rankEntry.Puuid,
		Tier:          rankEntry.Tier,
		Rank:          rankEntry.Rank,
		LeaguePoints:  rankEntry.LeaguePoints,
		Wins:          rankEntry.Wins,
		Losses:        rankEntry.Losses,
	})
}

func (service Service) CollectRankEntries(ctx context.Context, region, tier, division string) error {
	page := 0
	for {
		page++
		rankEntries, err := riot.GetRankEntries(region, tier, division, page)
		if err != nil {
			return err
		}
		if len(rankEntries) == 0 {
			break
		}

		for _, rankEntry := range rankEntries {
			service.Queries.InsertPuuid(ctx, db.InsertPuuidParams{
				Puuid:  rankEntry.Puuid,
				Region: region,
			})
			service.storeSummonerRank(ctx, rankEntry)
		}
	}

	log.Printf("Rank entries collected for %v %v %v\n", region, tier, division)

	return nil
}
