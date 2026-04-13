<template>
    <div
      class="group relative bg-white dark:bg-slate-800 rounded-2xl overflow-hidden shadow-lg hover:shadow-2xl transition-all duration-300 border border-gray-100 dark:border-slate-700 hover:scale-[1.02]"
      :class="{ 'opacity-70': !item.isAvailable }"
    >
      
      <PreOrderBadge v-if="item.isPreOrder" />
  
      
      <div class="relative h-48 overflow-hidden">
        <img
          :src="item.imageUrl"
          :alt="item.name"
          class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
          loading="lazy"
        />
        
        <div v-if="!item.isAvailable" class="absolute inset-0 bg-black/50 flex items-center justify-center">
          <span class="text-white font-bold text-lg bg-red-500/80 px-4 py-2 rounded-full">Currently Unavailable</span>
        </div>
      </div>
  
      
      <div class="p-6">
        <div class="flex justify-between items-start mb-3">
          <div>
            <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-1">{{ item.name }}</h3>
            <div class="flex items-center space-x-2 mb-2">
              <span class="text-sm text-gray-500 dark:text-gray-400">{{ item.preparationTime }} mins</span>
              <span class="text-xs text-brand-500 bg-brand-50 dark:bg-brand-900/30 px-2 py-1 rounded-full" 
                    v-if="item.tags && item.tags.length">
                {{ item.tags[0] }}
              </span>
            </div>
          </div>
          <span class="text-2xl font-bold text-brand-500">${{ item.price.toFixed(2) }}</span>
        </div>
  
        <p class="text-gray-600 dark:text-gray-300 mb-4 line-clamp-2">{{ item.description }}</p>
  
        
        <div v-if="item.nutritionalInfo" class="mb-4">
          <button
            @click="showNutrition = !showNutrition"
            class="text-sm text-gray-500 hover:text-brand-500 flex items-center"
          >
            <ChevronDown class="w-4 h-4 mr-1 transition-transform" :class="{ 'rotate-180': showNutrition }" />
            Nutritional Info
          </button>
          <div v-if="showNutrition" class="mt-2 grid grid-cols-4 gap-2 text-sm">
            <div class="text-center p-2 bg-gray-50 dark:bg-slate-700 rounded">
              <div class="font-bold">{{ item.nutritionalInfo.calories }}</div>
              <div class="text-gray-500 text-xs">Cal</div>
            </div>
            <div class="text-center p-2 bg-gray-50 dark:bg-slate-700 rounded">
              <div class="font-bold">{{ item.nutritionalInfo.protein }}g</div>
              <div class="text-gray-500 text-xs">Protein</div>
            </div>
            <div class="text-center p-2 bg-gray-50 dark:bg-slate-700 rounded">
              <div class="font-bold">{{ item.nutritionalInfo.carbs }}g</div>
              <div class="text-gray-500 text-xs">Carbs</div>
            </div>
            <div class="text-center p-2 bg-gray-50 dark:bg-slate-700 rounded">
              <div class="font-bold">{{ item.nutritionalInfo.fat }}g</div>
              <div class="text-gray-500 text-xs">Fat</div>
            </div>
          </div>
        </div>
  
        
        <div class="flex space-x-3">
          <button
            @click="$emit('addToCart', item)"
            :disabled="!item.isAvailable"
            class="flex-1 btn-primary text-center justify-center dark:text-white"
            :class="{ 'opacity-50 cursor-not-allowed': !item.isAvailable }"
          >
            <ShoppingCart class="w-4 h-4 inline mr-2 dark:text-white" />
            Add to Cart
          </button>
          <button
            @click="$emit('viewDetails', item)"
            class="px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
          >
            <Eye class="w-4 h-4 dark:text-white" />
          </button>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import type { MenuItem } from '../../types/menu'
  import { ShoppingCart, Eye, ChevronDown } from 'lucide-vue-next'
  import PreOrderBadge from './PreOrderBadge.vue'
  import { ref } from 'vue'
  
  defineProps<{
    item: MenuItem
  }>()
  
  defineEmits<{
    addToCart: [item: MenuItem]
    viewDetails: [item: MenuItem]
  }>()
  
  const showNutrition = ref(false)
  </script>
  
  <style scoped>
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    line-clamp: 2; 
    display: -moz-box;
    -moz-box-orient: vertical;
  }
  </style>