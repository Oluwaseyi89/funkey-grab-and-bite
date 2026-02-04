<template>
    <div class="flex items-center gap-6 p-4 border border-gray-200 dark:border-slate-700 rounded-2xl hover:bg-gray-50 dark:hover:bg-slate-800 transition-colors">
      <!-- Image -->
      <div class="w-24 h-24 flex-shrink-0">
        <img 
          :src="item.imageUrl" 
          :alt="item.name"
          class="w-full h-full object-cover rounded-xl"
        />
      </div>
      
      <!-- Details -->
      <div class="flex-1 min-w-0">
        <div class="flex justify-between items-start">
          <div>
            <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-2">{{ item.name }}</h3>
            <p class="text-gray-600 dark:text-gray-400 text-sm mb-3 line-clamp-2">{{ item.description }}</p>
          </div>
          <div class="text-right">
            <div class="text-2xl font-bold text-brand-500 mb-2">${{ item.price.toFixed(2) }}</div>
            <div class="flex items-center text-sm text-gray-500">
              <Clock class="w-4 h-4 mr-1" />
              {{ item.preparationTime }} mins
            </div>
          </div>
        </div>
        
        <!-- Tags -->
        <div v-if="item.tags?.length" class="flex flex-wrap gap-2 mb-4">
          <span
            v-for="tag in item.tags"
            :key="tag"
            class="px-3 py-1 bg-gray-100 dark:bg-slate-700 text-gray-700 dark:text-gray-300 rounded-full text-xs"
          >
            {{ tag }}
          </span>
        </div>
        
        <!-- Actions -->
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-500">
            <span v-if="item.isPreOrder" class="text-amber-600 dark:text-amber-400">
              ⏰ Pre-order required
            </span>
            <span v-else class="text-green-600 dark:text-green-400">
              ✅ Available now
            </span>
          </div>
          <div class="flex gap-3">
            <button
              @click="$emit('viewDetails', item)"
              class="px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
            >
              Details
            </button>
            <button
              @click="$emit('addToCart', item)"
              :disabled="!item.isAvailable"
              class="px-4 py-2 bg-brand-500 text-white rounded-lg hover:bg-brand-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              Add to Cart
            </button>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import type { MenuItem } from '../../types/menu'
  import { Clock } from 'lucide-vue-next'
  
  defineProps<{
    item: MenuItem
  }>()
  
  defineEmits<{
    addToCart: [item: MenuItem]
    viewDetails: [item: MenuItem]
  }>()
  </script>