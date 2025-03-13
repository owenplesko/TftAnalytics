package riot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Riot struct {
	apiKey      string
	rateLimiter RateLimiter
}

func New(apiKey string, rateDuration time.Duration) *Riot {
	return &Riot{
		apiKey:      apiKey,
		rateLimiter: newRateLimiter(rateDuration),
	}
}

func (riot *Riot) request(ctx context.Context, server string, route string, target any) error {
	url := fmt.Sprintf("https://%v.api.riotgames.com/%v", strings.ToLower(server), route)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", riot.apiKey)

	err := riot.rateLimiter.Wait(ctx, server)
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
