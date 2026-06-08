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

type OpenMeteoResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Temperature float64 `json:"temperature_2m"`
		Humidity    int     `json:"relative_humidity_2m"`
		WeatherCode int     `json:"weather_code"`
		WindSpeed   float64 `json:"wind_speed_10m"`
	} `json:"current"`
}

func FetchWeather() (*models.WeatherResponse, bool, error) {
	cacheKey := "weather:all"
	ttl := 30 * time.Minute

	if val, found := cache.Get(cacheKey); found {
		if resp, ok := val.(*models.WeatherResponse); ok {
			return resp, true, nil
		}
	}

	citiesData := []struct {
		name string
		lat  float64
		lon  float64
	}{
		{"Jakarta", -6.2088, 106.8456},
		{"Surabaya", -7.2575, 112.7521},
		{"Bandung", -6.9175, 107.6191},
		{"Medan", 3.5952, 98.6722},
		{"Semarang", -6.9932, 110.4203},
		{"Makassar", -5.1477, 119.4327},
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	list := make([]*models.CityWeather, len(citiesData))

	client := resty.New()
	client.SetTimeout(10 * time.Second)

	for i, city := range citiesData {
		wg.Add(1)
		go func(idx int, cName string, cLat, cLon float64) {
			defer wg.Done()
			url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m&timezone=Asia/Jakarta", cLat, cLon)

			var apiResp OpenMeteoResponse
			resp, err := client.R().
				SetResult(&apiResp).
				Get(url)

			if err != nil || resp.IsError() {
				return
			}

			weather := &models.CityWeather{
				City:         cName,
				Lat:          cLat,
				Lon:          cLon,
				Temperature:  apiResp.Current.Temperature,
				Humidity:     apiResp.Current.Humidity,
				WindSpeed:    apiResp.Current.WindSpeed,
				WeatherCode:  apiResp.Current.WeatherCode,
				WeatherDesc:  models.WeatherCodeToDesc(apiResp.Current.WeatherCode),
				LastUpdated:  time.Now(),
			}

			mu.Lock()
			list[idx] = weather
			mu.Unlock()
		}(i, city.name, city.lat, city.lon)
	}

	wg.Wait()

	var filteredList []models.CityWeather
	for _, item := range list {
		if item != nil {
			filteredList = append(filteredList, *item)
		}
	}

	if len(filteredList) == 0 {
		return nil, false, errors.New("failed to fetch weather data for all cities")
	}

	weatherResp := &models.WeatherResponse{
		Cities:      filteredList,
		LastUpdated: time.Now(),
	}

	cache.Set(cacheKey, weatherResp, ttl)
	return weatherResp, false, nil
}
