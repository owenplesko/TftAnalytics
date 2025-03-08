package main

import (
	"TFTAnalyticsServer/api"
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/leaderboard"
	"TFTAnalyticsServer/riot"
	"TFTAnalyticsServer/services"
	"net/http"
	"strconv"
	"time"

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
		Addr: os.Getenv("LEADERBOARD_URL"),
	})
	if err = rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	log.Println("Leaderboard connection successful")

	leaderboard := leaderboard.New(rdb)

	riotApiKey := os.Getenv("RIOT_KEY")

	rate, err := strconv.ParseFloat(os.Getenv("RIOT_RATE"), 32)
	if err != nil {
		panic("RIOT_RATE not set properly")
	}
	rateDuration := time.Duration(rate) * time.Millisecond
	riotClient := riot.New(riotApiKey, rateDuration)

	service := services.Service{
		Pool:        pool,
		Queries:     queries,
		Leaderboard: leaderboard,
		Riot:        riotClient,
	}

	for region := range riot.RegionToCluster {
		go service.RankEntryCollectionLoop(context.Background(), region)
		go service.MatchCollectionLoop(context.Background(), region)
	}

	service.Riot = riotClient.WithRequestPriority(1)
	apiEnv := api.Controller{
		Service: service,
	}
	http.ListenAndServe(":8080", apiEnv.New())
}
