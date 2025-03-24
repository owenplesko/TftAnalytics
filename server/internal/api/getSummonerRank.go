package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func (controller Api) getSummonerRank(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	region := r.PathValue("region")
	summonerId := r.PathValue("summonerId")

	rankEntry, err := controller.Service.GetSummonerRank(ctx, region, summonerId)
	if errors.Is(err, redis.Nil) {
		http.Error(w, "rank not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println(err.Error())
		log.Printf("Service.GetSummonerRank failed with err: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(rankEntry)
	w.Write(bytes)
}
