<template>
  <div class="mt-[-30px]">
    <PageHeader
      title-before="Your"
      highlight-text="Catering"
      title-after="Quote"
      subtitle="Review the quote and terms for your event"
      :narrow="true"
      :alignment="'center'"
      variant="gradient"
    />

    <div class="container-narrow px-4 sm:px-6 lg:px-8 py-8 max-w-2xl">
      <div v-if="isLoading" class="flex justify-center py-12 text-gray-500">
        <Loader2 class="w-6 h-6 animate-spin" />
      </div>

      <div v-else-if="error" class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 space-y-4">
        <p class="text-gray-800 dark:text-gray-200" role="alert">{{ error }}</p>
        <NuxtLink to="/catering/request" class="btn-primary dark:text-white inline-block">Find my booking</NuxtLink>
      </div>

      <CateringBookingView v-else-if="booking && token" :booking="booking" :auth="{ token }" @updated="booking = $event" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Loader2 } from 'lucide-vue-next'
import PageHeader from '../../components/layout/PageHeader.vue'
import CateringBookingView from '../../components/catering/CateringBookingView.vue'
import type { CateringBooking } from '../../types/catering'
import { useApi } from '../../utils/api'

// The emailed quote link. The token only identifies the booking; the page
// is read on the client so it never ends up in a pre-rendered page.
const route = useRoute()
const token = ref('')
const booking = ref<CateringBooking | null>(null)
const error = ref('')
const isLoading = ref(true)

onMounted(async () => {
  token.value = typeof route.query.token === 'string' ? route.query.token.trim() : ''
  if (!token.value) {
    error.value = 'This link is missing its code. Look up your booking with its reference and phone number instead.'
    isLoading.value = false
    return
  }
  try {
    booking.value = await useApi().getCateringBooking({ token: token.value })
    if (!booking.value) error.value = 'We could not find this booking.'
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Something went wrong. Please try again.'
  } finally {
    isLoading.value = false
  }
})
</script>
