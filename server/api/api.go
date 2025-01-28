package api

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"TFTAnalyticsServer/services"
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type ApiEnv struct {
	ServiceEnv services.ServiceEnv
}

func (env ApiEnv) New() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}", env.getSummonerByPuuid)
	router.HandleFunc("/v1/summoner/by-name-tag/{name}/{tag}", env.getSummonerByNameTag)
	router.HandleFunc("/v1/summoner/by-puuid/{puuid}/matches", env.getSummonerMatches)
	return router
}

func (env ApiEnv) getSummonerByPuuid(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")

	puuid := r.PathValue("puuid")

	summoner, err := env.ServiceEnv.GetSummonerByPuuid(ctx, puuid)
	if err == pgx.ErrNoRows {
		http.Error(w, "summoner not found", 404)
		return
	}
	if err != nil {
		http.Error(w, "something went wrong", 500)
		return
	}

	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}

func (env ApiEnv) getSummonerByNameTag(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")

	cluster := "americas"
	name := r.PathValue("name")
	tag := r.PathValue("tag")

	summoner, err := env.ServiceEnv.GetOrCollectSummonerByNameTag(ctx, cluster, name, tag)
	if err == riot.NotFoundError {
		http.Error(w, "summoner not found", 404)
		return
	}
	if err != nil {
		http.Error(w, "something went wrong", 500)
		return
	}

	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}

func (env ApiEnv) getSummonerMatches(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")

	puuid := r.PathValue("puuid")

	matches, err := env.ServiceEnv.GetMatchHistory(ctx, puuid)
	if err != nil {
		http.Error(w, "something went wrong", 500)
		return
	}

	// prevent returning null when list is empty
	if matches == nil {
		matches = []db.SummonerMatchHistoryRow{}
	}

	bytes, _ := json.Marshal(matches)
	w.Write(bytes)
}
