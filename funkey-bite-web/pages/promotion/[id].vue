<!-- pages/promotions/[id].vue -->
<template>
    <div>
      <!-- Back Button -->
      <div class="container-narrow py-6">
        <button @click="goBack" class="flex items-center text-gray-600 dark:text-gray-400 hover:text-brand-500">
          ← Back to Promotions
        </button>
      </div>
  
      <!-- Promotion Header -->
      <PageHeader
        :title="promotion?.title"
        :subtitle="promotion?.description"
        variant="gradient"
        alignment="center"
        header-class="pb-12"
      >
        <!-- Promotion Badge -->
        <div class="mt-8">
          <span class="inline-block bg-white/20 dark:bg-slate-800/50 px-6 py-2 rounded-full text-white font-bold text-lg">
            {{ formatDiscountValue(promotion?.discountValue, promotion?.discountType) }}
          </span>
        </div>
      </PageHeader>
  
      <!-- Promotion Details -->
      <div class="section-padding">
        <div class="container-narrow">
          <div v-if="promotion" class="bg-white dark:bg-slate-800 rounded-2xl shadow-xl overflow-hidden">
            <!-- Hero Image -->
            <div class="relative h-64 md:h-96">
              <img 
                :src="getPromotionImage(promotion)"
                :alt="promotion.title"
                class="w-full h-full object-cover"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent"></div>
              <div class="absolute bottom-6 left-6 right-6 text-white">
                <div class="text-4xl md:text-5xl font-bold mb-2">{{ promotion.title }}</div>
                <div class="text-xl opacity-90">{{ promotion.description }}</div>
              </div>
            </div>
  
            <!-- Promotion Details Grid -->
            <div class="p-8">
              <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-8">
                <!-- Discount Info -->
                <div class="bg-gray-50 dark:bg-slate-700/50 rounded-xl p-6">
                  <h3 class="text-xl font-bold mb-4">🎁 Discount Details</h3>
                  <div class="space-y-4">
                    <div>
                      <div class="text-sm text-gray-600 dark:text-gray-400">Type</div>
                      <div class="text-2xl font-bold text-brand-500">
                        {{ formatDiscountType(promotion.discountType) }}
                      </div>
                    </div>
                    <div>
                      <div class="text-sm text-gray-600 dark:text-gray-400">Value</div>
                      <div class="text-3xl font-bold">
                        {{ formatDiscountValue(promotion.discountValue, promotion.discountType) }}
                      </div>
                    </div>
                  </div>
                </div>
  
                <!-- Validity -->
                <div class="bg-gray-50 dark:bg-slate-700/50 rounded-xl p-6">
                  <h3 class="text-xl font-bold mb-4">📅 Validity Period</h3>
                  <div class="space-y-4">
                    <div>
                      <div class="text-sm text-gray-600 dark:text-gray-400">Starts</div>
                      <div class="text-xl font-bold">{{ formatDateFull(promotion.validFrom) }}</div>
                    </div>
                    <div>
                      <div class="text-sm text-gray-600 dark:text-gray-400">Ends</div>
                      <div class="text-xl font-bold">{{ formatDateFull(promotion.validUntil) }}</div>
                    </div>
                    <div class="pt-4 border-t">
                      <div class="text-sm text-gray-600 dark:text-gray-400">Time Remaining</div>
                      <div class="text-lg font-bold text-green-600 dark:text-green-400">
                        {{ getTimeRemaining(promotion.validUntil) }}
                      </div>
                    </div>
                  </div>
                </div>
  
                <!-- Status -->
                <div class="bg-gray-50 dark:bg-slate-700/50 rounded-xl p-6">
                  <h3 class="text-xl font-bold mb-4">⚡ Quick Actions</h3>
                  <div class="space-y-4">
                    <div class="flex items-center gap-2">
                      <div class="w-3 h-3 rounded-full bg-green-500 animate-pulse"></div>
                      <span class="font-medium">Active</span>
                    </div>
                    <div class="space-y-3">
                      <button
                        @click="claimPromotion(promotion)"
                        class="w-full btn-primary py-3 font-bold"
                      >
                        🎁 Claim This Offer
                      </button>
                      <button
                        @click="viewApplicableItems"
                        class="w-full btn-secondary py-3"
                      >
                        👀 View All Items
                      </button>
                    </div>
                  </div>
                </div>
              </div>
  
              <!-- Applicable Items Section -->
              <div v-if="applicableItems.length > 0" class="mb-12">
                <h3 class="text-2xl font-bold mb-6">🎯 Applicable Items</h3>
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                  <div
                    v-for="item in applicableItems"
                    :key="item.id"
                    class="border border-gray-200 dark:border-slate-700 rounded-xl p-4 hover:shadow-lg transition-shadow"
                  >
                    <div class="flex gap-4">
                      <div class="w-24 h-24 rounded-lg overflow-hidden flex-shrink-0">
                        <img 
                          :src="item.imageUrl" 
                          :alt="item.name"
                          class="w-full h-full object-cover"
                        />
                      </div>
                      <div class="flex-1">
                        <h4 class="font-bold text-lg mb-1">{{ item.name }}</h4>
                        <p class="text-gray-600 dark:text-gray-400 text-sm mb-2 line-clamp-2">
                          {{ item.description }}
                        </p>
                        <div class="flex items-center justify-between">
                          <div>
                            <div class="text-sm text-gray-500">Original Price</div>
                            <div class="font-bold text-gray-900 dark:text-white">
                              ${{ item.price.toFixed(2) }}
                            </div>
                          </div>
                          <div>
                            <div class="text-sm text-green-600 dark:text-green-400">Discounted</div>
                            <div class="font-bold text-brand-500">
                              {{ getDiscountedPrice(promotion, item.price) }}
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
  
              <!-- Terms & Conditions -->
              <div class="bg-gray-50 dark:bg-slate-700/30 rounded-xl p-6">
                <h3 class="text-xl font-bold mb-4">📋 Terms & Conditions</h3>
                <ul class="space-y-2 text-gray-700 dark:text-gray-300">
                  <li class="flex items-start gap-2">
                    <span class="text-green-500 mt-1">✓</span>
                    <span>Offer valid from {{ formatDate(promotion.validFrom) }} to {{ formatDate(promotion.validUntil) }}</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <span class="text-green-500 mt-1">✓</span>
                    <span>Discount applies before tax and delivery fees</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <span class="text-green-500 mt-1">✓</span>
                    <span>Cannot be combined with other promotions</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <span class="text-green-500 mt-1">✓</span>
                    <span>Limited to one use per customer</span>
                  </li>
                  <li class="flex items-start gap-2">
                    <span class="text-green-500 mt-1">✓</span>
                    <span>Management reserves the right to modify or cancel promotion</span>
                  </li>
                </ul>
              </div>
            </div>
          </div>
  
          <!-- Not Found State -->
          <div v-else class="text-center py-16">
            <div class="text-6xl mb-6">😕</div>
            <h3 class="text-2xl font-bold mb-4">Promotion Not Found</h3>
            <p class="text-gray-600 dark:text-gray-400 mb-6">
              This promotion may have expired or been removed.
            </p>
            <button @click="goBack" class="btn-primary">
              Back to All Promotions
            </button>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { computed } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
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
  
  // Components
  import PageHeader from '../../components/layout/PageHeader.vue'
  
  const route = useRoute()
  const router = useRouter()
  const promotionId = route.params.id as string
  
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
  


  const getPromotionImage = (promo: Promotion) => {
    if (applicableItems.value.length > 0) {
      return applicableItems.value[0].imageUrl
    }
    return '/images/promotion-default.jpg'
  }
  

  
  // Navigation methods
  const goBack = () => {
    router.push('/promotion')
  }
  
  const claimPromotion = (promo: Promotion) => {
    // Implement claim logic
    console.log('Claiming promotion:', promo.id)
    // navigateTo(`/order?promotion=${promo.id}`)
  }
  
  const viewApplicableItems = () => {
    if (applicableItems.value.length > 0) {
      // Navigate to menu with filtered items
      router.push(`/menu?promotion=${promotionId}`)
    }
  }
  </script>