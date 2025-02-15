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
	Id    string
	Score int
}

func (leaderboard Leaderboard) SetLeaderboardScores(ctx context.Context, params ...SetLeaderboardScoresParams) error {
	zParams := make([]redis.Z, len(params))

	for i, param := range params {
		zParams[i] = redis.Z{
			Member: param.Id,
			Score:  float64(param.Score),
		}
	}

	err := leaderboard.rdb.ZAdd(ctx, "leaderboard:global", zParams...).Err()
	return err
}
