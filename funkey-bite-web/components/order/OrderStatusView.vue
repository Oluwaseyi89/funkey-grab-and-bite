<template>
  <div class="space-y-4">
    <!-- Where the order is now -->
    <section class="card overflow-hidden" aria-labelledby="order-state-title">
      <div class="relative px-5 pb-5 pt-6" :class="hero.band">
        <div class="flex items-start gap-4">
          <span class="flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl" :class="hero.icon">
            <component :is="hero.iconComponent" class="h-7 w-7" aria-hidden="true" />
          </span>
          <div class="min-w-0">
            <h2 id="order-state-title" class="font-display text-xl font-bold leading-tight text-ink-900 dark:text-white">{{ hero.title }}</h2>
            <p class="mt-1 text-sm text-ink-600 dark:text-ink-300">{{ hero.text }}</p>
          </div>
        </div>
        <dl class="mt-5 grid grid-cols-2 gap-3 text-sm">
          <div class="col-span-2 rounded-xl bg-white/70 px-3 py-2 dark:bg-ink-900/60 sm:col-span-1">
            <dt class="text-xs text-ink-500 dark:text-ink-400">Order number</dt>
            <dd class="flex items-center gap-1 font-mono font-semibold text-ink-900 dark:text-white">
              <span class="break-all">{{ order.orderNumber }}</span>
              <button type="button" class="icon-btn -my-2 h-9 w-9 shrink-0" :aria-label="copiedField === 'order' ? 'Order number copied' : 'Copy order number'" @click="copy('order', order.orderNumber)">
                <component :is="copiedField === 'order' ? Check : Copy" class="h-4 w-4" aria-hidden="true" />
              </button>
            </dd>
          </div>
          <div class="col-span-2 rounded-xl bg-white/70 px-3 py-2 dark:bg-ink-900/60 sm:col-span-1">
            <dt class="text-xs text-ink-500 dark:text-ink-400">{{ showReadyTime ? 'Expected ready' : 'Placed' }}</dt>
            <dd class="font-semibold text-ink-900 dark:text-white">
              {{ showReadyTime ? formatTime(order.estimatedReadyTime) : formatDateTime(order.createdAt) }}
            </dd>
          </div>
        </dl>
      </div>

      <ol v-if="order.status !== 'cancelled'" class="grid grid-cols-5 gap-1.5 border-t border-ink-100 px-5 py-4 dark:border-ink-800" aria-label="Order progress">
        <li v-for="(step, index) in STEPS" :key="step" class="min-w-0">
          <div class="h-1.5 rounded-full transition-colors duration-500" :class="index <= currentStep ? 'bg-brand-500' : 'bg-ink-200 dark:bg-ink-700'" />
          <p
            class="mt-1.5 truncate text-[11px] sm:text-xs"
            :class="index === currentStep ? 'font-semibold text-ink-900 dark:text-white' : 'text-ink-500 dark:text-ink-400'"
            :aria-current="index === currentStep ? 'step' : undefined"
          >
            {{ STEP_LABELS[step] }}
          </p>
        </li>
      </ol>
    </section>

    <!-- Paying -->
    <section class="card p-5" aria-labelledby="payment-title">
      <div class="flex flex-wrap items-center justify-between gap-2">
        <h2 id="payment-title" class="font-display text-lg font-semibold text-ink-900 dark:text-white">Payment</h2>
        <span v-if="order.paymentStatus" class="badge" :class="paymentBadge">{{ PAYMENT_STATUS_LABELS[order.paymentStatus] ?? order.paymentStatus }}</span>
      </div>
      <p class="mt-1 text-sm text-ink-500 dark:text-ink-400">
        {{ order.paymentMethod ? ORDER_PAYMENT_METHOD_LABELS[order.paymentMethod] : '' }} · {{ naira(order.totalAmount) }}
      </p>

      <!-- Bank transfer, waiting -->
      <div v-if="awaitingTransfer && order.paymentAccountNumber && !windowClosed" class="mt-4 rounded-2xl border border-brand-100 bg-brand-50/60 p-4 dark:border-brand-900/40 dark:bg-brand-900/15">
        <p class="text-sm text-ink-700 dark:text-ink-200">
          Transfer exactly <strong class="text-ink-900 dark:text-white">{{ naira(order.totalAmount) }}</strong> to this account. Your order confirms by itself once it arrives.
        </p>
        <div class="mt-3 flex items-center justify-between gap-3 rounded-xl bg-white px-4 py-3 dark:bg-ink-900">
          <div>
            <p class="text-xs text-ink-500 dark:text-ink-400">Account number</p>
            <p class="font-mono text-xl font-bold tracking-wider text-ink-900 dark:text-white">{{ order.paymentAccountNumber }}</p>
          </div>
          <button type="button" class="btn-secondary min-h-[40px] px-3" @click="copy('account', order.paymentAccountNumber)">
            <component :is="copiedField === 'account' ? Check : Copy" class="h-4 w-4" aria-hidden="true" />
            {{ copiedField === 'account' ? 'Copied' : 'Copy' }}
          </button>
        </div>
        <dl class="mt-3 grid grid-cols-2 gap-2 text-sm">
          <div v-if="order.paymentBankName"><dt class="text-xs text-ink-500 dark:text-ink-400">Bank</dt><dd class="font-medium text-ink-900 dark:text-white">{{ order.paymentBankName }}</dd></div>
          <div v-if="order.paymentAccountName"><dt class="text-xs text-ink-500 dark:text-ink-400">Account name</dt><dd class="font-medium text-ink-900 dark:text-white">{{ order.paymentAccountName }}</dd></div>
          <div v-if="order.paymentExpiresAt" class="col-span-2"><dt class="text-xs text-ink-500 dark:text-ink-400">Pay before</dt><dd class="font-medium text-ink-900 dark:text-white">{{ formatTime(order.paymentExpiresAt) }} <span class="text-ink-500">({{ minutesLeft }})</span></dd></div>
        </dl>
      </div>

      <!-- Bank transfer, no details or window closed -->
      <div v-else-if="awaitingTransfer" class="mt-4 rounded-2xl bg-accent-50 p-4 text-sm text-accent-800 dark:bg-accent-700/15 dark:text-accent-200">
        <p class="font-semibold">{{ order.paymentAccountNumber ? 'The time to pay with those details has passed' : 'Your payment details aren’t ready yet' }}</p>
        <p class="mt-1">
          {{ order.paymentAccountNumber ? 'If you haven’t sent the transfer yet, get new details and don’t use the old account.' : 'Your order is saved. Get the details to pay.' }}
        </p>
        <button type="button" class="btn-primary mt-3 w-full" :disabled="busy" @click="getDetails">
          <Loader2 v-if="busy === 'details'" class="h-4 w-4 animate-spin" aria-hidden="true" />
          Get payment details
        </button>
      </div>

      <!-- Card, waiting -->
      <div v-else-if="awaitingCard" class="mt-4">
        <p class="text-sm text-ink-600 dark:text-ink-300">Your order is saved. Pay in Paystack’s secure window; your card details never reach us.</p>
        <button v-if="cardPaymentsEnabled" type="button" class="btn-primary btn-lg mt-3 w-full" :disabled="!!busy" @click="payByCard">
          <Loader2 v-if="busy === 'card'" class="h-5 w-5 animate-spin" aria-hidden="true" />
          <CreditCard v-else class="h-5 w-5" aria-hidden="true" />
          Pay {{ naira(order.totalAmount) }} by card
        </button>
      </div>

      <!-- Paid online -->
      <p v-else-if="order.paymentStatus === 'paid'" class="mt-4 flex items-center gap-2 rounded-2xl bg-leaf-50 px-4 py-3 text-sm font-medium text-leaf-800 dark:bg-leaf-800/25 dark:text-leaf-200">
        <BadgeCheck class="h-5 w-5 shrink-0" aria-hidden="true" />
        Paid{{ order.paymentPaidAt ? ` on ${formatDateTime(order.paymentPaidAt)}` : '' }}{{ order.paymentCollectedVia ? ` · ${COLLECTED_VIA[order.paymentCollectedVia]}` : '' }}. Thank you!
      </p>

      <!-- Pay on delivery -->
      <p v-else-if="order.paymentMethod === 'cash' && order.paymentStatus === 'pending'" class="mt-4 rounded-2xl bg-accent-50 px-4 py-3 text-sm text-accent-800 dark:bg-accent-700/15 dark:text-accent-200">
        Pay <strong>{{ naira(order.totalAmount) }}</strong> {{ order.orderType === 'delivery' ? 'when your order arrives' : 'when you collect it' }}: cash, POS or transfer.
      </p>

      <p v-if="message" class="mt-3 text-sm text-ink-700 dark:text-ink-200" role="status">{{ message }}</p>

      <div v-if="awaitingOnline" class="mt-3">
        <button type="button" class="btn-ghost w-full" :disabled="!!busy" @click="checkPayment">
          <Loader2 v-if="busy === 'check'" class="h-4 w-4 animate-spin" aria-hidden="true" />
          <RefreshCw v-else class="h-4 w-4" aria-hidden="true" />
          I’ve paid, check now
        </button>
      </div>
    </section>

    <!-- What was ordered -->
    <section class="card p-5" aria-labelledby="items-title">
      <h2 id="items-title" class="font-display text-lg font-semibold text-ink-900 dark:text-white">Your order</h2>
      <ul class="mt-2 divide-y divide-ink-100 dark:divide-ink-800">
        <li v-for="(item, index) in order.items" :key="`${item.menuItemId}-${index}`" class="flex justify-between gap-4 py-3">
          <div class="min-w-0">
            <p class="text-ink-900 dark:text-white"><span class="font-semibold">{{ item.quantity }}×</span> {{ item.name }}</p>
            <p v-if="item.specialInstructions" class="mt-0.5 text-xs italic text-ink-500 dark:text-ink-400">“{{ item.specialInstructions }}”</p>
          </div>
          <p class="shrink-0 tabular-nums text-ink-900 dark:text-white">{{ naira(item.unitPrice * item.quantity) }}</p>
        </li>
      </ul>
      <dl class="mt-2 space-y-1.5 border-t border-ink-100 pt-3 text-sm dark:border-ink-800">
        <div class="flex justify-between text-ink-600 dark:text-ink-300"><dt>Subtotal</dt><dd class="tabular-nums">{{ naira(subtotal) }}</dd></div>
        <div v-if="order.deliveryFee" class="flex justify-between text-ink-600 dark:text-ink-300"><dt>Delivery</dt><dd class="tabular-nums">{{ naira(order.deliveryFee) }}</dd></div>
        <div v-if="order.taxAmount" class="flex justify-between text-ink-600 dark:text-ink-300"><dt>Tax</dt><dd class="tabular-nums">{{ naira(order.taxAmount) }}</dd></div>
        <div v-if="order.discountAmount" class="flex justify-between text-leaf-700 dark:text-leaf-300"><dt>Discount</dt><dd class="tabular-nums">−{{ naira(order.discountAmount) }}</dd></div>
        <div class="flex justify-between pt-1 text-base font-bold text-ink-900 dark:text-white"><dt>Total</dt><dd class="tabular-nums">{{ naira(order.totalAmount) }}</dd></div>
      </dl>
    </section>

    <!-- Pickup or delivery -->
    <section class="card p-5" aria-labelledby="fulfil-title">
      <h2 id="fulfil-title" class="font-display text-lg font-semibold text-ink-900 dark:text-white">{{ order.orderType === 'delivery' ? 'Delivery' : 'Pickup' }}</h2>
      <dl class="mt-3 space-y-3 text-sm">
        <div v-if="order.orderType === 'delivery' && order.deliveryAddress" class="flex gap-3">
          <MapPin class="mt-0.5 h-4 w-4 shrink-0 text-ink-400" aria-hidden="true" />
          <div><dt class="sr-only">Address</dt><dd class="text-ink-900 dark:text-white">{{ order.deliveryAddress }}</dd></div>
        </div>
        <div v-else-if="order.orderType === 'pickup' && fullAddress" class="flex gap-3">
          <Store class="mt-0.5 h-4 w-4 shrink-0 text-ink-400" aria-hidden="true" />
          <div>
            <dt class="sr-only">Pick up from</dt>
            <dd class="text-ink-900 dark:text-white">{{ fullAddress }}</dd>
            <dd v-if="mapsUrl"><a :href="mapsUrl" target="_blank" rel="noopener noreferrer" class="link text-sm">Get directions</a></dd>
          </div>
        </div>
        <div class="flex gap-3">
          <User class="mt-0.5 h-4 w-4 shrink-0 text-ink-400" aria-hidden="true" />
          <div><dt class="sr-only">Contact</dt><dd class="text-ink-900 dark:text-white">{{ order.customerName }} · {{ order.customerPhone }}</dd></div>
        </div>
        <div v-if="order.notes" class="flex gap-3">
          <StickyNote class="mt-0.5 h-4 w-4 shrink-0 text-ink-400" aria-hidden="true" />
          <div><dt class="sr-only">Notes</dt><dd class="text-ink-900 dark:text-white">{{ order.notes }}</dd></div>
        </div>
      </dl>
    </section>

    <!-- Invoice, for pay-on-delivery orders -->
    <section v-if="invoice" class="card p-5">
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 class="font-display text-lg font-semibold text-ink-900 dark:text-white">Invoice</h2>
          <p class="text-sm text-ink-500 dark:text-ink-400">{{ invoice.invoiceNumber }}</p>
        </div>
        <div class="flex gap-2">
          <button type="button" class="btn-secondary" :aria-expanded="showInvoice" @click="showInvoice = !showInvoice">
            <FileText class="h-4 w-4" aria-hidden="true" /> {{ showInvoice ? 'Hide' : 'View' }}
          </button>
          <button v-if="showInvoice" type="button" class="btn-primary" @click="printInvoice">
            <Printer class="h-4 w-4" aria-hidden="true" /> Print
          </button>
        </div>
      </div>
      <div v-if="showInvoice" class="print-area mt-4">
        <InvoiceView :invoice="invoice" />
      </div>
    </section>

    <p v-if="phoneHref" class="py-2 text-center text-sm text-ink-500 dark:text-ink-400">
      Need help with this order? <a :href="phoneHref" class="link">Call us</a>
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import {
  BadgeCheck, Ban, Check, ChefHat, CircleCheckBig, Clock, Copy, CreditCard, FileText, Loader2, MapPin, PackageCheck,
  Printer, RefreshCw, StickyNote, Store, User, Wallet,
} from 'lucide-vue-next'
import InvoiceView from './InvoiceView.vue'
import { useApi } from '../../utils/api'
import { usePaystack } from '../../composables/usePaystack'
import { useBusinessSettings } from '../../composables/useBusinessSettings'
import { ORDER_PAYMENT_METHOD_LABELS, PAYMENT_STATUS_LABELS } from '../../utils/paymentLabels'
import { formatDateTime, formatTime, naira } from '../../utils/format'
import type { Invoice, Order, OrderStatus, PaymentCollectedVia } from '../../types/order'
import type { PaymentRefresh } from '../../types/account'

