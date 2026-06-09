import { defineStore } from 'pinia'

export interface PriceHistory {
  timestamp: string
  close: number
}

export interface Quote {
  symbol: string
  name: string
  price: number
  change: number
  changePercent: number
  currency: string
  history: PriceHistory[]
  lastUpdated?: string
}

export interface CommoditiesResponse {
  commodities: Quote[]
  lastUpdated: string
}

export const useCommoditiesStore = defineStore('commodities', {
  state: () => ({
    data: null as CommoditiesResponse | null,
    loading: false,
    error: null as string | null
  }),

  actions: {
    async fetchCommodities() {
      this.loading = true
      this.error = null
      try {
        const config = useRuntimeConfig()
        const apiBase = config.public.apiBase

        interface RawPriceHistory {
          timestamp: string
          close: number
        }

        interface RawQuote {
          symbol: string
          name: string
          price: number
          change: number
          change_percent: number
          currency: string
          history: RawPriceHistory[]
          last_updated: string
        }

        interface RawCommoditiesResponse {
          commodities: RawQuote[]
          last_updated: string
        }

        const raw = await $fetch<RawCommoditiesResponse>(`${apiBase}/api/v1/commodities`)

        this.data = {
          commodities: raw.commodities.map(c => ({
            symbol: c.symbol,
            name: c.name,
            price: c.price,
            change: c.change,
            changePercent: c.change_percent,
            currency: c.currency,
            history: c.history.map(h => ({
              timestamp: h.timestamp,
              close: h.close
            })),
            lastUpdated: c.last_updated
          })),
          lastUpdated: raw.last_updated
        }
      } catch (err: any) {
        console.error('Error fetching commodities:', err)
        this.error = err.message || 'Failed to fetch commodities'
      } finally {
        this.loading = false
      }
    }
  }
})
