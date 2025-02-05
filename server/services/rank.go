package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
)

func (env ServiceEnv) CollectSummonerRank(ctx context.Context, region, summonerId string) error {
	rankEntry, err := riot.GetRank(region, summonerId)
	if err != nil {
		return err
	}

	return env.storeSummonerRank(ctx, rankEntry)
}

func (env ServiceEnv) storeSummonerRank(ctx context.Context, rankEntry *riot.RankEntry) error {
	return env.Queries.InsertSummonerRank(ctx, db.InsertSummonerRankParams{
		SummonerPuuid: rankEntry.Puuid,
		QueueType:     rankEntry.QueueType,
		Tier:          rankEntry.Tier,
		Rank:          rankEntry.Rank,
		LeaguePoints:  rankEntry.LeaguePoints,
		Wins:          rankEntry.Wins,
		Losses:        rankEntry.Losses,
	})
}
