package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
	"time"
)

func (service Service) SummonerDataCollectionLoop(ctx context.Context, region string) {
	backoffTicker := time.NewTicker(time.Second * 5)

	for range backoffTicker.C {
		puuids, _ := service.Queries.GetPuuidsWithNullSummonerData(ctx, db.GetPuuidsWithNullSummonerDataParams{
			Limit:  100,
			Region: region,
		})

		for _, puuid := range puuids {
			service.CollectSummonerDetails(ctx, region, puuid)
		}
	}
}

func (service Service) SummonerRegionCollectionLoop(ctx context.Context) {
	backoffTicker := time.NewTicker(time.Second * 5)

	for range backoffTicker.C {
		puuids, _ := service.Queries.GetPuuidsWithNullRegion(ctx, 100)

		for _, puuid := range puuids {
			service.CollectSummonerRegion(ctx, puuid)
		}
	}
}

func (service Service) AccountDataCollectionLoop(ctx context.Context, cluster string) {
	backoffTicker := time.NewTicker(time.Second * 5)

	for range backoffTicker.C {
		puuids, _ := service.Queries.GetPuuidsWithNullAccountData(ctx, 100)
		for _, puuid := range puuids {
			service.CollectAccountByPuuid(ctx, cluster, puuid)
		}
	}
}

func (service Service) MatchHistoryCollectionLoop(ctx context.Context, cluster string) {
	backoffTicker := time.NewTicker(time.Second * 5)

	for range backoffTicker.C {
		for _, region := range riot.ClusterToRegions[cluster] {
			rows, _ := service.Queries.GetOldestMatchHistories(ctx, db.GetOldestMatchHistoriesParams{
				Region: region,
				Limit:  100,
			})

			for _, row := range rows {
				service.CollectMatchHistory(ctx, cluster, row.Puuid, time.Now().Add(-time.Hour*24*3))
			}
		}

	}
}

func (service Service) RankEntryCollectionLoop(ctx context.Context) {
	for {
		for region, _ := range riot.RegionToCluster {
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
}
