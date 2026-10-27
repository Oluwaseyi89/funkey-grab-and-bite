<template>
  <article class="space-y-6 print:space-y-4">
    <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 space-y-4">
      <div class="flex items-start justify-between gap-4">
        <div>
          <p class="text-sm text-gray-500">Reference</p>
          <p class="text-2xl font-mono font-bold text-brand-500">{{ request.reference }}</p>
        </div>
        <span class="px-3 py-1 rounded-full text-sm font-medium" :class="statusClass">{{ CATERING_STATUS_LABELS[request.status] }}</span>
      </div>

      <dl class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-3 text-sm">
        <div><dt class="text-gray-500">Event</dt><dd class="text-gray-900 dark:text-white">{{ request.eventName || eventTypeName }}</dd></div>
        <div><dt class="text-gray-500">Date</dt><dd class="text-gray-900 dark:text-white">{{ eventDate }}</dd></div>
        <div><dt class="text-gray-500">Guests</dt><dd class="text-gray-900 dark:text-white">{{ request.guestCount }}</dd></div>
        <div><dt class="text-gray-500">Venue</dt><dd class="text-gray-900 dark:text-white">{{ [request.venueArea, request.venueState].filter(Boolean).join(', ') }}</dd></div>
        <div><dt class="text-gray-500">Catering</dt><dd class="text-gray-900 dark:text-white">{{ request.packageName || 'Custom request' }}</dd></div>
        <div v-if="request.selectedItems.length"><dt class="text-gray-500">Dishes</dt><dd class="text-gray-900 dark:text-white">{{ request.selectedItems.map((item) => item.name).join(', ') }}</dd></div>
      </dl>

      <p v-if="!quote && request.status === 'pending'" class="p-4 rounded-lg bg-gray-50 dark:bg-slate-700/50 text-sm text-gray-700 dark:text-gray-300">
        <template v-if="request.estimatedTotal != null">Estimate: {{ naira(request.estimatedTotal) }}. </template>
        We're preparing your quote and will email it to you.
      </p>
    </div>

    <!-- Quote: shown to accept, and kept as the agreed price afterwards. -->
    <section v-if="quote" class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 space-y-4">
      <div class="flex items-baseline justify-between gap-4">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ quote.status === 'accepted' ? 'Your agreed quote' : 'Your quote' }}
        </h2>
        <span class="text-xs text-gray-500">Version {{ quote.version }}</span>
      </div>

      <p v-if="quote.message" class="text-sm text-gray-700 dark:text-gray-300 whitespace-pre-line">{{ quote.message }}</p>

      <table class="w-full text-sm">
        <tbody>
          <tr v-for="(line, index) in quote.lineItems" :key="index" class="border-b border-gray-100 dark:border-slate-700">
            <td class="py-2 text-gray-700 dark:text-gray-300">{{ line.description }}</td>
            <td class="py-2 text-right text-gray-900 dark:text-white">{{ naira(line.amount) }}</td>
          </tr>
          <tr>
            <td class="pt-3 font-semibold text-gray-900 dark:text-white">Total</td>
            <td class="pt-3 text-right text-lg font-bold text-gray-900 dark:text-white">{{ naira(quote.amount) }}</td>
          </tr>
          <tr>
            <td class="py-1 text-gray-600 dark:text-gray-400">Deposit to confirm ({{ quote.depositPercent }}%)</td>
            <td class="py-1 text-right text-gray-900 dark:text-white">{{ naira(quote.depositAmount) }}</td>
          </tr>
        </tbody>
      </table>

      <details class="rounded-lg border border-gray-200 dark:border-slate-700" :open="canAccept">
        <summary class="px-4 py-3 cursor-pointer text-sm font-medium text-gray-900 dark:text-white">Terms and conditions</summary>
        <p class="px-4 pb-4 text-sm text-gray-600 dark:text-gray-400 whitespace-pre-line max-h-72 overflow-y-auto">{{ quote.terms }}</p>
      </details>

      <template v-if="canAccept">
        <p class="text-sm text-gray-600 dark:text-gray-400">This quote is valid until {{ formatDay(quote.validUntil) }}.</p>
        <label class="flex items-start gap-3 text-sm text-gray-800 dark:text-gray-200">
          <input v-model="agreed" type="checkbox" class="mt-1 h-4 w-4" />
          <span>I have read and agree to the terms, and I accept this quote of {{ naira(quote.amount) }} for my event.</span>
        </label>
        <p v-if="actionError" class="text-sm text-red-600" role="alert">{{ actionError }}</p>
        <button class="btn-primary dark:text-white" :disabled="!agreed || busy" @click="accept">
          <Loader2 v-if="busy" class="w-4 h-4 animate-spin inline mr-2" />
          Accept quote
        </button>
      </template>
      <p v-else-if="quote.status === 'sent' && quote.isExpired" class="p-3 rounded-lg bg-amber-50 dark:bg-amber-900/20 text-sm text-amber-800 dark:text-amber-200">
        This quote expired on {{ formatDay(quote.validUntil) }}. Please contact us for an updated quote.
      </p>
      <p v-else-if="quote.acceptedAt" class="text-sm text-gray-600 dark:text-gray-400">Accepted {{ formatDateTime(quote.acceptedAt) }}.</p>
    </section>

    <!-- Invoice and payments, once the quote is accepted. -->
    <section v-if="booking.account" class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 space-y-4">
      <div class="flex items-baseline justify-between gap-4">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ booking.receipt ? `Receipt ${booking.receipt.receiptNumber}` : `Invoice ${request.invoiceNumber}` }}
        </h2>
        <button class="text-sm text-brand-500 font-medium print:hidden" @click="print">Print</button>
      </div>

      <dl class="text-sm space-y-1">
        <div class="flex justify-between"><dt class="text-gray-600 dark:text-gray-400">Total</dt><dd class="text-gray-900 dark:text-white">{{ naira(booking.account.total) }}</dd></div>
        <div class="flex justify-between"><dt class="text-gray-600 dark:text-gray-400">Deposit to confirm</dt><dd class="text-gray-900 dark:text-white">{{ naira(booking.account.depositDue) }}</dd></div>
        <div class="flex justify-between"><dt class="text-gray-600 dark:text-gray-400">Paid</dt><dd class="text-gray-900 dark:text-white">{{ naira(booking.account.amountPaid) }}</dd></div>
        <div class="flex justify-between font-semibold"><dt class="text-gray-900 dark:text-white">Balance</dt><dd class="text-gray-900 dark:text-white">{{ naira(booking.account.balance) }}</dd></div>
      </dl>

      <ul v-if="booking.payments.length" class="text-sm divide-y divide-gray-100 dark:divide-slate-700">
        <li v-for="payment in booking.payments" :key="payment.id" class="py-2 flex justify-between">
          <span class="text-gray-700 dark:text-gray-300">{{ formatDay(payment.paidAt) }} · {{ PAYMENT_METHOD_LABELS[payment.method] }}</span>
          <span class="text-gray-900 dark:text-white">{{ naira(payment.amount) }}</span>
        </li>
      </ul>

      <p v-if="booking.receipt" class="p-3 rounded-lg bg-green-50 dark:bg-green-900/20 text-sm text-green-800 dark:text-green-200">
        Paid in full on {{ formatDateTime(booking.receipt.issuedAt) }}. Thank you! We look forward to your event.
      </p>

      <div v-else-if="paymentDue" class="space-y-3 print:hidden">
        <p class="text-sm text-gray-700 dark:text-gray-300">
          <template v-if="booking.account.depositRemaining > 0">
            Pay {{ naira(booking.account.depositRemaining) }} more to confirm your booking.
          </template>
          <template v-else>Your booking is confirmed. Please pay the balance before your event.</template>
          You can pay by transfer online, or pay us in cash or by direct transfer and we'll record it.
        </p>

        <div v-if="charge" class="p-4 rounded-lg border border-gray-200 dark:border-slate-600 bg-gray-50 dark:bg-slate-700/50 text-sm space-y-1">
          <p class="font-medium text-gray-900 dark:text-white">Transfer exactly {{ naira(charge.amount) }} to:</p>
          <p v-if="charge.bankName" class="text-gray-700 dark:text-gray-300">Bank: {{ charge.bankName }}</p>
          <p class="text-gray-700 dark:text-gray-300">Account number: <span class="font-mono font-semibold">{{ charge.accountNumber }}</span></p>
          <p v-if="charge.accountName" class="text-gray-700 dark:text-gray-300">Account name: {{ charge.accountName }}</p>
          <p v-if="charge.expiresAt" class="text-xs text-gray-500">Use this account before {{ formatDateTime(charge.expiresAt) }}.</p>
          <p class="text-xs text-gray-500">Your payment is recorded automatically once it arrives; refresh this page to see it.</p>
        </div>

        <p v-if="actionError" class="text-sm text-red-600" role="alert">{{ actionError }}</p>
        <div class="flex flex-wrap gap-3">
          <button v-if="booking.account.depositRemaining > 0" class="btn-primary dark:text-white" :disabled="busy" @click="pay('deposit')">
            Pay deposit ({{ naira(booking.account.depositRemaining) }})
          </button>
          <button class="btn-secondary" :disabled="busy" @click="pay('balance')">
            Pay full balance ({{ naira(booking.account.balance) }})
          </button>
        </div>
      </div>
    </section>

    <p class="text-sm text-gray-600 dark:text-gray-400 print:hidden">
      Questions about your booking?
      <a v-if="phoneHref" :href="phoneHref" class="text-brand-500 font-medium">Call {{ settings?.phoneNumber }}</a>
      <template v-else>Contact us.</template>
    </p>
  </article>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Loader2 } from 'lucide-vue-next'