// Everything about one order: where it is, how to pay, what was in it.
// Shared by the guest tracking page and the account's order page; each
// passes how to reload the order and how to ask the server to check payment.
const props = defineProps<{
  order: Order
  invoice: Invoice | null
  reload: () => Promise<{ order: Order; invoice: Invoice | null }>
  checkPaymentStatus: () => Promise<{ order: Order } & PaymentRefresh>
}>()

const emit = defineEmits<{ update: [result: { order: Order; invoice?: Invoice | null }] }>()

const STEPS: OrderStatus[] = ['pending', 'confirmed', 'preparing', 'ready', 'completed']
const STEP_LABELS: Record<OrderStatus, string> = {
  pending: 'Received', confirmed: 'Confirmed', preparing: 'Cooking', ready: 'Ready', completed: 'Done', cancelled: 'Cancelled',
}
const COLLECTED_VIA: Record<PaymentCollectedVia, string> = { cash: 'Cash', pos: 'POS', transfer: 'Transfer' }

const api = useApi()
const { cardPaymentsEnabled, payWithCard } = usePaystack()
const { fullAddress, mapsUrl, phoneHref, load } = useBusinessSettings()

const busy = ref<'' | 'details' | 'card' | 'check'>('')
const message = ref('')
const showInvoice = ref(false)
const copiedField = ref('')

