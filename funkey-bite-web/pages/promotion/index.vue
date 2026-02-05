<!-- pages/promotions/index.vue -->
<template>
    <div>
      <PageHeader
        title="All Promotions"
        subtitle="Browse all active and upcoming promotions"
        variant="gradient"
        alignment="center"
      />
  
      <div class="section-padding">
        <div class="container-narrow">
          <!-- Filter Tabs -->
          <div class="flex flex-wrap gap-4 mb-8">
            <button
              v-for="filter in filters"
              :key="filter.id"
              @click="activeFilter = filter.id"
              class="px-4 py-2 rounded-full font-medium transition-colors"
              :class="[
                activeFilter === filter.id
                  ? 'bg-brand-500 text-white'
                  : 'bg-gray-100 dark:bg-slate-800 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-slate-700'
              ]"
            >
              {{ filter.label }} ({{ filteredPromotions(filter.id).length }})
            </button>
          </div>
  
          <!-- Promotions Grid -->
          <div v-if="filteredPromotions(activeFilter).length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div
              v-for="promotion in filteredPromotions(activeFilter)"
              :key="promotion.id"
              class="bg-white dark:bg-slate-800 rounded-xl shadow-lg overflow-hidden hover:shadow-xl transition-shadow"
            >
              <!-- Promotion Image -->
              <div class="relative h-48 overflow-hidden">
                <img 
                  :src="getPromotionImage(promotion)"
                  :alt="promotion.title"
                  class="w-full h-full object-cover"
                />
                <!-- Discount Badge -->
                <div class="absolute top-4 right-4 bg-red-500 text-white px-3 py-1 rounded-full font-bold">
                  {{ formatDiscountValue(promotion.discountValue, promotion.discountType) }}
                </div>
                <!-- Time Badge -->
                <div class="absolute bottom-4 left-4 bg-black/70 text-white px-3 py-1 rounded-full text-sm">
                  {{ getTimeRemaining(promotion.validUntil) }}
                </div>
              </div>
  
              <!-- Promotion Details -->
              <div class="p-6">
                <div class="flex items-center justify-between mb-3">
                  <span class="px-3 py-1 bg-brand-100 dark:bg-brand-900/30 text-brand-700 dark:text-brand-300 rounded-full text-xs font-medium">
                    {{ formatDiscountType(promotion.discountType) }}
                  </span>
                  <span class="text-sm text-gray-500">
                    {{ formatDate(promotion.validFrom) }} - {{ formatDate(promotion.validUntil) }}
                  </span>
                </div>
  
                <h3 class="text-xl font-bold mb-2">{{ promotion.title }}</h3>
                <p class="text-gray-600 dark:text-gray-400 mb-4 line-clamp-2">
                  {{ promotion.description }}
                </p>
  
                <!-- Applicable Items -->
                <div v-if="promotion.applicableItems?.length" class="mb-4">
                  <div class="text-sm text-gray-500 mb-2">Applies to {{ promotion.applicableItems.length }} items</div>
                  <div class="flex flex-wrap gap-2">
                    <span
                      v-for="itemId in promotion.applicableItems.slice(0, 3)"
                      :key="itemId"
                      class="px-2 py-1 bg-gray-100 dark:bg-slate-700 rounded text-xs"
                    >
                      Item #{{ itemId }}
                    </span>
                    <span v-if="promotion.applicableItems.length > 3" class="text-xs text-gray-500">
                      +{{ promotion.applicableItems.length - 3 }} more
                    </span>
                  </div>
                </div>
  
                <!-- Actions -->
                <div class="flex gap-3">
                  <NuxtLink
                    :to="`/promotions/${promotion.id}`"
                    class="flex-1 text-center px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
                  >
                    View Details
                  </NuxtLink>
                  <button
                    @click="claimPromotion(promotion)"
                    class="flex-1 px-4 py-2 bg-brand-500 text-white rounded-lg hover:bg-brand-600 transition-colors"
                  >
                    Claim Offer
                  </button>
                </div>
              </div>
            </div>
          </div>
  
          <!-- Empty State -->
          <div v-else class="text-center py-16">
            <div class="text-6xl mb-6">🎁</div>
            <h3 class="text-2xl font-bold mb-4">No Promotions Available</h3>
            <p class="text-gray-600 dark:text-gray-400 mb-6">
              Check back later for new promotions!
            </p>
            <NuxtLink to="/menu" class="btn-primary">
              Browse Menu
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed } from 'vue'
   import { mockPromotions, mockMenuItems } from '../../utils/mockData'
  import type { Promotion, MenuItem } from '../../types/menu'

  import { 
    formatDiscountType, 
    formatDiscountValue,
    formatDate,
    formatDateFull,
    getTimeRemaining,
    getDiscountedPrice  
} from '../../utils/helperFunctions'
import { useRoute, useRouter } from 'vue-router'
  
  // Use the same helper functions from the details page...
  // (Copy the formatDiscountType, formatDiscountValue, etc. functions)

  const route = useRoute()
  const router = useRouter()
  const promotionId = route.params.id as string
  
  const filters = [
    { id: 'all', label: 'All Promotions' },
    { id: 'active', label: 'Active Now' },
    { id: 'upcoming', label: 'Upcoming' },
    { id: 'ending', label: 'Ending Soon' }
  ]


  // Find promotion
  const promotion = computed(() => {
    return mockPromotions.find(p => p.id === promotionId)
  })
  
  // Get applicable items
  const applicableItems = computed(() => {
    if (!promotion.value?.applicableItems?.length) {
      return mockMenuItems.slice(0, 6) // Show first 6 items if no specific items
    }
    
    return mockMenuItems.filter(item => 
      promotion.value?.applicableItems?.includes(item.id)
    )
  })
  
  const activeFilter = ref('all')
  
  const filteredPromotions = (filterId: string) => {
    // Implement filtering logic based on dates
    return mockPromotions
  }

  const claimPromotion = (promo: Promotion) => {
    // Implement claim logic
    console.log('Claiming promotion:', promo.id)
    // navigateTo(`/order?promotion=${promo.id}`)
  }

  const getPromotionImage = (promo: Promotion) => {
    if (applicableItems.value.length > 0) {
      return applicableItems.value[0].imageUrl
    }
    return '/images/promotion-default.jpg'
  }
  </script>