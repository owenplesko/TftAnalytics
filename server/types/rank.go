package types

type RankData struct {
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int32  `json:"leaguePoints"`
	Wins         int32  `json:"wins"`
	Losses       int32  `json:"losses"`
}
