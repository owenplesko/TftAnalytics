package api

import (
	"TFTAnalyticsServer/services"
	"net/http"
)

type Controller struct {
	Service services.Service
}

func (controller Controller) New() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}", controller.getSummonerByPuuid)
	router.HandleFunc("/v1/summoner/by-name-tag/{name}/{tag}", controller.getSummonerByNameTag)
	router.HandleFunc("/v1/summoner/by-summoner-id/{puuid}/rank", controller.getSummonerRank)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/update", controller.updateSummoner)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/matches", controller.getSummonerMatches)
	router.HandleFunc("/v1/match/{matchid}/comps", controller.getMatchComps)
	return router
}
