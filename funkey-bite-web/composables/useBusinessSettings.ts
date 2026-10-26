import { computed } from 'vue'
import { useState } from 'nuxt/app'
import type { BusinessSettings, OpeningHours } from '~/types/settings'
import { useApi } from '~/utils/api'

// The API judges open/closed in this zone (BUSINESS_TIMEZONE); "today" is
// highlighted in it too, whatever the visitor's own timezone is.
export const BUSINESS_TIMEZONE = 'Africa/Lagos'

const WEEK = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday']

export interface HoursGroup {
  label: string
  hours: string
  days: string[]
}

// "08:00" -> "8:00 AM"
export const formatClock = (value: string) => {
  const [hour = 0, minute = 0] = value.split(':').map(Number)
  const period = hour < 12 ? 'AM' : 'PM'
  const displayHour = hour % 12 === 0 ? 12 : hour % 12
  return `${displayHour}:${String(minute).padStart(2, '0')} ${period}`
}

export const formatDayHours = (h: OpeningHours) => {
  if (!h.isOpen) return 'Closed'
  const range = `${formatClock(h.openTime)} - ${formatClock(h.closeTime)}`
  return h.closeTime <= h.openTime ? `${range} (next day)` : range
}

const shortDay = (day: string) => day.slice(0, 3)

// Collapses consecutive days with identical hours, e.g. "Mon - Fri".
const groupHours = (hours: OpeningHours[]): HoursGroup[] => {
  const groups: HoursGroup[] = []
  for (const h of hours) {
    const label = formatDayHours(h)
    const last = groups[groups.length - 1]
    if (last && last.hours === label) {
      last.days.push(h.day)
    } else {
      groups.push({ label: '', hours: label, days: [h.day] })
    }
  }
  for (const group of groups) {
    const first = group.days[0] ?? ''
    const last = group.days[group.days.length - 1] ?? ''
    group.label = group.days.length === 1 ? first : `${shortDay(first)} - ${shortDay(last)}`
  }
  return groups
}

let pending: Promise<void> | null = null

// Shared, client-side business settings. The site is statically generated, so
// fetching in the browser keeps contact details and hours live without a rebuild.
export function useBusinessSettings() {
  const settings = useState<BusinessSettings | null>('business-settings', () => null)
  const loaded = useState<boolean>('business-settings-loaded', () => false)
  const api = process.client ? useApi() : null

  const load = async (force = false) => {
    if (!api || (loaded.value && !force)) return
    if (!pending) {
      pending = api
        .getBusinessSettings()
        .then((data) => {
          if (data) settings.value = data
          loaded.value = true
        })
        .finally(() => {
          pending = null
        })
    }
    await pending
  }

  const status = computed<'loading' | 'ready' | 'unavailable'>(() =>
    settings.value ? 'ready' : loaded.value ? 'unavailable' : 'loading',
  )

  const openingHours = computed<OpeningHours[]>(() =>
    WEEK.map((day) => {
      const found = settings.value?.openingHours?.find((h) => h.day === day)
      return found ?? { day, openTime: '', closeTime: '', isOpen: false }
    }),
  )

  const hoursGroups = computed(() => groupHours(openingHours.value))

  const today = computed(() =>
    new Intl.DateTimeFormat('en-US', { weekday: 'long', timeZone: BUSINESS_TIMEZONE }).format(new Date()),
  )

  const addressLines = computed(() => {
    const s = settings.value
    if (!s) return []
    const locality = [s.city, s.state, s.postalCode].filter(Boolean).join(', ')
    return [s.address, locality, s.country].filter(Boolean)
  })

  const fullAddress = computed(() => addressLines.value.join(', '))

  const mapsUrl = computed(() =>
    fullAddress.value
      ? `https://www.google.com/maps/search/?api=1&query=${encodeURIComponent(fullAddress.value)}`
      : '',
  )

  const phoneHref = computed(() =>
    settings.value?.phoneNumber ? `tel:${settings.value.phoneNumber.replace(/[^\d+]/g, '')}` : '',
  )

  return {
    settings,
    status,
    load,
    openingHours,
    hoursGroups,
    today,
    addressLines,
    fullAddress,
    mapsUrl,
    phoneHref,
  }
}
