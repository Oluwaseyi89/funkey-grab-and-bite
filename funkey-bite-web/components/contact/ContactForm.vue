<template>
    <form @submit.prevent="submitForm" class="space-y-6">
      <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Send us a Message</h3>
      
      <!-- Name & Email -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Your Name *
          </label>
          <input
            v-model="formData.name"
            type="text"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
            placeholder="John Doe"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Email Address *
          </label>
          <input
            v-model="formData.email"
            type="email"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
            placeholder="john@example.com"
          />
        </div>
      </div>
  
      <!-- Phone -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Phone Number
        </label>
        <input
          v-model="formData.phone"
          type="tel"
          class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
          placeholder="+1 (555) 123-4567"
        />
      </div>
  
      <!-- Subject -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Subject *
        </label>
        <select
          v-model="formData.subject"
          required
          class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
        >
          <option value="" disabled selected>Select a subject</option>
          <option value="general">General Inquiry</option>
          <option value="catering">Catering Request</option>
          <option value="feedback">Feedback & Suggestions</option>
          <option value="partnership">Business Partnership</option>
          <option value="careers">Careers</option>
          <option value="other">Other</option>
        </select>
      </div>
  
      <!-- Message -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Message *
        </label>
        <textarea
          v-model="formData.message"
          rows="5"
          required
          class="w-full px-4 py-3 border border-gray-300 dark:border-slate-600 rounded-lg focus:ring-2 focus:ring-brand-500 focus:border-transparent"
          placeholder="How can we help you?"
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
          Sending Message...
        </template>
        <template v-else>
          Send Message
        </template>
      </button>
  
      <!-- Privacy Note -->
      <p class="text-sm text-gray-500 dark:text-gray-400 text-center">
        We'll respond within 24 hours. Your information is secure and won't be shared.
      </p>
    </form>
  </template>
  
  <script setup lang="ts">
  import { Loader2 } from 'lucide-vue-next'
  import { ref } from 'vue'
  
  interface ContactFormData {
    name: string
    email: string
    phone: string
    subject: string
    message: string
  }
  
  const props = withDefaults(defineProps<{
    isSubmitting?: boolean
  }>(), {
    isSubmitting: false
  })
  
  const emit = defineEmits<{
    submit: [data: ContactFormData]
  }>()
  
  const formData = ref<ContactFormData>({
    name: '',
    email: '',
    phone: '',
    subject: '',
    message: ''
  })
  
  const submitForm = () => {
    emit('submit', formData.value)
  }
  </script>