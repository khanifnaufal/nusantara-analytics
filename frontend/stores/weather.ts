import { defineStore } from 'pinia'

export interface CityWeather {
  city: string
  temperature: number
  humidity: number
  windSpeed: number
  weatherCode: number
  weatherDesc: string
  lat?: number
  lon?: number
}

export interface WeatherResponse {
  cities: CityWeather[]
  lastUpdated: string
}

export const useWeatherStore = defineStore('weather', {
  state: () => ({
    data: null as WeatherResponse | null,
    loading: false,
    error: null as string | null
  }),

  actions: {
    async fetchWeather() {
      this.loading = true
      this.error = null
      try {
        const config = useRuntimeConfig()
        const apiBase = config.public.apiBase

        interface RawCityWeather {
          city: string
          lat: number
          lon: number
          temperature: number
          humidity: number
          wind_speed: number
          weather_code: number
          weather_desc: string
          last_updated: string
        }

        interface RawWeatherResponse {
          cities: RawCityWeather[]
          last_updated: string
        }

        const raw = await $fetch<RawWeatherResponse>(`${apiBase}/api/v1/weather`)

        this.data = {
          cities: raw.cities.map(c => ({
            city: c.city,
            lat: c.lat,
            lon: c.lon,
            temperature: c.temperature,
            humidity: c.humidity,
            windSpeed: c.wind_speed,
            weatherCode: c.weather_code,
            weatherDesc: c.weather_desc
          })),
          lastUpdated: raw.last_updated
        }
      } catch (err: any) {
        console.error('Error fetching weather:', err)
        this.error = err.message || 'Failed to fetch weather'
      } finally {
        this.loading = false
      }
    }
  }
})
