package services

import (
	"TFTAnalyticsServer/internal/db"
	"TFTAnalyticsServer/internal/leaderboard"
	"TFTAnalyticsServer/pkg/riot"
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
		matchesAfterCutoff: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
	}
}
