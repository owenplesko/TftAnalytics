package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
	"log"
	"time"
)

func (service *Service) MatchCollectionLoop(ctx context.Context, region string) {
	backoffTime := time.NewTicker(time.Second * 5)

	for range backoffTime.C {
		summoners, err := service.Queries.GetOldestMatchesAfter(ctx, db.GetOldestMatchesAfterParams{
			Limit:  100,
			Region: region,
		})
		if err != nil {
			log.Printf("in MatchCollectionLoop Queries.GetOldestMatchesAfter failed with err: %v", err)
			continue
		}

		for _, summoner := range summoners {
			var after time.Time
			if summoner.MatchesAfterTimestamp.Valid {
				after = summoner.MatchesAfterTimestamp.Time
			} else {
				after = time.Now().UTC().Add(-time.Hour * 24 * 30)
			}
			err = service.CollectMatchHistory(ctx, region, summoner.Puuid, after)
			if err != nil {
				log.Printf("in MatchCollectionLoop service.CollectMatchHistory failed with err: %v", err)
			}
		}
	}
}

func (service *Service) RankEntryCollectionLoop(ctx context.Context, region string) {
	backoffTicker := time.NewTicker(time.Minute * 20)

	for range backoffTicker.C {
		for _, tier := range riot.ApexTiers {
			err := service.CollectApexRankEntries(ctx, region, tier)
			if err != nil {
				log.Printf("in RankEntryCollectionLoop CollectApexRankEntries failed with err: %v", err)
			}
		}
		for _, tier := range riot.Tiers {
			for _, division := range riot.Divisions {
				err := service.CollectRankEntries(ctx, region, tier, division)
				if err != nil {
					log.Printf("in RankEntryCollectionLoop CollectRankEntries failed with err: %v", err)
				}
			}
		}
	}
}
