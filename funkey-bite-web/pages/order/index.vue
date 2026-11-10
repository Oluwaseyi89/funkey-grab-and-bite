<template>
  <div class="container-narrow py-5 md:py-10">
    <div class="flex items-center gap-2">
      <NuxtLink to="/menu" class="icon-btn -ml-2" aria-label="Back to the menu">
        <ArrowLeft class="h-5 w-5" aria-hidden="true" />
      </NuxtLink>
      <h1 class="font-display text-2xl font-bold text-ink-900 dark:text-white md:text-3xl">Checkout</h1>
    </div>

    <!-- Empty cart -->
    <div v-if="mounted && cart.isEmpty" class="mx-auto max-w-sm py-16 text-center">
      <div class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-brand-50 dark:bg-brand-900/30">
        <ShoppingBag class="h-9 w-9 text-brand-500" aria-hidden="true" />
      </div>
      <h2 class="mt-4 text-xl font-semibold text-ink-900 dark:text-white">Your cart is empty</h2>
      <p class="mt-1 text-ink-500 dark:text-ink-400">Add a few things from the menu, then come back here to order.</p>
      <NuxtLink to="/menu" class="btn-primary btn-lg mt-6">Browse the menu</NuxtLink>
      <p class="mt-6 text-sm text-ink-500 dark:text-ink-400">Looking for an order you placed? <NuxtLink to="/order/track" class="link">Track it</NuxtLink></p>
    </div>

    <form v-else-if="mounted" class="mt-4 grid gap-6 lg:grid-cols-[1fr_380px] lg:items-start" novalidate @submit.prevent="placeOrder">
      <div class="min-w-0 space-y-4">
        <p v-if="blockedReason && settingsStatus !== 'loading'" role="alert" class="rounded-2xl bg-accent-50 px-4 py-3 text-sm text-accent-800 dark:bg-accent-700/15 dark:text-accent-200">
          {{ blockedReason }}
        </p>

        <!-- 1. How -->
        <section class="card p-5" aria-labelledby="how-title">
          <h2 id="how-title" class="mb-3 font-display text-lg font-semibold text-ink-900 dark:text-white">How would you like it?</h2>
          <OrderTypeSelector
            v-model="orderType"
            :delivery-fee="settings?.deliveryFee ?? null"
            :pickup-open="settings ? settings.isPickupOpen : undefined"
            :delivery-open="settings ? settings.isDeliveryOpen : undefined"
          />
          <p class="mt-3 text-sm text-ink-500 dark:text-ink-400">
            Feeding a crowd? <NuxtLink to="/catering" class="link">Request catering</NuxtLink>
          </p>
        </section>

        <!-- 2. Items -->
        <section class="card p-5" aria-labelledby="items-title">
          <div class="flex items-center justify-between">
            <h2 id="items-title" class="font-display text-lg font-semibold text-ink-900 dark:text-white">Your items</h2>
            <NuxtLink to="/menu" class="link inline-flex min-h-[44px] items-center gap-1 text-sm"><Plus class="h-4 w-4" aria-hidden="true" />Add more</NuxtLink>
          </div>
          <ul class="divide-y divide-ink-100 dark:divide-ink-800">
            <CartLine v-for="item in cart.items" :key="item.menuItem.id" :item="item" />
          </ul>
        </section>

        <!-- 3. Details -->
        <section class="card p-5" aria-labelledby="details-title">
          <div class="mb-4 flex items-start justify-between gap-3">
            <h2 id="details-title" class="font-display text-lg font-semibold text-ink-900 dark:text-white">Your details</h2>
            <NuxtLink v-if="!signedIn" :to="{ path: '/account/login', query: { next: '/order' } }" class="link shrink-0 text-sm">Sign in to fill these in</NuxtLink>
          </div>
          <CustomerInfoForm ref="detailsForm" v-model="details" :order-type="orderType" />
        </section>

        <!-- 4. Payment -->
        <section class="card p-5" aria-labelledby="pay-title">
          <h2 id="pay-title" class="mb-3 font-display text-lg font-semibold text-ink-900 dark:text-white">Payment</h2>
          <fieldset>
            <legend class="sr-only">How would you like to pay?</legend>
            <div class="space-y-2">
              <label
                v-for="method in paymentMethods"
                :key="method.id"
                class="flex cursor-pointer items-center gap-3 rounded-2xl border-2 p-4 transition has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-brand-500"
                :class="payment === method.id ? 'border-brand-500 bg-brand-50/60 dark:bg-brand-900/15' : 'border-ink-200 hover:border-ink-300 dark:border-ink-700'"
              >
                <input v-model="payment" type="radio" name="payment" :value="method.id" class="sr-only" />
                <span class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl" :class="method.tint">
                  <component :is="method.icon" class="h-5 w-5" aria-hidden="true" />
                </span>
                <span class="min-w-0 flex-1">
                  <span class="block font-semibold text-ink-900 dark:text-white">{{ method.name }}</span>
                  <span class="block text-sm text-ink-500 dark:text-ink-400">{{ method.description }}</span>
                </span>
                <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full border-2" :class="payment === method.id ? 'border-brand-600 bg-brand-600' : 'border-ink-300 dark:border-ink-600'" aria-hidden="true">
                  <span v-if="payment === method.id" class="h-2 w-2 rounded-full bg-white" />
                </span>
              </label>
            </div>
          </fieldset>
        </section>
      </div>

      <!-- Summary: under the form on phones, alongside on desktop. -->
      <aside class="space-y-4 lg:sticky lg:top-24" aria-labelledby="summary-title">
        <section class="card p-5">
          <h2 id="summary-title" class="font-display text-lg font-semibold text-ink-900 dark:text-white">Summary</h2>

          <!-- Promo code, checked by the server -->
          <div class="mt-4">
            <div v-if="promo.applied" class="flex items-center justify-between gap-3 rounded-xl bg-leaf-50 px-3 py-2.5 text-sm dark:bg-leaf-800/25">
              <span class="inline-flex items-center gap-2 font-medium text-leaf-800 dark:text-leaf-200">
                <TicketPercent class="h-4 w-4" aria-hidden="true" />
                <span class="font-mono">{{ promo.applied }}</span> applied
              </span>
              <button type="button" class="min-h-[36px] px-2 font-medium text-ink-600 hover:text-brand-600 dark:text-ink-300" @click="removePromo">Remove</button>
            </div>
            <div v-else>
              <label for="promo-code" class="form-label">Promo code</label>
              <div class="flex gap-2">
                <input
                  id="promo-code"
                  v-model="promo.input"
                  autocomplete="off"
                  autocapitalize="characters"
                  class="form-input font-mono"
                  placeholder="Enter code"
                  :aria-invalid="promo.error ? 'true' : undefined"
                  :aria-describedby="promo.error ? 'promo-error' : undefined"
                  @keydown.enter.prevent="applyPromo"
                />
                <button type="button" class="btn-secondary shrink-0" :disabled="!promo.input.trim() || promo.checking" @click="applyPromo">
                  <Loader2 v-if="promo.checking" class="h-4 w-4 animate-spin" aria-hidden="true" />
                  Apply
                </button>
              </div>
              <p v-if="promo.error" id="promo-error" class="form-error" role="alert">{{ promo.error }}</p>
            </div>
          </div>

          <dl class="mt-4 space-y-2 text-sm">
            <div class="flex justify-between text-ink-600 dark:text-ink-300"><dt>Items ({{ cart.totalItems }})</dt><dd class="tabular-nums">{{ naira(cart.subtotal) }}</dd></div>
            <template v-if="settings">
              <div v-if="orderType === 'delivery'" class="flex justify-between text-ink-600 dark:text-ink-300"><dt>Delivery</dt><dd class="tabular-nums">{{ deliveryFee ? naira(deliveryFee) : 'Free' }}</dd></div>
              <div v-if="tax > 0" class="flex justify-between text-ink-600 dark:text-ink-300"><dt>Tax ({{ settings.taxRate }}%)</dt><dd class="tabular-nums">{{ naira(tax) }}</dd></div>
              <div v-if="discount > 0" class="flex justify-between text-leaf-700 dark:text-leaf-300"><dt>Discount</dt><dd class="tabular-nums">−{{ naira(discount) }}</dd></div>
              <div class="flex justify-between border-t border-ink-100 pt-3 text-lg font-bold text-ink-900 dark:border-ink-800 dark:text-white"><dt>Total</dt><dd class="tabular-nums">{{ naira(total) }}</dd></div>
            </template>
            <div v-else class="flex justify-between border-t border-ink-100 pt-3 text-ink-500 dark:border-ink-800">
              <dt>Fees and total</dt><dd>{{ settingsStatus === 'loading' ? 'Loading…' : 'Unavailable' }}</dd>
            </div>
          </dl>
          <p v-if="belowMinimum" class="mt-3 text-sm text-accent-700 dark:text-accent-300">
            The minimum order is {{ naira(settings?.minOrderAmount ?? 0) }}. Add {{ naira((settings?.minOrderAmount ?? 0) - cart.subtotal) }} more to order.
          </p>

          <button type="submit" class="btn-primary btn-lg mt-5 hidden w-full lg:flex" :disabled="!canPlace">
            <Loader2 v-if="placing" class="h-5 w-5 animate-spin" aria-hidden="true" />
            {{ placeLabel }}
          </button>
          <p class="mt-3 flex items-center justify-center gap-1.5 text-xs text-ink-500 dark:text-ink-400">
            <Lock class="h-3.5 w-3.5" aria-hidden="true" /> Card payments go through Paystack. We never see your card.
          </p>
        </section>
      </aside>

      <!-- Phones: the place-order button stays in reach. -->
      <div class="fixed inset-x-0 z-40 border-t border-ink-200/70 bg-white/95 px-4 pt-3 backdrop-blur-xl dark:border-ink-800 dark:bg-ink-950/95 lg:hidden" style="bottom: calc(4rem + env(safe-area-inset-bottom)); padding-bottom: 0.75rem">
        <button type="submit" class="btn-primary btn-lg w-full justify-between" :disabled="!canPlace">
          <span class="inline-flex items-center gap-2">
            <Loader2 v-if="placing" class="h-5 w-5 animate-spin" aria-hidden="true" />
            {{ placing ? 'Placing order…' : 'Place order' }}
          </span>
          <span v-if="settings" class="tabular-nums">{{ naira(total) }}</span>
        </button>
      </div>
      <!-- Room so the fixed button never covers the summary. -->
      <div class="h-20 lg:hidden" aria-hidden="true" />
    </form>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useNuxtApp, useState } from 'nuxt/app'
