package models

import "time"

type PriceHistory struct {
	Timestamp time.Time `json:"timestamp"`
	Close     float64   `json:"close"`
}

type Quote struct {
	Symbol        string         `json:"symbol"`
	Name          string         `json:"name"`
	Price         float64        `json:"price"`
	Change        float64        `json:"change"`
	ChangePercent float64        `json:"change_percent"`
	Currency      string         `json:"currency"`
	History       []PriceHistory `json:"history"`
	LastUpdated   time.Time      `json:"last_updated"`
}

type CommoditiesResponse struct {
	Commodities []Quote   `json:"commodities"`
	LastUpdated time.Time `json:"last_updated"`
}
