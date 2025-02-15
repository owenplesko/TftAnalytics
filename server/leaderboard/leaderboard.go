package leaderboard

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func New(rdb *redis.Client) *Leaderboard {
	return &Leaderboard{rdb}
}

type Leaderboard struct {
	rdb *redis.Client
}

type SetLeaderboardScoresParams struct {
	SummonerId string
	RankScore  int
}

func (leaderboard Leaderboard) SetLeaderboardScore(ctx context.Context, param SetLeaderboardScoresParams) error {
	zParam := redis.Z{
		Member: param.SummonerId,
		Score:  float64(param.RankScore),
	}

	err := leaderboard.rdb.ZAdd(ctx, "leaderboard:global", zParam).Err()
	return err
}

func (leaderboard Leaderboard) BatchSetLeaderboardScores(ctx context.Context, params []SetLeaderboardScoresParams) error {
	zParams := make([]redis.Z, len(params))

	for i, param := range params {
		zParams[i] = redis.Z{
			Member: param.SummonerId,
			Score:  float64(param.RankScore),
		}
	}

	err := leaderboard.rdb.ZAdd(ctx, "leaderboard:global", zParams...).Err()
	return err
}
