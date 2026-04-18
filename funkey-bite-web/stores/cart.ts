import { defineStore } from 'pinia'
import { useNuxtApp } from 'nuxt/app'
import { ref, computed } from 'vue'
import type { CartItem, MenuItem } from '~/types/menu'

export const useCartStore = defineStore('cart', () => {
  const items = ref<CartItem[]>([])
  const isOpen = ref(false)

  const totalItems = computed(() => items.value.reduce((sum: number, item: CartItem) => sum + item.quantity, 0))
  const subtotal = computed(() => items.value.reduce((sum: number, item: CartItem) => sum + item.menuItem.price * item.quantity, 0))
  const tax = computed(() => subtotal.value * 0.08)
  const total = computed(() => subtotal.value + tax.value)
  const totalPrice = computed(() => {
    const sub = items.value.reduce((sum: number, item: CartItem) => sum + (item.menuItem.price * item.quantity), 0)
    const t = sub * 0.08
    return sub + t
  })
  const itemCount = (menuItemId: string) => items.value.find((item: CartItem) => item.menuItem.id === menuItemId)?.quantity ?? 0

  function addItem(menuItem: MenuItem, quantity = 1, specialInstructions?: string) {
    const existingItem = items.value.find((item: CartItem) => item.menuItem.id === menuItem.id)
    if (existingItem) {
      existingItem.quantity += quantity
      if (specialInstructions !== undefined) {
        existingItem.specialInstructions = specialInstructions
      }
    } else {
      items.value.push({ menuItem, quantity, specialInstructions })
    }
    if (import.meta.client) {
      const { $toast } = useNuxtApp()
      const toast = $toast as { success?: (message: string) => void } | undefined
      toast?.success?.(`${menuItem.name} added to cart`)
    }
  }

  function removeItem(menuItemId: string) {
    items.value = items.value.filter((item: CartItem) => item.menuItem.id !== menuItemId)
  }

  function updateQuantity(menuItemId: string, quantity: number) {
    const item = items.value.find((item: CartItem) => item.menuItem.id === menuItemId)
    if (!item) return
    if (quantity <= 0) {
      removeItem(menuItemId)
    } else {
      item.quantity = quantity
    }
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
    tax,
    total,
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

