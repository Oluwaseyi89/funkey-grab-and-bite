<template>
  <div class="container-narrow max-w-2xl py-5 md:py-10">
    <!-- Just placed -->
    <div v-if="justPlaced && result" class="mb-4 animate-rise-in rounded-3xl bg-gradient-to-br from-leaf-500 to-leaf-600 p-5 text-white shadow-lift" role="status">
      <p class="flex items-center gap-2 font-display text-2xl font-bold">
        <PartyPopper class="h-7 w-7" aria-hidden="true" /> Thank you, {{ firstName }}!
      </p>
      <p class="mt-1 text-white/90">Your order is in. We’ve emailed you the details, and this page keeps itself up to date.</p>
    </div>

    <template v-if="result">
      <div class="mb-3 flex items-center justify-between gap-3">
        <h1 class="font-display text-2xl font-bold text-ink-900 dark:text-white md:text-3xl">{{ justPlaced ? 'Your order' : 'Order status' }}</h1>
        <button type="button" class="btn-ghost text-sm" @click="reset">Find another</button>
      </div>
      <OrderStatusView
        :order="result.order"
        :invoice="result.invoice"
        :reload="reload"
        :check-payment-status="checkPaymentStatus"
        @update="onUpdate"
      />
      <NuxtLink to="/menu" class="btn-secondary mt-4 w-full">Order something else</NuxtLink>
    </template>

    <template v-else>
      <p class="eyebrow">Orders</p>
      <h1 class="mt-1 font-display text-3xl font-bold text-ink-900 dark:text-white md:text-4xl">Track your order</h1>
      <p class="mt-2 text-ink-600 dark:text-ink-300">See where your order is, pay for it, or get its invoice.</p>

      <div v-if="loading && !result" class="mt-6 space-y-3" aria-busy="true">
        <div class="skeleton h-40 rounded-2xl" />
        <div class="skeleton h-24 rounded-2xl" />
      </div>

      <template v-else>
        <section v-if="recent.length" class="mt-6" aria-labelledby="recent-title">
          <h2 id="recent-title" class="text-sm font-semibold text-ink-700 dark:text-ink-200">On this phone</h2>
          <ul class="mt-2 space-y-2">
            <li v-for="entry in recent" :key="entry.orderNumber">
              <button type="button" class="card flex w-full items-center gap-3 p-4 text-left transition hover:shadow-lift" @click="open(entry.orderNumber, entry.phone)">
                <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-brand-50 text-brand-700 dark:bg-brand-900/30 dark:text-brand-300">
                  <Receipt class="h-5 w-5" aria-hidden="true" />
                </span>
                <span class="min-w-0 flex-1">
                  <span class="block truncate font-mono text-sm font-semibold text-ink-900 dark:text-white">{{ entry.orderNumber }}</span>
                  <span class="block text-xs text-ink-500 dark:text-ink-400">{{ formatDateTime(entry.placedAt) }} · {{ naira(entry.total) }}</span>
                </span>
                <ChevronRight class="h-5 w-5 text-ink-400" aria-hidden="true" />
              </button>
            </li>
          </ul>
        </section>

        <form class="card mt-6 space-y-4 p-5" novalidate @submit.prevent="submit">
          <h2 class="font-display text-lg font-semibold text-ink-900 dark:text-white">{{ recent.length ? 'Another order' : 'Find your order' }}</h2>
          <div>
            <label for="track-number" class="form-label">Order number</label>
            <input
              id="track-number"
              v-model="orderNumber"
              autocomplete="off"
              autocapitalize="characters"
              class="form-input font-mono uppercase"
              placeholder="FG-2026-11-…"
              :aria-invalid="!!error || undefined"
              aria-describedby="track-number-hint"
            />
            <p id="track-number-hint" class="form-hint">It’s in your confirmation email.</p>
          </div>
          <div>
            <label for="track-phone" class="form-label">Phone number used for the order</label>
            <input id="track-phone" v-model="phone" type="tel" inputmode="tel" autocomplete="tel" class="form-input" placeholder="0801 234 5678" :aria-invalid="!!error || undefined" />
          </div>
          <p v-if="error" class="form-error" role="alert">{{ error }}</p>
          <button type="submit" class="btn-primary btn-lg w-full" :disabled="loading">
            <Loader2 v-if="loading" class="h-5 w-5 animate-spin" aria-hidden="true" />
            Track order
          </button>
        </form>

        <p v-if="signedIn" class="mt-4 text-center text-sm text-ink-600 dark:text-ink-300">
          Signed in? <NuxtLink to="/account" class="link">See all your orders</NuxtLink>
        </p>
        <p v-else class="mt-4 text-center text-sm text-ink-600 dark:text-ink-300">
          Have an account? <NuxtLink :to="{ path: '/account/login', query: { next: '/account' } }" class="link">Sign in to see every order</NuxtLink>
        </p>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useState } from 'nuxt/app'
