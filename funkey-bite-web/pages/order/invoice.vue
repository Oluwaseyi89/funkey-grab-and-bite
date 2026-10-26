<template>
  <div class="mt-[-30px]">
    <PageHeader
      class="no-print"
      title-before="Your"
      highlight-text="Invoice"
      title-after=""
      subtitle="View or print the invoice for a pay-on-delivery order"
      :narrow="true"
      :alignment="'center'"
      variant="gradient"
    />

    <div class="container-narrow px-4 sm:px-6 lg:px-8 py-8 max-w-3xl">
      <form
        v-if="!result"
        class="no-print bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 space-y-4"
        @submit.prevent="lookup"
      >
        <p class="text-sm text-gray-600 dark:text-gray-400">
          Enter your order number and the phone number you ordered with.
        </p>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label for="invoice-order-number" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Order Number</label>
            <input
              id="invoice-order-number"
              v-model="orderNumber"
              type="text"
              required
              autocomplete="off"
              class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white"
              placeholder="FG-2026-10-…"
            />
          </div>
          <div>
            <label for="invoice-phone" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Phone Number</label>
            <input
              id="invoice-phone"
              v-model="phone"
              type="tel"
              required
              autocomplete="tel"
              class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white"
              placeholder="0801 234 5678"
            />
          </div>
        </div>
        <p v-if="error" class="text-sm text-red-600" role="alert">{{ error }}</p>
        <button type="submit" class="btn-primary w-full sm:w-auto" :disabled="isLoading">
          <Loader2 v-if="isLoading" class="w-4 h-4 animate-spin inline mr-2" />
          Show Invoice
        </button>
      </form>

      <div v-else>
        <InvoiceView v-if="result.invoice" :invoice="result.invoice" />
        <p v-else class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 text-gray-700 dark:text-gray-300">
          Order {{ result.order.orderNumber }} was paid online, so it has no pay-on-delivery invoice.
        </p>

        <div class="no-print flex flex-col sm:flex-row gap-3 mt-6">
          <button v-if="result.invoice" type="button" class="btn-primary" @click="printInvoice">
            <Printer class="w-4 h-4 inline mr-2" />
            Print / Save as PDF
          </button>
          <button type="button" class="btn-secondary" @click="reset">Look up another order</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useState } from 'nuxt/app'
import { Loader2, Printer } from 'lucide-vue-next'
import PageHeader from '../../components/layout/PageHeader.vue'
import InvoiceView from '../../components/order/InvoiceView.vue'
import { useApi } from '../../utils/api'
import type { InvoiceLookup, OrderLookupResult } from '../../types/order'

// Set by the order page right after checkout, so the customer isn't asked to
// type their details again. Kept in app state, never in the URL.
const pendingLookup = useState<InvoiceLookup | null>('invoice-lookup', () => null)

const orderNumber = ref('')
const phone = ref('')
const result = ref<OrderLookupResult | null>(null)
const error = ref('')
const isLoading = ref(false)

const lookup = async () => {
  error.value = ''
  isLoading.value = true
  try {
    const found = await useApi().lookupOrder(orderNumber.value.trim(), phone.value.trim())
    if (found) {
      result.value = found
    } else {
      error.value = "We couldn't find an order with that number and phone. Please check both and try again."
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Something went wrong. Please try again.'
  } finally {
    isLoading.value = false
  }
}

const reset = () => {
  result.value = null
  orderNumber.value = ''
  phone.value = ''
}

const printInvoice = () => window.print()

onMounted(() => {
  if (pendingLookup.value) {
    orderNumber.value = pendingLookup.value.orderNumber
    phone.value = pendingLookup.value.phone
    pendingLookup.value = null
    lookup()
  }
})
</script>

<style>
/* Print only the invoice: hide the site chrome and the page controls. */
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
