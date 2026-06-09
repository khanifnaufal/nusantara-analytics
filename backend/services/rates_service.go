package services

import (
	"errors"
	"time"

	"indo-stats-backend/cache"
	"indo-stats-backend/models"

	"github.com/go-resty/resty/v2"
)

type FrankfurterResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

func FetchRates() (*models.RatesResponse, bool, error) {
	cacheKey := "rates:latest"
	ttl := 1 * time.Hour
	staleTtl := 24 * time.Hour // Keep as stale fallback

	if val, found := cache.Get(cacheKey); found {
		if resp, ok := val.(*models.RatesResponse); ok {
			// Fresh hit
			if time.Since(resp.LastUpdated) < ttl {
				return resp, true, nil
			}

			// Stale hit, try to fetch fresh
			freshResp, err := fetchRatesFromAPI()
			if err == nil {
				cache.Set(cacheKey, freshResp, staleTtl)
				return freshResp, false, nil
			}

			// Fetch failed, return stale fallback with stale = true
			staleResp := *resp
			staleResp.Stale = true
			return &staleResp, true, nil
		}
	}

	// No cache
	freshResp, err := fetchRatesFromAPI()
	if err != nil {
		return nil, false, err
	}

	cache.Set(cacheKey, freshResp, staleTtl)
	return freshResp, false, nil
}

func fetchRatesFromAPI() (*models.RatesResponse, error) {
	client := resty.New()
	client.SetTimeout(10 * time.Second)
	var apiResp FrankfurterResponse
	resp, err := client.R().
		SetResult(&apiResp).
		Get("https://api.frankfurter.app/latest?from=IDR&to=USD,EUR,SGD,JPY,MYR,SAR,AUD")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New("failed to fetch rates from API: " + resp.Status())
	}

	ratesResp := &models.RatesResponse{
		Base:        apiResp.Base,
		Date:        apiResp.Date,
		Rates:       apiResp.Rates,
		LastUpdated: time.Now(),
	}
	return ratesResp, nil
}