const isOnline = computed(() => props.order.paymentMethod === 'transfer' || props.order.paymentMethod === 'card')
const awaitingOnline = computed(() => isOnline.value && props.order.paymentStatus === 'pending' && props.order.status !== 'cancelled')
const awaitingTransfer = computed(() => awaitingOnline.value && props.order.paymentMethod === 'transfer')
const awaitingCard = computed(() => awaitingOnline.value && props.order.paymentMethod === 'card')
const currentStep = computed(() => STEPS.indexOf(props.order.status))
const finished = computed(() => props.order.status === 'completed' || props.order.status === 'cancelled')
const showReadyTime = computed(() =>
  Boolean(props.order.estimatedReadyTime) && !finished.value && !awaitingOnline.value && props.order.status !== 'ready',
)

const subtotal = computed(() =>
  props.order.subtotalAmount ?? props.order.items.reduce((sum, item) => sum + item.unitPrice * item.quantity, 0),
)

// The headline at the top, in plain words.
const hero = computed(() => {
  const o = props.order
  if (o.status === 'cancelled') {
    return {
      title: 'Order cancelled',
      text: o.paymentStatus === 'expired' ? 'The payment wasn’t made in time, so we cancelled it.' : 'This order was cancelled. Call us if that’s unexpected.',
      iconComponent: Ban, icon: 'bg-ink-200 text-ink-700 dark:bg-ink-700 dark:text-ink-200', band: 'bg-ink-100/70 dark:bg-ink-800/50',
    }
  }
  if (awaitingOnline.value) {
    return {
      title: 'Waiting for your payment',
      text: 'We’ll start cooking as soon as your payment comes through.',
      iconComponent: Wallet, icon: 'bg-accent-100 text-accent-700 dark:bg-accent-700/30 dark:text-accent-200', band: 'bg-accent-50/80 dark:bg-accent-700/10',
    }
  }
  const map: Record<string, { title: string; text: string; iconComponent: unknown }> = {
    pending: { title: 'Order received', text: 'We’ve got it and will confirm it shortly.', iconComponent: Clock },
    confirmed: { title: 'Order confirmed', text: 'The kitchen has your order.', iconComponent: CircleCheckBig },
    preparing: { title: 'Cooking now', text: 'Your food is being prepared.', iconComponent: ChefHat },
    ready: {
      title: o.orderType === 'delivery' ? 'Ready, heading your way' : 'Ready for pickup',
      text: o.orderType === 'delivery' ? 'Your order is ready and will be with you soon.' : 'Come and collect it whenever you’re ready.',
      iconComponent: PackageCheck,
    },
    completed: { title: 'Enjoy your meal!', text: 'This order is complete. Thanks for ordering with us.', iconComponent: BadgeCheck },
  }
  return {
    ...(map[o.status] ?? map.pending)!,
    icon: 'bg-leaf-100 text-leaf-700 dark:bg-leaf-800/40 dark:text-leaf-200',
    band: 'bg-leaf-50/80 dark:bg-leaf-800/10',
  }
})

