package models

import "time"

type MarketResponse struct {
	Stocks      []Quote   `json:"stocks"`
	LastUpdated time.Time `json:"last_updated"`
}
