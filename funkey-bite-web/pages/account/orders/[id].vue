<template>
  <div class="mt-[-30px]">
    <PageHeader
      class="no-print"
      title-before="Order"
      highlight-text="Details"
      :subtitle="order ? `#${order.orderNumber}` : 'Your order'"
      :narrow="true"
      :alignment="'center'"
      variant="gradient"
    />

    <div class="container-narrow px-4 sm:px-6 lg:px-8 py-8 max-w-3xl space-y-6">
      <NuxtLink to="/account" class="no-print inline-flex items-center gap-1 text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white">
        <ArrowLeft class="w-4 h-4" /> Back to my orders
      </NuxtLink>

      <div v-if="loading" class="flex justify-center py-10 text-gray-500"><Loader2 class="w-6 h-6 animate-spin" /></div>

      <div v-else-if="loadError" class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 text-center space-y-3">
        <p class="text-gray-700 dark:text-gray-300" role="alert">{{ loadError }}</p>
        <button v-if="!notFound" class="btn-secondary" @click="load">Try again</button>
      </div>

      <template v-else-if="order">
        <!-- Summary and progress -->
        <section class="no-print bg-white dark:bg-slate-800 rounded-xl p-5 border border-gray-200 dark:border-slate-700">
          <div class="flex flex-wrap items-start justify-between gap-3">
            <div>
              <p class="font-semibold text-gray-900 dark:text-white">#{{ order.orderNumber }}</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">Placed {{ formatDate(order.createdAt) }}</p>
            </div>
            <span class="px-2 py-1 rounded-full text-xs font-medium" :class="order.status === 'cancelled' ? 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300' : 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300'">
              {{ ORDER_STATUS_LABELS[order.status] ?? order.status }}
            </span>
          </div>

          <p v-if="order.status === 'cancelled'" class="mt-4 text-sm text-gray-700 dark:text-gray-300">
            This order was cancelled.<span v-if="order.paymentStatus === 'expired'"> The payment wasn’t made in time.</span>
          </p>
          <ol v-else class="mt-5 grid grid-cols-5 gap-1" aria-label="Order progress">
            <li v-for="(step, index) in STEPS" :key="step" class="text-center">
              <div class="h-1.5 rounded-full" :class="index <= currentStep ? 'bg-brand-500' : 'bg-gray-200 dark:bg-slate-700'" />
              <p class="mt-2 text-[11px] sm:text-xs" :class="index <= currentStep ? 'text-gray-900 dark:text-white font-medium' : 'text-gray-500 dark:text-gray-400'" :aria-current="index === currentStep ? 'step' : undefined">
                {{ ORDER_STATUS_LABELS[step] }}
              </p>
            </li>
          </ol>
        </section>

        <!-- Items and amounts -->
        <section class="no-print bg-white dark:bg-slate-800 rounded-xl p-5 border border-gray-200 dark:border-slate-700">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">Items</h2>
          <ul class="divide-y divide-gray-100 dark:divide-slate-700">
            <li v-for="(item, index) in order.items" :key="`${item.menuItemId}-${index}`" class="py-3 flex justify-between gap-4">
              <div>
                <p class="text-gray-900 dark:text-white">{{ item.quantity }} × {{ item.name }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ formatNaira(item.unitPrice) }} each</p>
                <p v-if="item.specialInstructions" class="text-xs text-gray-600 dark:text-gray-400 mt-1">Note: {{ item.specialInstructions }}</p>
              </div>
              <p class="font-medium text-gray-900 dark:text-white whitespace-nowrap">{{ formatNaira(item.unitPrice * item.quantity) }}</p>
            </li>
          </ul>
          <dl class="mt-3 pt-3 border-t border-gray-200 dark:border-slate-700 space-y-1 text-sm">
            <div class="flex justify-between text-gray-600 dark:text-gray-400"><dt>Subtotal</dt><dd>{{ formatNaira(subtotal) }}</dd></div>
            <div v-if="order.deliveryFee" class="flex justify-between text-gray-600 dark:text-gray-400"><dt>Delivery fee</dt><dd>{{ formatNaira(order.deliveryFee) }}</dd></div>
            <div v-if="order.taxAmount" class="flex justify-between text-gray-600 dark:text-gray-400"><dt>Tax</dt><dd>{{ formatNaira(order.taxAmount) }}</dd></div>
            <div v-if="order.discountAmount" class="flex justify-between text-green-700 dark:text-green-400"><dt>Discount</dt><dd>−{{ formatNaira(order.discountAmount) }}</dd></div>
            <div class="flex justify-between font-bold text-gray-900 dark:text-white pt-1"><dt>Total</dt><dd>{{ formatNaira(order.totalAmount) }}</dd></div>
          </dl>
        </section>

        <!-- Fulfilment -->
        <section class="no-print bg-white dark:bg-slate-800 rounded-xl p-5 border border-gray-200 dark:border-slate-700">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">{{ order.orderType === 'delivery' ? 'Delivery' : 'Pickup' }}</h2>
          <dl class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-3 text-sm">
            <div v-if="order.orderType === 'delivery' && order.deliveryAddress" class="sm:col-span-2">
              <dt class="text-gray-500 dark:text-gray-400">Address</dt><dd class="text-gray-900 dark:text-white">{{ order.deliveryAddress }}</dd>
            </div>
            <div v-if="order.pickupTime">
              <dt class="text-gray-500 dark:text-gray-400">Pickup time</dt><dd class="text-gray-900 dark:text-white">{{ formatDate(order.pickupTime) }}</dd>
            </div>
            <div v-if="order.estimatedReadyTime && !isFinished">
              <dt class="text-gray-500 dark:text-gray-400">Expected ready</dt><dd class="text-gray-900 dark:text-white">{{ formatDate(order.estimatedReadyTime) }}</dd>
            </div>
            <div>
              <dt class="text-gray-500 dark:text-gray-400">Contact</dt>
              <dd class="text-gray-900 dark:text-white">{{ order.customerName }} · {{ order.customerPhone }}<span v-if="order.customerEmail"> · {{ order.customerEmail }}</span></dd>
            </div>
            <div v-if="order.notes" class="sm:col-span-2">
              <dt class="text-gray-500 dark:text-gray-400">Notes</dt><dd class="text-gray-900 dark:text-white">{{ order.notes }}</dd>
            </div>
          </dl>
        </section>

        <!-- Payment -->
        <section class="no-print bg-white dark:bg-slate-800 rounded-xl p-5 border border-gray-200 dark:border-slate-700">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">Payment</h2>
          <div class="flex flex-wrap gap-2 text-xs font-medium">
            <span v-if="order.paymentMethod" class="px-2 py-1 rounded-full bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300">{{ ORDER_PAYMENT_METHOD_LABELS[order.paymentMethod] }}</span>
            <span v-if="order.paymentStatus" class="px-2 py-1 rounded-full" :class="paymentStatusClass(order.paymentStatus)">{{ PAYMENT_STATUS_LABELS[order.paymentStatus] ?? order.paymentStatus }}</span>
          </div>

          <dl class="mt-3 grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-3 text-sm">
            <div v-if="order.paymentPaidAt">
              <dt class="text-gray-500 dark:text-gray-400">Paid</dt>
              <dd class="text-gray-900 dark:text-white">{{ formatDate(order.paymentPaidAt) }}<span v-if="order.paymentCollectedVia"> · {{ COLLECTED_VIA[order.paymentCollectedVia] }}</span></dd>
            </div>
            <div v-if="order.paymentReference && isOnline">
              <dt class="text-gray-500 dark:text-gray-400">Payment reference</dt><dd class="text-gray-900 dark:text-white font-mono break-all">{{ order.paymentReference }}</dd>
            </div>
          </dl>

          <div v-if="awaitingTransfer" class="mt-4 rounded-lg bg-yellow-50 dark:bg-yellow-900/20 p-4 text-sm text-gray-800 dark:text-gray-200">
            Transfer <strong>{{ formatNaira(order.totalAmount) }}</strong> to
            <strong class="font-mono">{{ order.paymentAccountNumber }}</strong>
            <span v-if="order.paymentBankName"> ({{ order.paymentBankName }})</span><span v-if="order.paymentAccountName">, {{ order.paymentAccountName }}</span>.
            <span v-if="order.paymentExpiresAt"> Pay before {{ formatDate(order.paymentExpiresAt) }}.</span>
          </div>

          <p v-if="message" class="mt-3 text-sm text-gray-700 dark:text-gray-300" role="status">{{ message }}</p>

          <div v-if="isOnline && order.paymentStatus === 'pending'" class="mt-4 flex flex-wrap gap-2">
            <button v-if="order.paymentMethod === 'card' && cardPaymentsEnabled" class="btn-primary text-sm dark:text-white" :disabled="busy" @click="payNow">
              Pay with card
            </button>
            <button class="btn-secondary text-sm inline-flex items-center gap-2" :disabled="busy" @click="refresh">
              <Loader2 v-if="busy" class="w-4 h-4 animate-spin" /><RefreshCw v-else class="w-4 h-4" />
              Refresh payment status
            </button>
          </div>
        </section>

        <!-- Invoice (pay-on-delivery orders) -->
        <section v-if="invoice">
          <div class="no-print flex flex-wrap gap-3">
            <button class="btn-secondary text-sm inline-flex items-center gap-2" :aria-expanded="showInvoice" @click="showInvoice = !showInvoice">
              <FileText class="w-4 h-4" /> {{ showInvoice ? 'Hide invoice' : 'View invoice' }}
            </button>
            <button v-if="showInvoice" class="btn-primary text-sm dark:text-white inline-flex items-center gap-2" @click="printInvoice">
              <Printer class="w-4 h-4" /> Print / Save as PDF
            </button>
          </div>
          <div v-if="showInvoice" class="mt-4">
            <InvoiceView :invoice="invoice" />
          </div>
        </section>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, FileText, Loader2, Printer, RefreshCw } from 'lucide-vue-next'
