import { defineStore } from 'pinia'
import { useNuxtApp } from 'nuxt/app'
import { ref, computed } from 'vue'
import type { CartItem, MenuItem } from '~/types/menu'
import { menuStock } from '~/utils/menuStock'

type Toast = { warning?: (message: string) => void }

function warn(message: string) {
  if (!import.meta.client) return
  const { $toast } = useNuxtApp()
  ;($toast as Toast | undefined)?.warning?.(message)
}

export const useCartStore = defineStore('cart', () => {
  const items = ref<CartItem[]>([])
  // A promotion code the customer chose (e.g. from the promotions page); the
  // server checks it at checkout.
  const promoCode = ref('')
  const isOpen = ref(false)
  // Read out by screen readers and used to animate the cart badge.
  const announcement = ref('')
  const bump = ref(0)

  const totalItems = computed(() => items.value.reduce((sum, item) => sum + item.quantity, 0))
  const subtotal = computed(() => items.value.reduce((sum, item) => sum + item.menuItem.price * item.quantity, 0))
  // Items only: tax and delivery come from the business settings at checkout.
  const totalPrice = computed(() => subtotal.value)
  const isEmpty = computed(() => items.value.length === 0)

  const find = (menuItemId: string) => items.value.find((item) => item.menuItem.id === menuItemId)
  const itemCount = (menuItemId: string) => find(menuItemId)?.quantity ?? 0

  // Whether one more of this item can be added.
  const canAddMore = (menuItem: MenuItem) => {
    const stock = menuStock(menuItem)
    if (!stock.orderable) return false
    return stock.maxQuantity == null || itemCount(menuItem.id) < stock.maxQuantity
  }

  const say = (message: string) => {
    announcement.value = message
    bump.value++
  }

  // Adds quantity of an item, never more than is in stock. Returns how many
  // were actually added.
  function addItem(menuItem: MenuItem, quantity = 1, specialInstructions?: string): number {
    const stock = menuStock(menuItem)
    if (!stock.orderable) {
      warn(`${menuItem.name} is ${stock.label.toLowerCase()}`)
      return 0
    }

    const existing = find(menuItem.id)
    const current = existing?.quantity ?? 0
    let adding = quantity
    if (stock.maxQuantity != null && current + adding > stock.maxQuantity) {
      adding = stock.maxQuantity - current
      if (adding <= 0) {
        warn(`Only ${stock.maxQuantity} ${menuItem.name} left, and they're all in your cart`)
        return 0
      }
      warn(`Only ${stock.maxQuantity} ${menuItem.name} left`)
    }

    if (existing) {
      existing.quantity += adding
      // Keep the latest price and stock with the item.
      existing.menuItem = menuItem
      if (specialInstructions !== undefined) existing.specialInstructions = specialInstructions || undefined
    } else {
      items.value.push({ menuItem, quantity: adding, specialInstructions: specialInstructions || undefined })
    }
    say(`${menuItem.name} added. ${totalItems.value} in your cart.`)
    return adding
  }

  function removeItem(menuItemId: string) {
    const item = find(menuItemId)
    items.value = items.value.filter((entry) => entry.menuItem.id !== menuItemId)
    if (item) say(`${item.menuItem.name} removed from your cart.`)
  }

  function updateQuantity(menuItemId: string, quantity: number) {
    const item = find(menuItemId)
    if (!item) return
    if (quantity <= 0) {
      removeItem(menuItemId)
      return
    }
    const { maxQuantity } = menuStock(item.menuItem)
    if (maxQuantity != null && maxQuantity > 0 && quantity > maxQuantity) {
      warn(`Only ${maxQuantity} ${item.menuItem.name} left`)
      quantity = maxQuantity
    }
    item.quantity = quantity
    say(`${item.menuItem.name}: ${quantity}.`)
  }

  const increase = (menuItem: MenuItem) => addItem(menuItem, 1)
  const decrease = (menuItemId: string) => updateQuantity(menuItemId, itemCount(menuItemId) - 1)

  function setInstructions(menuItemId: string, text: string) {
    const item = find(menuItemId)
    if (item) item.specialInstructions = text.trim() || undefined
  }

  function clearCart() {
    items.value = []
    promoCode.value = ''
  }

  const openCart = () => { isOpen.value = true }
  const closeCart = () => { isOpen.value = false }
  const toggleCart = () => { isOpen.value = !isOpen.value }

  return {
    items,
    promoCode,
    isOpen,
    announcement,
    bump,
    totalItems,
    subtotal,
    totalPrice,
    isEmpty,
    itemCount,
    canAddMore,
    addItem,
    removeItem,
    updateQuantity,
    increase,
    decrease,
    setInstructions,
    clearCart,
    openCart,
    closeCart,
    toggleCart,
  }
}, {
  // The cart and chosen code survive a reload; whether the sheet is open doesn't.
  persist: { pick: ['items', 'promoCode'] },
})
