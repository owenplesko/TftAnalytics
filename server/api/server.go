package api

import (
	"TFTAnalyticsServer/services"
	"net/http"
)

type Controller struct {
	Service *services.Service
}

func New(controller Controller) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}", controller.getSummonerByPuuid)
	router.HandleFunc("/v1/summoner/by-name-tag/{name}/{tag}", controller.getSummonerByNameTag)
	router.HandleFunc("/v1/leaderboard/{region}/{summonerId}", controller.getSummonerRank)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/update", controller.updateSummoner)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/matches", controller.getSummonerMatches)
	router.HandleFunc("/v1/match/{matchid}/comps", controller.getMatchComps)
	return router
}
