import { defineStore } from 'pinia'
import type { Quote } from './commodities'

export interface MarketResponse {
  stocks: Quote[]
  lastUpdated: string
}

export const useMarketStore = defineStore('market', {
  state: () => ({
    data: null as MarketResponse | null,
    loading: false,
    error: null as string | null
  }),

  actions: {
    async fetchMarket() {
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

        interface RawMarketResponse {
          stocks: RawQuote[]
          last_updated: string
        }

        const raw = await $fetch<RawMarketResponse>(`${apiBase}/api/v1/market`)

        this.data = {
          stocks: raw.stocks.map(c => ({
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
        console.error('Error fetching market:', err)
        this.error = err.message || 'Failed to fetch market'
      } finally {
        this.loading = false
      }
    }
  }
})
