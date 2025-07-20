package services

import (
	"context"

	"github.com/owenplesko/TftAnalytics/internal/db"
)


func (service *Service) GetUnitPlacments(ctx context.Context) ([]db.UnitPlacement, error) {
	return service.queries.GetUnitPlacements(ctx)
}

