package main

import (
	"TFTAnalyticsServer/api"
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/leaderboard"
	"TFTAnalyticsServer/riot"
	"TFTAnalyticsServer/services"
	"net/http"

	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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
	if err = pool.Ping(context.Background()); err != nil {
		panic(err)
	}
	log.Println("Db connection successful")

	queries := db.New(pool)

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if err = rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	log.Println("Leaderboard connection successful")

	leaderboard := leaderboard.New(rdb)

	service := services.Service{
		Pool:        pool,
		Queries:     queries,
		Leaderboard: leaderboard,
	}

	for region, _ := range riot.RegionToCluster {
		go service.RankEntryCollectionLoop(context.Background(), region)
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
