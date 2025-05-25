package main

import (
	"strconv"
	"time"

	"context"
	"embed"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/owenplesko/TftAnalytics/internal/api"
	"github.com/owenplesko/TftAnalytics/internal/db"
	"github.com/owenplesko/TftAnalytics/internal/leaderboard"
	"github.com/owenplesko/TftAnalytics/internal/schedule"
	"github.com/owenplesko/TftAnalytics/internal/services"
	"github.com/owenplesko/TftAnalytics/pkg/riot"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
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
	riotRate, err := strconv.ParseFloat(os.Getenv("RIOT_RATE"), 32)
	if err != nil {
		panic("RIOT_RATE not set properly")
	}
	limit := rate.Every(time.Duration(riotRate) * time.Millisecond)
	riotClient := riot.New(os.Getenv("RIOT_KEY"), limit)

	// create service
	service := services.New(pool, leaderboard, riotClient)

	schedule := schedule.New(service)
	schedule.Start()

	// start collection loops
	for region := range riot.RegionToCluster {
		go service.MatchCollectionLoop(context.Background(), region)
	}

	// start api
	apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		apiPort = 9001
	}
	api.Api{Service: service, Queries: db.New(pool)}.ListenAndServe(apiPort)
}
