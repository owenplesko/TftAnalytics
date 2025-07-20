package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/owenplesko/TftAnalytics/internal/db"
)

func (controller Api) getSummonerMatches(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := r.PathValue("puuid")

	// validate query params
	var err error
	queryParams := r.URL.Query()

	limit := 20 // default value
	if limitParamArr, ok := queryParams["limit"]; ok {

		limit, err = strconv.Atoi(limitParamArr[0])

		if err != nil || limit < 0 {
			http.Error(w, "bad limit param", http.StatusBadRequest)
			return
		}
	}

	before := time.Now() // default value
	if beforeParamArr, ok := queryParams["before"]; ok {
		layout := "2006-01-02T15:04:05.999999Z"

		before, err = time.Parse(layout, beforeParamArr[0])

		if err != nil {
			http.Error(w, "bad before param", http.StatusBadRequest)
			return
		}
	}

	var queue int
	if queueParamArr, ok := queryParams["queue"]; ok {
		queue, err = strconv.Atoi(queueParamArr[0])

		if err != nil {
			http.Error(w, "bad queue param", http.StatusBadRequest)
			return
		}
	}

	// TODO: extract this constant somewhere
	setNumber := 14

	matches, err := controller.Queries.GetSummonerMatchHistory(ctx, db.GetSummonerMatchHistoryParams{
		SummonerPuuid: puuid,
		Limit:         int64(limit),
		Before: pgtype.Timestamp{
			Time:  before,
			Valid: true,
		},
		QueueID: pgtype.Int4{
			Int32: int32(queue),
			Valid: queue != 0,
		},
		SetNumber: int32(setNumber),
	})
	if err != nil {
		log.Printf("Queries.GetSummonerMatchHistory failed with err: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// prevent returning nil when list is empty
	if matches == nil {
		matches = []db.GetSummonerMatchHistoryRow{}
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(matches)
	w.Write(bytes)
}
