import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { formatClock, useBusinessSettings } from './useBusinessSettings'

const WEEK = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']

const minutes = (clock: string) => {
  const [h = 0, m = 0] = clock.split(':').map(Number)
  return h * 60 + m
}

// "Open now · until 10:00 PM" or "Closed · opens 8:00 AM tomorrow", in the
// business's own timezone. Whether ordering is allowed still comes from the
// server's isOpenNow; this only explains it.
export function useStoreStatus() {
  const { settings, status, load } = useBusinessSettings()
  const now = ref(Date.now())
  let timer: ReturnType<typeof setInterval> | undefined

  onMounted(() => {
    load()
    timer = setInterval(() => { now.value = Date.now() }, 60_000)
  })
  onBeforeUnmount(() => clearInterval(timer))

  // The weekday and minute of the day where the business is.
  const local = computed(() => {
    const timeZone = settings.value?.timezone || undefined
    const parts = new Intl.DateTimeFormat('en-US', { timeZone, weekday: 'long', hour: '2-digit', minute: '2-digit', hourCycle: 'h23' })
      .formatToParts(new Date(now.value))
    const get = (type: string) => parts.find((p) => p.type === type)?.value ?? ''
    return { day: get('weekday'), minute: Number(get('hour')) * 60 + Number(get('minute')) }
  })

  const hoursFor = (day: string) => settings.value?.openingHours?.find((h) => h.day === day)

  const isOpen = computed(() => Boolean(settings.value?.isOpenNow))

  const detail = computed(() => {
    const s = settings.value
    if (!s) return ''
    const today = local.value.day
    if (s.isOpenNow) {
      const hours = hoursFor(today)
      return hours?.closeTime ? `until ${formatClock(hours.closeTime)}` : ''
    }
    // The next opening: later today, or the next open day this week.
    const startIndex = WEEK.indexOf(today)
    for (let offset = 0; offset < 7; offset++) {
      const day = WEEK[(startIndex + offset) % 7]!
      const hours = hoursFor(day)
      if (!hours?.isOpen || !hours.openTime) continue
      if (offset === 0 && minutes(hours.openTime) <= local.value.minute) continue
      const when = offset === 0 ? 'today' : offset === 1 ? 'tomorrow' : day
      return `opens ${formatClock(hours.openTime)} ${when}`
    }
    return ''
  })

  const label = computed(() => {
    if (status.value === 'loading') return ''
    if (status.value === 'unavailable') return ''
    return `${isOpen.value ? 'Open now' : 'Closed'}${detail.value ? ` · ${detail.value}` : ''}`
  })

  return { isOpen, label, detail, ready: computed(() => status.value === 'ready') }
}
