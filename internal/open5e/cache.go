package open5e

import "time"

type cacheEntry struct {
	data      MonstersResponse
	expiresAt time.Time
}

func (c *Client) getCached(url string) (MonstersResponse, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.cache[url]
	if !exists {
		return MonstersResponse{}, false
	}

	if time.Now().After(entry.expiresAt) {
		delete(c.cache, url)
		return MonstersResponse{}, false
	}

	return entry.data, true
}

func (c *Client) setCached(url string, data MonstersResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[url] = cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(c.cacheTTL),
	}
}

func (c *Client) reapExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()

			for url, entry := range c.cache {
				if time.Now().After(entry.expiresAt) {
					delete(c.cache, url)
				}
			}

			c.mu.Unlock()

		case <-c.stopReaper:
			return
		}
	}
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		close(c.stopReaper)
	})
}
