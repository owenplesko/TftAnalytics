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

	return service.Queries.UpsertSummonerRank(ctx, db.UpsertSummonerRankParams{
		SummonerPuuid: rankEntry.Puuid,
		Tier:          rankEntry.Tier,
		Rank:          rankEntry.Rank,
		LeaguePoints:  rankEntry.LeaguePoints,
		Wins:          rankEntry.Wins,
		Losses:        rankEntry.Losses,
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

		// transform rank entries to upsert puuid params
		upsertPuuidParams := make([]db.BatchUpsertPuuidsParams, len(rankEntries))

		for i, rankEntry := range rankEntries {
			upsertPuuidParams[i] = db.BatchUpsertPuuidsParams{
				Puuid:  rankEntry.Puuid,
				Region: region,
			}
		}

		// batch insert puuids
		service.Queries.BatchUpsertPuuids(ctx, upsertPuuidParams).Exec(nil)

		// transform rank entries to upsert rank entry params
		upsertRankEntryParams := make([]db.BatchUpsertSummonerRankParams, len(rankEntries))

		for i, rankEntry := range rankEntries {
			upsertRankEntryParams[i] = db.BatchUpsertSummonerRankParams{
				SummonerPuuid: rankEntry.Puuid,
				Tier:          rankEntry.Tier,
				Rank:          rankEntry.Rank,
				LeaguePoints:  rankEntry.LeaguePoints,
				Wins:          rankEntry.Wins,
				Losses:        rankEntry.Losses,
			}
		}

		// batch upsert rank entries
		service.Queries.BatchUpsertSummonerRank(ctx, upsertRankEntryParams).Exec(nil)
	}

	log.Printf("Rank entries collected for %v %v %v\n", region, tier, division)

	return nil
}

func (service Service) CollectApexRankEntries(ctx context.Context, region, tier string) error {
	rankPage, err := riot.GetApexRankPage(region, tier)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// transform rank entries to upsert puuid params
	upsertPuuidParams := make([]db.BatchUpsertPuuidsParams, len(rankPage.Entries))

	for i, rankEntry := range rankPage.Entries {
		upsertPuuidParams[i] = db.BatchUpsertPuuidsParams{
			Puuid:  rankEntry.Puuid,
			Region: region,
		}
	}

	// batch insert puuids
	service.Queries.BatchUpsertPuuids(ctx, upsertPuuidParams).Exec(nil)

	// transform rank entries to upsert rank entry params
	upsertRankEntryParams := make([]db.BatchUpsertSummonerRankParams, len(rankPage.Entries))

	for i, rankEntry := range rankPage.Entries {
		upsertRankEntryParams[i] = db.BatchUpsertSummonerRankParams{
			SummonerPuuid: rankEntry.Puuid,
			Tier:          rankEntry.Tier,
			Rank:          rankEntry.Rank,
			LeaguePoints:  rankEntry.LeaguePoints,
			Wins:          rankEntry.Wins,
			Losses:        rankEntry.Losses,
		}
	}

	// batch upsert rank entries
	service.Queries.BatchUpsertSummonerRank(ctx, upsertRankEntryParams).Exec(nil)

	log.Printf("Rank entries collected for %v %v\n", region, tier)
	return nil
}
