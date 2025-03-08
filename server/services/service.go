package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/leaderboard"
	"TFTAnalyticsServer/riot"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool        *pgxpool.Pool
	queries     *db.Queries
	leaderboard *leaderboard.Leaderboard
	riot        *riot.Riot
}

func New(pool *pgxpool.Pool, leaderboard *leaderboard.Leaderboard, riot *riot.Riot) *Service {
	return &Service{
		pool:        pool,
		queries:     db.New(pool),
		leaderboard: leaderboard,
		riot:        riot,
	}
}

func (service *Service) WithRiotRequestPriority(priority int) *Service {
	return &Service{
		pool:        service.pool,
		queries:     service.queries,
		leaderboard: service.leaderboard,
		riot:        service.riot.WithRequestPriority(priority),
	}
}
