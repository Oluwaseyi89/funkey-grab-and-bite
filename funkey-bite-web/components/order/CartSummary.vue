<template>
    <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 lg:sticky lg:top-24">
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Order Summary</h3>
      
      
      <div class="space-y-2 max-h-48 sm:max-h-64 overflow-y-auto mb-4 pr-1">
        <div v-for="item in items" :key="item.menuItem.id" class="flex justify-between gap-3 text-sm py-1">
          <span class="text-gray-600 dark:text-gray-400 min-w-0 leading-snug">{{ item.menuItem.name }} <span class="text-gray-400 dark:text-gray-500">×{{ item.quantity }}</span></span>
          <span class="font-semibold flex-shrink-0 text-gray-900 dark:text-white">${{ (item.menuItem.price * item.quantity).toFixed(2) }}</span>
        </div>
      </div>
  
      
      <div class="space-y-2 border-t border-gray-200 dark:border-slate-700 pt-4">
        <div class="flex justify-between text-gray-600 dark:text-gray-400">
          <span>Subtotal</span>
          <span>${{ subtotal.toFixed(2) }}</span>
        </div>
        <div class="flex justify-between text-gray-600 dark:text-gray-400">
          <span>Tax (8%)</span>
          <span>${{ tax.toFixed(2) }}</span>
        </div>
        <div v-if="deliveryFee > 0" class="flex justify-between text-gray-600 dark:text-gray-400">
          <span>Delivery Fee</span>
          <span>${{ deliveryFee.toFixed(2) }}</span>
        </div>
        <div class="flex justify-between text-lg font-bold text-gray-900 dark:text-white pt-2 border-t">
          <span>Total</span>
          <span class="text-brand-500">${{ total.toFixed(2) }}</span>
        </div>
      </div>
  
      
      
      <div class="mt-4">
        <div class="flex flex-col sm:flex-row gap-2">
          <input
            v-model="promoCode"
            type="text"
            placeholder="Promo code"
            class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:border-transparent bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
          />
          <button
            @click="applyPromo"
            class="w-full sm:w-auto px-4 py-2 bg-gray-100 dark:bg-slate-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors font-medium whitespace-nowrap"
          >
            Apply
          </button>
        </div>
        <div v-if="discount > 0" class="mt-2 text-green-600 dark:text-green-400 text-sm">
            -${{ discount.toFixed(2) }} discount applied
        </div>
    </div>
  
      
      <button
        @click="$emit('checkout')"
        :disabled="items.length === 0 || isLoading"
        class="w-full btn-primary mt-6 py-4 text-lg font-bold disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <template v-if="isLoading">
          <Loader2 class="w-5 h-5 animate-spin inline mr-2" />
          Processing...
        </template>
        <template v-else>
          {{ items.length === 0 ? 'Cart is Empty' : `Checkout · $${total.toFixed(2)}` }}
        </template>
      </button>
  
      
      <div class="mt-4 text-center text-sm text-gray-500 dark:text-gray-400">
        🔒 Secure payment • No card details stored
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import type { CartItem } from '../../types/menu'
  import { Loader2 } from 'lucide-vue-next'
  import { ref, computed } from 'vue'
  

  const props = withDefaults(defineProps<{
  items: CartItem[]
  isLoading?: boolean
  deliveryFee?: number
}>(), {
  deliveryFee: 0 // ✅ Default value
})

  
  const emit = defineEmits<{
    checkout: []
  }>()
  
  const promoCode = ref('')
  const discount = ref(0)
  
  const subtotal = computed(() => {
    return props.items.reduce((sum, item) => sum + (item.menuItem.price * item.quantity), 0)
  })
  
  const tax = computed(() => subtotal.value * 0.08)
  
  const total = computed(() => {
    return subtotal.value + tax.value + (props.deliveryFee || 0) - discount.value
  })
  
  const applyPromo = () => {
    if (promoCode.value.toUpperCase() === 'WELCOME10') {
      discount.value = subtotal.value * 0.1
    } else {
      discount.value = 0
    }
  }
  </script>