package services

import "context"

func (service *Service) RefreshUnitPlacement(ctx context.Context) error {
	return service.queries.RefreshUnitPlacement(ctx)
}
