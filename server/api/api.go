package api

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/services"
	"context"
	"encoding/json"
	"net/http"
)

type ApiEnv struct {
	ServiceEnv services.ServiceEnv
}

func (env ApiEnv) New() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/summoner/{puuid}", env.getSummonerByPuuid)
	router.HandleFunc("/account/{cluster}/{name}/{tag}", env.getOrCollectSummonerByNameTag)
	router.HandleFunc("/summoner/{puuid}/matches", env.getSummonerMatches)
	return router
}

func (env ApiEnv) getSummonerByPuuid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	puuid := r.PathValue("puuid")
	ctx := context.Background()

	summoner, err := env.ServiceEnv.GetSummonerByPuuid(ctx, puuid)
	if err != nil {
		http.Error(w, "summoner not found", 404)
		return
	}

	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}

func (env ApiEnv) getOrCollectSummonerByNameTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := context.Background()
	cluster := r.PathValue("cluster")
	name := r.PathValue("name")
	tag := r.PathValue("tag")

	summoner, err := env.ServiceEnv.GetOrCollectSummonerByNameTag(ctx, cluster, name, tag)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}

func (env ApiEnv) getSummonerMatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ctx := context.Background()
	puuid := r.PathValue("puuid")

	matches, err := env.ServiceEnv.GetMatchHistory(ctx, puuid)
	if err != nil {
		http.Error(w, "summoner not found", 404)
		return
	}

	// prevent returning null when list is empty
	if matches == nil {
		matches = []db.SummonerMatchHistoryRow{}
	}

	bytes, _ := json.Marshal(matches)
	w.Write(bytes)
}
