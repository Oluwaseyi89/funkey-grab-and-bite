// store/menu.ts
import { defineStore } from 'pinia'
import type { MenuCategory, MenuItem } from '~/types/menu'
import { api } from '~/utils/api'

export const useMenuStore = defineStore('menu', {
  state: () => ({
    categories: [] as MenuCategory[],
    items: [] as MenuItem[],
    selectedCategory: null as string | null,
    isLoading: false,
    error: null as string | null,
  }),

  getters: {
    menuItems: (state) => state.items,
    getItemsByCategory: (state) => {
      return (categoryId: string) => {
        return state.items.filter(item => item.categoryId === categoryId)
      }
    },

    featuredItems: (state) =>
      state.items.filter(i => i.tags?.includes('best seller') || i.tags?.includes('popular')).slice(0, 6),

    categoryItems: (state) => (categoryId: string) =>
      state.items.filter(i => i.categoryId === categoryId),

    preOrderItems: (state) =>
      state.items.filter(i => i.isPreOrder === true),
  },

  actions: {
    async fetchCategories() {
      this.isLoading = true
      try {
        this.categories = await api.getMenuCategories()
        this.error = null
      } catch (err) {
        this.error = 'Failed to load categories'
        console.error(err)
      } finally {
        this.isLoading = false
      }
    },

    async fetchMenuItems(categoryId?: string) {
      this.isLoading = true
      try {
        this.items = await api.getMenuItems(categoryId)
        this.selectedCategory = categoryId || null
        this.error = null
      } catch (err) {
        this.error = 'Failed to load menu items'
        console.error(err)
      } finally {
        this.isLoading = false
      }
    },

    async fetchMenuItem(id: string) {
      this.isLoading = true
      try {
        return await api.getMenuItem(id)
      } catch (err) {
        this.error = 'Failed to load menu item'
        console.error(err)
        return null
      } finally {
        this.isLoading = false
      }
    },
  },
})























// import { defineStore } from 'pinia'
// import type { MenuCategory, MenuItem } from '~/types/menu'
// import { api } from '~/utils/api'

// export const useMenuStore = defineStore('menu', {
//   state: () => ({
//     categories: [] as MenuCategory[],
//     items: [] as MenuItem[],
//     selectedCategory: null as string | null,
//     isLoading: false,
//     error: null as string | null,
//   }),

//   getters: {
//     featuredItems: (state) => 
//       state.items.filter(item => item.tags?.includes('best seller') || item.tags?.includes('popular')).slice(0, 6),
    
//     categoryItems: (state) => (categoryId: string) =>
//       state.items.filter(item => item.categoryId === categoryId),
    
//     preOrderItems: (state) =>
//       state.items.filter(item => item.isPreOrder),
//   },

//   actions: {
//     async fetchCategories() {
//       this.isLoading = true
//       try {
//         this.categories = await api.getMenuCategories()
//         this.error = null
//       } catch (err) {
//         this.error = 'Failed to load categories'
//         console.error(err)
//       } finally {
//         this.isLoading = false
//       }
//     },

//     async fetchMenuItems(categoryId?: string) {
//       this.isLoading = true
//       try {
//         this.items = await api.getMenuItems(categoryId)
//         this.selectedCategory = categoryId || null
//         this.error = null
//       } catch (err) {
//         this.error = 'Failed to load menu items'
//         console.error(err)
//       } finally {
//         this.isLoading = false
//       }
//     },

//     async fetchMenuItem(id: string) {
//       this.isLoading = true
//       try {
//         const item = await api.getMenuItem(id)
//         this.error = null
//         return item
//       } catch (err) {
//         this.error = 'Failed to load menu item'
//         console.error(err)
//         return null
//       } finally {
//         this.isLoading = false
//       }
//     },
//   },
// })