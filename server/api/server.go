package api

import (
	"TFTAnalyticsServer/riot/scheduler"
	"TFTAnalyticsServer/services"
	"net/http"
)

type Controller struct {
	Service *services.Service
}

func setSchedulerPriority(next http.Handler, priority int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = scheduler.WithPriority(ctx, priority)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (controller Controller) New() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}", controller.getSummonerByPuuid)
	router.HandleFunc("/v1/summoner/by-name-tag/{name}/{tag}", controller.getSummonerByNameTag)
	router.HandleFunc("/v1/leaderboard/{region}/{summonerId}", controller.getSummonerRank)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/update", controller.updateSummoner)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/matches", controller.getSummonerMatches)
	router.HandleFunc("/v1/match/{matchid}/comps", controller.getMatchComps)

	root := http.NewServeMux()
	root.Handle("/", setSchedulerPriority(router, 1))

	return root
}
