import { defineStore } from 'pinia'

export interface Quake {
  id: string
  place: string
  magnitude: number
  depth: number
  lat: number
  lon: number
  time: string
  magnitudeLevel: string
  tsunami: number
  url: string
}

export interface QuakesResponse {
  total: number
  quakes: Quake[]
  lastUpdated: string
}

export const useQuakesStore = defineStore('quakes', {
  state: () => ({
    data: null as QuakesResponse | null,
    loading: false,
    error: null as string | null
  }),

  actions: {
    async fetchQuakes() {
      this.loading = true
      this.error = null
      try {
        const config = useRuntimeConfig()
        const apiBase = config.public.apiBase

        interface RawQuake {
          id: string
          place: string
          magnitude: number
          depth: number
          lat: number
          lon: number
          time: string
          alert: string
          tsunami: number
          magnitude_level: string
          url: string
        }

        interface RawQuakesResponse {
          total: number
          quakes: RawQuake[]
          last_updated: string
        }

        const raw = await $fetch<RawQuakesResponse>(`${apiBase}/api/v1/quakes`)

        this.data = {
          total: raw.total,
          quakes: raw.quakes.map(q => ({
            id: q.id,
            place: q.place,
            magnitude: q.magnitude,
            depth: q.depth,
            lat: q.lat,
            lon: q.lon,
            time: q.time,
            magnitudeLevel: q.magnitude_level,
            tsunami: q.tsunami,
            url: q.url
          })),
          lastUpdated: raw.last_updated
        }
      } catch (err: any) {
        console.error('Error fetching quakes:', err)
        this.error = err.message || 'Failed to fetch quakes'
      } finally {
        this.loading = false
      }
    }
  }
})
