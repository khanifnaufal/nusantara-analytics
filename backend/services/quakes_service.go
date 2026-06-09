package services

import (
	"errors"
	"time"

	"indo-stats-backend/cache"
	"indo-stats-backend/models"

	"github.com/go-resty/resty/v2"
)

type USGSGeoJSONResponse struct {
	Features []struct {
		ID         string `json:"id"`
		Properties struct {
			Mag     float64 `json:"mag"`
			Place   string  `json:"place"`
			Time    int64   `json:"time"`
			Alert   string  `json:"alert"`
			Tsunami int     `json:"tsunami"`
			URL     string  `json:"url"`
		} `json:"properties"`
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

func FetchQuakes() (*models.QuakesResponse, bool, error) {
	cacheKey := "quakes:indonesia"
	ttl := 10 * time.Minute

	if val, found := cache.Get(cacheKey); found {
		if resp, ok := val.(*models.QuakesResponse); ok {
			return resp, true, nil
		}
	}

	client := resty.New()
	client.SetTimeout(10 * time.Second)
	var apiResp USGSGeoJSONResponse

	url := "https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson&minlatitude=-11&maxlatitude=6&minlongitude=95&maxlongitude=141&minmagnitude=3.0&orderby=time&limit=20"
	resp, err := client.R().
		SetResult(&apiResp).
		Get(url)

	if err != nil {
		return nil, false, err
	}
	if resp.IsError() {
		return nil, false, errors.New("failed to fetch earthquakes from USGS API: " + resp.Status())
	}

	var list []models.Quake
	for _, feat := range apiResp.Features {
		var lat, lon, depth float64
		coords := feat.Geometry.Coordinates
		if len(coords) > 0 {
			lon = coords[0]
		}
		if len(coords) > 1 {
			lat = coords[1]
		}
		if len(coords) > 2 {
			depth = coords[2]
		}

		quake := models.Quake{
			ID:             feat.ID,
			Place:          feat.Properties.Place,
			Magnitude:      feat.Properties.Mag,
			Depth:          depth,
			Lat:            lat,
			Lon:            lon,
			Time:           time.UnixMilli(feat.Properties.Time),
			Alert:          feat.Properties.Alert,
			Tsunami:        feat.Properties.Tsunami,
			MagnitudeLevel: models.GetMagnitudeLevel(feat.Properties.Mag),
			URL:            feat.Properties.URL,
		}
		list = append(list, quake)
	}

	quakesResp := &models.QuakesResponse{
		Total:       len(list),
		Quakes:      list,
		LastUpdated: time.Now(),
	}

	cache.Set(cacheKey, quakesResp, ttl)
	return quakesResp, false, nil
}
