import { defineStore } from 'pinia'
import { useNuxtApp } from 'nuxt/app'
import { ref, computed } from 'vue'
import type { CartItem, MenuItem } from '~/types/menu'
import { menuStock } from '~/utils/menuStock'

type Toast = { success?: (message: string) => void; warning?: (message: string) => void }

function notify(kind: keyof Toast, message: string) {
  if (!import.meta.client) return
  const { $toast } = useNuxtApp()
  ;($toast as Toast | undefined)?.[kind]?.(message)
}

export const useCartStore = defineStore('cart', () => {
  const items = ref<CartItem[]>([])
  const isOpen = ref(false)

  const totalItems = computed(() => items.value.reduce((sum: number, item: CartItem) => sum + item.quantity, 0))
  const subtotal = computed(() => items.value.reduce((sum: number, item: CartItem) => sum + item.menuItem.price * item.quantity, 0))
  // Items only: tax and delivery come from the business settings at checkout.
  const totalPrice = computed(() => subtotal.value)
  const itemCount = (menuItemId: string) => items.value.find((item: CartItem) => item.menuItem.id === menuItemId)?.quantity ?? 0

  function addItem(menuItem: MenuItem, quantity = 1, specialInstructions?: string) {
    const stock = menuStock(menuItem)
    if (!stock.orderable) {
      notify('warning', `${menuItem.name} is ${stock.label.toLowerCase()}`)
      return
    }

    const existingItem = items.value.find((item: CartItem) => item.menuItem.id === menuItem.id)
    const wanted = (existingItem?.quantity ?? 0) + quantity
    // Never put more in the cart than there is stock for.
    if (stock.maxQuantity != null && wanted > stock.maxQuantity) {
      if ((existingItem?.quantity ?? 0) >= stock.maxQuantity) {
        notify('warning', `Only ${stock.maxQuantity} ${menuItem.name} left, and they're all in your cart`)
        return
      }
      quantity = stock.maxQuantity - (existingItem?.quantity ?? 0)
      notify('warning', `Only ${stock.maxQuantity} ${menuItem.name} left`)
    }

    if (existingItem) {
      existingItem.quantity += quantity
      // Keep the latest stock with the item for later quantity changes.
      existingItem.menuItem = menuItem
      if (specialInstructions !== undefined) {
        existingItem.specialInstructions = specialInstructions
      }
    } else {
      items.value.push({ menuItem, quantity, specialInstructions })
    }
    notify('success', `${menuItem.name} added to cart`)
  }

  function removeItem(menuItemId: string) {
    items.value = items.value.filter((item: CartItem) => item.menuItem.id !== menuItemId)
  }

  function updateQuantity(menuItemId: string, quantity: number) {
    const item = items.value.find((item: CartItem) => item.menuItem.id === menuItemId)
    if (!item) return
    if (quantity <= 0) {
      removeItem(menuItemId)
      return
    }
    const { maxQuantity } = menuStock(item.menuItem)
    if (maxQuantity != null && maxQuantity > 0 && quantity > maxQuantity) {
      notify('warning', `Only ${maxQuantity} ${item.menuItem.name} left`)
      item.quantity = maxQuantity
      return
    }
    item.quantity = quantity
  }

  function clearCart() {
    items.value = []
  }

  function toggleCart() {
    isOpen.value = !isOpen.value
  }

  return {
    items,
    isOpen,
    totalItems,
    subtotal,
    totalPrice,
    itemCount,
    addItem,
    removeItem,
    updateQuantity,
    clearCart,
    toggleCart,
  }
}, {
  persist: true,
})

