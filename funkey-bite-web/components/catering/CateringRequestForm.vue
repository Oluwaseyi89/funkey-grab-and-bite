<template>
    <form @submit.prevent="submitForm" class="space-y-6">
      <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Request Catering Quote</h3>
      
      <!-- Contact Information -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="catering-name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Contact Name *
          </label>
          <input
            v-model="formData.contactName"
            id="catering-name"
            type="text"
            autocomplete="name"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
            placeholder="John Doe"
          />
        </div>
        <div>
          <label for="catering-phone" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Phone Number *
          </label>
          <input
            v-model="formData.contactPhone"
            id="catering-phone"
            type="tel"
            autocomplete="tel"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
            placeholder="+1 (555) 123-4567"
          />
        </div>
      </div>
      
      <div>
        <label for="catering-email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Email Address *
        </label>
        <input
          v-model="formData.contactEmail"
          id="catering-email"
          type="email"
          autocomplete="email"
          required
          class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
          placeholder="john@example.com"
        />
      </div>
      
      <!-- Event Details -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="catering-event-date" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Event Date *
          </label>
          <input
            v-model="formData.eventDate"
            id="catering-event-date"
            type="date"
            required
            :min="minDate"
            class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white"
          />
        </div>
        <div>
          <label for="catering-event-time" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Event Time *
          </label>
          <input
            v-model="formData.eventTime"
            id="catering-event-time"
            type="time"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white"
          />
        </div>
      </div>
      
      <div>
        <label for="catering-event-name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Event Name / Occasion
        </label>
        <input
          v-model="formData.eventName"
          id="catering-event-name"
          type="text"
          class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
          placeholder="e.g., Company Annual Meeting, Sarah's Wedding"
        />
      </div>
      
      <!-- Budget -->
      <div>
        <label for="catering-budget" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Estimated Budget
        </label>
        <div class="relative">
          <span class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-500 dark:text-gray-400">$</span>
          <input
            v-model="formData.budget"
            id="catering-budget"
            type="number"
            min="0"
            step="100"
            class="w-full pl-10 pr-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
            placeholder="Enter your budget (optional)"
          />
        </div>
      </div>
      
      <!-- Special Requests -->
      <div>
        <label for="catering-special-requests" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Special Requests / Notes
        </label>
        <textarea
          v-model="formData.specialRequests"
          id="catering-special-requests"
          rows="4"
          class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
          placeholder="Dietary restrictions, venue details, theme, etc."
        ></textarea>
      </div>
      
      <!-- Submit Button -->
      <button
        type="submit"
        :disabled="props.isSubmitting"
        class="w-full btn-primary py-4 text-lg font-bold disabled:opacity-50 disabled:cursor-not-allowed dark:text-white"
      >
        <template v-if="props.isSubmitting">
          <Loader2 class="w-5 h-5 animate-spin inline mr-2" />
          Sending Request...
        </template>
        <template v-else>
          Request Quote
        </template>
      </button>
      
      <!-- Privacy Note -->
      <p class="text-sm text-gray-500 dark:text-gray-400 text-center">
        We'll contact you within 24 hours to discuss your event. No commitment required.
      </p>
    </form>
  </template>
  
  <script setup lang="ts">
  import type { CateringRequest } from '../../types/order'
  import { Loader2 } from 'lucide-vue-next'
  import { ref, computed } from 'vue'
  
  const props = withDefaults(defineProps<{
    isSubmitting?: boolean
    guestCount?: number
    selectedPackage?: string
    selectedEvent?: string
  }>(), {
    isSubmitting: false
  })
  
  const emit = defineEmits<{
    submit: [data: Partial<CateringRequest>]
  }>()
  
  const formData = ref({
    contactName: '',
    contactPhone: '',
    contactEmail: '',
    eventName: '',
    eventDate: '',
    eventTime: '',
    budget: undefined as number | undefined,
    specialRequests: ''
  })
  
  const minDate = computed(() => {
    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 1)
    return tomorrow.toISOString().split('T')[0]
  })
  
  const submitForm = () => {
    const requestData = {
      ...formData.value,
      guestCount: props.guestCount || 0,
      eventType: props.selectedEvent || '',
      package: props.selectedPackage || ''
    }
    emit('submit', requestData)
  }
  </script>