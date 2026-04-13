<template>

  <section 
    v-if="activePromotions.length > 0"
    class="section-padding mt-8 mb-8 md:mt-12 md:mb-12 mx-5"
  >
    <div class="container-narrow">
      <div class="bg-gradient-to-r from-brand-500 to-accent-500 rounded-3xl p-6 md:p-10 relative overflow-hidden shadow-xl">

        <div class="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full -translate-y-32 translate-x-32"></div>
        <div class="absolute bottom-0 left-0 w-96 h-96 bg-white/5 rounded-full translate-y-48 -translate-x-48"></div>

        <div class="relative z-10">

          <div class="flex flex-col md:flex-row md:items-center justify-between mb-8 text-white">
            <div>
              <h2 class="text-2xl md:text-3xl font-bold">🔥 Hot Promotions</h2>
              <p class="opacity-90 mt-2">Limited time offers for you!</p>
            </div>
            <div class="mt-4 md:mt-0">
              <span class="bg-white/20 px-3 py-1 rounded-full text-sm">
                {{ currentIndex + 1 }} of {{ activePromotions.length }}
              </span>
            </div>
          </div>

          <div 
            v-if="currentPromotion"
            class="bg-white/15 backdrop-blur-sm rounded-2xl p-6 md:p-8 mb-8 transition-all duration-500"
            :key="currentIndex"
          >
            <div class="text-white">

              <div class="flex flex-col lg:flex-row gap-8 mb-8">

                <div class="lg:w-1/3">
                  <div class="relative h-48 lg:h-64 rounded-xl overflow-hidden shadow-2xl border-2 border-white/20">
                    <img 
                      :src="getPromotionImage(currentPromotion)"
                      :alt="currentPromotion.title"
                      class="w-full h-full object-cover"
                    />

                    <div class="absolute top-4 right-4 bg-red-500 text-white px-3 py-1 rounded-full font-bold shadow-lg">
                      {{ formatDiscountValue(currentPromotion.discountValue, currentPromotion.discountType) }}
                    </div>
                  </div>
                </div>
                

                <div class="lg:w-2/3">
                  <div class="inline-block bg-white/30 px-4 py-1 rounded-full text-sm font-semibold mb-4">
                    {{ formatDiscountType(currentPromotion.discountType) }}
                  </div>
                  
                  <h3 class="text-2xl md:text-3xl font-bold mb-3">
                    {{ currentPromotion.title }}
                  </h3>
                  
                  <p class="text-lg opacity-90 mb-6">
                    {{ currentPromotion.description }}
                  </p>
                  

                  <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
                    <div class="bg-white/10 rounded-xl p-4 text-center">
                      <div class="text-2xl font-bold">
                        {{ getApplicableItems(currentPromotion).length }}
                      </div>
                      <div class="text-sm opacity-90 mt-1">
                        {{ getApplicableItems(currentPromotion).length === 1 ? 'Item' : 'Items' }}
                      </div>
                    </div>
                    
                    <div class="bg-white/10 rounded-xl p-4 text-center">
                      <div class="text-xl font-bold">
                        {{ formatDate(currentPromotion.validFrom) }}
                      </div>
                      <div class="text-sm opacity-90 mt-1">Starts</div>
                    </div>
                    
                    <div class="bg-white/10 rounded-xl p-4 text-center">
                      <div class="text-xl font-bold">
                        {{ formatDate(currentPromotion.validUntil) }}
                      </div>
                      <div class="text-sm opacity-90 mt-1">Ends</div>
                    </div>
                    
                    <div class="bg-white/10 rounded-xl p-4 text-center">
                      <div class="text-xl font-bold">
                        {{ getPromotionCategory(currentPromotion) }}
                      </div>
                      <div class="text-sm opacity-90 mt-1">Category</div>
                    </div>
                  </div>
                </div>
              </div>
              

              <div v-if="getApplicableItems(currentPromotion).length > 0" class="mt-8 pt-8 border-t border-white/20">
                <h4 class="font-bold mb-4 text-xl">🎯 Applicable Items:</h4>
                <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                  <div 
                    v-for="item in getApplicableItems(currentPromotion)" 
                    :key="item.id"
                    class="bg-white/10 rounded-xl p-4 hover:bg-white/15 transition-colors cursor-pointer group"
                    @click="viewMenuItem(item.id)"
                  >
                    <div class="flex items-center gap-3">
                      <div class="w-12 h-12 rounded-lg overflow-hidden flex-shrink-0">
                        <img 
                          :src="item.imageUrl"
                          :alt="item.name"
                          class="w-full h-full object-cover group-hover:scale-110 transition-transform"
                        />
                      </div>
                      <div class="min-w-0">
                        <h5 class="font-semibold truncate">{{ item.name }}</h5>
                        <p class="text-sm opacity-90">${{ item.price.toFixed(2) }}</p>
                        <p class="text-xs opacity-75 mt-1">
                          {{ getDiscountAmount(currentPromotion, item.price) }}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              

              <div class="mt-8 pt-8 border-t border-white/20">
                <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
                  <div>
                    <div class="flex items-center gap-2 mb-2">
                      <div class="w-2 h-2 rounded-full bg-green-400 animate-pulse"></div>
                      <span class="font-medium">
                        {{ getTimeRemaining(currentPromotion.validUntil) }}
                      </span>
                    </div>
                    <p class="text-sm opacity-90">
                      ⚡ {{ currentPromotion.applicableItems?.length ? 'Limited items only!' : 'Applies to all items in category' }}
                    </p>
                  </div>
                  
                  <div class="flex gap-4">
                    <button
                      @click="navigateToPromotionDetails(currentPromotion)"
                      class="px-6 py-3 bg-white/20 text-white rounded-xl hover:bg-white/30 transition-colors border border-white/30 font-medium"
                    >
                      📖 Full Details
                    </button>
                    <button
                      @click="claimOffer(currentPromotion)"
                      class="px-8 py-3 bg-white text-brand-600 font-bold rounded-xl hover:bg-gray-100 transition-colors shadow-lg flex items-center gap-2"
                    >
                      <span>🎁 Claim Offer</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="flex flex-col sm:flex-row items-center justify-between gap-4">

            <div class="flex items-center gap-4">
              <button
                @click="prevPromotion"
                class="p-2 rounded-full bg-white/20 hover:bg-white/30 text-white transition-colors disabled:opacity-50"
                :disabled="activePromotions.length <= 1"
              >
                ←
              </button>
              

              <div class="flex gap-2">
                <button
                  v-for="(_, index) in activePromotions"
                  :key="index"
                  @click="currentIndex = index"
                  class="w-3 h-3 rounded-full transition-all"
                  :class="[
                    index === currentIndex 
                      ? 'bg-white w-8' 
                      : 'bg-white/40 hover:bg-white/60'
                  ]"
                  aria-label="Go to promotion"
                />
              </div>
              
              <button
                @click="nextPromotion"
                class="p-2 rounded-full bg-white/20 hover:bg-white/30 text-white transition-colors disabled:opacity-50"
                :disabled="activePromotions.length <= 1"
              >
                →
              </button>
            </div>
            

            <NuxtLink
              to="/promotion"
              class="text-white hover:text-gray-200 transition-colors underline flex items-center gap-2"
            >
              👀 View All Promotions
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { navigateTo } from 'nuxt/app' 
import type { Promotion, MenuItem } from '../../types/menu'
import { mockMenuItems, mockCategories } from '../../utils/mockData'

