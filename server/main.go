package main

import (
	"TFTAnalyticsServer/internal/api"
	"TFTAnalyticsServer/internal/db"
	"TFTAnalyticsServer/internal/leaderboard"
	"TFTAnalyticsServer/internal/services"
	"TFTAnalyticsServer/pkg/riot"
	"strconv"
	"time"

	"context"
	"embed"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

//go:embed sql/migrations/*.sql
var embedMigrations embed.FS

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found")
	} else {
		log.Println(".env file loaded")
	}

	// create database connection
	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	if err = pool.Ping(context.Background()); err != nil {
		panic(err)
	}
	log.Println("DB connection successful")

	// run database migrations
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	dbConn := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(dbConn, "sql/migrations"); err != nil {
		panic(err)
	}
	if err := dbConn.Close(); err != nil {
		panic(err)
	}
	log.Println("DB migrations successful")

	// create leaderboard connection
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("LEADERBOARD_URL"),
	})
	if err = rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	leaderboard := leaderboard.New(rdb)
	log.Println("Leaderboard connection successful")

	// create riot client
	rate, err := strconv.ParseFloat(os.Getenv("RIOT_RATE"), 32)
	if err != nil {
		panic("RIOT_RATE not set properly")
	}
	rateDuration := time.Duration(rate) * time.Millisecond
	riotClient := riot.New(os.Getenv("RIOT_KEY"), rateDuration)

	// create service
	service := services.New(pool, leaderboard, riotClient)

	// start collection loops
	for region := range riot.RegionToCluster {
		go service.RankEntryCollectionLoop(context.Background(), region)
		go service.MatchCollectionLoop(context.Background(), region)
	}

	// start api
	apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		apiPort = 9001
	}
	api.Api{Service: service, Queries: db.New(pool)}.ListenAndServe(apiPort)
}
