import { defineStore } from 'pinia'
import type { MenuCategory, MenuItem } from '~/types/menu'
import { useApi } from '~/utils/api'

let pending: Promise<void> | null = null

export const useMenuStore = defineStore('menu', {
  state: () => ({
    categories: [] as MenuCategory[],
    items: [] as MenuItem[],
    isLoading: false,
    loaded: false,
    error: null as string | null,
  }),

  getters: {
    menuItems: (state) => state.items,

    // Active categories that have something on them, in the admin's order.
    visibleCategories: (state) =>
      [...state.categories]
        .filter((category) => category.isActive && state.items.some((item) => item.categoryId === category.id))
        .sort((a, b) => a.displayOrder - b.displayOrder),

    getItemsByCategory: (state) => (categoryId: string) =>
      state.items.filter((item) => item.categoryId === categoryId),

    // Items tagged as favourites, or the first few on the menu when none are,
    // with anything that can't be ordered left out.
    featuredItems: (state) => {
      const available = state.items.filter((item) => item.isAvailable)
      const tagged = available.filter((item) =>
        item.tags?.some((tag) => ['best seller', 'popular', 'featured'].includes(tag.toLowerCase())),
      )
      return (tagged.length ? tagged : available).slice(0, 8)
    },
  },

  actions: {
    // Loads the menu once per visit; force reloads it (e.g. "Try again").
    async load(force = false) {
      if (this.loaded && !force) return
      if (!pending) {
        pending = (async () => {
          this.isLoading = true
          this.error = null
          try {
            const api = useApi()
            const [categories, items] = await Promise.all([api.getMenuCategories(), api.getMenuItems()])
            this.categories = categories
            this.items = items
            this.loaded = true
          } catch (err) {
            this.error = err instanceof Error ? err.message : 'We couldn’t load the menu.'
          } finally {
            this.isLoading = false
            pending = null
          }
        })()
      }
      await pending
    },

    async fetchCategories() {
      await this.load()
    },

    async fetchMenuItems() {
      await this.load()
    },
  },
})
