package riot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Riot struct {
	apiKey          string
	rateLimiter     RateLimiter
	requestPriority int
}

func New(apiKey string, rateDuration time.Duration) *Riot {
	return &Riot{
		apiKey:          apiKey,
		rateLimiter:     newRateLimiter(rateDuration),
		requestPriority: 0,
	}
}

func (r *Riot) WithRequestPriority(priority int) *Riot {
	return &Riot{
		apiKey:          r.apiKey,
		rateLimiter:     r.rateLimiter,
		requestPriority: priority,
	}
}

func (riot *Riot) request(server string, route string, target any) error {
	url := fmt.Sprintf("https://%v.api.riotgames.com/%v", strings.ToLower(server), route)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Riot-Token", riot.apiKey)

	err := riot.rateLimiter.Wait(server, riot.requestPriority)
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