import { ArrowLeft, Banknote, CreditCard, Landmark, Loader2, Lock, Plus, ShoppingBag, TicketPercent } from 'lucide-vue-next'
import CartLine from '../../components/cart/CartLine.vue'
import OrderTypeSelector from '../../components/order/OrderTypeSelector.vue'
import CustomerInfoForm from '../../components/order/CustomerInfoForm.vue'
import { useCartStore } from '../../stores/cart'
import { useAuthStore } from '../../stores/auth'
import { useApi } from '../../utils/api'
import { checkoutKeyFor, clearCheckoutKey } from '../../utils/checkoutKey'
import { rememberOrder } from '../../utils/recentOrders'
import { naira } from '../../utils/format'
import { usePaystack } from '../../composables/usePaystack'
import { useBusinessSettings } from '../../composables/useBusinessSettings'
import { usePageSeo } from '../../composables/usePageSeo'
import type { CustomerInfo, InvoiceLookup, Order, OrderCreatePayload, OrderType, PaymentMethod } from '../../types/order'

usePageSeo({ title: 'Checkout', description: 'Review your order and pay.', path: '/order', noindex: true })

const cart = useCartStore()
const auth = useAuthStore()
const api = useApi()
const router = useRouter()
const { $toast } = useNuxtApp()
const toast = $toast as { error?: (message: string) => void } | undefined
const { cardPaymentsEnabled, payWithCard } = usePaystack()
const { settings, status: settingsStatus, load: loadSettings } = useBusinessSettings()

