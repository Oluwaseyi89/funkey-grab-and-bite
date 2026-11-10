<template>
  <div class="mt-[-30px]">
    <PageHeader title-before="My" highlight-text="Account" subtitle="Your orders and payments" :narrow="true" :alignment="'center'" variant="gradient" />

    <div class="container-narrow px-4 sm:px-6 lg:px-8 py-8 max-w-3xl space-y-6">
      <div v-if="auth.user" class="bg-white dark:bg-slate-800 rounded-xl p-5 border border-gray-200 dark:border-slate-700 flex flex-wrap items-center justify-between gap-4">
        <div>
          <p class="font-semibold text-gray-900 dark:text-white">{{ auth.user.fullName }}</p>
          <p class="text-sm text-gray-600 dark:text-gray-400">{{ auth.user.phone }}<span v-if="auth.user.email"> · {{ auth.user.email }}</span></p>
        </div>
        <div class="flex gap-3">
          <NuxtLink to="/menu" class="btn-primary dark:text-white">Order food</NuxtLink>
          <button class="btn-secondary" @click="signOut">Sign out</button>
        </div>
      </div>

      <section>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">Your orders</h2>
        <div v-if="loading" class="flex justify-center py-10 text-gray-500"><Loader2 class="w-6 h-6 animate-spin" /></div>
        <p v-else-if="loadError" class="text-red-600" role="alert">{{ loadError }}</p>
        <div v-else-if="orders.length === 0" class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 text-center text-gray-600 dark:text-gray-400">
          No orders yet. Orders you place while signed in, and earlier ones placed with your number and email, show up here.
        </div>
        <ul v-else class="space-y-3">
          <li v-for="order in orders" :key="order.id" class="bg-white dark:bg-slate-800 rounded-xl p-5 border border-gray-200 dark:border-slate-700">
            <div class="flex flex-wrap items-start justify-between gap-3">
              <div>
                <NuxtLink :to="`/account/orders/${order.id}`" class="font-semibold text-gray-900 dark:text-white hover:underline">#{{ order.orderNumber }}</NuxtLink>
                <p class="text-sm text-gray-500 dark:text-gray-400">{{ formatDate(order.createdAt) }} · {{ order.items.length }} item{{ order.items.length === 1 ? '' : 's' }}</p>
              </div>
              <p class="font-bold text-gray-900 dark:text-white">{{ formatNaira(order.totalAmount) }}</p>
            </div>

            <div class="mt-3 flex flex-wrap gap-2 text-xs font-medium">
              <span class="px-2 py-1 rounded-full bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300">{{ ORDER_STATUS_LABELS[order.status] ?? order.status }}</span>
              <span v-if="order.paymentMethod" class="px-2 py-1 rounded-full bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300">{{ ORDER_PAYMENT_METHOD_LABELS[order.paymentMethod] }}</span>
              <span v-if="order.paymentStatus" class="px-2 py-1 rounded-full" :class="paymentStatusClass(order.paymentStatus)">{{ PAYMENT_STATUS_LABELS[order.paymentStatus] ?? order.paymentStatus }}</span>
            </div>

            <div v-if="order.paymentMethod === 'transfer' && order.paymentStatus === 'pending' && order.paymentAccountNumber" class="mt-3 text-sm text-gray-700 dark:text-gray-300">
              Transfer {{ formatNaira(order.totalAmount) }} to <strong class="font-mono">{{ order.paymentAccountNumber }}</strong>
              <span v-if="order.paymentBankName"> ({{ order.paymentBankName }})</span>
            </div>

            <p v-if="messages[order.id]" class="mt-3 text-sm text-gray-700 dark:text-gray-300" role="status">{{ messages[order.id] }}</p>

            <div class="mt-3 flex flex-wrap gap-2">
              <NuxtLink :to="`/account/orders/${order.id}`" class="btn-secondary text-sm">View details</NuxtLink>
              <button
                v-if="order.paymentMethod === 'card' && order.paymentStatus === 'pending' && cardPaymentsEnabled"
                class="btn-primary text-sm dark:text-white"
                :disabled="busy[order.id]"
                @click="payNow(order)"
              >
                Pay with card
              </button>
              <button
                v-if="isOnline(order) && order.paymentStatus === 'pending'"
                class="btn-secondary text-sm inline-flex items-center gap-2"
                :disabled="busy[order.id]"
                @click="refresh(order)"
              >
                <Loader2 v-if="busy[order.id]" class="w-4 h-4 animate-spin" /><RefreshCw v-else class="w-4 h-4" />
                Refresh payment status
              </button>
            </div>
          </li>
        </ul>
        <div v-if="totalPages > 1" class="mt-4 flex justify-between items-center text-sm">
          <button class="btn-secondary" :disabled="page <= 1 || loading" @click="load(page - 1)">Newer</button>
          <span class="text-gray-500 dark:text-gray-400">Page {{ page }} of {{ totalPages }}</span>
          <button class="btn-secondary" :disabled="page >= totalPages || loading" @click="load(page + 1)">Older</button>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Loader2, RefreshCw } from 'lucide-vue-next'
