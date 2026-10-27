<template>
  <div class="mt-[-30px]">
    <PageHeader
      title-before="Your"
      highlight-text="Catering"
      title-after="Booking"
      subtitle="Check your request, accept your quote and pay"
      :narrow="true"
      :alignment="'center'"
      variant="gradient"
    />

    <div class="container-narrow px-4 sm:px-6 lg:px-8 py-8 max-w-2xl">
      <form
        v-if="!booking"
        class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 space-y-4"
        @submit.prevent="lookup"
      >
        <p class="text-sm text-gray-600 dark:text-gray-400">
          Enter the reference from your confirmation (e.g. CAT-7K3M9Q) and the phone number you used.
        </p>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label for="cat-ref" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Reference</label>
            <input id="cat-ref" v-model="reference" required autocomplete="off" class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white uppercase" placeholder="CAT-…" />
          </div>
          <div>
            <label for="cat-phone" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Phone number</label>
            <input id="cat-phone" v-model="phone" type="tel" required autocomplete="tel" class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white" placeholder="0801 234 5678" />
          </div>
        </div>
        <p v-if="error" class="text-sm text-red-600" role="alert">{{ error }}</p>
        <button type="submit" class="btn-primary dark:text-white" :disabled="isLoading">
          <Loader2 v-if="isLoading" class="w-4 h-4 animate-spin inline mr-2" />
          Find my booking
        </button>
      </form>

      <template v-else-if="auth">
        <CateringBookingView :booking="booking" :auth="auth" @updated="booking = $event" />
        <button class="btn-secondary mt-6 print:hidden" @click="reset">Check another request</button>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Loader2 } from 'lucide-vue-next'
import PageHeader from '../../components/layout/PageHeader.vue'
import CateringBookingView from '../../components/catering/CateringBookingView.vue'
import type { CateringBooking, CateringBookingAuth } from '../../types/catering'
import { useApi } from '../../utils/api'

const reference = ref('')
const phone = ref('')
const booking = ref<CateringBooking | null>(null)
const auth = ref<CateringBookingAuth | null>(null)
const error = ref('')
const isLoading = ref(false)

const lookup = async () => {
  error.value = ''
  isLoading.value = true
  try {
    const credentials = { reference: reference.value.trim().toUpperCase(), phone: phone.value.trim() }
    const found = await useApi().getCateringBooking(credentials)
    if (found) {
      auth.value = credentials
      booking.value = found
    } else {
      error.value = "We couldn't find a request with that reference and phone number."
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Something went wrong. Please try again.'
  } finally {
    isLoading.value = false
  }
}

const reset = () => {
  booking.value = null
  auth.value = null
  reference.value = ''
  phone.value = ''
}
</script>
