package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
	"log"
)

func (service Service) GetSummonerRank(ctx context.Context, puuid string) (db.TftRank, error) {
	return service.Queries.GetSummonerRank(ctx, puuid)
}

func (service Service) CollectSummonerRank(ctx context.Context, region, summonerId string) error {
	rankEntry, err := riot.GetRank(region, summonerId)
	if err != nil {
		return err
	}

	return service.Queries.UpsertRank(ctx, db.UpsertRankParams{
		SummonerID:   rankEntry.SummonerId,
		Tier:         rankEntry.Tier,
		Rank:         rankEntry.Rank,
		LeaguePoints: rankEntry.LeaguePoints,
		Wins:         rankEntry.Wins,
		Losses:       rankEntry.Losses,
	})
}

func (service Service) CollectRankEntries(ctx context.Context, region, tier, division string) error {
	page := 0
	for {
		page++
		rankEntries, err := riot.GetRankEntries(region, tier, division, page)
		if err != nil {
			return err
		}
		if len(rankEntries) == 0 {
			break
		}

		// batch upsert puuids
		upsertPuuidParams := make([]db.BatchUpsertPuuidParams, len(rankEntries))

		for i, rankEntry := range rankEntries {
			upsertPuuidParams[i] = db.BatchUpsertPuuidParams{
				Puuid:  rankEntry.Puuid,
				Region: region,
			}
		}

		service.Queries.BatchUpsertPuuid(ctx, upsertPuuidParams).Exec(nil)

		// transform rank entries to upsert rank entry params
		upsertRankEntryParams := make([]db.BatchUpsertRankParams, len(rankEntries))

		for i, rankEntry := range rankEntries {
			upsertRankEntryParams[i] = db.BatchUpsertRankParams{
				SummonerID:   rankEntry.SummonerId,
				Tier:         rankEntry.Tier,
				Rank:         rankEntry.Rank,
				LeaguePoints: rankEntry.LeaguePoints,
				Wins:         rankEntry.Wins,
				Losses:       rankEntry.Losses,
			}
		}

		// batch upsert rank entries
		service.Queries.BatchUpsertRank(ctx, upsertRankEntryParams).Exec(nil)

		log.Printf("Rank entries collected for %v %v %v page %v\n", region, tier, division, page)
	}

	return nil
}

func (service Service) CollectApexRankEntries(ctx context.Context, region, tier string) error {
	rankPage, err := riot.GetApexRankPage(region, tier)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// transform rank entries to upsert rank entry params
	upsertRankEntryParams := make([]db.BatchUpsertRankParams, len(rankPage.Entries))

	for i, rankEntry := range rankPage.Entries {
		upsertRankEntryParams[i] = db.BatchUpsertRankParams{
			SummonerID:   rankEntry.SummonerId,
			Tier:         rankPage.Tier,
			Rank:         rankEntry.Rank,
			LeaguePoints: rankEntry.LeaguePoints,
			Wins:         rankEntry.Wins,
			Losses:       rankEntry.Losses,
		}
	}

	// batch upsert rank entries
	service.Queries.BatchUpsertRank(ctx, upsertRankEntryParams).Exec(nil)

	log.Printf("Rank entries collected for %v %v\n", region, tier)
	return nil
}
