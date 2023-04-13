package api

import (
	"TheCollectorDG/collection/matchCollection"
	"TheCollectorDG/collection/summonerCollection"
	"os"

	"github.com/gofiber/fiber/v2"
)

var summonerCollectionRegionRouter map[string]*summonerCollection.RegionalSummonerCollectionQueue
var matchCollectionRegionRouter map[string]*matchCollection.RegionalMatchCollectionQueue

func Setup(summonerCQR map[string]*summonerCollection.RegionalSummonerCollectionQueue,
	matchCQR map[string]*matchCollection.RegionalMatchCollectionQueue) {
	summonerCollectionRegionRouter = summonerCQR
	matchCollectionRegionRouter = matchCQR
}

func Start() {
	app := fiber.New()

	app.Get("summoner/:region/:name", GetSummoner)
	app.Get("comps/:puuid", GetCompHistory)
	app.Get("update/profile/:puuid", UpdateProfile)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}
	app.Listen("0.0.0.0:" + port)
}
