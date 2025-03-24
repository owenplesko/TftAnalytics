package leaderboard

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type SetRankParams struct {
	SummonerId string
	RankData   RankData
}

func (leaderboard Leaderboard) SetRank(ctx context.Context, leaderboardName string, params ...SetRankParams) error {
	err := leaderboard.setLeaderboardScores(ctx, leaderboardName, params...)
	if err != nil {
		return err
	}

	err = leaderboard.setRankData(ctx, params...)

	return err
}

func (leaderboard Leaderboard) setLeaderboardScores(ctx context.Context, leaderboardName string, params ...SetRankParams) error {
	zParams := make([]redis.Z, len(params))

	for i, param := range params {
		zParams[i] = redis.Z{
			Member: param.SummonerId,
			Score:  float64(getRankScore(param.RankData)),
		}
	}

	err := leaderboard.rdb.ZAdd(ctx, "leaderboard:global", zParams...).Err()
	if err != nil {
		return err
	}

	err = leaderboard.rdb.ZAdd(ctx, "leaderboard:"+leaderboardName, zParams...).Err()

	return err
}

func (leaderboard Leaderboard) setRankData(ctx context.Context, params ...SetRankParams) error {
	pipe := leaderboard.rdb.Pipeline()
	for _, param := range params {
		bytes, _ := json.Marshal(param.RankData)
		pipe.Do(ctx, "JSON.SET", "rank:"+param.SummonerId, "$", bytes)
	}
	_, err := pipe.Exec(ctx)
	return err
}
