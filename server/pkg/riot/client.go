package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
)

type Riot struct {
	apiKey   string
	limiters map[string]*rate.Limiter
}

func New(apiKey string, limit rate.Limit) *Riot {
	limiters := make(map[string]*rate.Limiter)

	for cluster := range ClusterToRegions {
		limiters[cluster] = rate.NewLimiter(limit, 1)
	}
	for region := range RegionToCluster {
		limiters[region] = rate.NewLimiter(limit, 1)
	}

	return &Riot{
		apiKey:   apiKey,
		limiters: limiters,
	}
}

func (riot *Riot) request(ctx context.Context, server string, route string, target any) error {
	url := fmt.Sprintf("https://%v.api.riotgames.com/%v", strings.ToLower(server), route)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", riot.apiKey)

	limiter, ok := riot.limiters[server]
	if !ok {
		return fmt.Errorf("no limiter found for server %s", server)
	}

	err := limiter.Wait(ctx)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to %v: %w", url, err)
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Printf("error closing response body of request to %v: %v", url, err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusNotFound:
			return fmt.Errorf("failed to get resource from %v: %w", url, ErrNotFound)
		default:
			return fmt.Errorf("failed to get resource from %v: status code %v", url, res.StatusCode)
		}
	}

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		return fmt.Errorf("error decoding body of request to %v: %w", url, err)
	}

	return nil
}
