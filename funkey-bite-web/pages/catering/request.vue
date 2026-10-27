<template>
  <div class="mt-[-30px]">
    <PageHeader
      title-before="Your"
      highlight-text="Catering"
      title-after="Request"
      subtitle="Check the status of your catering request"
      :narrow="true"
      :alignment="'center'"
      variant="gradient"
    />

    <div class="container-narrow px-4 sm:px-6 lg:px-8 py-8 max-w-2xl">
      <form
        v-if="!request"
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
          Check status
        </button>
      </form>

      <article v-else class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 space-y-4">
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

        <div class="p-4 rounded-lg bg-gray-50 dark:bg-slate-700/50 text-sm">
          <p v-if="request.quotedAmount != null" class="text-gray-900 dark:text-white">
            Your quote: <strong class="text-lg">{{ naira(request.quotedAmount) }}</strong>
          </p>
          <p v-else-if="request.estimatedTotal != null" class="text-gray-700 dark:text-gray-300">
            Estimate: {{ naira(request.estimatedTotal) }}. We'll confirm the final price in your quote.
          </p>
          <p v-else class="text-gray-700 dark:text-gray-300">We're preparing your quote.</p>
          <p class="mt-2 text-gray-600 dark:text-gray-400">
            Questions about your request?
            <a v-if="phoneHref" :href="phoneHref" class="text-brand-500 font-medium">Call {{ settings?.phoneNumber }}</a>
            <template v-else>Contact us.</template>
          </p>
        </div>

        <button class="btn-secondary" @click="reset">Check another request</button>
      </article>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Loader2 } from 'lucide-vue-next'
import PageHeader from '../../components/layout/PageHeader.vue'
import type { CateringRequest } from '../../types/catering'
import { CATERING_STATUS_LABELS, EVENT_TYPES } from '../../types/catering'
import { useApi } from '../../utils/api'
import { useBusinessSettings } from '../../composables/useBusinessSettings'

const { settings, phoneHref, load: loadBusinessSettings } = useBusinessSettings()

const reference = ref('')
const phone = ref('')
const request = ref<CateringRequest | null>(null)
const error = ref('')
const isLoading = ref(false)

const naira = (amount: number) => `₦${amount.toLocaleString('en-NG')}`

const eventTypeName = computed(() => {
  if (!request.value) return ''
  if (request.value.eventType === 'other' && request.value.eventTypeOther) return request.value.eventTypeOther
  return EVENT_TYPES.find((type) => type.id === request.value?.eventType)?.name ?? request.value.eventType
})

const eventDate = computed(() => {
  if (!request.value) return ''
  const date = new Date(`${request.value.eventDate}T${request.value.eventTime ?? '00:00'}`)
  const day = date.toLocaleDateString('en-NG', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
  return request.value.eventTime ? `${day}, ${request.value.eventTime}` : day
})

const statusClass = computed(() => {
  switch (request.value?.status) {
    case 'confirmed':
    case 'completed':
      return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300'
    case 'quoted':
      return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300'
    case 'declined':
    case 'cancelled':
      return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300'
    default:
      return 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-300'
  }
})

const lookup = async () => {
  error.value = ''
  isLoading.value = true
  try {
    const found = await useApi().lookupCateringRequest(reference.value.trim().toUpperCase(), phone.value.trim())
    if (found) {
      request.value = found
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
  request.value = null
  reference.value = ''
  phone.value = ''
}

onMounted(() => loadBusinessSettings())
</script>
