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

func FetchMarket() (*models.MarketResponse, bool, error) {
	cacheKey := "market:all"
	ttl := 15 * time.Minute

	if val, found := cache.Get(cacheKey); found {
		if resp, ok := val.(*models.MarketResponse); ok {
			return resp, true, nil
		}
	}

	symbols := []string{"^JKSE", "BBCA.JK", "BBRI.JK", "TLKM.JK"}
	names := map[string]string{
		"^JKSE":   "IHSG",
		"BBCA.JK": "BCA",
		"BBRI.JK": "BRI",
		"TLKM.JK": "Telkom",
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

			// If Yahoo Finance returns 0 or no change, calculate it from history
			if (quote.Change == 0 || quote.ChangePercent == 0) && len(history) >= 1 {
				latest := quote.Price
				prev := history[len(history)-1].Close
				if len(history) >= 2 && latest == prev {
					prev = history[len(history)-2].Close
				}
				if prev != 0 {
					quote.Change = latest - prev
					quote.ChangePercent = (quote.Change / prev) * 100
				}
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
		return nil, false, errors.New("failed to fetch market data for all symbols")
	}

	response := &models.MarketResponse{
		Stocks:      filteredQuotes,
		LastUpdated: time.Now(),
	}

	cache.Set(cacheKey, response, ttl)
	return response, false, nil
}