const mounted = ref(false)
const signedIn = ref(false)
const placing = ref(false)
const orderType = ref<OrderType>('pickup')
const payment = ref<PaymentMethod>('transfer')
const details = ref<CustomerInfo>({ customerName: '', customerPhone: '', customerEmail: '', deliveryAddress: '', notes: '' })
const detailsForm = ref<InstanceType<typeof CustomerInfoForm>>()

// ---- Money ----
const deliveryFee = computed(() => (orderType.value === 'delivery' ? settings.value?.deliveryFee ?? 0 : 0))
const tax = computed(() => Math.round(cart.subtotal * ((settings.value?.taxRate ?? 0) / 100) * 100) / 100)
const beforeDiscount = computed(() => cart.subtotal + tax.value + deliveryFee.value)

// The server decides the discount; this shows what it said for this total.
const promo = reactive({ input: '', applied: '', discount: 0, checkedFor: 0, checking: false, error: '' })
const discount = computed(() => (promo.applied ? Math.min(promo.discount, beforeDiscount.value) : 0))
const total = computed(() => Math.max(0, beforeDiscount.value - discount.value))
const belowMinimum = computed(() => Boolean(settings.value?.minOrderAmount) && cart.subtotal < (settings.value?.minOrderAmount ?? 0))

const checkCode = async (code: string) => {
  promo.checking = true
  promo.error = ''
  try {
    const result = await api.checkPromotion(code, beforeDiscount.value)
    if (result.isValid) {
      promo.applied = code
      promo.discount = result.discount ?? 0
      promo.checkedFor = beforeDiscount.value
      promo.input = ''
      cart.promoCode = promo.applied
    } else {
      promo.applied = ''
      promo.error = result.message || 'That code can’t be used on this order.'
      cart.promoCode = ''
    }
  } catch (err) {
    promo.error = err instanceof Error ? err.message : 'We couldn’t check that code.'
  } finally {
    promo.checking = false
  }
}

