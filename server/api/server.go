package api

import (
	"TFTAnalyticsServer/services"
	"net/http"
)

type ApiEnv struct {
	ServiceEnv services.ServiceEnv
}

func (env ApiEnv) New() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}", env.getSummonerByPuuid)
	router.HandleFunc("/v1/summoner/by-name-tag/{name}/{tag}", env.getSummonerByNameTag)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/update", env.updateSummoner)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/matches", env.getSummonerMatches)
	router.HandleFunc("/v1/match/{matchid}/comps", env.getMatchComps)
	return router
}
