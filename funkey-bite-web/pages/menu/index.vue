<template>
  <div>
    <!-- Enhanced Page Header -->
    <PageHeader
      title-before="Our"
      highlight-text="Delicious"
      title-after="Menu"
      subtitle="Explore our wide range of mouth-watering dishes, each crafted with fresh ingredients and passion."
      alignment="center"
      :narrow="true"
      header-class="pb-12"
    >
      <!-- Search Bar in the actions slot -->
      <div class="mt-8 max-w-2xl mx-auto w-full">
        <SearchBar @search="handleSearch" />
        
        <!-- Quick Filter Chips (Optional) -->
        <div class="flex flex-wrap justify-center gap-2 mt-4">
          <button
            v-for="category in topCategories"
            :key="category.id"
            @click="handleCategoryFilter(category.id)"
            class="px-4 py-2 rounded-full text-sm font-medium transition-all hover:scale-105"
            :class="[
              selectedCategory === category.id
                ? 'bg-brand-500 text-white shadow-lg'
                : 'bg-white/80 dark:bg-slate-800/80 text-gray-700 dark:text-gray-300 hover:bg-white dark:hover:bg-slate-800'
            ]"
          >
            {{ category.name }}
            <span class="ml-1 text-xs opacity-75">({{ category.itemsCount }})</span>
          </button>
          
          <button
            v-if="selectedCategory"
            @click="clearAllFilters"
            class="px-4 py-2 rounded-full text-sm font-medium bg-gray-100 dark:bg-slate-700 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-slate-600 transition-all"
          >
            Clear Filter
          </button>
        </div>
      </div>
    </PageHeader>

    <!-- Rest of your content remains largely the same -->
    <div  class="section-padding my-5 mx-5 md:my-8 md:mx-8 px-3 py-3 md:py-8 md:px-8">
      <div class="container-narrow">
        <!-- Active Filters Bar -->
        <div v-if="activeFilters.length" class="mb-8 p-4 bg-white/50 dark:bg-slate-800/50 backdrop-blur-sm rounded-xl border border-gray-100 dark:border-slate-700">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Filter class="w-4 h-4 text-gray-500" />
              <span class="text-sm text-gray-600 dark:text-gray-400">Active Filters:</span>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="filter in activeFilters"
                  :key="filter"
                  class="inline-flex items-center bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 px-3 py-1 rounded-full text-sm"
                >
                  {{ filter }}
                  <button @click="removeFilter(filter)" class="ml-2 hover:text-red-500 transition-colors">
                    <X class="w-3 h-3" />
                  </button>
                </span>
              </div>
            </div>
            <button 
              @click="clearAllFilters" 
              class="text-sm text-brand-500 hover:text-brand-600 dark:hover:text-brand-400 font-medium transition-colors"
            >
              Clear All
            </button>
          </div>
        </div>

        <div class="flex flex-col lg:flex-row gap-8">
          <!-- Sidebar Filters -->
          <div class="lg:w-1/4 space-y-6">
            <!-- Category Filter -->
            <div class="bg-white/50 dark:bg-slate-800/50 backdrop-blur-sm rounded-xl p-5 border border-gray-100 dark:border-slate-700">
              <div class="flex items-center justify-between mb-4">
                <h3 class="font-bold text-gray-900 dark:text-white">Categories</h3>
                <Filter class="w-4 h-4 text-gray-500" />
              </div>
              <CategoryFilter
                :categories="categoriesWithCount"
                :selected-category="selectedCategory"
                @filter="handleCategoryFilter"
              />
            </div>

            <!-- Price Filter -->
            <div class="bg-white/50 dark:bg-slate-800/50 backdrop-blur-sm rounded-xl p-5 border border-gray-100 dark:border-slate-700">
              <div class="flex items-center justify-between mb-4">
                <h3 class="font-bold text-gray-900 dark:text-white">Price Range</h3>
                <DollarSign class="w-4 h-4 text-gray-500" />
              </div>
              <PriceFilter
                :min-price="0"
                :max-price="100"
                :current-price="maxPrice"
                @price-change="handlePriceFilter"
              />
              <div class="mt-4 text-sm text-gray-600 dark:text-gray-400">
                <span class="font-medium">Max: ${{ maxPrice }}</span>
                <span class="float-right">{{ priceFilteredCount }} items</span>
              </div>
            </div>

            <!-- Popular Items -->
            <div v-if="popularItems.length" class="bg-white/50 dark:bg-slate-800/50 backdrop-blur-sm rounded-xl p-5 border border-gray-100 dark:border-slate-700">
              <div class="flex items-center justify-between mb-4">
                <h3 class="font-bold text-gray-900 dark:text-white">Popular Picks</h3>
                <TrendingUp class="w-4 h-4 text-brand-500" />
              </div>
              <div class="space-y-3">
                <div
                  v-for="item in popularItems"
                  :key="item.id"
                  @click="handleViewDetails(item)"
                  class="flex items-center gap-3 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 cursor-pointer transition-colors"
                >
                  <img
                    :src="item.imageUrl"
                    :alt="item.name"
                    class="w-12 h-12 object-cover rounded-lg"
                  />
                  <div class="flex-1 min-w-0">
                    <p class="font-medium text-gray-900 dark:text-white truncate">{{ item.name }}</p>
                    <p class="text-sm text-brand-500 font-bold">${{ item.price.toFixed(2) }}</p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Tip Section -->
            <div class="bg-gradient-to-r from-brand-50/80 to-accent-50/80 dark:from-brand-900/20 dark:to-accent-900/20 rounded-xl p-5 border border-brand-100 dark:border-brand-700/30">
              <div class="flex items-center space-x-3">
                <div class="w-10 h-10 rounded-full bg-brand-100 dark:bg-brand-900/30 flex items-center justify-center">
                  <Clock class="w-5 h-5 text-brand-600 dark:text-brand-400" />
                </div>
                <div>
                  <h4 class="font-bold text-brand-700 dark:text-brand-300 text-sm mb-1">Pre-order Tip</h4>
                  <p class="text-xs text-gray-600 dark:text-gray-400">
                    Pre-order items need 30+ minutes preparation time. Order early to ensure availability!
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Main Content -->
          <div class="lg:w-3/4">
            <!-- Results Summary -->
            <div class="mb-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
              <div>
                <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
                  {{ filteredItems.length }} Items Available
                </h2>
                <p v-if="searchQuery" class="text-gray-600 dark:text-gray-400 mt-1">
                  Showing results for "<span class="font-medium">{{ searchQuery }}</span>"
                </p>
              </div>
              
              <!-- Sort Options -->
              <div class="flex items-center gap-4 dark:text-white">
                <select
                  v-model="sortBy"
                  class="bg-white dark:bg-slate-800 border border-gray-200 dark:border-slate-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-transparent"
                >
                  <option value="name">Sort by Name</option>
                  <option value="price-low">Price: Low to High</option>
                  <option value="price-high">Price: High to Low</option>
                  <option value="popular">Most Popular</option>
                </select>
                
                <!-- View Toggle (Optional) -->
                <div class="hidden sm:flex items-center gap-1 p-1 bg-gray-100 dark:bg-slate-800 rounded-lg dark:text-white">
                  <button
                    @click="gridView = 'grid'"
                    class="p-2 rounded transition-colors"
                    :class="gridView === 'grid' ? 'bg-white dark:bg-slate-700 shadow' : 'hover:bg-gray-200 dark:hover:bg-slate-700'"
                  >
                    <Grid class="w-4 h-4" />
                  </button>
                  <button
                    @click="gridView = 'list'"
                    class="p-2 rounded transition-colors"
                    :class="gridView === 'list' ? 'bg-white dark:bg-slate-700 shadow' : 'hover:bg-gray-200 dark:hover:bg-slate-700'"
                  >
                    <List class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>

            <!-- Menu Grid -->
            <MenuGrid
              :filtered-items="paginatedItems"
              :is-loading="menuStore.isLoading"
              :current-page="currentPage"
              :total-pages="totalPages"
              :view="gridView"
              @add-to-cart="handleAddToCart"
              @view-details="handleViewDetails"
              @page-change="handlePageChange"
            />

            <!-- Empty State -->
            <div v-if="filteredItems.length === 0 && !menuStore.isLoading" class="text-center py-16">
              <Search class="w-24 h-24 text-gray-300 dark:text-gray-700 mx-auto mb-6" />
              <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">No items found</h3>
              <p class="text-gray-600 dark:text-gray-400 mb-6 max-w-md mx-auto">
                No menu items match your current filters. Try adjusting your search or filters.
              </p>
              <button
                @click="clearAllFilters"
                class="btn-primary"
              >
                Clear All Filters
              </button>
            </div>

            <!-- Cart Summary (Mobile) -->
            <div class="lg:hidden fixed bottom-0 left-0 right-0 bg-white dark:bg-slate-800 border-t border-gray-200 dark:border-slate-700 shadow-lg p-4 z-50">
              <div class="flex justify-between items-center">
                <div>
                  <div class="text-sm text-gray-600 dark:text-gray-400">Total</div>
                  <div class="font-bold text-xl text-brand-500">${{ cart.totalPrice.toFixed(2) }}</div>
                  <div class="text-xs text-gray-500">{{ cart.totalItems }} items</div>
                </div>
                <div class="flex gap-3">
                  <button
                    @click="cart.toggleCart()"
                    class="btn-primary px-6"
                  >
                    View Cart
                  </button>
                  <button
                    v-if="cart.totalItems > 0"
                    @click="navigateTo('/order')"
                    class="btn-secondary px-6"
                  >
                    Checkout
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Item Details Modal -->
    <Teleport to="body">
      <ItemDetailsModal
        v-if="selectedItem"
        :item="selectedItem"
        @close="selectedItem = null"
        @add-to-cart="handleAddToCart"
      />
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMenuStore } from '../../stores/menu'
import { useCartStore } from '../../stores/cart'
import { navigateTo } from 'nuxt/app'
import { 
  X, 
  Filter, 
  DollarSign, 
  TrendingUp, 
  Clock,
  Grid,
  List,
  Search 
} from 'lucide-vue-next'
import type { MenuItem } from '../../types/menu'
import { useNuxtApp } from 'nuxt/app'
import { useRoute } from 'vue-router'


