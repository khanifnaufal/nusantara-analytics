package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"indo-stats-backend/cache"
	"indo-stats-backend/models"

	"github.com/go-resty/resty/v2"
)

func FetchCommodities() (*models.CommoditiesResponse, bool, error) {
	cacheKey := "commodities:all"
	ttl := 15 * time.Minute

	if val, found := cache.Get(cacheKey); found {
		if resp, ok := val.(*models.CommoditiesResponse); ok {
			return resp, true, nil
		}
	}

	symbols := []string{"GC=F", "CL=F", "PALM.KL"}
	names := map[string]string{
		"GC=F":    "Emas",
		"CL=F":    "Minyak Mentah",
		"PALM.KL": "CPO",
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	quotes := make([]*models.Quote, len(symbols))

	client := resty.New()
	client.SetTimeout(10 * time.Second)
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	for i, symbol := range symbols {
		wg.Add(1)
		go func(idx int, sym string) {
			defer wg.Done()
			url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=7d", sym)

			var apiResp YahooFinanceResponse
			resp, err := client.R().
				SetResult(&apiResp).
				Get(url)

			if err != nil || resp.IsError() {
				return
			}

			if len(apiResp.Chart.Result) == 0 {
				return
			}

			res := apiResp.Chart.Result[0]
			meta := res.Meta

			history := make([]models.PriceHistory, 0)
			if len(res.Timestamp) > 0 && len(res.Indicators.Quote) > 0 && len(res.Indicators.Quote[0].Close) > 0 {
				closes := res.Indicators.Quote[0].Close
				for k, ts := range res.Timestamp {
					if k >= len(closes) {
						break
					}
					if closes[k] == nil {
						continue
					}
					history = append(history, models.PriceHistory{
						Timestamp: time.Unix(ts, 0),
						Close:     *closes[k],
					})
				}
			}

			quote := &models.Quote{
				Symbol:        sym,
				Name:          names[sym],
				Price:         meta.RegularMarketPrice,
				Change:        meta.RegularMarketChange,
				ChangePercent: meta.RegularMarketChangePercent,
				Currency:      meta.Currency,
				History:       history,
				LastUpdated:   time.Now(),
			}

			mu.Lock()
			quotes[idx] = quote
			mu.Unlock()
		}(i, symbol)
	}

	wg.Wait()

	var filteredQuotes []models.Quote
	for _, q := range quotes {
		if q != nil {
			filteredQuotes = append(filteredQuotes, *q)
		}
	}

	if len(filteredQuotes) == 0 {
		return nil, false, errors.New("failed to fetch commodities data for all symbols")
	}

	response := &models.CommoditiesResponse{
		Commodities: filteredQuotes,
		LastUpdated: time.Now(),
	}

	cache.Set(cacheKey, response, ttl)
	return response, false, nil
}
