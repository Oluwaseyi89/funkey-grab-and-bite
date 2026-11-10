<template>
  <div class="container-narrow max-w-2xl py-5 md:py-10">
    <section v-if="auth.user" class="card flex items-center gap-4 p-5">
      <span class="flex h-14 w-14 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-brand-500 to-accent-500 font-display text-xl font-bold text-white" aria-hidden="true">
        {{ initials }}
      </span>
      <div class="min-w-0 flex-1">
        <h1 class="truncate font-display text-xl font-bold text-ink-900 dark:text-white">Hi, {{ firstName }}</h1>
        <p class="truncate text-sm text-ink-500 dark:text-ink-400">{{ auth.user.phone }}<span v-if="auth.user.email"> · {{ auth.user.email }}</span></p>
      </div>
      <button type="button" class="btn-ghost shrink-0 text-sm" @click="signOut">
        <LogOut class="h-4 w-4" aria-hidden="true" /> Sign out
      </button>
    </section>

    <section class="mt-6" aria-labelledby="orders-title">
      <div class="flex items-center justify-between">
        <h2 id="orders-title" class="font-display text-lg font-semibold text-ink-900 dark:text-white">Your orders</h2>
        <NuxtLink to="/menu" class="link inline-flex min-h-[44px] items-center gap-1 text-sm"><Plus class="h-4 w-4" aria-hidden="true" />New order</NuxtLink>
      </div>

      <div v-if="loading" class="mt-3 space-y-3" aria-busy="true" aria-label="Loading your orders">
        <div v-for="i in 3" :key="i" class="skeleton h-24 rounded-2xl" />
      </div>
      <p v-else-if="loadError" class="card mt-3 p-5 text-brand-700 dark:text-brand-300" role="alert">{{ loadError }}</p>
      <div v-else-if="orders.length === 0" class="card mt-3 p-8 text-center">
        <Receipt class="mx-auto h-10 w-10 text-ink-300" aria-hidden="true" />
        <p class="mt-3 font-semibold text-ink-900 dark:text-white">No orders yet</p>
        <p class="mt-1 text-sm text-ink-500 dark:text-ink-400">Orders you place while signed in, and earlier ones with your number and email, show up here.</p>
        <NuxtLink to="/menu" class="btn-primary mt-5">Browse the menu</NuxtLink>
      </div>

      <ul v-else class="mt-3 space-y-3">
        <li v-for="order in orders" :key="order.id">
          <NuxtLink :to="`/account/orders/${order.id}`" class="card flex items-center gap-4 p-4 transition hover:shadow-lift">
            <span class="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl" :class="needsPayment(order) ? 'bg-accent-50 text-accent-700 dark:bg-accent-700/20 dark:text-accent-300' : 'bg-ink-100 text-ink-600 dark:bg-ink-800 dark:text-ink-300'">
              <component :is="order.orderType === 'delivery' ? Bike : Store" class="h-6 w-6" aria-hidden="true" />
            </span>
            <span class="min-w-0 flex-1">
              <span class="flex items-baseline justify-between gap-3">
                <span class="truncate font-mono text-sm font-semibold text-ink-900 dark:text-white">{{ order.orderNumber }}</span>
                <span class="shrink-0 font-semibold tabular-nums text-ink-900 dark:text-white">{{ naira(order.totalAmount) }}</span>
              </span>
              <span class="block text-xs text-ink-500 dark:text-ink-400">{{ formatDateTime(order.createdAt) }} · {{ plural(itemCount(order), 'item') }}</span>
              <span class="mt-2 flex flex-wrap gap-1.5">
                <span class="badge" :class="order.status === 'cancelled' ? 'bg-ink-100 text-ink-600 dark:bg-ink-800 dark:text-ink-300' : 'bg-leaf-50 text-leaf-800 dark:bg-leaf-800/30 dark:text-leaf-200'">
                  {{ ORDER_STATUS_LABELS[order.status] ?? order.status }}
                </span>
                <span v-if="needsPayment(order)" class="badge bg-accent-100 text-accent-800 dark:bg-accent-700/30 dark:text-accent-200">Payment needed</span>
                <span v-else-if="order.paymentStatus === 'paid'" class="badge bg-ink-100 text-ink-600 dark:bg-ink-800 dark:text-ink-300">Paid</span>
              </span>
            </span>
            <ChevronRight class="h-5 w-5 shrink-0 text-ink-300" aria-hidden="true" />
          </NuxtLink>
        </li>
      </ul>

      <nav v-if="totalPages > 1" class="mt-4 flex items-center justify-between text-sm" aria-label="Order pages">
        <button type="button" class="btn-secondary" :disabled="page <= 1 || loading" @click="load(page - 1)">Newer</button>
        <span class="text-ink-500 dark:text-ink-400">Page {{ page }} of {{ totalPages }}</span>
        <button type="button" class="btn-secondary" :disabled="page >= totalPages || loading" @click="load(page + 1)">Older</button>
      </nav>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Bike, ChevronRight, LogOut, Plus, Receipt, Store } from 'lucide-vue-next'
import { ApiRequestError, useApi } from '../../utils/api'
import { useAuthStore } from '../../stores/auth'
import { ORDER_STATUS_LABELS } from '../../utils/paymentLabels'
import { formatDateTime, naira, plural } from '../../utils/format'
import { usePageSeo } from '../../composables/usePageSeo'
import type { Order } from '../../types/order'

usePageSeo({ title: 'My account', description: 'Your orders and payments.', path: '/account', noindex: true })

const router = useRouter()
const auth = useAuthStore()
const api = useApi()

const orders = ref<Order[]>([])
const page = ref(1)
const totalPages = ref(0)
const loading = ref(true)
const loadError = ref('')

const firstName = computed(() => auth.user?.fullName.split(' ')[0] ?? '')
const initials = computed(() =>
  (auth.user?.fullName ?? '').split(' ').filter(Boolean).slice(0, 2).map((part) => part[0]?.toUpperCase()).join(''),
)

const itemCount = (order: Order) => order.items.reduce((sum, item) => sum + item.quantity, 0)
const needsPayment = (order: Order) =>
  (order.paymentMethod === 'transfer' || order.paymentMethod === 'card') && order.paymentStatus === 'pending' && order.status !== 'cancelled'

// An expired session sends the customer to sign in again.
const handleAuthError = (err: unknown) => {
  if (err instanceof ApiRequestError && err.status === 401) {
    auth.signOut()
    router.push({ path: '/account/login', query: { next: '/account' } })
    return true
  }
  return false
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
