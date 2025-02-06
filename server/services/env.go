package services

import (
	"TFTAnalyticsServer/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}
