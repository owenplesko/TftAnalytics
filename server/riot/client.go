package riot

import (
	"TFTAnalyticsServer/riot/limiter"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Riot struct {
	apiKey          string
	rateLimiter     RateLimiter
	requestPriority int
}

func New(apiKey string, rateLimiter RateLimiter, requestPriority int) *Riot {
	return &Riot{
		apiKey:          apiKey,
		rateLimiter:     rateLimiter,
		requestPriority: requestPriority,
	}
}

func (riot *Riot) request(server string, route string, target interface{}) error {
	url := fmt.Sprintf("https://%v.api.riotgames.com/%v", strings.ToLower(server), route)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", riot.apiKey)

	limiter, err := riot.getLimiter(server)
	if err != nil {
		return err
	}
	<-limiter.Wait(riot.requestPriority)

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

func (riot *Riot) getLimiter(server string) (*limiter.Limiter, error) {
	cluster, ok := RegionToCluster[server]
	if !ok {
		cluster = server
	}

	limiter, ok := riot.rateLimiter[cluster]
	if !ok {
		return limiter, NoLimiterError
	}

	return limiter, nil
}
