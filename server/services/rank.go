package services

import (
	"TFTAnalyticsServer/leaderboard"
	"TFTAnalyticsServer/riot"
	"TFTAnalyticsServer/types"
	"context"
	"log"
)

var tierScoreMap = map[string]int{
	"CHALLENGER":  2800,
	"GRANDMASTER": 2800,
	"MASTER":      2800,
	"DIAMOND":     2400,
	"EMERALD":     2000,
	"PLATINUM":    1600,
	"GOLD":        1200,
	"SILVER":      800,
	"BRONZE":      400,
	"IRON":        0,
}

var divisionScoreMap = map[string]int{
	"I":   300,
	"II":  200,
	"III": 100,
	"IV":  0,
}

func getRankScore(lp int, tier, division string) int {
	return lp + tierScoreMap[tier] + tierScoreMap[division]
}

func (service Service) GetSummonerRank(ctx context.Context, summonerId string) (types.RankData, error) {
	return service.Leaderboard.GetRankData(ctx, summonerId)
}

func (service Service) CollectSummonerRank(ctx context.Context, region, summonerId string) error {
	rankEntry, err := riot.GetRank(region, summonerId)
	if err != nil {
		return err
	}

	_ = service.Leaderboard.SetLeaderboardScores(ctx, "global", leaderboard.SetLeaderboardScoresParams{
		Id:    rankEntry.SummonerId,
		Score: getRankScore(int(rankEntry.LeaguePoints), rankEntry.Tier, rankEntry.Rank),
	})

	_ = service.Leaderboard.SetLeaderboardScores(ctx, region, leaderboard.SetLeaderboardScoresParams{
		Id:    rankEntry.SummonerId,
		Score: getRankScore(int(rankEntry.LeaguePoints), rankEntry.Tier, rankEntry.Rank),
	})

	err = service.Leaderboard.SetRankData(ctx, leaderboard.SetRankDataParams{
		Id: rankEntry.SummonerId,
		Data: types.RankData{
			Tier:         rankEntry.Tier,
			Rank:         rankEntry.Rank,
			LeaguePoints: rankEntry.LeaguePoints,
			Wins:         rankEntry.Wins,
			Losses:       rankEntry.Losses,
		},
	})

	return err
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

		// create datastore params
		//upsertPuuidParams := make([]db.BatchUpsertPuuidParams, len(rankEntries))
		setLeaderboardScoresParams := make([]leaderboard.SetLeaderboardScoresParams, len(rankEntries))
		setRankDataParams := make([]leaderboard.SetRankDataParams, len(rankEntries))

		for i, rankEntry := range rankEntries {
			//upsertPuuidParams[i] = db.BatchUpsertPuuidParams{
			//	Puuid:  rankEntry.Puuid,
			//	Region: region,
			//}

			setLeaderboardScoresParams[i] = leaderboard.SetLeaderboardScoresParams{
				Id:    rankEntry.SummonerId,
				Score: getRankScore(int(rankEntry.LeaguePoints), rankEntry.Tier, rankEntry.Rank),
			}

			setRankDataParams[i] = leaderboard.SetRankDataParams{
				Id: rankEntry.SummonerId,
				Data: types.RankData{
					Tier:         tier,
					Rank:         rankEntry.Rank,
					LeaguePoints: rankEntry.LeaguePoints,
					Wins:         rankEntry.Wins,
					Losses:       rankEntry.Losses,
				},
			}
		}

		//service.Queries.BatchUpsertPuuid(ctx, upsertPuuidParams).Exec(nil)
		_ = service.Leaderboard.SetLeaderboardScores(ctx, "global", setLeaderboardScoresParams...)
		_ = service.Leaderboard.SetLeaderboardScores(ctx, region, setLeaderboardScoresParams...)
		_ = service.Leaderboard.SetRankData(ctx, setRankDataParams...)

		//log.Printf("Rank entries collected for %v %v %v page %v\n", region, tier, division, page)
	}

	return nil
}

func (service Service) CollectApexRankEntries(ctx context.Context, region, tier string) error {
	rankPage, err := riot.GetApexRankPage(region, tier)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// create datastore params
	setLeaderboardScoresParams := make([]leaderboard.SetLeaderboardScoresParams, len(rankPage.Entries))
	setRankDataParams := make([]leaderboard.SetRankDataParams, len(rankPage.Entries))

	for i, rankEntry := range rankPage.Entries {
		setLeaderboardScoresParams[i] = leaderboard.SetLeaderboardScoresParams{
			Id:    rankEntry.SummonerId,
			Score: getRankScore(int(rankEntry.LeaguePoints), rankPage.Tier, rankEntry.Rank),
		}

		setRankDataParams[i] = leaderboard.SetRankDataParams{
			Id: rankEntry.SummonerId,
			Data: types.RankData{
				Tier:         tier,
				Rank:         rankEntry.Rank,
				LeaguePoints: rankEntry.LeaguePoints,
				Wins:         rankEntry.Wins,
				Losses:       rankEntry.Losses,
			},
		}
	}

	// batch upsert rank entries
	_ = service.Leaderboard.SetLeaderboardScores(ctx, "global", setLeaderboardScoresParams...)
	_ = service.Leaderboard.SetLeaderboardScores(ctx, region, setLeaderboardScoresParams...)
	_ = service.Leaderboard.SetRankData(ctx, setRankDataParams...)

	//log.Printf("Rank entries collected for %v %v\n", region, tier)
	return nil
}