const paymentBadge = computed(() => {
  switch (props.order.paymentStatus) {
    case 'paid': return 'bg-leaf-100 text-leaf-800 dark:bg-leaf-800/40 dark:text-leaf-200'
    case 'pending': return 'bg-accent-100 text-accent-800 dark:bg-accent-700/30 dark:text-accent-200'
    case 'failed':
    case 'needs_review':
    case 'refund_due': return 'bg-brand-100 text-brand-800 dark:bg-brand-900/40 dark:text-brand-200'
    default: return 'bg-ink-100 text-ink-700 dark:bg-ink-800 dark:text-ink-200'
  }
})

// A clock for the transfer window and the regular status check.
const now = ref(Date.now())
const windowClosed = computed(() => Boolean(props.order.paymentExpiresAt) && new Date(props.order.paymentExpiresAt!).getTime() <= now.value)
const minutesLeft = computed(() => {
  if (!props.order.paymentExpiresAt) return ''
  const minutes = Math.max(0, Math.round((new Date(props.order.paymentExpiresAt).getTime() - now.value) / 60_000))
  return minutes <= 1 ? 'about a minute left' : `${minutes} minutes left`
})

const update = (order: Order, invoice?: Invoice | null) => emit('update', { order, invoice })

const copy = async (field: string, text?: string) => {
  if (!text) return
  try {
    await navigator.clipboard.writeText(text)
    copiedField.value = field
    setTimeout(() => { if (copiedField.value === field) copiedField.value = '' }, 2000)
  } catch {
    // Clipboard blocked: the text is on screen.
  }
}

