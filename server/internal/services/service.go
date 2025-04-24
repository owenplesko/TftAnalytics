package services

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/owenplesko/TftAnalytics/internal/db"
	"github.com/owenplesko/TftAnalytics/internal/leaderboard"
	"github.com/owenplesko/TftAnalytics/internal/services/internal/dedupe"
	"github.com/owenplesko/TftAnalytics/pkg/riot"
)

type Service struct {
	pool               *pgxpool.Pool
	queries            *db.Queries
	leaderboard        *leaderboard.Leaderboard
	riot               *riot.Riot
	deduplicator       *dedupe.Deduplicator
	matchesAfterCutoff time.Time
}

func New(pool *pgxpool.Pool, leaderboard *leaderboard.Leaderboard, riot *riot.Riot) *Service {
	return &Service{
		pool:               pool,
		queries:            db.New(pool),
		leaderboard:        leaderboard,
		riot:               riot,
		deduplicator:       dedupe.New(),
		matchesAfterCutoff: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
	}
}
