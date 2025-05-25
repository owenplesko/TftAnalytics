package services

import (
	"context"

	"github.com/owenplesko/TftAnalytics/internal/db"
)

func (service *Service) GetOldestMatchHistories(ctx context.Context, region string) ([]db.GetOldestMatchesAfterRow, error) {
	return service.queries.GetOldestMatchesAfter(ctx, db.GetOldestMatchesAfterParams{
		Limit:  100,
		Region: region,
	})
}
