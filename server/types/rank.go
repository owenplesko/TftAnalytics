package types

type Rank struct {
	Data                RankData            `json:"rankData"`
	LeaderboardPosition LeaderboardPosition `json:"leaderboardPosition"`
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
