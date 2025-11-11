package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/pkg/riot"
)

func (service *Service) CollectAllRankEntries(ctx context.Context, region string) error {
	key := fmt.Sprintf("ALL_RANK_ENTRIES_%s", region)

	_, err, _ := service.group.Do(key, func() (interface{}, error) {
		for _, tier := range riot.ApexTiers {
			err := service.CollectApexRankEntries(ctx, region, tier)
			if err != nil {
				log.Printf("CollectApexRankEntries failed with err: %v", err)
			}
		}
		for _, tier := range riot.Tiers {
			for _, division := range riot.Divisions {
				err := service.CollectRankEntries(ctx, region, tier, division)
				if err != nil {
					log.Printf("CollectRankEntries failed with err: %v", err)
				}
			}
		}

		return nil, nil
	})

	return err
}
