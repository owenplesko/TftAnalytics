package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (controller Controller) getSummonerRank(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	region := r.PathValue("region")
	summonerId := r.PathValue("summonerId")

	rankEntry, err := controller.Service.GetSummonerRank(ctx, region, summonerId)
	if err == pgx.ErrNoRows {
		http.Error(w, "rank not found", 404)
		return
	}
	if err != nil {
		log.Println(err.Error())
		log.Printf("Service.GetSummonerRank failed with err: %v\n", err)
		http.Error(w, "something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(rankEntry)
	w.Write(bytes)
}
