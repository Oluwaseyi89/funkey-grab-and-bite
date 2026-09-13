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

        <div
          v-if="paymentMethod === 'transfer' && paymentStatus === 'pending' && paymentAccountNumber"
          class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-5 mb-6 text-left"
        >
          <p class="font-bold text-blue-900 dark:text-blue-300 mb-3">Complete your payment</p>
          <p class="text-sm text-blue-800 dark:text-blue-400 mb-4">
            Transfer the exact order total to this account. Your order will confirm automatically once the transfer is received.
          </p>
          <div class="space-y-2 text-sm">
            <div class="flex items-center justify-between gap-3">
              <span class="text-blue-700 dark:text-blue-400">Account Number</span>
              <div class="flex items-center gap-2">
                <span class="font-mono font-bold text-blue-900 dark:text-blue-200">{{ paymentAccountNumber }}</span>
                <button
                  type="button"
                  @click="copyAccountNumber"
                  class="text-blue-600 dark:text-blue-400 hover:text-blue-800"
                  title="Copy account number"
                >
                  <component :is="accountNumberCopied ? Check : Copy" class="w-4 h-4" />
                </button>
              </div>
            </div>
            <div v-if="paymentAccountName" class="flex items-center justify-between">
              <span class="text-blue-700 dark:text-blue-400">Account Name</span>
              <span class="font-semibold text-blue-900 dark:text-blue-200">{{ paymentAccountName }}</span>
            </div>
            <div v-if="paymentBankName" class="flex items-center justify-between">
              <span class="text-blue-700 dark:text-blue-400">Bank</span>
              <span class="font-semibold text-blue-900 dark:text-blue-200">{{ paymentBankName }}</span>
            </div>
            <div v-if="paymentExpiresAt" class="flex items-center justify-between">
              <span class="text-blue-700 dark:text-blue-400">Expires</span>
              <span class="font-semibold text-blue-900 dark:text-blue-200">{{ formattedExpiry }}</span>
            </div>
          </div>
        </div>

        <div
          v-else-if="paymentMethod === 'cash'"
          class="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-xl p-4 mb-6 text-sm text-amber-800 dark:text-amber-400"
        >
          Pay in cash when your order arrives or is picked up.
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
  import { CheckCircle, Clock, Loader2, Copy, Check } from 'lucide-vue-next'
  import { useApi } from '../../utils/api'
  import type { PaymentMethod, PaymentStatus } from '../../types/order'

  const props = defineProps<{
    orderNumber: string
    estimatedTime: string
    customerPhone: string
    paymentMethod?: PaymentMethod
    paymentStatus?: PaymentStatus
    paymentAccountNumber?: string
    paymentAccountName?: string
    paymentBankName?: string
    paymentExpiresAt?: string
  }>()

  defineEmits<{
    close: []
  }>()

  const api = useApi()
  const isTracking = ref(false)
  const liveStatus = ref('')
  const accountNumberCopied = ref(false)

  const formattedExpiry = computed(() => {
    if (!props.paymentExpiresAt) return ''
    return new Date(props.paymentExpiresAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  })

  const copyAccountNumber = async () => {
    if (!props.paymentAccountNumber) return
    try {
      await navigator.clipboard.writeText(props.paymentAccountNumber)
      accountNumberCopied.value = true
      setTimeout(() => { accountNumberCopied.value = false }, 2000)
    } catch {
      // clipboard access can be denied by the browser — the number is still visible to copy manually
    }
  }

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