<template>
    <div>
      
      <PageHeader>
        <SearchBar @search="handleSearch" />
      </PageHeader>
  
      <div class="section-padding">
        <div class="container-narrow">
          <div class="flex flex-col lg:flex-row gap-8">
            
            <div class="lg:w-1/4 space-y-8">
              <CategoryFilter
                :categories="categoriesWithCount"
                :selected-category="selectedCategory"
                @filter="handleCategoryFilter"
              />
              <PriceFilter
                :min-price="0"
                :max-price="100"
                @price-change="handlePriceFilter"
              />
              <div class="bg-brand-50 dark:bg-brand-900/20 p-4 rounded-xl">
                <h4 class="font-bold text-brand-700 dark:text-brand-300 mb-2">💡 Tip</h4>
                <p class="text-sm text-gray-600 dark:text-gray-400">
                  Pre-order items require 30+ minutes preparation. Order early for best experience!
                </p>
              </div>
            </div>
  
            
            <div class="lg:w-3/4">
              
              <div v-if="activeFilters.length" class="mb-6 flex flex-wrap gap-2">
                <span
                  v-for="filter in activeFilters"
                  :key="filter"
                  class="inline-flex items-center bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 px-3 py-1 rounded-full text-sm"
                >
                  {{ filter }}
                  <button @click="removeFilter(filter)" class="ml-2 hover:text-red-500">
                    <X class="w-3 h-3" />
                  </button>
                </span>
                <button @click="clearAllFilters" class="text-sm text-brand-500 hover:text-brand-600">
                  Clear All
                </button>
              </div>
  
              
              <MenuGrid
                :filtered-items="paginatedItems"
                :is-loading="menuStore.isLoading"
                :current-page="currentPage"
                :total-pages="totalPages"
                @add-to-cart="handleAddToCart"
                @view-details="handleViewDetails"
                @reset-filters="clearAllFilters"
                @page-change="handlePageChange"
              />
  
              
              <div class="lg:hidden fixed bottom-0 left-0 right-0 bg-white dark:bg-slate-800 border-t shadow-lg p-4">
                <div class="flex justify-between items-center">
                  <div>
                    <div class="font-bold">{{ cart.totalItems }} items</div>
                    <div class="text-brand-500 font-bold">${{ cart.totalPrice.toFixed(2) }}</div>
                  </div>
                  <button @click="cart.toggleCart()" class="btn-primary">
                    View Cart
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
  
      
      <ItemDetailsModal
        v-if="selectedItem"
        :item="selectedItem"
        @close="selectedItem = null"
        @add-to-cart="handleAddToCart"
      />
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import { useMenuStore } from '../../stores/menu'
  import { useCartStore } from '../../stores/cart'
  import { X } from 'lucide-vue-next'
  import type { MenuItem } from '../../types/menu'
  import { useNuxtApp } from 'nuxt/app'
  
  import PageHeader from '../layout/PageHeader.vue'
  import SearchBar from './SearchBar.vue'
  import CategoryFilter from './CategoryFilter.vue'
  import PriceFilter from './PriceFilter.vue'
  import MenuGrid from './MenuGrid.vue'
  import ItemDetailsModal from './ItemDetailsModal.vue'
  
  const menuStore = useMenuStore()
  const cart = useCartStore()
  
  const selectedCategory = ref<string | null>(null)
  const searchQuery = ref('')
  const maxPrice = ref(100)
  const selectedItem = ref<MenuItem | null>(null)
  const currentPage = ref(1)
  const itemsPerPage = 12
  
  const categoriesWithCount = computed(() => {
    return menuStore.categories.map(category => ({
      ...category,
      itemsCount: menuStore.getItemsByCategory(category.id).length
    }))
  })
  
  const filteredItems = computed(() => {
    let items = menuStore.menuItems
  
    if (selectedCategory.value) {
      items = items.filter(item => item.categoryId === selectedCategory.value)
    }
  
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      items = items.filter(item =>
        item.name.toLowerCase().includes(query) ||
        item.description.toLowerCase().includes(query) ||
        item.tags?.some(tag => tag.toLowerCase().includes(query))
      )
    }
  
    items = items.filter(item => item.price <= maxPrice.value)
  
    return items
  })
  
  const paginatedItems = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage
    const end = start + itemsPerPage
    return filteredItems.value.slice(start, end)
  })
  
  const totalPages = computed(() => {
    return Math.ceil(filteredItems.value.length / itemsPerPage)
  })
  
  const activeFilters = computed(() => {
    const filters: string[] = []
    if (selectedCategory.value) {
      const category = menuStore.categories.find(c => c.id === selectedCategory.value)
      if (category) filters.push(category.name)
    }
    if (searchQuery.value) filters.push(`Search: "${searchQuery.value}"`)
    if (maxPrice.value < 100) filters.push(`Under $${maxPrice.value}`)
    return filters
  })
  
  const handleCategoryFilter = (categoryId: string | null) => {
    selectedCategory.value = categoryId
    currentPage.value = 1 // Reset to first page on filter change
  }
  
  const handleSearch = (query: string) => {
    searchQuery.value = query
    currentPage.value = 1 // Reset to first page on search
  }
  
  const handlePriceFilter = (price: number) => {
    maxPrice.value = price
    currentPage.value = 1 // Reset to first page on price filter
  }
  

const handleAddToCart = (item: MenuItem) => {
  cart.addItem(item)
  
  if (import.meta.client) {
    const { $toast } = useNuxtApp()
    
    const toast = $toast as {
      success?: (msg: string) => void
      error?: (msg: string) => void
      warning?: (msg: string) => void
      info?: (msg: string) => void
    }
    
    if (toast.success) {
      toast.success(`Added ${item.name} to cart`)
    }
  }
}
  
  const handleViewDetails = (item: MenuItem) => {
    selectedItem.value = item
  }
  
  const handlePageChange = (page: number) => {
    currentPage.value = page
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
  
  const removeFilter = (filter: string) => {
    if (filter.startsWith('Search:')) {
      searchQuery.value = ''
    } else if (filter.startsWith('Under $')) {
      maxPrice.value = 100
    } else {
      const category = menuStore.categories.find(c => c.name === filter)
      if (category) selectedCategory.value = null
    }
    currentPage.value = 1
  }
  
  const clearAllFilters = () => {
    selectedCategory.value = null
    searchQuery.value = ''
    maxPrice.value = 100
    currentPage.value = 1
  }
  
  onMounted(async () => {
    await Promise.all([
      menuStore.fetchCategories(),
      menuStore.fetchMenuItems()
    ])
  })
  </script>