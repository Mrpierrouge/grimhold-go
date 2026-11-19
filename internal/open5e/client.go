package open5e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const baseURL = "https://api.open5e.com"
const monstersEndpoint = "/v2/creatures"

type Client struct {
	HTTPClient *http.Client

	mu         sync.Mutex
	cache      map[string]cacheEntry
	cacheTTL   time.Duration
	stopReaper chan struct{}
	closeOnce  sync.Once
}

func NewClient() *Client {
	c := &Client{
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		cache:      make(map[string]cacheEntry),
		cacheTTL:   5 * time.Minute,
		stopReaper: make(chan struct{}),
	}

	go c.reapExpired()

	return c
}

func (c *Client) GetMonsters(pageURL string) (MonstersResponse, error) {
	if pageURL == "" {
		pageURL = baseURL + monstersEndpoint + "?limit=10"
	}

	if data, found := c.getCached(pageURL); found {
		fmt.Println("[CACHE] Returning cached page")
		return data, nil
	}

	var monstersResponse MonstersResponse

	resp, err := c.HTTPClient.Get(pageURL)
	if err != nil {
		return MonstersResponse{}, fmt.Errorf("fetching monsters: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return MonstersResponse{}, fmt.Errorf("Open5e returned status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&monstersResponse); err != nil {
		return MonstersResponse{}, fmt.Errorf("decoding monsters: %w", err)
	}
	c.setCached(pageURL, monstersResponse)

	return monstersResponse, nil
}

