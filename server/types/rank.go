package types

type Rank struct {
	Data                 RankData              `json:"rankData"`
	LeaderboardPositions []LeaderboardPosition `json:"leaderboardPositions"`
}

type LeaderboardPosition struct {
	Leaderboard string `json:"leaderboard"`
	Position    int    `json:"position"`
}

type RankData struct {
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
}
