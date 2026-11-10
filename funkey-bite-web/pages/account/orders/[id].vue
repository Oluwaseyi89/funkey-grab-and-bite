<template>
  <div class="container-narrow max-w-2xl py-5 md:py-10">
    <div class="mb-3 flex items-center gap-2">
      <NuxtLink to="/account" class="icon-btn -ml-2" aria-label="Back to my orders">
        <ArrowLeft class="h-5 w-5" aria-hidden="true" />
      </NuxtLink>
      <h1 class="font-display text-2xl font-bold text-ink-900 dark:text-white md:text-3xl">Order details</h1>
    </div>

    <div v-if="loading" class="space-y-3" aria-busy="true" aria-label="Loading your order">
      <div class="skeleton h-44 rounded-2xl" />
      <div class="skeleton h-28 rounded-2xl" />
      <div class="skeleton h-40 rounded-2xl" />
    </div>

    <div v-else-if="loadError" class="card p-8 text-center" role="alert">
      <p class="text-ink-700 dark:text-ink-200">{{ loadError }}</p>
      <div class="mt-4 flex justify-center gap-3">
        <button v-if="!notFound" type="button" class="btn-primary" @click="load">Try again</button>
        <NuxtLink to="/account" class="btn-secondary">My orders</NuxtLink>
      </div>
    </div>

    <OrderStatusView
      v-else-if="result"
      :order="result.order"
      :invoice="result.invoice"
      :reload="fetchOrder"
      :check-payment-status="checkPaymentStatus"
      @update="onUpdate"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'
import OrderStatusView from '../../../components/order/OrderStatusView.vue'
import { ApiRequestError, useApi } from '../../../utils/api'
import { useAuthStore } from '../../../stores/auth'
import { usePageSeo } from '../../../composables/usePageSeo'
import type { Invoice, Order, OrderLookupResult } from '../../../types/order'

usePageSeo({ title: 'Order details', description: 'Your order and its payment.', path: '/account', noindex: true })

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const api = useApi()

const orderId = computed(() => String(route.params.id ?? ''))
const result = ref<OrderLookupResult | null>(null)
const loading = ref(true)
const loadError = ref('')
const notFound = ref(false)

const sendToSignIn = () => {
  auth.signOut()
  router.replace({ path: '/account/login', query: { next: route.fullPath } })
}

const fetchOrder = () => api.getMyOrder(orderId.value)

// The server checks an unpaid online payment with Paystack.
const checkPaymentStatus = async () => {
  try {
    return await api.refreshMyPaymentStatus(orderId.value)
  } catch (err) {
    if (err instanceof ApiRequestError && err.status === 401) sendToSignIn()
    throw err
  }
}

const load = async () => {
  loading.value = true
  loadError.value = ''
  notFound.value = false
  try {
    result.value = await fetchOrder()
  } catch (err) {
    if (err instanceof ApiRequestError && err.status === 401) {
      sendToSignIn()
      return
    }
    notFound.value = err instanceof ApiRequestError && (err.status === 404 || err.status === 400)
    loadError.value = notFound.value ? 'We couldn’t find that order in your account.' : 'We couldn’t load this order. Please try again.'
  } finally {
    loading.value = false
  }
}

const onUpdate = ({ order, invoice }: { order: Order; invoice?: Invoice | null }) => {
  if (!result.value) return
  result.value = { order, invoice: invoice === undefined ? result.value.invoice : invoice }
}

onMounted(() => {
  if (!auth.isSignedIn) {
    router.replace({ path: '/account/login', query: { next: route.fullPath } })
    return
  }
  load()
})
</script>