import { ChevronRight, Loader2, PartyPopper, Receipt } from 'lucide-vue-next'
import OrderStatusView from '../../components/order/OrderStatusView.vue'
import { useApi } from '../../utils/api'
import { useAuthStore } from '../../stores/auth'
import { readRecentOrders, type RecentOrder } from '../../utils/recentOrders'
import { formatDateTime, naira } from '../../utils/format'
import { usePageSeo } from '../../composables/usePageSeo'
import type { Invoice, Order, OrderLookupResult, InvoiceLookup } from '../../types/order'

usePageSeo({ title: 'Track your order', description: 'Check your order’s progress and payment.', path: '/order/track', noindex: true })

const api = useApi()
const auth = useAuthStore()

// Set by checkout so the new order opens straight away. Kept in app state,
// never in the address bar, since it includes the phone number.
const handoff = useState<InvoiceLookup | null>('track-lookup', () => null)

const orderNumber = ref('')
const phone = ref('')
const result = ref<OrderLookupResult | null>(null)
const loading = ref(false)
const error = ref('')
const justPlaced = ref(false)
const recent = ref<RecentOrder[]>([])
const signedIn = ref(false)

// The number and phone of the order on screen, for refreshing it.
const current = ref<{ orderNumber: string; phone: string } | null>(null)

const firstName = computed(() => result.value?.order.customerName.split(' ')[0] ?? '')

const open = async (number: string, customerPhone: string) => {
  loading.value = true
  error.value = ''
  try {
    const found = await api.lookupOrder(number.trim(), customerPhone.trim())
    if (found) {
      result.value = found
      current.value = { orderNumber: found.order.orderNumber, phone: customerPhone.trim() }
      window.scrollTo({ top: 0 })
    } else {
      error.value = 'We couldn’t find an order with that number and phone. Check both and try again.'
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Something went wrong. Please try again.'
  } finally {
    loading.value = false
  }
}

const submit = () => {
  if (!orderNumber.value.trim() || !phone.value.trim()) {
    error.value = 'Enter both the order number and the phone number.'
    return
  }
  open(orderNumber.value.toUpperCase(), phone.value)
}

const reload = async () => {
  const found = current.value ? await api.lookupOrder(current.value.orderNumber, current.value.phone) : null
  if (!found) throw new Error('Order not found')
  return found
}

const checkPaymentStatus = () => {
  if (!current.value) return Promise.reject(new Error('No order open'))
  return api.refreshPaymentStatus(current.value.orderNumber, current.value.phone)
}

const onUpdate = ({ order, invoice }: { order: Order; invoice?: Invoice | null }) => {
  if (!result.value) return
  result.value = { order, invoice: invoice === undefined ? result.value.invoice : invoice }
}

const reset = () => {
  result.value = null
  current.value = null
  justPlaced.value = false
  orderNumber.value = ''
  phone.value = ''
  recent.value = readRecentOrders()
}

onMounted(() => {
  signedIn.value = auth.isSignedIn
  recent.value = readRecentOrders()
  if (handoff.value) {
    const { orderNumber: number, phone: customerPhone } = handoff.value
    handoff.value = null
    justPlaced.value = true
    open(number, customerPhone)
  }
})
</script>
