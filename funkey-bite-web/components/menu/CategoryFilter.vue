<template>
    <div class="space-y-4">
      <h3 class="font-bold text-gray-900 dark:text-white">Categories</h3>
      <div class="space-y-2">
        <button
          v-for="category in categories"
          :key="category.id"
          @click="$emit('filter', category.id)"
          class="w-full text-left px-4 py-3 rounded-lg transition-colors"
          :class="[
            selectedCategories && selectedCategories.includes(category.id)
              ? 'bg-brand-500 text-white'
              : 'bg-gray-100 dark:bg-slate-800 hover:bg-gray-200 dark:hover:bg-slate-700 text-gray-700 dark:text-gray-300'
          ]"
        >
          <div class="flex justify-between items-center">
            <span>{{ category.name }}</span>
            <span v-if="category.itemsCount" class="text-sm opacity-80">
              ({{ category.itemsCount }})
            </span>
          </div>
        </button>
      </div>
      <button
        v-if="selectedCategories && selectedCategories.length > 0"
        @click="$emit('clearAll')"
        class="w-full text-center text-brand-500 hover:text-brand-600 text-sm font-medium py-2"
      >
        Clear Filters
      </button>
    </div>
  </template>
  
  <script setup lang="ts">
  import type { MenuCategory } from '../../types/menu'
  
  interface ExtendedCategory extends MenuCategory {
    itemsCount?: number
  }
  
  defineProps<{
    categories: ExtendedCategory[]
    selectedCategories?: string[]
  }>()

  defineEmits<{
    filter: [categoryId: string]
    clearAll: []
  }>()
  </script>