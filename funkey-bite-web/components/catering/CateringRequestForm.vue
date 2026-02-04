<template>
    <form @submit.prevent="submitForm" class="space-y-6">
      <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Request Catering Quote</h3>
      
      <!-- Contact Information -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Contact Name *
          </label>
          <input
            v-model="formData.contactName"
            type="text"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
            placeholder="John Doe"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Phone Number *
          </label>
          <input
            v-model="formData.contactPhone"
            type="tel"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
            placeholder="+1 (555) 123-4567"
          />
        </div>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Email Address *
        </label>
        <input
          v-model="formData.contactEmail"
          type="email"
          required
          class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
          placeholder="john@example.com"
        />
      </div>
      
      <!-- Event Details -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Event Date *
          </label>
          <input
            v-model="formData.eventDate"
            type="date"
            required
            :min="minDate"
            class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Event Time *
          </label>
          <input
            v-model="formData.eventTime"
            type="time"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
          />
        </div>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Event Name / Occasion
        </label>
        <input
          v-model="formData.eventName"
          type="text"
          class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
          placeholder="e.g., Company Annual Meeting, Sarah's Wedding"
        />
      </div>
      
      <!-- Budget -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Estimated Budget
        </label>
        <div class="relative">
          <span class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-500">$</span>
          <input
            v-model="formData.budget"
            type="number"
            min="0"
            step="100"
            class="w-full pl-10 pr-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
            placeholder="Enter your budget (optional)"
          />
        </div>
      </div>
      
      <!-- Special Requests -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Special Requests / Notes
        </label>
        <textarea
          v-model="formData.specialRequests"
          rows="4"
          class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
          placeholder="Dietary restrictions, venue details, theme, etc."
        ></textarea>
      </div>
      
      <!-- Submit Button -->
      <button
        type="submit"
        :disabled="isSubmitting"
        class="w-full btn-primary py-4 text-lg font-bold disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <template v-if="isSubmitting">
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