// Orders placed on this device, so a guest can open them again without
// typing the order number. Kept only in this browser; nothing is sent anywhere.

export interface RecentOrder {
  orderNumber: string
  phone: string
  placedAt: string
  total: number
}

const KEY = 'funkey-recent-orders'
const LIMIT = 5

export function readRecentOrders(): RecentOrder[] {
  if (!import.meta.client) return []
  try {
    const parsed = JSON.parse(localStorage.getItem(KEY) || '[]')
    return Array.isArray(parsed) ? parsed.filter((order) => order?.orderNumber && order?.phone) : []
  } catch {
    return []
  }
}

export function rememberOrder(order: RecentOrder) {
  if (!import.meta.client) return
  const others = readRecentOrders().filter((entry) => entry.orderNumber !== order.orderNumber)
  try {
    localStorage.setItem(KEY, JSON.stringify([order, ...others].slice(0, LIMIT)))
  } catch {
    // Storage full or blocked: the order can still be found by its number.
  }
}

export function forgetRecentOrders() {
  if (!import.meta.client) return
  try {
    localStorage.removeItem(KEY)
  } catch {
    // ignore
  }
}
