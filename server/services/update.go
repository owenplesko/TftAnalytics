package services

import (
	"TFTAnalyticsServer/riot"
	"context"
)

func (env ServiceEnv) UpdateAllSummonerInfo(ctx context.Context, puuid string) error {
	summoner, err := env.GetSummonerByPuuid(ctx, puuid)
	if err != nil {
		return err
	}

	go env.CollectSummonerRank(ctx, summoner.Region.String, summoner.SummonerID.String)
	go env.CollectSummonerDetails(ctx, summoner.Region.String, puuid)
	go env.CollectAccountByPuuid(ctx, "americas", puuid)
	go env.CollectMatchHistory(ctx, riot.RegionToCluster[summoner.Region.String], puuid, summoner.FullUpdateTimestamp.Time)

	return nil
}