import type { CateringBooking, CateringBookingAuth, CateringPaymentCharge } from '../../types/catering'
import { CATERING_STATUS_LABELS, EVENT_TYPES, PAYMENT_METHOD_LABELS } from '../../types/catering'
import { useApi } from '../../utils/api'
import { useBusinessSettings } from '../../composables/useBusinessSettings'

const props = defineProps<{ booking: CateringBooking; auth: CateringBookingAuth }>()
const emit = defineEmits<{ updated: [booking: CateringBooking] }>()

const { settings, phoneHref, load: loadBusinessSettings } = useBusinessSettings()

const agreed = ref(false)
const busy = ref(false)
const actionError = ref('')
const charge = ref<CateringPaymentCharge | null>(props.booking.pendingCharge ?? null)

const request = computed(() => props.booking.request)
const quote = computed(() => props.booking.quote)
const canAccept = computed(
  () => request.value.status === 'quoted' && quote.value?.status === 'sent' && !quote.value.isExpired,
)
const paymentDue = computed(() => request.value.status === 'accepted' || request.value.status === 'confirmed')

const naira = (amount: number) => `₦${amount.toLocaleString('en-NG', { maximumFractionDigits: 2 })}`
const formatDay = (value: string) =>
  new Date(value.length === 10 ? `${value}T00:00` : value).toLocaleDateString('en-NG', { day: 'numeric', month: 'long', year: 'numeric' })