const applyPromo = () => {
  const code = promo.input.trim()
  if (code) checkCode(code)
}

const removePromo = () => {
  promo.applied = ''
  promo.discount = 0
  cart.promoCode = ''
}

// A percentage discount changes with the total, so ask again when it does.
let recheck: ReturnType<typeof setTimeout> | undefined
watch(beforeDiscount, (amount) => {
  if (!promo.applied || amount === promo.checkedFor) return
  clearTimeout(recheck)
  recheck = setTimeout(() => checkCode(promo.applied), 600)
})

// ---- Payment methods ----
const paymentMethods = computed(() => {
  const list: { id: PaymentMethod; name: string; description: string; icon: unknown; tint: string }[] = []
  if (cardPaymentsEnabled) {
    list.push({ id: 'card', name: 'Card', description: 'Pay now in Paystack’s secure window.', icon: CreditCard, tint: 'bg-leaf-50 text-leaf-700 dark:bg-leaf-800/30 dark:text-leaf-300' })
  }
  list.push({ id: 'transfer', name: 'Bank transfer', description: 'We’ll give you an account number to pay into.', icon: Landmark, tint: 'bg-brand-50 text-brand-700 dark:bg-brand-900/30 dark:text-brand-300' })
  if (settings.value?.payOnDeliveryEnabled) {
    const delivery = orderType.value === 'delivery'
    list.push({
      id: 'cash',
      name: delivery ? 'Pay on delivery' : 'Pay at pickup',
      description: `Cash, POS or transfer when you ${delivery ? 'receive' : 'collect'} it.`,
      icon: Banknote,
      tint: 'bg-accent-50 text-accent-700 dark:bg-accent-700/20 dark:text-accent-300',
    })
  }
  return list
})

