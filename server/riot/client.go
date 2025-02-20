package riot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Riot struct {
	apiKey   string
	limiters map[string]*time.Ticker
}

func NewClient(apiKey string, rateDuration time.Duration) *Riot {
	return &Riot{
		apiKey:   apiKey,
		limiters: newLimiters(rateDuration),
	}
}

func newLimiters(rateDuration time.Duration) map[string]*time.Ticker {
	limiters := make(map[string]*time.Ticker)

	for region := range RegionToCluster {
		limiters[region] = time.NewTicker(rateDuration)
	}
	for cluster := range ClusterToRegions {
		limiters[cluster] = time.NewTicker(rateDuration)
	}

	return limiters
}

func (riot *Riot) request(server string, limitServer string, route string, target interface{}) error {
	url := fmt.Sprintf("https://%v.api.riotgames.com/%v", strings.ToLower(server), route)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", riot.apiKey)

	limiter, ok := riot.limiters[limitServer]
	if !ok {
		return fmt.Errorf("no limiter defined %v", server)
	}
	<-limiter.C

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return NotFoundError
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("Status: %v URL: %v\n", res.Status, url)
		return errors.New(res.Status)
	}

	return json.NewDecoder(res.Body).Decode(target)
}
