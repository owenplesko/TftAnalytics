package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/owenplesko/TftAnalytics/internal/db"
)

func (controller Api) getSummonerStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := r.PathValue("puuid")

	// validate query params
	var err error
	queryParams := r.URL.Query()

	var setNumber int
	if paramArr, ok := queryParams["set"]; ok {
		setNumber, err = strconv.Atoi(paramArr[0])

		if err != nil {
			http.Error(w, "failed to parse set to int", http.StatusBadRequest)
			return
		}
	}

	var queueID int
	queueIdValid := false
	if paramArr, ok := queryParams["queue"]; ok {
		queueID, err = strconv.Atoi(paramArr[0])

		if err != nil {
			http.Error(w, "failed to parse queue to int", http.StatusBadRequest)
			return
		}

		queueIdValid = true
	}

	stats, err := controller.Queries.GetSummonerStats(ctx, db.GetSummonerStatsParams{
		SummonerPuuid: puuid,
		SetNumber:     int32(setNumber),
		QueueID: pgtype.Int4{
			Int32: int32(queueID),
			Valid: queueIdValid,
		},
	})
	if err == pgx.ErrNoRows {
		http.Error(w, "stats not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Queries.GetSummonerStats failed with err: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(stats)
	w.Write(bytes)
}
