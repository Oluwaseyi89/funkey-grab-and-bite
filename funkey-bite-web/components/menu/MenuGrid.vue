<template>
  <div class="space-y-8">
    
    <div v-if="isLoading">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div v-for="i in 6" :key="i" class="animate-pulse">
          <div class="h-48 bg-gray-200 dark:bg-slate-700 rounded-2xl mb-4"></div>
          <div class="h-4 bg-gray-200 dark:bg-slate-700 rounded mb-2"></div>
          <div class="h-4 bg-gray-200 dark:bg-slate-700 rounded w-2/3"></div>
        </div>
      </div>
    </div>

    
    <div v-else-if="!filteredItems.length" class="text-center py-12 px-4">
      <div class="max-w-md mx-auto">
        <Package class="w-20 h-20 text-gray-300 dark:text-gray-600 mx-auto mb-6" />
        <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">No items found</h3>
        <p class="text-gray-600 dark:text-gray-400 mb-8">
          Try adjusting your search criteria or filters to find what you're looking for.
        </p>
        <button 
          @click="$emit('resetFilters')" 
          class="btn-primary px-6 py-3"
        >
          Reset All Filters
        </button>
      </div>
    </div>

    
    <div v-else>
      
      <div class="mb-8 flex justify-between items-center">
        <div class="text-gray-600 dark:text-gray-400">
          Showing <span class="font-bold text-gray-900 dark:text-white">{{ filteredItems.length }}</span> 
          {{ filteredItems.length === 1 ? 'item' : 'items' }}
        </div>
        <div v-if="totalPages > 1" class="text-sm text-gray-500 dark:text-gray-400">
          Page {{ currentPage }} of {{ totalPages }}
        </div>
      </div>

       
       <div 
        v-if="view === 'grid'"
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-8 mr-5"
      >
        <MenuCard
          v-for="item in filteredItems"
          :key="item.id"
          :item="item"
          @add-to-cart="$emit('addToCart', item)"
          @view-details="$emit('viewDetails', item)"
        />
      </div>

      
      <div 
        v-if="view === 'list'"
        class="space-y-6"
      >
        <MenuItemList
          v-for="item in filteredItems"
          :key="item.id"
          :item="item"
          @add-to-cart="$emit('addToCart', item)"
          @view-details="$emit('viewDetails', item)"
        />
      </div>
    </div>

    
    <div v-if="totalPages > 1 && !isLoading && filteredItems.length" class="pt-8 border-t border-gray-200 dark:border-slate-700">
      <div class="flex flex-col sm:flex-row justify-between items-center gap-4">
        <div class="text-sm text-gray-600 dark:text-gray-400">
          Showing {{ Math.min(filteredItems.length, itemsPerPage) }} of {{ filteredItems.length }} items
        </div>
        
        <div class="flex items-center space-x-2">
          <button
            @click="$emit('pageChange', currentPage - 1)"
            :disabled="currentPage === 1"
            class="px-4 py-2 rounded-lg border border-gray-300 dark:border-slate-600 hover:bg-gray-50 dark:hover:bg-slate-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            ← Previous
          </button>
          
          <div class="flex items-center space-x-2">
            <button
              v-for="page in visiblePages"
              :key="page"
              @click="$emit('pageChange', page)"
              :class="[
                'w-10 h-10 rounded-lg transition-colors',
                currentPage === page
                  ? 'bg-brand-500 text-white'
                  : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-slate-700'
              ]"
            >
              {{ page }}
            </button>
            <span v-if="hasMorePages" class="text-gray-500">...</span>
          </div>
          
          <button
            @click="$emit('pageChange', currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="px-4 py-2 rounded-lg border border-gray-300 dark:border-slate-600 hover:bg-gray-50 dark:hover:bg-slate-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            Next →
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { MenuItem } from '../../types/menu'
import { Package } from 'lucide-vue-next'
import MenuCard from './MenuCard.vue'
import MenuItemList from './MenuItemList.vue';
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  filteredItems: MenuItem[]
  isLoading: boolean
  currentPage: number
  totalPages: number
  itemsPerPage?: number
  view?: 'grid' | 'list' // Add this
}>(), {
  itemsPerPage: 12
})

const emit = defineEmits<{
  addToCart: [item: MenuItem]
  viewDetails: [item: MenuItem]
  resetFilters: []
  pageChange: [page: number]
}>()

const visiblePages = computed(() => {
  const pages: number[] = []
  const maxVisible = 5
  
  if (props.totalPages <= maxVisible) {
    for (let i = 1; i <= props.totalPages; i++) {
      pages.push(i)
    }
  } else {
    let start = Math.max(1, props.currentPage - 2)
    let end = Math.min(props.totalPages, start + maxVisible - 1)
    
    if (end - start + 1 < maxVisible) {
      start = Math.max(1, end - maxVisible + 1)
    }
    
    for (let i = start; i <= end; i++) {
      pages.push(i)
    }
  }
  
  return pages
})

const hasMorePages = computed(() => {
  return props.totalPages > visiblePages.value.length && 
         visiblePages.value[visiblePages.value.length - 1] < props.totalPages
})
</script>