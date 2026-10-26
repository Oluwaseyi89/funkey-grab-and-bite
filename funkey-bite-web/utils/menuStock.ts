import type { MenuItem } from '~/types/menu'

// Below this, customers see "Only N left"; above it no count is shown.
export const LOW_STOCK_THRESHOLD = 5

export interface StockState {
  orderable: boolean
  // The most that can be ordered, or null when stock isn't tracked.
  maxQuantity: number | null
  // Shown on the item when it needs calling out; empty otherwise.
  label: string
  tone: 'out' | 'low' | null
}

// What a customer may order of an item, from its availability and stock. The
// API checks stock again at checkout; this keeps customers from building a
// cart that would be rejected.
export function menuStock(item: MenuItem): StockState {
  if (!item.isAvailable) {
    return { orderable: false, maxQuantity: 0, label: 'Currently unavailable', tone: 'out' }
  }

  const stock = item.stockQuantity
  if (stock == null) {
    return { orderable: true, maxQuantity: null, label: '', tone: null }
  }
  if (stock <= 0) {
    return { orderable: false, maxQuantity: 0, label: 'Out of stock', tone: 'out' }
  }
  if (stock < LOW_STOCK_THRESHOLD) {
    return { orderable: true, maxQuantity: stock, label: `Only ${stock} left`, tone: 'low' }
  }
  return { orderable: true, maxQuantity: stock, label: '', tone: null }
}
