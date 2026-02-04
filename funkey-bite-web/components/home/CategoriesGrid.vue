<template>
    <section class="section-padding mt-4 mb-4 md:mt-8 md:mb-8 mx-3 px-3 py-3">
      <div class="container-narrow">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white mb-4">
            Our <span class="text-gradient">Delicious</span> Categories
          </h2>
          <p class="text-gray-600 dark:text-gray-300 max-w-2xl mx-auto">
            Explore our wide range of mouth-watering categories, each crafted with fresh ingredients and passion.
          </p>
        </div>
  
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4 md:gap-6">
          <div
            v-for="category in categories"
            :key="category.id"
            class="group relative bg-white dark:bg-slate-800 rounded-xl p-6 text-center shadow-lg hover:shadow-2xl transition-all duration-300 hover:scale-105 cursor-pointer border border-gray-100 dark:border-slate-700"
            @click="navigateToMenu(category.id)"
          >
            <div
              class="w-16 h-16 mx-auto mb-4 bg-gradient-to-br from-brand-100 to-brand-50 dark:from-brand-900/30 dark:to-brand-800/20 rounded-full flex items-center justify-center group-hover:from-brand-200 group-hover:to-brand-100 transition-all duration-300"
            >
              <div
                class="w-10 h-10 bg-brand-500 rounded-full flex items-center justify-center transform group-hover:scale-110 transition-transform"
              >
                <UtensilsCrossed class="w-6 h-6 text-white" />
              </div>
            </div>
            <h3 class="font-bold text-gray-900 dark:text-white mb-2">{{ category.name }}</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 line-clamp-2">{{ category.description }}</p>
            <!-- Show active status badge -->
            <span 
              v-if="!category.isActive"
              class="absolute top-2 right-2 text-xs bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300 px-2 py-1 rounded-full"
            >
              Coming Soon
            </span>
          </div>
        </div>
      </div>
    </section>
  </template>
  
  <script setup lang="ts">
  import { UtensilsCrossed } from 'lucide-vue-next'
  import { navigateTo } from 'nuxt/app'
  import type { MenuCategory } from '../../types/menu'
  
  defineProps<{
    categories: MenuCategory[]
  }>()
  
  const navigateToMenu = (categoryId: string) => {
    navigateTo(`/menu?category=${categoryId}`)
  }
  </script>
  
  <style scoped>
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    line-clamp: 2;
  }
  </style>