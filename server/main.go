package main

import (
	"TFTAnalyticsServer/api"
	"TFTAnalyticsServer/db"
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

	serviceEnv := services.ServiceEnv{
		Pool:    pool,
		Queries: queries,
	}

	go serviceEnv.ClusterCollectionLoop(context.Background(), "americas")
	go serviceEnv.RegionCollectionLoop(context.Background(), "na1")

	apiEnv := api.ApiEnv{
		ServiceEnv: serviceEnv,
	}
	http.ListenAndServe(":8080", apiEnv.New())
}