const run = async (kind: typeof busy.value, task: () => Promise<void>) => {
  busy.value = kind
  message.value = ''
  try {
    await task()
  } catch (err) {
    message.value = err instanceof Error ? err.message : 'Something went wrong. Please try again.'
  } finally {
    busy.value = ''
  }
}

const checkPayment = () => run('check', async () => {
  const result = await props.checkPaymentStatus()
  update(result.order)
  message.value = result.message || 'Status updated.'
})

const getDetails = () => run('details', async () => {
  const renewed = await api.getPaymentDetails(props.order.orderNumber, props.order.customerPhone)
  update({ ...props.order, ...renewed, items: props.order.items })
  now.value = Date.now()
})

const payByCard = () => run('card', async () => {
  const started = await api.getPaymentDetails(props.order.orderNumber, props.order.customerPhone)
  if (!started.paymentAccessCode) throw new Error('We couldn’t start the card payment. Please try again.')
  if ((await payWithCard(started.paymentAccessCode)) === 'paid') {
    const result = await props.checkPaymentStatus()
    update(result.order)
    message.value = result.message || 'Payment received. Thank you!'
  } else {
    message.value = 'Payment cancelled. You can pay whenever you’re ready.'
  }
})

const printInvoice = () => window.print()

// Keeps the page current without the customer pressing anything: every
// 20 seconds while the order is moving, quietly.
let timer: ReturnType<typeof setInterval> | undefined
onMounted(() => {
  load()
  timer = setInterval(async () => {
    now.value = Date.now()
    if (finished.value || busy.value || document.visibilityState !== 'visible') return
    try {
      const latest = await props.reload()
      update(latest.order, latest.invoice)
    } catch {
      // A missed refresh is fine; the next one will try again.
    }
  }, 20_000)
})
onBeforeUnmount(() => clearInterval(timer))
</script>

<style>
/* Printing shows only the invoice. */
@media print {
  body * {
    visibility: hidden;
  }
  .print-area,
  .print-area * {
    visibility: visible;
  }
  .print-area {
    position: absolute;
    inset: 0;
  }
}
</style>
