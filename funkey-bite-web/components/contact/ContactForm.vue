<template>
    <form @submit.prevent="submitForm" class="space-y-6">
      <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">Send us a Message</h3>
      
      <!-- Name & Email -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="contact-name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Your Name *
          </label>
          <input
            v-model="formData.name"
            id="contact-name"
            type="text"
            autocomplete="name"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
            placeholder="John Doe"
          />
        </div>
        <div>
          <label for="contact-email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Email Address *
          </label>
          <input
            v-model="formData.email"
            id="contact-email"
            type="email"
            autocomplete="email"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
            placeholder="john@example.com"
          />
        </div>
      </div>
  
      <!-- Phone -->
      <div>
        <label for="contact-phone" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Phone Number
        </label>
        <input
          v-model="formData.phone"
          id="contact-phone"
          type="tel"
          autocomplete="tel"
          class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
          placeholder="+1 (555) 123-4567"
        />
      </div>
  
      <!-- Subject -->
      <div>
        <label for="contact-subject" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Subject *
        </label>
        <select
          v-model="formData.subject"
          id="contact-subject"
          required
          class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white"
        >
          <option value="" disabled selected class="text-gray-500 dark:text-gray-400">Select a subject</option>
          <option value="general" class="text-gray-900 dark:text-white">General Inquiry</option>
          <option value="catering" class="text-gray-900 dark:text-white">Catering Request</option>
          <option value="feedback" class="text-gray-900 dark:text-white">Feedback & Suggestions</option>
          <option value="partnership" class="text-gray-900 dark:text-white">Business Partnership</option>
          <option value="careers" class="text-gray-900 dark:text-white">Careers</option>
          <option value="other" class="text-gray-900 dark:text-white">Other</option>
        </select>
      </div>
  
      <!-- Message -->
      <div>
        <label for="contact-message" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Message *
        </label>
        <textarea
          v-model="formData.message"
          id="contact-message"
          rows="5"
          required
          class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
          placeholder="How can we help you?"
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