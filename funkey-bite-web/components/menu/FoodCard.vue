<template>
    <div
      ref="cardRef"
      class="group bg-white dark:bg-slate-800 rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-500 overflow-hidden card-hover border border-gray-100 dark:border-slate-700"
    >
      <!-- Image Section -->
      <div class="relative h-48 overflow-hidden">
        <div class="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent z-10"></div>
        <img
          :src="imageUrl"
          :alt="item.name"
          class="w-full h-full object-cover transform group-hover:scale-110 transition-transform duration-700"
          loading="lazy"
        />
        
        <!-- Badges -->
        <div class="absolute top-3 left-3 z-20 flex flex-wrap gap-2">
          <span 
            v-if="item.isPreOrder"
            class="px-3 py-1 bg-purple-500 text-white text-xs font-bold rounded-full"
          >
            Pre-order
          </span>
          <span 
            v-if="item.tags?.includes('best seller')"
            class="px-3 py-1 bg-yellow-500 text-white text-xs font-bold rounded-full"
          >
            Best Seller
          </span>
          <span 
            v-if="item.tags?.includes('spicy')"
            class="px-3 py-1 bg-red-500 text-white text-xs font-bold rounded-full"
          >
            Spicy
          </span>
        </div>
        
        <!-- Quick Add Button -->
        <button
          @click="addToCart"
          class="absolute top-3 right-3 z-20 w-10 h-10 bg-white/90 dark:bg-slate-800/90 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300 hover:bg-brand-500 hover:text-white"
          aria-label="Add to cart"
        >
          <Plus class="w-5 h-5" />
        </button>
      </div>
      
      <!-- Content Section -->
      <div class="p-6">
        <div class="flex justify-between items-start mb-2">
          <h3 class="text-xl font-bold text-gray-900 dark:text-white group-hover:text-brand-500 dark:group-hover:text-brand-400 transition-colors">
            {{ item.name }}
          </h3>
          <div class="text-2xl font-bold text-brand-500">
            ${{ item.price.toFixed(2) }}
          </div>
        </div>
        
        <p class="text-gray-600 dark:text-gray-300 mb-4 line-clamp-2">
          {{ item.description }}
        </p>
        
        <!-- Nutritional Info -->
        <div v-if="item.nutritionalInfo" class="mb-4">
          <div class="flex items-center text-sm text-gray-500 dark:text-gray-400 mb-2">
            <Nutrient class="w-4 h-4 mr-2" />
            Nutritional Info
          </div>
          <div class="grid grid-cols-4 gap-2 text-xs">
            <div class="text-center p-2 bg-gray-50 dark:bg-slate-700 rounded-lg">
              <div class="font-bold">{{ item.nutritionalInfo.calories }}</div>
              <div class="text-gray-500 dark:text-gray-400">Cal</div>
            </div>
            <div class="text-center p-2 bg-gray-50 dark:bg-slate-700 rounded-lg">
              <div class="font-bold">{{ item.nutritionalInfo.protein }}g</div>
              <div class="text-gray-500 dark:text-gray-400">Protein</div>
            </div>
            <div class="text-center p-2 bg-gray-50 dark:bg-slate-700 rounded-lg">
              <div class="font-bold">{{ item.nutritionalInfo.carbs }}g</div>
              <div class="text-gray-500 dark:text-gray-400">Carbs</div>
            </div>
            <div class="text-center p-2 bg-gray-50 dark:bg-slate-700 rounded-lg">
              <div class="font-bold">{{ item.nutritionalInfo.fat }}g</div>
              <div class="text-gray-500 dark:text-gray-400">Fat</div>
            </div>
          </div>
        </div>
        
        <!-- Footer -->
        <div class="flex items-center justify-between pt-4 border-t border-gray-100 dark:border-slate-700">
          <div class="flex items-center text-sm text-gray-500 dark:text-gray-400">
            <Clock class="w-4 h-4 mr-1" />
            {{ item.preparationTime }} min
          </div>
          
          <div class="flex space-x-2">
            <button
              @click="showDetails"
              class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 hover:text-brand-500 dark:hover:text-brand-400 transition-colors"
            >
              Details
            </button>
            <button
              @click="addToCart"
              class="btn-primary px-4 py-2 text-sm"
            >
              Add to Cart
            </button>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { Plus, Clock } from 'lucide-vue-next'
  import type { MenuItem } from '../../types/menu'
  import { onMounted, computed, ref } from 'vue'
  import { useNuxtApp } from 'nuxt/app'
  import { useCartStore } from '../../stores/cart'
  import { useRuntimeConfig } from 'nuxt/app'
  import { gsap } from 'gsap'
  
  interface Props {
    item: MenuItem
  }
  
  const props = defineProps<Props>()
  const cardRef = ref<HTMLElement>()
  const cart = useCartStore()
  const { $toast } = useNuxtApp()
  
  const imageUrl = computed(() => {
    const config = useRuntimeConfig()
    return props.item.imageUrl.startsWith('http') 
      ? props.item.imageUrl 
      : `${config.public.s3BucketUrl}/${props.item.imageUrl}`
  })
  
  const addToCart = () => {
    cart.addItem(props.item)
    // Animate button feedback
    if (import.meta.client && cardRef.value) {
      const button = cardRef.value.querySelector('.btn-primary')
      gsap.to(button, {
        scale: 1.2,
        duration: 0.2,
        yoyo: true,
        repeat: 1,
      })
    }
  }
  
  const showDetails = () => {
    // Emit event to parent to show modal
    const modalEvent = new CustomEvent('show-item-details', { detail: props.item })
    window.dispatchEvent(modalEvent)
  }
  
  onMounted(() => {
    if (import.meta.client && cardRef.value) {
      // Entrance animation
      gsap.from(cardRef.value, {
        scrollTrigger: {
          trigger: cardRef.value,
          start: 'top 80%',
          toggleActions: 'play none none reverse',
        },
        opacity: 0,
        y: 50,
        duration: 0.6,
        ease: 'power2.out',
      })
    }
  })
  </script>