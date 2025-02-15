package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/leaderboard"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	Pool        *pgxpool.Pool
	Queries     *db.Queries
	Leaderboard *leaderboard.Leaderboard
}
