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

          <div v-if="liveStatus" class="mt-4 flex items-center justify-center gap-2">
            <span class="text-sm text-gray-500 dark:text-gray-400">Live status:</span>
            <span class="px-3 py-1 rounded-full text-xs font-semibold capitalize" :class="statusClass">
              {{ liveStatus }}
            </span>
          </div>
        </div>
        
        <div class="space-y-3">
          <button
            @click="trackOrder"
            :disabled="isTracking"
            class="w-full btn-primary flex items-center justify-center gap-2"
          >
            <Loader2 v-if="isTracking" class="w-4 h-4 animate-spin" />
            <span>{{ liveStatus ? 'Refresh Status' : 'Track Order' }}</span>
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
  import { ref, computed } from 'vue'
  import { CheckCircle, Clock, Loader2 } from 'lucide-vue-next'
  import { useApi } from '../../utils/api'

  const props = defineProps<{
    orderNumber: string
    estimatedTime: string
    customerPhone: string
  }>()

  defineEmits<{
    close: []
  }>()

  const api = useApi()
  const isTracking = ref(false)
  const liveStatus = ref('')

  const statusClass = computed(() => {
    switch (liveStatus.value) {
      case 'pending':   return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
      case 'confirmed': return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
      case 'preparing': return 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400'
      case 'ready':     return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
      case 'completed': return 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
      case 'cancelled': return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
      default:          return 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-400'
    }
  })

  const trackOrder = async () => {
    if (!props.customerPhone || !props.orderNumber) return
    isTracking.value = true
    try {
      const order = await api.getOrder(props.customerPhone, props.orderNumber)
      if (order) {
        liveStatus.value = order.status
      }
    } catch {
      // silent — status badge simply won't appear on error
    } finally {
      isTracking.value = false
    }
  }
  </script>