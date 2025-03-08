package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/leaderboard"
	"TFTAnalyticsServer/riot"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	Pool        *pgxpool.Pool
	Queries     *db.Queries
	Leaderboard *leaderboard.Leaderboard
	Riot        *riot.Riot
}

func New(pool *pgxpool.Pool, leaderboard *leaderboard.Leaderboard, riot *riot.Riot) *Service {
	return &Service{
		Pool:        pool,
		Queries:     db.New(pool),
		Leaderboard: leaderboard,
		Riot:        riot,
	}
}
