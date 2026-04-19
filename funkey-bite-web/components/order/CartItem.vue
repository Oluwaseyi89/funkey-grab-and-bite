<template>
    <div class="flex flex-row w-full items-start gap-3 p-3 sm:p-4 bg-white dark:bg-slate-800 rounded-xl border border-gray-200 dark:border-slate-700">
      <!-- Item Image -->
      <div class="w-20 h-20 sm:w-24 sm:h-24 flex-shrink-0">
        <img 
          :src="item.menuItem.imageUrl" 
          :alt="item.menuItem.name"
          class="w-full h-full object-cover rounded-lg"
        />
      </div>
  
      <!-- Item Details -->
      <div class="flex-1 min-w-0">
        <div class="flex justify-between items-start gap-2">
          <div class="min-w-0 flex-1">
            <h4 class="font-bold text-gray-900 dark:text-white leading-snug">{{ item.menuItem.name }}</h4>
            <p class="text-xs sm:text-sm text-gray-600 dark:text-gray-400 line-clamp-2 mt-0.5">{{ item.menuItem.description }}</p>
            <div v-if="item.specialInstructions" class="mt-1 text-xs sm:text-sm text-brand-500 truncate">
              📝 {{ item.specialInstructions }}
            </div>
          </div>
          <div class="text-right flex-shrink-0">
            <div class="font-bold text-base sm:text-lg text-brand-500">₦{{ (item.menuItem.price * item.quantity).toFixed(2) }}</div>
            <div class="text-xs text-gray-500">${{ item.menuItem.price.toFixed(2) }} each</div>
          </div>
        </div>
  
        <!-- Quantity Controls -->
        <div class="flex items-center justify-between mt-3">
          <div class="flex items-center space-x-2">
            <button
              @click="updateQuantity(-1)"
              class="w-8 h-8 flex items-center justify-center bg-gray-100 dark:bg-slate-700 rounded-full hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
              :disabled="item.quantity <= 1"
              :aria-label="`Decrease quantity of ${item.menuItem.name}`"
            >
              <Minus class="w-4 h-4" />
            </button>
            <span class="w-10 text-center font-bold" aria-live="polite" :aria-label="`Quantity: ${item.quantity}`">{{ item.quantity }}</span>
            <button
              @click="updateQuantity(1)"
              class="w-8 h-8 flex items-center justify-center bg-gray-100 dark:bg-slate-700 rounded-full hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
              :aria-label="`Increase quantity of ${item.menuItem.name}`"
            >
              <Plus class="w-4 h-4" />
            </button>
          </div>
          <button
            @click="$emit('remove')"
            class="text-red-500 hover:text-red-600 text-sm font-medium flex items-center"
            :aria-label="`Remove ${item.menuItem.name} from cart`"
          >
            <Trash2 class="w-4 h-4" />
            <span class="ml-1">Remove</span>
          </button>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import type { CartItem } from '../../types/menu'
  import { Minus, Plus, Trash2 } from 'lucide-vue-next'
  
  const { item } = defineProps<{
    item: CartItem
  }>()
  
  const emit = defineEmits<{
    updateQuantity: [quantity: number]
    remove: []
  }>()
  
  const updateQuantity = (change: number) => {
    emit('updateQuantity', change)
  }
  </script>