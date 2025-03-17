package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/leaderboard"
	"TFTAnalyticsServer/riot"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool               *pgxpool.Pool
	queries            *db.Queries
	leaderboard        *leaderboard.Leaderboard
	riot               *riot.Riot
	serviceStatus      map[string]string
	matchesAfterCutoff time.Time
}

func New(pool *pgxpool.Pool, leaderboard *leaderboard.Leaderboard, riot *riot.Riot) *Service {
	return &Service{
		pool:               pool,
		queries:            db.New(pool),
		leaderboard:        leaderboard,
		riot:               riot,
		serviceStatus:      make(map[string]string),
		matchesAfterCutoff: time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC),
	}
}