// Components
import PageHeader from '../../components/layout/PageHeader.vue'
import SearchBar from '../../components/menu/SearchBar.vue'
import CategoryFilter from '../../components/menu/CategoryFilter.vue'
import PriceFilter from '../../components/menu/PriceFilter.vue'
import MenuGrid from '../../components/menu/MenuGrid.vue'
import ItemDetailsModal from '../../components/menu/ItemDetailsModal.vue'

const route = useRoute()


// Stores
const menuStore = useMenuStore()
const cart = useCartStore()

// State
const selectedCategory = ref<string | null>(null)
const searchQuery = ref('')
const maxPrice = ref(100)
const selectedItem = ref<MenuItem | null>(null)
const currentPage = ref(1)
const sortBy = ref('name')
const gridView = ref<'grid' | 'list'>('grid')
const itemsPerPage = 12

// Computed
const categoriesWithCount = computed(() => {
  return menuStore.categories.map(category => ({
    ...category,
    itemsCount: menuStore.getItemsByCategory(category.id).length
  }))
})

const topCategories = computed(() => {
  return categoriesWithCount.value
    .sort((a, b) => b.itemsCount - a.itemsCount)
    .slice(0, 5)
})

const filteredItems = computed(() => {
  let items = menuStore.menuItems

  // Filter by category
  if (selectedCategory.value) {
    items = items.filter(item => item.categoryId === selectedCategory.value)
  }

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    items = items.filter(item =>
      item.name.toLowerCase().includes(query) ||
      item.description.toLowerCase().includes(query) ||
      item.tags?.some(tag => tag.toLowerCase().includes(query))
    )
  }

  // Filter by price
  items = items.filter(item => item.price <= maxPrice.value)

  // Apply sorting
  items = [...items].sort((a, b) => {
    switch (sortBy.value) {
      case 'price-low':
        return a.price - b.price
      case 'price-high':
        return b.price - a.price
      case 'popular':
        return (b.popularity || 0) - (a.popularity || 0)
      default: // 'name'
        return a.name.localeCompare(b.name)
    }
  })

  return items
})

