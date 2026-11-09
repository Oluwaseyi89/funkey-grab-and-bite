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
        
        <h2 id="order-confirmation-title" class="text-2xl font-bold text-gray-900 dark:text-white mb-2">
          {{ awaitingOnline ? 'Order received' : 'Order Confirmed!' }}
        </h2>
        <p class="text-gray-600 dark:text-gray-400 mb-6">
          {{ awaitingOnline
            ? "Thank you. We'll start preparing it as soon as your payment comes through."
            : "Thank you for your order. We're preparing it now." }}
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
          v-if="awaitingTransfer && accountNumber && !paymentWindowClosed"
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
                <span class="font-mono font-bold text-blue-900 dark:text-blue-200">{{ accountNumber }}</span>
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
            <div v-if="accountName" class="flex items-center justify-between">
              <span class="text-blue-700 dark:text-blue-400">Account Name</span>
              <span class="font-semibold text-blue-900 dark:text-blue-200">{{ accountName }}</span>
            </div>
            <div v-if="bankName" class="flex items-center justify-between">
              <span class="text-blue-700 dark:text-blue-400">Bank</span>
              <span class="font-semibold text-blue-900 dark:text-blue-200">{{ bankName }}</span>
            </div>
            <div v-if="expiresAt" class="flex items-center justify-between">
              <span class="text-blue-700 dark:text-blue-400">Expires</span>
              <span class="font-semibold text-blue-900 dark:text-blue-200">{{ formattedExpiry }}</span>
            </div>
          </div>
        </div>

        <div
          v-else-if="awaitingTransfer"
          class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-5 mb-6 text-left"
        >
          <p class="font-bold text-blue-900 dark:text-blue-300 mb-2">
            {{ accountNumber ? 'Your payment window has closed' : 'Payment details are not ready yet' }}
          </p>
          <p class="text-sm text-blue-800 dark:text-blue-400 mb-4">
            {{ accountNumber
              ? "If you haven't sent the transfer yet, get new details. Don't use the old account."
              : "We couldn't get an account for your transfer. Your order is saved; get the details to pay." }}
          </p>
          <p v-if="renewError" class="text-sm text-red-600 dark:text-red-400 mb-3" role="alert">{{ renewError }}</p>
          <button
            type="button"
            class="btn-primary w-full flex items-center justify-center gap-2"
            :disabled="isRenewing"
            @click="getPaymentDetails"
          >
            <Loader2 v-if="isRenewing" class="w-4 h-4 animate-spin" />
            <span>Get payment details</span>
          </button>
        </div>

        <div
          v-else-if="awaitingCard"
          class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-5 mb-6 text-left"
        >
          <p class="font-bold text-blue-900 dark:text-blue-300 mb-2">Complete your card payment</p>
          <p class="text-sm text-blue-800 dark:text-blue-400 mb-4">
            Your order is saved. Pay securely in Paystack's window; your card details never reach us.
          </p>
          <p v-if="renewError" class="text-sm text-red-600 dark:text-red-400 mb-3" role="alert">{{ renewError }}</p>
          <button
            v-if="cardPaymentsEnabled"
            type="button"
            class="btn-primary w-full flex items-center justify-center gap-2"
            :disabled="isRenewing"
            @click="payByCard"
          >
            <Loader2 v-if="isRenewing" class="w-4 h-4 animate-spin" />
            <span>Pay with card</span>
          </button>
        </div>

        <div
          v-else-if="isOnline && currentPaymentStatus === 'paid'"
          class="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl p-4 mb-6 text-sm text-green-800 dark:text-green-300 text-left"
          role="status"
        >
          Payment received. Thank you!
        </div>

        <div
          v-else-if="paymentMethod === 'cash'"
          class="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-xl p-4 mb-6 text-sm text-amber-800 dark:text-amber-400 text-left"
        >
          <p>
            Pay <strong v-if="totalAmount != null">₦{{ totalAmount.toFixed(2) }}</strong>
            {{ orderType === 'delivery' ? 'when your order is delivered' : 'when you collect your order' }}
            — cash, POS or transfer.
          </p>
          <p v-if="invoiceNumber" class="mt-2">
            Invoice <strong>{{ invoiceNumber }}</strong> has been emailed to you.
            <NuxtLink to="/order/invoice" class="underline font-medium" @click="$emit('close')">View or print invoice</NuxtLink>
          </p>
        </div>

        <div class="space-y-3">
          <button
            @click="trackOrder"
            :disabled="isTracking"
            class="w-full btn-primary flex items-center justify-center gap-2"
          >
            <Loader2 v-if="isTracking" class="w-4 h-4 animate-spin" />
            <span>{{ awaitingOnline ? 'Refresh payment status' : liveStatus ? 'Refresh Status' : 'Track Order' }}</span>
          </button>
          <p v-if="refreshMessage" class="text-sm text-gray-600 dark:text-gray-400" role="status">{{ refreshMessage }}</p>
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
  import { ref, computed, onMounted, onUnmounted } from 'vue'
  import { CheckCircle, Clock, Loader2, Copy, Check } from 'lucide-vue-next'
  import { useApi } from '../../utils/api'
  import { usePaystack } from '../../composables/usePaystack'
  import type { Order, OrderType, PaymentMethod, PaymentStatus } from '../../types/order'

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
    invoiceNumber?: string
    orderType?: OrderType
    totalAmount?: number
  }>()

  defineEmits<{
    close: []
  }>()

  const api = useApi()
  const isTracking = ref(false)
  const liveStatus = ref('')
  const accountNumberCopied = ref(false)

  // New payment details fetched on this screen replace the ones from checkout.
  const renewed = ref<Order | null>(null)
  const isRenewing = ref(false)
  const renewError = ref('')
  const accountNumber = computed(() => renewed.value?.paymentAccountNumber ?? props.paymentAccountNumber)
  const accountName = computed(() => renewed.value?.paymentAccountName ?? props.paymentAccountName)
  const bankName = computed(() => renewed.value?.paymentBankName ?? props.paymentBankName)
  const expiresAt = computed(() => renewed.value?.paymentExpiresAt ?? props.paymentExpiresAt)
  const currentPaymentStatus = computed(() => renewed.value?.paymentStatus ?? props.paymentStatus)
  const isOnline = computed(() => props.paymentMethod === 'transfer' || props.paymentMethod === 'card')
  const awaitingOnline = computed(() => isOnline.value && currentPaymentStatus.value === 'pending')
  const awaitingTransfer = computed(() => props.paymentMethod === 'transfer' && awaitingOnline.value)
  const awaitingCard = computed(() => props.paymentMethod === 'card' && awaitingOnline.value)
  const { cardPaymentsEnabled, payWithCard } = usePaystack()

  // Opens Paystack's popup with a fresh card payment, then has the server
  // confirm it.
  const payByCard = async () => {
    if (!props.customerPhone || !props.orderNumber) return
    isRenewing.value = true
    renewError.value = ''
    try {
      const started = await api.getPaymentDetails(props.orderNumber, props.customerPhone)
      if (!started.paymentAccessCode) throw new Error("We couldn't start the card payment. Please try again.")
      if ((await payWithCard(started.paymentAccessCode)) === 'paid') {
        await refreshStatus()
      } else {
        renewError.value = 'Payment cancelled. You can pay whenever you’re ready.'
      }
    } catch (err) {
      renewError.value = err instanceof Error ? err.message : 'Something went wrong. Please try again.'
    } finally {
      isRenewing.value = false
    }
  }

  // Re-checked every half minute so the screen notices the window closing.
  const now = ref(Date.now())
  let clock: ReturnType<typeof setInterval> | undefined
  onMounted(() => { clock = setInterval(() => { now.value = Date.now() }, 30_000) })
  onUnmounted(() => clearInterval(clock))
  const paymentWindowClosed = computed(() =>
    !!expiresAt.value && new Date(expiresAt.value).getTime() <= now.value)

  const getPaymentDetails = async () => {
    if (!props.customerPhone || !props.orderNumber) return
    isRenewing.value = true
    renewError.value = ''
    try {
      renewed.value = await api.getPaymentDetails(props.orderNumber, props.customerPhone)
      now.value = Date.now()
    } catch (err) {
      renewError.value = err instanceof Error ? err.message : 'Something went wrong. Please try again.'
    } finally {
      isRenewing.value = false
    }
  }

  const formattedExpiry = computed(() => {
    if (!expiresAt.value) return ''
    return new Date(expiresAt.value).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  })

  const copyAccountNumber = async () => {
    if (!accountNumber.value) return
    try {
      await navigator.clipboard.writeText(accountNumber.value)
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

  // While an online payment is outstanding, the server checks it with
  // Paystack, so the customer sees a payment as soon as it lands.
  const refreshMessage = ref('')
  const refreshStatus = async () => {
    const result = await api.refreshPaymentStatus(props.orderNumber, props.customerPhone)
    renewed.value = result.order
    liveStatus.value = result.order.status
    refreshMessage.value = result.message
  }

  const trackOrder = async () => {
    if (!props.customerPhone || !props.orderNumber) return
    isTracking.value = true
    refreshMessage.value = ''
    try {
      if (awaitingOnline.value) {
        await refreshStatus()
      } else {
        const order = await api.getOrder(props.customerPhone, props.orderNumber)
        if (order) {
          liveStatus.value = order.status
        }
      }
    } catch (err) {
      refreshMessage.value = err instanceof Error ? err.message : ''
    } finally {
      isTracking.value = false
    }
  }
  </script>