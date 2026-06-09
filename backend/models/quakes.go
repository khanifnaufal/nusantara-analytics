package models

import "time"

type Quake struct {
	ID             string    `json:"id"`
	Place          string    `json:"place"`
	Magnitude      float64   `json:"magnitude"`
	Depth          float64   `json:"depth"`
	Lat            float64   `json:"lat"`
	Lon            float64   `json:"lon"`
	Time           time.Time `json:"time"`
	Alert          string    `json:"alert"`
	Tsunami        int       `json:"tsunami"`
	MagnitudeLevel string    `json:"magnitude_level"`
	URL            string    `json:"url"`
}

type QuakesResponse struct {
	Total       int       `json:"total"`
	Quakes      []Quake   `json:"quakes"`
	LastUpdated time.Time `json:"last_updated"`
}

func GetMagnitudeLevel(mag float64) string {
	if mag >= 6.0 {
		return "major"
	} else if mag >= 5.0 {
		return "strong"
	} else if mag >= 4.0 {
		return "moderate"
	} else if mag >= 3.0 {
		return "minor"
	}
	return "unknown"
}
