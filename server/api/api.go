package api

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"TFTAnalyticsServer/services"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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
	router.HandleFunc("/v1/match/{matchid}/comps", env.getMatchComps)
	return router
}

func (env ApiEnv) getSummonerByPuuid(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	puuid := r.PathValue("puuid")

	summoner, err := env.ServiceEnv.GetSummonerByPuuid(ctx, puuid)
	if err == pgx.ErrNoRows {
		http.Error(w, "summoner not found", 404)
		return
	}
	if err != nil {
		http.Error(w, "the database is sad :(", 500)
		return
	}

	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}

func (env ApiEnv) getSummonerByNameTag(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	cluster := "americas"
	name := r.PathValue("name")
	tag := r.PathValue("tag")

	summoner, err := env.ServiceEnv.GetOrCollectSummonerByNameTag(ctx, cluster, name, tag)
	if err == riot.NotFoundError {
		http.Error(w, "summoner not found", 404)
		return
	}
	if err != nil {
		http.Error(w, "the database is sad :(", 500)
		return
	}

	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}

func (env ApiEnv) getSummonerMatches(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	puuid := r.PathValue("puuid")

	// validate query params
	queryParams := r.URL.Query()

	limit := int64(20) // default value
	if limitParamArr, ok := queryParams["limit"]; ok {
		var err error

		limit, err = strconv.ParseInt(limitParamArr[0], 10, 32)

		if err != nil || limit < 0 {
			http.Error(w, "bad limit param", 400)
			return
		}
	}

	after := time.Now() // default value
	if afterParamArr, ok := queryParams["after"]; ok {
		var err error
		layout := "2006-01-02 15:04:05.999999"

		after, err = time.Parse(layout, afterParamArr[0])

		if err != nil {
			http.Error(w, "bad after param", 400)
			return
		}
	}

	matches, err := env.ServiceEnv.GetMatchHistory(ctx, puuid, int32(limit), after)
	if err != nil {
		http.Error(w, "the database is sad :(", 500)
		return
	}

	// prevent returning null when list is empty
	if matches == nil {
		matches = []db.SummonerMatchHistoryRow{}
	}

	bytes, _ := json.Marshal(matches)
	w.Write(bytes)
}

func (env ApiEnv) getMatchComps(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	matchId := r.PathValue("matchid")

	comps, err := env.ServiceEnv.GetMatchComps(ctx, matchId)
	if err != nil {
		http.Error(w, "the database is sad :(", 500)
		return
	}

	// prevent returning null when list is empty
	if comps == nil {
		comps = []db.GetMatchCompsRow{}
	}

	bytes, _ := json.Marshal(comps)
	w.Write(bytes)
}
