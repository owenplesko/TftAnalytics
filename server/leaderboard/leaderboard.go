package leaderboard

import (
	"TFTAnalyticsServer/types"
	"context"
	"encoding/json"

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

func (leaderboard Leaderboard) SetLeaderboardScores(ctx context.Context, leaderboardName string, params ...SetLeaderboardScoresParams) error {
	zParams := make([]redis.Z, len(params))

	for i, param := range params {
		zParams[i] = redis.Z{
			Member: param.Id,
			Score:  float64(param.Score),
		}
	}

	err := leaderboard.rdb.ZAdd(ctx, "leaderboard:"+leaderboardName, zParams...).Err()
	return err
}

func (leaderboard Leaderboard) GetLeaderboardRank(ctx context.Context, leaderboardName string, id string) (int, error) {
	res := leaderboard.rdb.ZRevRank(ctx, "leaderboard:"+leaderboardName, id)
	rank := int(res.Val()) + 1

	return rank, res.Err()
}

type SetRankDataParams struct {
	Id   string
	Data types.RankData
}

func (leaderboard Leaderboard) SetRankData(ctx context.Context, params ...SetRankDataParams) error {
	pipe := leaderboard.rdb.Pipeline()
	for _, param := range params {
		bytes, _ := json.Marshal(param.Data)
		pipe.Do(ctx, "JSON.SET", "rank:"+param.Id, "$", bytes)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (leaderboard Leaderboard) GetRankData(ctx context.Context, id string) (types.RankData, error) {
	var rankData types.RankData

	jsonRaw, err := leaderboard.rdb.JSONGet(ctx, "rank:"+id).Result()
	if err != nil {
		return rankData, err
	}

	err = json.Unmarshal([]byte(jsonRaw), &rankData)

	return rankData, err
}
