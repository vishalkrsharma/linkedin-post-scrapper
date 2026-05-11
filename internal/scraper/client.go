package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"linkedin-hunter/internal/jobs"
)

type Client interface {
	Scrape(keyword, location string) (jobs.ScraperResult, error)
}

type client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) Client {
	return &client{
		baseURL: baseURL,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

type ScrapeRequest struct {
	Keyword  string `json:"keyword"`
	Location string `json:"location"`
}

func (c *client) Scrape(keyword, location string) (jobs.ScraperResult, error) {
	url := c.baseURL + "/scrape"

	body := ScrapeRequest{Keyword: keyword, Location: location}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return jobs.ScraperResult{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return jobs.ScraperResult{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return jobs.ScraperResult{}, fmt.Errorf("failed to call scraper: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return jobs.ScraperResult{}, fmt.Errorf("scraper returned status %d", resp.StatusCode)
	}

	var result struct {
		Jobs  []jobs.JobPost `json:"jobs"`
		Count int            `json:"count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return jobs.ScraperResult{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return jobs.ScraperResult{Jobs: result.Jobs, Count: result.Count}, nil
}