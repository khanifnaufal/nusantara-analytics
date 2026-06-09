import { ref } from 'vue'
import { useIntervalFn } from '@vueuse/core'

/**
 * usePolling
 * Reusable polling composable that runs fetchFn immediately on call,
 * then polls it at the given interval.
 * 
 * @param fetchFn - The asynchronous function to execute periodically
 * @param intervalMs - The polling interval in milliseconds (default: 60000)
 */
export const usePolling = (
  fetchFn: () => Promise<any> | any,
  intervalMs = 60000,
  options: { immediate?: boolean } = {}
) => {
  const isPolling = ref(true)
  const { immediate = true } = options

  // Auto-fetch saat dipanggil pertama kali jika immediate true
  if (immediate) {
    fetchFn()
  }

  const { pause: pauseInterval, resume: resumeInterval } = useIntervalFn(
    async () => {
      await fetchFn()
    },
    intervalMs,
    { immediate: false }
  )

  const pause = () => {
    pauseInterval()
    isPolling.value = false
  }

  const resume = () => {
    resumeInterval()
    isPolling.value = true
  }

  return {
    isPolling,
    pause,
    resume
  }
}