// Never leave a method selected that's no longer offered.
watch(paymentMethods, (methods) => {
  if (!methods.some((m) => m.id === payment.value)) payment.value = methods[0]?.id ?? 'transfer'
}, { immediate: true })

// ---- Can the order be placed? (mirrors the server's checks) ----
const blockedReason = computed(() => {
  const s = settings.value
  if (settingsStatus.value === 'loading') return ''
  if (!s) return 'We couldn’t load our prices and opening hours. Please refresh the page.'
  if (!s.isOpenNow) return 'We’re closed right now, so orders can’t be placed until we open. Your cart is saved.'
  if (orderType.value === 'delivery' && !s.isDeliveryOpen) return 'Delivery isn’t available right now. Choose pickup instead.'
  if (orderType.value === 'pickup' && !s.isPickupOpen) return 'Pickup isn’t available right now. Choose delivery instead.'
  return ''
})

const canPlace = computed(() => !placing.value && !!settings.value && !blockedReason.value && !belowMinimum.value && !cart.isEmpty)
const placeLabel = computed(() => (placing.value ? 'Placing order…' : settings.value ? `Place order · ${naira(total.value)}` : 'Place order'))

const fail = (message: string) => toast?.error?.(message)

const placeOrder = async () => {
  if (!canPlace.value) return
  if (!detailsForm.value?.validate()) {
    fail('Please check your details.')
    return
  }
  placing.value = true
  try {
    // The page may have been open a while: check hours and fees again.
    await loadSettings(true)
    if (blockedReason.value) {
      fail(blockedReason.value)
      return
    }

    const d = details.value
    const payload: OrderCreatePayload = {
      customerName: d.customerName.trim(),
      customerPhone: d.customerPhone.trim(),
      customerEmail: d.customerEmail.trim(),
      orderType: orderType.value,
      paymentMethod: payment.value,
      deliveryAddress: orderType.value === 'delivery' ? d.deliveryAddress.trim() : undefined,
      notes: d.notes.trim() || undefined,
      promotionCode: promo.applied || undefined,
      items: cart.items.map((item) => ({
        menuItemId: Number(item.menuItem.id),
        name: item.menuItem.name,
        quantity: item.quantity,
        unitPrice: item.menuItem.price,
        specialInstructions: item.specialInstructions,
      })),
    }

    // The same key is reused if this exact order is sent again, so it's placed once.
    let order: Order = await api.createOrder(payload, checkoutKeyFor(payload))
    clearCheckoutKey()
    rememberOrder({ orderNumber: order.orderNumber, phone: payload.customerPhone, placedAt: order.createdAt, total: order.totalAmount })
    cart.clearCart()

    // Card orders open Paystack straight away; closing it leaves a pay button
    // on the order page.
    if (order.paymentMethod === 'card' && order.paymentAccessCode) {
      try {
        if ((await payWithCard(order.paymentAccessCode)) === 'paid') {
          order = (await api.refreshPaymentStatus(order.orderNumber, payload.customerPhone)).order
        }
      } catch (err) {
        fail(err instanceof Error ? err.message : 'The card payment window could not open.')
      }
    }

    useState<InvoiceLookup | null>('track-lookup', () => null).value = { orderNumber: order.orderNumber, phone: payload.customerPhone }
    await router.push('/order/track')
  } catch (err) {
    fail(err instanceof Error ? err.message : 'Your order didn’t go through. Please try again.')
  } finally {
    placing.value = false
  }
}

onMounted(() => {
  mounted.value = true
  signedIn.value = auth.isSignedIn
  loadSettings()
  // Signed-in customers start with their details filled in.
  const user = auth.user
  if (user) {
    details.value = { ...details.value, customerName: user.fullName, customerPhone: user.phone, customerEmail: user.email ?? '' }
  }
  // A code chosen on the deals page.
  if (cart.promoCode) {
    promo.input = cart.promoCode
    watch(settings, (s) => { if (s && promo.input && !promo.applied) checkCode(promo.input) }, { immediate: true })
  }
})
</script>
