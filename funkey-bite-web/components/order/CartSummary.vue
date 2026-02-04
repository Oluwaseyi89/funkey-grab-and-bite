<template>
    <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 sticky top-24">
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Order Summary</h3>
      
      <!-- Items List -->
      <div class="space-y-3 max-h-64 overflow-y-auto mb-4">
        <div v-for="item in items" :key="item.menuItem.id" class="flex justify-between text-sm">
          <span class="text-gray-600 dark:text-gray-400">
            {{ item.menuItem.name }} × {{ item.quantity }}
          </span>
          <span class="font-medium">${{ (item.menuItem.price * item.quantity).toFixed(2) }}</span>
        </div>
      </div>
  
      <!-- Totals -->
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
  
      <!-- Promo Code -->
      <div class="mt-4">
        <div class="flex space-x-2">
          <input
            v-model="promoCode"
            type="text"
            placeholder="Promo code"
            class="flex-1 px-4 py-2 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
          />
          <button
            @click="applyPromo"
            class="px-4 py-2 bg-gray-100 dark:bg-slate-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
          >
            Apply
          </button>
        </div>
        <div v-if="discount > 0" class="mt-2 text-green-600 dark:text-green-400 text-sm">
          -${{ discount.toFixed(2) }} discount applied
        </div>
      </div>
  
      <!-- Checkout Button -->
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
          {{ items.length === 0 ? 'Cart is Empty' : `Proceed to Checkout - $${total.toFixed(2)}` }}
        </template>
      </button>
  
      <!-- Secure Payment -->
      <div class="mt-4 text-center text-sm text-gray-500 dark:text-gray-400">
        🔒 Secure payment • No card details stored
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import type { CartItem } from '../../types/menu'
  import { Loader2 } from 'lucide-vue-next'
  import { ref, computed } from 'vue'
  
//   defineProps<{
//     items: CartItem[]
//     isLoading?: boolean
//     deliveryFee?: number
//   }>()

  const props = withDefaults(defineProps<{
  items: CartItem[]
  isLoading?: boolean
  deliveryFee?: number
}>(), {
  deliveryFee: 0 // ✅ Default value
})

// OR

// const { items, isLoading, deliveryFee = 0 } = defineProps<{
//   items: CartItem[]
//   isLoading?: boolean
//   deliveryFee?: number
// }>()
  
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
      // Show success toast
    } else {
      discount.value = 0
      // Show error toast
    }
  }
  </script>