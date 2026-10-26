<template>
  <article class="invoice bg-white text-gray-900 rounded-xl border border-gray-200 p-6 sm:p-8">
    <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 pb-6 border-b border-gray-200">
      <div>
        <h2 class="text-2xl font-bold">Invoice</h2>
        <p class="text-sm text-gray-600 mt-1">{{ invoice.invoiceNumber }}</p>
        <p class="text-sm text-gray-600">Order {{ invoice.orderNumber }} · {{ issuedAt }}</p>
      </div>
      <span
        class="self-start px-3 py-1 rounded-full text-xs font-semibold"
        :class="isPaid ? 'bg-green-100 text-green-800' : 'bg-amber-100 text-amber-800'"
      >
        {{ isPaid ? 'Paid' : 'Payment due on ' + (invoice.orderType === 'delivery' ? 'delivery' : 'pickup') }}
      </span>
    </div>

    <section class="grid grid-cols-1 sm:grid-cols-2 gap-6 py-6 border-b border-gray-200 text-sm">
      <div>
        <h3 class="font-semibold text-gray-500 uppercase text-xs mb-2">From</h3>
        <p class="font-medium">{{ invoice.business.name }}</p>
        <p v-if="invoice.business.address">{{ invoice.business.address }}</p>
        <p>{{ invoice.business.phone }}</p>
        <p v-if="invoice.business.email">{{ invoice.business.email }}</p>
      </div>
      <div>
        <h3 class="font-semibold text-gray-500 uppercase text-xs mb-2">Bill to</h3>
        <p class="font-medium">{{ invoice.customer.name }}</p>
        <p>{{ invoice.customer.phone }}</p>
        <p v-if="invoice.customer.email">{{ invoice.customer.email }}</p>
        <p v-if="invoice.customer.address" class="mt-1">Deliver to: {{ invoice.customer.address }}</p>
      </div>
    </section>

    <div class="overflow-x-auto">
      <table class="w-full text-sm mt-4">
        <thead>
          <tr class="text-left text-gray-500 border-b border-gray-200">
            <th class="py-2 font-medium">Item</th>
            <th class="py-2 font-medium text-right">Qty</th>
            <th class="py-2 font-medium text-right">Price</th>
            <th class="py-2 font-medium text-right">Amount</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(line, index) in invoice.items" :key="index" class="border-b border-gray-100">
            <td class="py-2 pr-2">{{ line.name }}</td>
            <td class="py-2 text-right">{{ line.quantity }}</td>
            <td class="py-2 text-right whitespace-nowrap">{{ naira(line.unitPrice) }}</td>
            <td class="py-2 text-right whitespace-nowrap">{{ naira(line.total) }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <dl class="mt-4 ml-auto max-w-xs space-y-1 text-sm">
      <div class="flex justify-between"><dt>Subtotal</dt><dd>{{ naira(invoice.subtotal) }}</dd></div>
      <div v-if="invoice.deliveryFee" class="flex justify-between"><dt>Delivery fee</dt><dd>{{ naira(invoice.deliveryFee) }}</dd></div>
      <div v-if="invoice.tax" class="flex justify-between"><dt>Tax</dt><dd>{{ naira(invoice.tax) }}</dd></div>
      <div v-if="invoice.discount" class="flex justify-between"><dt>Discount</dt><dd>-{{ naira(invoice.discount) }}</dd></div>
      <div class="flex justify-between pt-2 border-t border-gray-200 text-base font-bold">
        <dt>{{ isPaid ? 'Total paid' : 'Total due' }}</dt><dd>{{ naira(invoice.total) }}</dd>
      </div>
    </dl>

    <p class="mt-6 p-3 rounded-lg text-sm" :class="isPaid ? 'bg-green-50 text-green-800' : 'bg-amber-50 text-amber-800'">
      <template v-if="isPaid">Payment received{{ paidVia }}. Thank you!</template>
      <template v-else>
        Please pay {{ naira(invoice.total) }} {{ invoice.orderType === 'delivery' ? 'when your order is delivered' : 'when you collect your order' }}.
        Cash, POS and bank transfer are accepted.
      </template>
    </p>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Invoice } from '../../types/order'

const props = defineProps<{ invoice: Invoice }>()

const naira = (amount: number) =>
  `₦${amount.toLocaleString('en-NG', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`

const isPaid = computed(() => props.invoice.paymentStatus === 'paid')

const issuedAt = computed(() =>
  new Date(props.invoice.issuedAt).toLocaleString('en-NG', { dateStyle: 'medium', timeStyle: 'short' }),
)

const paidVia = computed(() => {
  switch (props.invoice.paymentCollectedVia) {
    case 'cash': return ' in cash'
    case 'pos': return ' by POS'
    case 'transfer': return ' by bank transfer'
    default: return ''
  }
})
</script>
