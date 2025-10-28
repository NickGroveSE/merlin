package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OverwatchService struct{}

type OverwatchFilters struct {
	Role     string `json:"role"`
	Input    string `json:"input"`
	GameMode string `json:"gameMode"`
	RankTier string `json:"rankTier"`
	Map      string `json:"map"`
	Region   string `json:"region"`
}

type OWHero struct {
	Name     string
	Color    string
	Image    string
	PickRate float32
	WinRate  float32
}

type OWDataResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Data      []OWHero  `json:"data"`
}

func NewOverwatchService() *OverwatchService {
	return &OverwatchService{}
}

func (o *OverwatchService) Scrape(filters OverwatchFilters) ([]OWHero, error) {
	// Build the URL with query parameters
	url := o.buildURL(filters)

	fmt.Printf("Making request to: %s\n", url)

	// Make HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON response
	var response OWDataResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return response.Data, nil
}

func (o *OverwatchService) buildURL(filters OverwatchFilters) string {
	baseURL := "https://api.metatrack.ing/overwatch"

	// Build query parameters
	url := fmt.Sprintf("%s?role=%s&input=%s&gameMode=%s&map=%s&region=%s",
		baseURL,
		filters.Role,
		filters.Input,
		filters.GameMode,
		filters.Map,
		filters.Region,
	)

	// Add rankTier if it exists
	if filters.RankTier != "0" {
		url += fmt.Sprintf("&rankTier=%s", filters.RankTier)
	}

	return url
}