import PageHeader from '../../../components/layout/PageHeader.vue'
import InvoiceView from '../../../components/order/InvoiceView.vue'
import { ApiRequestError, useApi } from '../../../utils/api'
import { useAuthStore } from '../../../stores/auth'
import { usePaystack } from '../../../composables/usePaystack'
import { ORDER_STATUS_LABELS, ORDER_PAYMENT_METHOD_LABELS, PAYMENT_STATUS_LABELS, paymentStatusClass } from '../../../utils/paymentLabels'
import type { Invoice, Order, OrderStatus, PaymentCollectedVia } from '../../../types/order'

const STEPS: OrderStatus[] = ['pending', 'confirmed', 'preparing', 'ready', 'completed']
const COLLECTED_VIA: Record<PaymentCollectedVia, string> = { cash: 'Cash', pos: 'POS', transfer: 'Transfer' }

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const api = useApi()
const { cardPaymentsEnabled, payWithCard } = usePaystack()

const orderId = computed(() => String(route.params.id ?? ''))
const order = ref<Order | null>(null)
const invoice = ref<Invoice | null>(null)
const loading = ref(true)
const loadError = ref('')
const notFound = ref(false)
const busy = ref(false)
const message = ref('')
const showInvoice = ref(false)

const formatNaira = (amount: number) =>
  `₦${amount.toLocaleString('en-NG', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
const formatDate = (value: string) => new Date(value).toLocaleString('en-NG', { dateStyle: 'medium', timeStyle: 'short' })

const isOnline = computed(() => order.value?.paymentMethod === 'transfer' || order.value?.paymentMethod === 'card')
const isFinished = computed(() => order.value?.status === 'completed' || order.value?.status === 'cancelled')
const currentStep = computed(() => (order.value ? STEPS.indexOf(order.value.status) : -1))
const awaitingTransfer = computed(() =>
  order.value?.paymentMethod === 'transfer' && order.value.paymentStatus === 'pending' && Boolean(order.value.paymentAccountNumber),
)
// Older orders have no stored breakdown; their items add up to the subtotal.
const subtotal = computed(() =>
  order.value?.subtotalAmount ?? (order.value?.items ?? []).reduce((sum, item) => sum + item.unitPrice * item.quantity, 0),
)

const sendToSignIn = () => {
  auth.signOut()
  router.replace({ path: '/account/login', query: { next: route.fullPath } })
}

// An expired session sends the customer to sign in again.
const handleAuthError = (err: unknown) => {
  if (err instanceof ApiRequestError && err.status === 401) {
    sendToSignIn()
    return true
  }
  return false
}

const load = async () => {
  loading.value = true
  loadError.value = ''
  notFound.value = false
  try {
    const result = await api.getMyOrder(orderId.value)
    order.value = result.order
    invoice.value = result.invoice
  } catch (err) {
    if (handleAuthError(err)) return
    if (err instanceof ApiRequestError && (err.status === 404 || err.status === 400)) {
      notFound.value = true
      loadError.value = 'We couldn’t find that order in your account.'
    } else {
      loadError.value = 'We couldn’t load this order. Please try again.'
    }
  } finally {
    loading.value = false
  }
}

// Asks the server to check the payment with Paystack.
const refresh = async () => {
  if (!order.value) return
  busy.value = true
  try {
    const result = await api.refreshMyPaymentStatus(order.value.id)
    order.value = result.order
    message.value = result.message || 'Status updated.'
  } catch (err) {
    if (!handleAuthError(err)) message.value = err instanceof Error ? err.message : 'Please try again.'
  } finally {
    busy.value = false
  }
}

const payNow = async () => {
  if (!order.value) return
  busy.value = true
  message.value = ''
  try {
    const renewed = await api.getPaymentDetails(order.value.orderNumber, order.value.customerPhone)
    if (!renewed.paymentAccessCode) throw new Error('We couldn’t start the card payment. Please try again.')
    const outcome = await payWithCard(renewed.paymentAccessCode)
    busy.value = false
    if (outcome === 'paid') {
      await refresh()
    } else {
      message.value = 'Payment cancelled. You can pay whenever you’re ready.'
    }
  } catch (err) {
    message.value = err instanceof Error ? err.message : 'Please try again.'
  } finally {
    busy.value = false
  }
}

const printInvoice = () => window.print()

onMounted(() => {
  if (!auth.isSignedIn) {
    router.replace({ path: '/account/login', query: { next: route.fullPath } })
    return
  }
  load()
})
</script>

<style>
/* Printing shows only the invoice: the site chrome and order panels are hidden. */
@media print {
  header,
  footer,
  .no-print {
    display: none !important;
  }

  main {
    padding: 0 !important;
  }

  .invoice {
    border: none !important;
  }
}
</style>
