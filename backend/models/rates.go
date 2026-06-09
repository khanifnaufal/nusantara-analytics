package models

import "time"

type RatesResponse struct {
	Base        string             `json:"base"`
	Date        string             `json:"date"`
	Rates       map[string]float64 `json:"rates"`
	LastUpdated time.Time          `json:"last_updated"`
	Stale       bool               `json:"stale,omitempty"`
}
