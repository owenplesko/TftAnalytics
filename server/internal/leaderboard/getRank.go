package leaderboard

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

func (leaderboard Leaderboard) GetRank(ctx context.Context, region string, puuid string) (Rank, error) {
	leaderboardPosition, err := leaderboard.getLeaderboardPosition(ctx, region, puuid)
	if err != nil {
		return Rank{}, err
	}

	rankData, err := leaderboard.getRankData(ctx, puuid)
	if err != nil {
		return Rank{}, err
	}

	rank := Rank{
		Data:                rankData,
		LeaderboardPosition: leaderboardPosition,
	}

	return rank, nil
}

func (leaderboard Leaderboard) getRankData(ctx context.Context, puuid string) (RankData, error) {
	var rankData RankData

	jsonRaw, err := leaderboard.rdb.JSONGet(ctx, "rank:"+puuid).Result()
	if err != nil {
		return rankData, err
	}
	// idk why JSONGet doesn't return redis.Nil on missing key!!!
	if jsonRaw == "" {
		return rankData, redis.Nil
	}

	err = json.Unmarshal([]byte(jsonRaw), &rankData)

	return rankData, err
}

func (leaderboard Leaderboard) getLeaderboardPosition(ctx context.Context, region string, puuid string) (LeaderboardPosition, error) {
	leaderboardName := "leaderboard:" + region

	position, err := leaderboard.rdb.ZRevRank(ctx, leaderboardName, puuid).Result()
	if err != nil {
		return LeaderboardPosition{}, err
	}

	cardinality, err := leaderboard.rdb.ZCard(ctx, leaderboardName).Result()
	if err != nil {
		return LeaderboardPosition{}, err
	}

	leaderboardPosition := LeaderboardPosition{
		Leaderboard: region,
		Position:    int(position + 1),
		Top:         float32(position+1) / float32(cardinality) * 100,
	}

	return leaderboardPosition, nil
}
