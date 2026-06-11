import { defineStore } from 'pinia'

export interface CanvasWidget {
  id: string
  title: string
  dataSource: 'rates' | 'weather' | 'commodities' | 'market' | 'quakes'
  metric: string
  chartType: string
  position: { x: number; y: number }
  size: { w: number; h: number }
  config: Record<string, any>
}

let localIdCounter = 0

// Simple unique ID generator
export function generateUniqueId(): string {
  if (typeof globalThis !== 'undefined' && globalThis.crypto && typeof globalThis.crypto.randomUUID === 'function') {
    return globalThis.crypto.randomUUID()
  }
  localIdCounter++
  return `id-${localIdCounter}`
}

export const useCanvasStore = defineStore('canvas', {
  state: () => ({
    widgets: [] as CanvasWidget[],
    selectedWidgetId: null as string | null
  }),

  actions: {
    addWidget(widget: Omit<CanvasWidget, 'id'> & { id?: string }) {
      const newWidget: CanvasWidget = {
        ...widget,
        id: widget.id || generateUniqueId()
      }
      this.widgets.push(newWidget)
      this.selectedWidgetId = newWidget.id
    },

    removeWidget(id: string) {
      this.widgets = this.widgets.filter(w => w.id !== id)
      if (this.selectedWidgetId === id) {
        this.selectedWidgetId = null
      }
    },

    updatePosition(id: string, x: number, y: number) {
      const widget = this.widgets.find(w => w.id === id)
      if (widget) {
        widget.position = { x, y }
      }
    },

    updateSize(id: string, w: number, h: number) {
      const widget = this.widgets.find(w => w.id === id)
      if (widget) {
        widget.size = { w, h }
      }
    },

    updateConfig(id: string, config: Record<string, any>) {
      const widget = this.widgets.find(w => w.id === id)
      if (widget) {
        widget.config = { ...widget.config, ...config }
      }
    },

    clearCanvas() {
      this.widgets = []
      this.selectedWidgetId = null
    },

    exportToBase64(): string {
      try {
        const jsonStr = JSON.stringify(this.widgets)
        // Encode UTF-8 string to base64
        return btoa(unescape(encodeURIComponent(jsonStr)))
      } catch (err) {
        console.error('Error encoding canvas to base64:', err)
        return ''
      }
    },

    importFromBase64(str: string): boolean {
      if (!str) return false
      try {
        const decodedStr = decodeURIComponent(escape(atob(str)))
        const parsed = JSON.parse(decodedStr)
        
        const isValidWidget = (obj: any): obj is CanvasWidget => {
          return obj && 
            typeof obj.id === 'string' &&
            typeof obj.title === 'string' &&
            ['rates', 'weather', 'commodities', 'market', 'quakes'].includes(obj.dataSource) &&
            typeof obj.metric === 'string' &&
            typeof obj.chartType === 'string' &&
            obj.position && typeof obj.position.x === 'number' && typeof obj.position.y === 'number' &&
            obj.size && typeof obj.size.w === 'number' && typeof obj.size.h === 'number' &&
            typeof obj.config === 'object'
        }

        if (Array.isArray(parsed) && parsed.every(isValidWidget)) {
          this.widgets = parsed
          this.selectedWidgetId = null
          return true
        }
      } catch (err) {
        console.error('Error importing canvas from base64:', err)
      }
      return false
    },

    exportToUrl(): string {
      const base64 = this.exportToBase64()
      if (typeof window !== 'undefined' && base64) {
        const url = new URL(window.location.href)
        url.searchParams.set('canvas', base64)
        return url.toString()
      }
      return ''
    }
  }
})
