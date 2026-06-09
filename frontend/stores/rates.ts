import { defineStore } from 'pinia'

export interface RatesResponse {
  base: string
  date: string
  rates: Record<string, number>
  lastUpdated: string
}

export const useRatesStore = defineStore('rates', {
  state: () => ({
    data: null as RatesResponse | null,
    loading: false,
    error: null as string | null,
    lastUpdated: null as Date | null
  }),

  getters: {
    ratesArray(state): { currency: string; value: number }[] {
      if (!state.data || !state.data.rates) return []
      return Object.entries(state.data.rates)
        .map(([currency, value]) => ({ currency, value }))
        .sort((a, b) => a.value - b.value)
    }
  },

  actions: {
    async fetchRates() {
      this.loading = true
      this.error = null
      try {
        const config = useRuntimeConfig()
        const apiBase = config.public.apiBase
        
        const raw = await $fetch<{
          base: string
          date: string
          rates: Record<string, number>
          last_updated: string
        }>(`${apiBase}/api/v1/rates`)

        this.data = {
          base: raw.base,
          date: raw.date,
          rates: raw.rates,
          lastUpdated: raw.last_updated
        }
        this.lastUpdated = raw.last_updated ? new Date(raw.last_updated) : null
      } catch (err: any) {
        console.error('Error fetching rates:', err)
        this.error = err.message || 'Failed to fetch rates'
      } finally {
        this.loading = false
      }
    }
  }
})
