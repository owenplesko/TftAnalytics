package leaderboard

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

func (leaderboard Leaderboard) GetRank(ctx context.Context, region string, summonerId string) (Rank, error) {
	leaderboardPosition, err := leaderboard.getLeaderboardPosition(ctx, region, summonerId)
	if err != nil {
		return Rank{}, err
	}

	rankData, err := leaderboard.getRankData(ctx, summonerId)
	if err != nil {
		return Rank{}, err
	}

	rank := Rank{
		Data:                rankData,
		LeaderboardPosition: leaderboardPosition,
	}

	return rank, nil
}

func (leaderboard Leaderboard) getRankData(ctx context.Context, summonerId string) (RankData, error) {
	var rankData RankData

	jsonRaw, err := leaderboard.rdb.JSONGet(ctx, "rank:"+summonerId).Result()
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

func (leaderboard Leaderboard) getLeaderboardPosition(ctx context.Context, region string, summonerId string) (LeaderboardPosition, error) {
	leaderboardName := "leaderboard:" + region

	position, err := leaderboard.rdb.ZRevRank(ctx, leaderboardName, summonerId).Result()
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
