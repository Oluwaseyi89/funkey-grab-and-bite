<template>
    <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4" role="presentation">
      <div
        class="bg-white dark:bg-slate-800 rounded-2xl max-w-md w-full p-8 text-center"
        role="dialog"
        aria-modal="true"
        aria-labelledby="order-confirmation-title"
      >
        <div class="w-20 h-20 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-6">
          <CheckCircle class="w-10 h-10 text-green-600 dark:text-green-400" />
        </div>
        
        <h2 id="order-confirmation-title" class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Order Confirmed!</h2>
        <p class="text-gray-600 dark:text-gray-400 mb-6">
          Thank you for your order. We're preparing it now.
        </p>
        
        <div class="bg-gray-50 dark:bg-slate-700 rounded-xl p-6 mb-6">
          <div class="text-sm text-gray-500 dark:text-gray-400 mb-2">Order Number</div>
          <div class="text-2xl font-bold text-brand-500 mb-4">{{ orderNumber }}</div>
          
          <div class="flex items-center justify-center space-x-2 text-gray-600 dark:text-gray-300">
            <Clock class="w-4 h-4" />
            <span>Estimated ready: {{ estimatedTime }}</span>
          </div>
        </div>
        
        <div class="space-y-3">
          <button @click="$emit('close')" class="w-full btn-primary">
            Track Order
          </button>
          <button @click="$emit('close')" class="w-full btn-secondary">
            Back to Menu
          </button>
        </div>
        
        <p class="mt-6 text-sm text-gray-500 dark:text-gray-400">
          We'll send updates to your email and phone.
        </p>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { CheckCircle, Clock } from 'lucide-vue-next'
  
  defineProps<{
    orderNumber: string
    estimatedTime: string
  }>()
  
  defineEmits<{
    close: []
  }>()
  </script>