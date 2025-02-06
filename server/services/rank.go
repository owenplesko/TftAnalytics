package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
)

func (env Service) GetSummonerRank(ctx context.Context, puuid string) (db.TftRank, error) {
	return env.Queries.GetSummonerRank(ctx, puuid)
}

func (env Service) CollectSummonerRank(ctx context.Context, region, summonerId string) error {
	rankEntry, err := riot.GetRank(region, summonerId)
	if err != nil {
		return err
	}

	return env.storeSummonerRank(ctx, rankEntry)
}

func (env Service) storeSummonerRank(ctx context.Context, rankEntry riot.RankEntry) error {
	return env.Queries.UpsertSummonerRank(ctx, db.UpsertSummonerRankParams{
		SummonerPuuid: rankEntry.Puuid,
		Tier:          rankEntry.Tier,
		Rank:          rankEntry.Rank,
		LeaguePoints:  rankEntry.LeaguePoints,
		Wins:          rankEntry.Wins,
		Losses:        rankEntry.Losses,
	})
}

func (env Service) CollectRankEntries(ctx context.Context, region, tier, division string) {
	page := 0
	for {
		page++
		rankEntries, err := riot.GetRankEntries(region, tier, division, page)
		if err != nil || len(rankEntries) == 0 {
			break
		}

	}
}
