package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/owenplesko/TftAnalytics/internal/db"
	"github.com/owenplesko/TftAnalytics/internal/services"
)

type Api struct {
	Service *services.Service
	Queries *db.Queries
}

func (controller Api) ListenAndServe(port int) {
	router := http.NewServeMux()
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}", controller.getSummonerByPuuid)
	router.HandleFunc("/v1/summoner/by-name-tag/{name}/{tag}", controller.getSummonerByNameTag)
	router.HandleFunc("/v1/leaderboard/{region}/{summonerId}", controller.getSummonerRank)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/update", controller.updateSummoner)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/matches", controller.getSummonerMatches)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/stats", controller.getSummonerStats)
	router.HandleFunc("/v1/match/{matchid}/comps", controller.getMatchComps)

	root := http.NewServeMux()
	root.Handle("/", setCors(setSchedulerPriority(router)))

	log.Printf("api listening on port: %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), root)
	log.Fatalf("api failed with err: %v", err)
}
