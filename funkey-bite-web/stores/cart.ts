import { defineStore } from 'pinia'
import type { CartItem, MenuItem } from '~/types/menu'
import { useToast } from 'vue-toastification'


export const useCartStore = defineStore('cart', {
  state: () => ({
    items: [] as CartItem[],
    isOpen: false,
  }),

  getters: {
    totalItems(state): number {
      return state.items.reduce((sum, item) => sum + item.quantity, 0)
    },

    subtotal(state): number {
      return state.items.reduce(
        (sum, item) => sum + item.menuItem.price * item.quantity,
        0
      )
    },

    tax(): number {
      return this.subtotal * 0.08
    },

    total(): number {
      return this.subtotal + this.tax
    },
    totalPrice: (state) => {
      const subtotal = state.items.reduce((sum, item) => sum + (item.menuItem.price * item.quantity), 0)
      const tax = subtotal * 0.08
      return subtotal + tax
    },

    itemCount(state) {
      return (menuItemId: string): number =>
        state.items.find(item => item.menuItem.id === menuItemId)?.quantity ?? 0
    },
  },

  actions: {
    addItem(
      menuItem: MenuItem,
      quantity = 1,
      specialInstructions?: string
    ) {
      const existingItem = this.items.find(
        item => item.menuItem.id === menuItem.id
      )

      if (existingItem) {
        existingItem.quantity += quantity
        if (specialInstructions !== undefined) {
          existingItem.specialInstructions = specialInstructions
        }
      } else {
        this.items.push({
          menuItem,
          quantity,
          specialInstructions,
        })
      }

      const toast = useToast()
      toast.success(`${menuItem.name} added to cart`)
    },

    removeItem(menuItemId: string) {
      this.items = this.items.filter(
        item => item.menuItem.id !== menuItemId
      )
    },

    updateQuantity(menuItemId: string, quantity: number) {
      const item = this.items.find(
        item => item.menuItem.id === menuItemId
      )

      if (!item) return

      if (quantity <= 0) {
        this.removeItem(menuItemId)
      } else {
        item.quantity = quantity
      }
    },

    clearCart() {
      this.items = []
    },

    toggleCart() {
      this.isOpen = !this.isOpen
    },
  },

  persist: true,
})

