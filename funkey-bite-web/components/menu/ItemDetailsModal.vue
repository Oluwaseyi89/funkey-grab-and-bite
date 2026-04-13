<template>
  
  <div class="modal-root" style="z-index: 2147483647;">
    
    
    <div 
      class="fixed inset-0"
      :class="[
        'md:flex md:items-center md:justify-center md:p-4',
        isMobile ? 'bg-white dark:bg-slate-900' : 'bg-black/50 backdrop-blur-sm'
      ]"
    >
      
      <div 
        ref="modalRef"
        class="w-full h-full md:h-auto md:max-w-2xl md:max-h-[90vh] md:rounded-2xl md:overflow-hidden"
        :class="[
          isMobile 
            ? 'overflow-y-auto bg-white dark:bg-slate-900' 
            : 'bg-white dark:bg-slate-800 shadow-2xl'
        ]"
      >
        
        <div v-if="isMobile" class="sticky top-0 z-10 bg-white dark:bg-slate-900 border-b border-gray-200 dark:border-slate-700 p-4 flex items-center justify-between">
          <button @click="handleClose" class="p-2">
            <X class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </button>
          <h2 class="text-lg font-bold text-gray-900 dark:text-white truncate px-2">{{ item.name }}</h2>
          <div class="w-10"></div>
        </div>

        
        <button 
          v-if="!isMobile"
          @click="handleClose" 
          class="absolute top-4 right-4 z-10 p-2 rounded-full bg-white/90 dark:bg-slate-800/90 text-gray-800 dark:text-gray-200 hover:bg-white dark:hover:bg-slate-700 transition-colors shadow-lg"
        >
          <X class="w-5 h-5" />
        </button>

        
        <div class="grid md:grid-cols-2 h-[75vh] md:h-auto px-3 md:px-8 py-3 md:py-8">
          
          <div class="relative h-64 md:h-auto md:min-h-[300px]">
            <img 
              :src="item.imageUrl" 
              :alt="item.name" 
              class="w-full h-full object-cover"
            />
          </div>

          
          <div class="p-6 md:p-8 overflow-y-auto">
            <div v-if="!isMobile" class="flex justify-between items-start mb-4">
              <div>
                <h2 class="text-2xl font-bold text-gray-900 dark:text-white">{{ item.name }}</h2>
                <div class="flex items-center space-x-2 mt-2">
                  <Clock class="w-4 h-4 text-gray-500" />
                  <span class="text-gray-600 dark:text-gray-400">{{ item.preparationTime }} mins</span>
                </div>
              </div>
              <span class="text-3xl font-bold text-brand-500">${{ item.price.toFixed(2) }}</span>
            </div>

            <div v-if="isMobile" class="flex justify-between items-center mb-4">
              <div class="flex items-center space-x-2">
                <Clock class="w-4 h-4 text-gray-500" />
                <span class="text-gray-600 dark:text-gray-400">{{ item.preparationTime }} mins</span>
              </div>
              <span class="text-2xl font-bold text-brand-500">${{ item.price.toFixed(2) }}</span>
            </div>

            <p class="text-gray-600 dark:text-gray-300 mb-6">{{ item.description }}</p>

            
            <div v-if="item.tags?.length" class="flex flex-wrap gap-2 mb-6">
              <span
                v-for="tag in item.tags"
                :key="tag"
                class="px-3 py-1 bg-gray-100 dark:bg-slate-700 text-gray-700 dark:text-gray-300 rounded-full text-sm"
              >
                {{ tag }}
              </span>
            </div>

            
            <div v-if="item.nutritionalInfo" class="mb-6">
              <h4 class="font-bold text-gray-900 dark:text-white mb-3">Nutritional Information</h4>
              <div class="grid grid-cols-4 gap-2">
                <div class="text-center p-3 bg-gray-50 dark:bg-slate-700 rounded-lg">
                  <div class="font-bold text-lg">{{ item.nutritionalInfo.calories }}</div>
                  <div class="text-gray-500 text-xs">Calories</div>
                </div>
                <div class="text-center p-3 bg-gray-50 dark:bg-slate-700 rounded-lg">
                  <div class="font-bold text-lg">{{ item.nutritionalInfo.protein }}g</div>
                  <div class="text-gray-500 text-xs">Protein</div>
                </div>
                <div class="text-center p-3 bg-gray-50 dark:bg-slate-700 rounded-lg">
                  <div class="font-bold text-lg">{{ item.nutritionalInfo.carbs }}g</div>
                  <div class="text-gray-500 text-xs">Carbs</div>
                </div>
                <div class="text-center p-3 bg-gray-50 dark:bg-slate-700 rounded-lg">
                  <div class="font-bold text-lg">{{ item.nutritionalInfo.fat }}g</div>
                  <div class="text-gray-500 text-xs">Fat</div>
                </div>
              </div>
            </div>

            
            <div v-if="item.isPreOrder" class="mb-6 p-4 bg-amber-50 dark:bg-amber-900/20 rounded-lg">
              <AlertTriangle class="w-5 h-5 text-amber-500 inline mr-2" />
              <span class="text-amber-700 dark:text-amber-300 text-sm">
                Pre-order: Order at least {{ item.preparationTime }} minutes in advance.
              </span>
            </div>

            
            <div class="flex space-x-4 mt-8 md:mt-0">
              <button
                @click="handleAddToCart"
                :disabled="!item.isAvailable"
                class="flex-1 btn-primary text-lg py-4 dark:text-white"
                :class="{ 'opacity-50 cursor-not-allowed': !item.isAvailable }"
              >
                <ShoppingCart class="w-5 h-5 inline mr-2" />
                Add to Cart
              </button>
              <button 
                v-if="!isMobile"
                @click="handleClose" 
                class="px-6 py-4 border border-gray-300 dark:border-slate-600 rounded-xl hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors dark:text-white"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { MenuItem } from '../../types/menu'
import { X, Clock, ShoppingCart, AlertTriangle } from 'lucide-vue-next'

const { item } = defineProps<{
  item: MenuItem
}>()

const emit = defineEmits<{
  close: []
  addToCart: [item: MenuItem]
}>()

const isMobile = ref(false)
const modalRef = ref<HTMLElement>()

const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
}

const handleClose = () => {
  emit('close')
}

const handleAddToCart = () => {
  emit('addToCart', item)
}

const handleEscape = (e: KeyboardEvent) => {
  if (e.key === 'Escape') handleClose()
}

const handleClickOutside = (e: MouseEvent) => {
  if (modalRef.value && !modalRef.value.contains(e.target as Node)) {
    handleClose()
  }
}

let savedScrollY = 0

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  document.addEventListener('keydown', handleEscape)
  document.addEventListener('mousedown', handleClickOutside)

  savedScrollY = window.scrollY

  document.body.style.overflow = 'hidden'
  document.body.style.position = 'fixed'
  document.body.style.top = `-${savedScrollY}px`
  document.body.style.width = '100%'

  modalRef.value?.focus()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
  document.removeEventListener('keydown', handleEscape)
  document.removeEventListener('mousedown', handleClickOutside)

  document.body.style.overflow = ''
  document.body.style.position = ''
  document.body.style.top = ''
  document.body.style.width = ''

  window.scrollTo(0, savedScrollY)
})

</script>

<style scoped>
.modal-root {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 2147483647; 
}

.fixed {
  position: fixed !important;
}

@media (max-width: 767px) {
  .modal-root {
    -webkit-overflow-scrolling: touch;
    overflow-y: auto;
  }
}
</style>