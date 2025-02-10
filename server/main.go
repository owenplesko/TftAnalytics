package main

import (
	"TFTAnalyticsServer/api"
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"TFTAnalyticsServer/services"
	"net/http"

	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	} else {
		log.Println(".env file loaded")
	}

	pool, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	log.Println("Db connection successful")

	queries := db.New(pool)

	service := services.Service{
		Pool:    pool,
		Queries: queries,
	}

	go service.RankEntryCollectionLoop(context.Background())
	for region, _ := range riot.RegionToCluster {
		go service.SummonerDataCollectionLoop(context.Background(), region)
	}
	for cluster, _ := range riot.ClusterToRegions {
		go service.MatchHistoryCollectionLoop(context.Background(), cluster)
		go service.AccountDataCollectionLoop(context.Background(), cluster)
	}
	go service.SummonerRegionCollectionLoop(context.Background())

	apiEnv := api.Controller{
		Service: service,
	}
	http.ListenAndServe(":8080", apiEnv.New())
}
