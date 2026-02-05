<template>
    <section class="section-padding bg-gradient-to-b mt-8 mb-8 md:mt-12 mx-5 py-3 px-3 md:mb-12 from-white to-gray-50 dark:from-slate-900 dark:to-slate-800">
      <div class="container-narrow">
        <div class="flex justify-between items-center mb-12">
          <div>
            <h2 class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white">
              <span class="text-gradient">Featured</span> Items
            </h2>
            <p class="text-gray-600 dark:text-gray-300 mt-2">Customer favorites you'll love</p>
          </div>
          <NuxtLink to="/menu" class="btn-secondary">
            View All Menu
            <ArrowRight class="w-4 h-4 inline ml-2" />
          </NuxtLink>
        </div>
  
        <!-- Loading Skeleton -->
        <div v-if="isLoading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <div v-for="i in 6" :key="i" class="animate-pulse">
            <div class="h-48 bg-gray-200 dark:bg-slate-700 rounded-2xl mb-4"></div>
            <div class="h-4 bg-gray-200 dark:bg-slate-700 rounded mb-2"></div>
            <div class="h-4 bg-gray-200 dark:bg-slate-700 rounded w-2/3"></div>
          </div>
        </div>
  
        <!-- Items Grid -->
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 px-3 md:px-8 py-3 md:py-8">
          <FoodCard v-for="item in items" :key="item.id" :item="item" />
        </div>
      </div>
    </section>
  </template>
  
  <script setup lang="ts">
  import { ArrowRight } from 'lucide-vue-next'
  import FoodCard from '../menu/FoodCard.vue'
  import type { MenuItem } from '../../types/menu'

  
  defineProps<{
    items: MenuItem[]
    isLoading: boolean
  }>()
  </script>