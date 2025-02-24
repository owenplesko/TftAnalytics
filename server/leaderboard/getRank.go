package leaderboard

import (
	"TFTAnalyticsServer/types"
	"context"
	"encoding/json"
)

func (leaderboard Leaderboard) GetRank(ctx context.Context, region string, summonerId string) (types.Rank, error) {
	rankData, err := leaderboard.getRankData(ctx, summonerId)
	if err != nil {
		return types.Rank{}, err
	}

	leaderboardPosition, err := leaderboard.getLeaderboardPosition(ctx, region, summonerId)
	if err != nil {
		return types.Rank{}, err
	}

	rank := types.Rank{
		Data:                rankData,
		LeaderboardPosition: leaderboardPosition,
	}

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

func (leaderboard Leaderboard) getLeaderboardPosition(ctx context.Context, region string, summonerId string) (types.LeaderboardPosition, error) {
	leaderboardName := "leaderboard:" + region

	position, err := leaderboard.rdb.ZRevRank(ctx, leaderboardName, summonerId).Result()
	if err != nil {
		return types.LeaderboardPosition{}, err
	}

	cardinality, err := leaderboard.rdb.ZCard(ctx, leaderboardName).Result()
	if err != nil {
		return types.LeaderboardPosition{}, err
	}

	leaderboardPosition := types.LeaderboardPosition{
		Leaderboard: region,
		Position:    int(position + 1),
		Top:         float32(position+1) / float32(cardinality) * 100,
	}

	return leaderboardPosition, nil
}
