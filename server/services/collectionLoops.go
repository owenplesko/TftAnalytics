package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
	"log"
	"time"
)

func (service Service) MatchCollectionLoop(ctx context.Context, region string) {
	backoffTime := time.NewTicker(time.Second * 5)

	for range backoffTime.C {
		summoners, err := service.Queries.GetOldestMatchesAfter(ctx, db.GetOldestMatchesAfterParams{
			Limit:  100,
			Region: region,
		})
		if err != nil {
			log.Println(err.Error())
			continue
		}

		for _, summoner := range summoners {
			_ = service.CollectMatchHistory(ctx, region, summoner.Puuid, summoner.MatchesAfterTimestamp.Time)
		}
	}
}

func (service Service) RankEntryCollectionLoop(ctx context.Context, region string) {
	//backoffTicker := time.NewTicker(time.Minute * 20)

	for /*range backoffTicker.C*/ {
		for _, tier := range riot.ApexTiers {
			service.CollectApexRankEntries(ctx, region, tier)
		}
		for _, tier := range riot.Tiers {
			for _, division := range riot.Divisions {
				service.CollectRankEntries(ctx, region, tier, division)
			}
		}
	}
}