import PageHeader from '../../components/layout/PageHeader.vue'
import { ApiRequestError, useApi } from '../../utils/api'
import { useAuthStore } from '../../stores/auth'
import { usePaystack } from '../../composables/usePaystack'
import { ORDER_STATUS_LABELS, ORDER_PAYMENT_METHOD_LABELS, PAYMENT_STATUS_LABELS, paymentStatusClass } from '../../utils/paymentLabels'
import type { Order } from '../../types/order'

const router = useRouter()
const auth = useAuthStore()
const api = useApi()
const { cardPaymentsEnabled, payWithCard } = usePaystack()

const orders = ref<Order[]>([])
const page = ref(1)
const totalPages = ref(0)
const loading = ref(true)
const loadError = ref('')
const busy = reactive<Record<string, boolean>>({})
const messages = reactive<Record<string, string>>({})

const formatNaira = (amount: number) =>
  `₦${amount.toLocaleString('en-NG', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
const formatDate = (value: string) => new Date(value).toLocaleString('en-NG', { dateStyle: 'medium', timeStyle: 'short' })
const isOnline = (order: Order) => order.paymentMethod === 'transfer' || order.paymentMethod === 'card'

// An expired session sends the customer to sign in again.
const handleAuthError = (err: unknown) => {
  if (err instanceof ApiRequestError && err.status === 401) {
    auth.signOut()
    router.push({ path: '/account/login', query: { next: '/account' } })
    return true
  }
  return false
}

const replace = (updated: Order) => {
  orders.value = orders.value.map((o) => (o.id === updated.id ? updated : o))
}

const load = async (target = 1) => {
  loading.value = true
  loadError.value = ''
  try {
    const result = await api.getMyOrders(target)
    orders.value = result.orders
    totalPages.value = result.totalPages
    page.value = target
  } catch (err) {
    if (!handleAuthError(err)) loadError.value = 'We couldn’t load your orders. Please try again.'
  } finally {
    loading.value = false
  }
}

// Asks the server to check the payment with Paystack.
const refresh = async (order: Order) => {
  busy[order.id] = true
  try {
    const result = await api.refreshMyPaymentStatus(order.id)
    replace(result.order)
    messages[order.id] = result.message || 'Status updated.'
  } catch (err) {
    if (!handleAuthError(err)) messages[order.id] = err instanceof Error ? err.message : 'Please try again.'
  } finally {
    busy[order.id] = false
  }
}

const payNow = async (order: Order) => {
  busy[order.id] = true
  messages[order.id] = ''
  try {
    const renewed = await api.getPaymentDetails(order.orderNumber, order.customerPhone)
    if (!renewed.paymentAccessCode) throw new Error('We couldn’t start the card payment. Please try again.')
    const outcome = await payWithCard(renewed.paymentAccessCode)
    busy[order.id] = false
    if (outcome === 'paid') {
      await refresh(order)
    } else {
      messages[order.id] = 'Payment cancelled. You can pay whenever you’re ready.'
    }
  } catch (err) {
    messages[order.id] = err instanceof Error ? err.message : 'Please try again.'
  } finally {
    busy[order.id] = false
  }
}

const signOut = () => {
  auth.signOut()
  router.push('/')
}

onMounted(async () => {
  if (!auth.isSignedIn) {
    router.replace({ path: '/account/login', query: { next: '/account' } })
    return
  }
  try {
    auth.setUser(await api.getProfile())
  } catch (err) {
    if (handleAuthError(err)) return
  }
  await load(1)
})
</script>
