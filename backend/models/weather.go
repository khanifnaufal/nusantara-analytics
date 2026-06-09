package models

import "time"

type CityWeather struct {
	City         string    `json:"city"`
	Lat          float64   `json:"lat"`
	Lon          float64   `json:"lon"`
	Temperature  float64   `json:"temperature"`
	Humidity     int       `json:"humidity"`
	WindSpeed    float64   `json:"wind_speed"`
	WeatherCode  int       `json:"weather_code"`
	WeatherDesc  string    `json:"weather_desc"`
	LastUpdated  time.Time `json:"last_updated"`
}

type WeatherResponse struct {
	Cities      []CityWeather `json:"cities"`
	LastUpdated time.Time     `json:"last_updated"`
}

func WeatherCodeToDesc(code int) string {
	switch code {
	case 0:
		return "Cerah"
	case 1, 2, 3:
		return "Berawan"
	case 45, 48:
		return "Berkabut"
	case 51, 53, 55:
		return "Gerimis"
	case 56, 57:
		return "Gerimis Membeku"
	case 61:
		return "Hujan Ringan"
	case 63:
		return "Hujan Sedang"
	case 65:
		return "Hujan Lebat"
	case 66, 67:
		return "Hujan Beku"
	case 71, 73, 75:
		return "Hujan Salju"
	case 77:
		return "Hujan Es"
	case 80, 81, 82:
		return "Hujan Pancaroba"
	case 85, 86:
		return "Hujan Salju Showers"
	case 95:
		return "Badai Guntur"
	case 96, 99:
		return "Badai Guntur dengan Hujan Es"
	default:
		return "Kondisi Tidak Diketahui"
	}
}
