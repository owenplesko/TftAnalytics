package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (env Controller) getSummonerRank(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")

	region := r.PathValue("region")
	summonerId := r.PathValue("summonerId")

	rankEntry, err := env.Service.GetSummonerRank(ctx, region, summonerId)
	if err == pgx.ErrNoRows {
		http.Error(w, "rank not found", 404)
		return
	}
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(rankEntry)
	w.Write(bytes)
}
