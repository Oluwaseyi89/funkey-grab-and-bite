import type { OrderCreatePayload } from '~/types/order'

// One Idempotency-Key per checkout attempt. Pressing "Place order" again for
// the same order, even after a reload, reuses it so the order can't be placed
// twice; a changed cart or changed details get a new key.
const STORAGE_KEY = 'checkout-idempotency'

interface StoredKey {
  fingerprint: string
  key: string
}

// Used when session storage is unavailable (e.g. some private windows), so
// repeated presses in this page view still share a key.
let inMemory: StoredKey | null = null

const newKey = (): string => {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return `checkout-${crypto.randomUUID()}`
  }
  return `checkout-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 12)}`
}

const read = (): StoredKey | null => {
  try {
    const raw = sessionStorage.getItem(STORAGE_KEY)
    return raw ? (JSON.parse(raw) as StoredKey) : inMemory
  } catch {
    return inMemory
  }
}

export const checkoutKeyFor = (payload: OrderCreatePayload): string => {
  const fingerprint = JSON.stringify(payload)
  const stored = read()
  if (stored && stored.fingerprint === fingerprint) return stored.key

  const key = newKey()
  inMemory = { fingerprint, key }
  try {
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(inMemory))
  } catch {
    // inMemory still covers this page view.
  }
  return key
}

// Called once the order is placed, so the next order gets a fresh key.
export const clearCheckoutKey = (): void => {
  inMemory = null
  try {
    sessionStorage.removeItem(STORAGE_KEY)
  } catch {
    // nothing stored
  }
}