const priceFilteredCount = computed(() => {
  return menuStore.menuItems.filter(item => item.price <= maxPrice.value).length
})

const popularItems = computed(() => {
  return [...menuStore.menuItems]
    .sort((a, b) => (b.popularity || 0) - (a.popularity || 0))
    .slice(0, 3)
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
  if (searchQuery.value) filters.push(`"${searchQuery.value}"`)
  if (maxPrice.value < 100) filters.push(`Under $${maxPrice.value}`)
  return filters
})

// Methods
const handleCategoryFilter = (categoryId: string | null) => {
  selectedCategory.value = categoryId
  currentPage.value = 1
}

const handleSearch = (query: string) => {
  searchQuery.value = query
  currentPage.value = 1
}

const handlePriceFilter = (price: number) => {
  maxPrice.value = price
  currentPage.value = 1
}

const handleAddToCart = (item: MenuItem) => {
  cart.addItem(item)
  
  if (import.meta.client) {
    const { $toast } = useNuxtApp()
    
    if ($toast && typeof ($toast as any).success === 'function') {
      ($toast as any).success(`Added ${item.name} to cart`)
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
  if (filter.startsWith('"')) {
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
  sortBy.value = 'name'
}

// Lifecycle
// onMounted(async () => {
//   await Promise.all([
//     menuStore.fetchCategories(),
//     menuStore.fetchMenuItems()
//   ])
  
//   // Reset scroll position
//   window.scrollTo({ top: 0, behavior: 'instant' })
// })
// Replace your current onMounted with this:
onMounted(async () => {
  // Read category from URL query parameter
  if (route.query.category) {
    selectedCategory.value = route.query.category as string
  }
  
  await Promise.all([
    menuStore.fetchCategories(),
    menuStore.fetchMenuItems()
  ])
  
  // Reset scroll position
  window.scrollTo({ top: 0, behavior: 'instant' })
})
</script>

<style scoped>
/* Custom scrollbar for price filter */
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #555;
}

/* Smooth transitions */
.smooth-transition {
  transition: all 0.3s ease-in-out;
}

/* Ensure the mobile cart doesn't overlap with modals */
.fixed {
  z-index: 40;
}
</style>