const formatDateTime = (value: string) =>
  new Date(value).toLocaleString('en-NG', { day: 'numeric', month: 'short', year: 'numeric', hour: 'numeric', minute: '2-digit' })

const eventTypeName = computed(() => {
  if (request.value.eventType === 'other' && request.value.eventTypeOther) return request.value.eventTypeOther
  return EVENT_TYPES.find((type) => type.id === request.value.eventType)?.name ?? request.value.eventType
})

const eventDate = computed(() => {
  const date = new Date(`${request.value.eventDate}T${request.value.eventTime ?? '00:00'}`)
  const day = date.toLocaleDateString('en-NG', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
  return request.value.eventTime ? `${day}, ${request.value.eventTime}` : day
})

const statusClass = computed(() => {
  switch (request.value.status) {
    case 'confirmed':
    case 'paid':
    case 'completed':
      return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300'
    case 'quoted':
    case 'accepted':
      return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300'
    case 'declined':
    case 'cancelled':
      return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300'
    default:
      return 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-300'
  }
})

const run = async (action: () => Promise<void>) => {
  actionError.value = ''
  busy.value = true
  try {
    await action()
  } catch (err) {
    actionError.value = err instanceof Error ? err.message : 'Something went wrong. Please try again.'
  } finally {
    busy.value = false
  }
}

const accept = () =>
  run(async () => {
    emit('updated', await useApi().acceptCateringQuote(props.auth))
  })

const pay = (portion: 'deposit' | 'balance') =>
  run(async () => {
    charge.value = await useApi().payCateringOnline(props.auth, portion)
  })

const print = () => window.print()

onMounted(() => loadBusinessSettings())
</script>
