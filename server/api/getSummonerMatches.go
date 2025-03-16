package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (controller Controller) getSummonerMatches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := r.PathValue("puuid")

	// validate query params
	queryParams := r.URL.Query()

	limit := int64(20) // default value
	if limitParamArr, ok := queryParams["limit"]; ok {
		var err error

		limit, err = strconv.ParseInt(limitParamArr[0], 10, 32)

		if err != nil || limit < 0 {
			http.Error(w, "bad limit param", http.StatusBadRequest)
			return
		}
	}

	before := time.Now() // default value
	if beforeParamArr, ok := queryParams["before"]; ok {
		var err error
		layout := "2006-01-02T15:04:05.999999Z"

		before, err = time.Parse(layout, beforeParamArr[0])

		if err != nil {
			http.Error(w, "bad before param", http.StatusBadRequest)
			return
		}
	}

	matches, err := controller.Service.GetMatchHistory(ctx, puuid, int32(limit), before)
	if err != nil {
		log.Printf("Service.GetMatchHistory failed with err: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(matches)
	w.Write(bytes)
}