interface Props {
  promotions?: Promotion[]
  menuItems?: MenuItem[]
}

const props = withDefaults(defineProps<Props>(), {
  promotions: () => [],
  menuItems: () => mockMenuItems
})

const emit = defineEmits<{
  claimOffer: [promotion: Promotion]
  viewDetails: [promotion: Promotion]
}>()

const currentIndex = ref(0)
let autoPlayInterval: NodeJS.Timeout

const activePromotions = computed(() => {
  return props.promotions.filter(promo => promo.isActive)
})

const currentPromotion = computed(() => {
  return activePromotions.value[currentIndex.value] || null
})

const getApplicableItems = (promotion: Promotion): MenuItem[] => {
  if (!promotion.applicableItems?.length) {
    return props.menuItems
  }
  
  return props.menuItems.filter(item => 
    promotion.applicableItems?.includes(item.id)
  )
}

const getPromotionImage = (promotion: Promotion): string => {
  const items = getApplicableItems(promotion)
  if (items.length > 0) {
    return items[0].imageUrl || '/images/promotion-default.jpg'
  }
  return '/images/promotion-default.jpg'
}

const getPromotionCategory = (promotion: Promotion): string => {
  const items = getApplicableItems(promotion)
  if (items.length > 0) {
    const category = mockCategories.find(cat => cat.id === items[0].categoryId)
    return category?.name || 'Various'
  }
  return 'Various'
}

const getDiscountAmount = (promotion: Promotion, originalPrice: number): string => {
  switch(promotion.discountType) {
    case 'percentage':
      const discount = originalPrice * (promotion.discountValue / 100)
      return `Save $${discount.toFixed(2)}`
    case 'fixed':
      return `Save $${promotion.discountValue.toFixed(2)}`
    case 'bogo':
      return 'Buy 1 Get 1 Free'
    default:
      return 'Special Offer'
  }
}

const formatDiscountType = (type: string) => {
  return type === 'bogo' ? 'BOGO' : type.toUpperCase()
}

const formatDiscountValue = (value: number, type: string) => {
  switch(type) {
    case 'percentage':
      return `${value}% OFF`
    case 'fixed':
      return `$${value.toFixed(2)} OFF`
    case 'bogo':
      return 'Buy 1 Get 1'
    default:
      return `${value} OFF`
  }
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric' 
  })
}

const getTimeRemaining = (validUntil: string) => {
  const now = new Date()
  const endDate = new Date(validUntil)
  const diffTime = endDate.getTime() - now.getTime()
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays <= 0) return 'Ending today!'
  if (diffDays === 1) return 'Ends tomorrow'
  if (diffDays <= 7) return `${diffDays} days left`
  return `Valid until ${formatDate(validUntil)}`
}

const navigateToPromotionDetails = (promotion: Promotion) => {
  navigateTo(`/promotion/${promotion.id}`)
}

const viewMenuItem = (itemId: string) => {
  navigateTo(`/menu?item=${itemId}`)
}

const claimOffer = (promotion: Promotion) => {
  emit('claimOffer', promotion)
}

const nextPromotion = () => {
  if (activePromotions.value.length > 1) {
    currentIndex.value = (currentIndex.value + 1) % activePromotions.value.length
  }
}

const prevPromotion = () => {
  if (activePromotions.value.length > 1) {
    currentIndex.value = currentIndex.value === 0 
      ? activePromotions.value.length - 1 
      : currentIndex.value - 1
  }
}

onMounted(() => {
  if (activePromotions.value.length > 1) {
    autoPlayInterval = setInterval(nextPromotion, 10000) // 10 seconds
  }
})

onUnmounted(() => {
  if (autoPlayInterval) {
    clearInterval(autoPlayInterval)
  }
})
</script>

<style scoped>

.promotion-enter-active,
.promotion-leave-active {
  transition: all 0.5s ease;
}

.promotion-enter-from,
.promotion-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

@keyframes pulse-subtle {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.8; }
}

.animate-pulse-subtle {
  animation: pulse-subtle 2s ease-in-out infinite;
}
</style>