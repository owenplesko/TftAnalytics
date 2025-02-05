package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func (env ApiEnv) getSummonerMatches(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")

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
		http.Error(w, "something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(matches)
	w.Write(bytes)
}
