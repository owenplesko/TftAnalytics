package riot

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var key string
var limiters map[string]*time.Ticker

func init() {
	for region, cluster := range RegionToCluster {
		ClusterToRegions[cluster] = append(ClusterToRegions[cluster], region)
	}

	godotenv.Load()
	key = os.Getenv("RIOT_KEY")

	rate, err := strconv.ParseFloat(os.Getenv("RIOT_RATE"), 32)
	if err != nil {
		panic("RIOT_RATE not set properly")
	}
	rateDuration := time.Duration(rate) * time.Millisecond

	limiters = make(map[string]*time.Ticker)

	for region := range RegionToCluster {
		limiters[region] = time.NewTicker(rateDuration)
	}
	for cluster := range ClusterToRegions {
		limiters[cluster] = time.NewTicker(rateDuration)
	}
}
