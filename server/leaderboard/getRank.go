package leaderboard

import (
	"TFTAnalyticsServer/types"
	"context"
	"encoding/json"
)

func (leaderboard Leaderboard) GetRank(ctx context.Context, leaderboardName string, summonerId string) (types.Rank, error) {
	rank := types.Rank{}

	rankData, err := leaderboard.getRankData(ctx, summonerId)
	if err != nil {
		return rank, err
	}
	rank.Data = rankData

	leaderboardPositions, err := leaderboard.getLeaderboardPositions(ctx, leaderboardName, summonerId)
	if err != nil {
		return rank, err
	}
	rank.LeaderboardPositions = leaderboardPositions

	return rank, nil
}

func (leaderboard Leaderboard) getRankData(ctx context.Context, summonerId string) (types.RankData, error) {
	var rankData types.RankData

	jsonRaw, err := leaderboard.rdb.JSONGet(ctx, "rank:"+summonerId).Result()
	if err != nil {
		return rankData, err
	}

	err = json.Unmarshal([]byte(jsonRaw), &rankData)

	return rankData, err
}

func (leaderboard Leaderboard) getLeaderboardPositions(ctx context.Context, learderboardName string, summonerId string) ([]types.LeaderboardPosition, error) {
	gloablPosition, err := leaderboard.rdb.ZRevRank(ctx, "leaderboard:global", summonerId).Result()
	if err != nil {
		return nil, err
	}

	regionPosition, err := leaderboard.rdb.ZRevRank(ctx, "leaderboard:"+learderboardName, summonerId).Result()
	if err != nil {
		return nil, err
	}

	leaderboardPositions := []types.LeaderboardPosition{
		types.LeaderboardPosition{
			Leaderboard: "global",
			Position:    int(gloablPosition) + 1,
		},
		types.LeaderboardPosition{
			Leaderboard: learderboardName,
			Position:    int(regionPosition) + 1,
		},
	}

	return leaderboardPositions, nil
}
