<template>
    <form @submit.prevent="submitForm" class="space-y-6">
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Customer Information</h3>
      
      <!-- Name & Phone -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Full Name *
          </label>
          <input
            v-model="formData.customerName"
            type="text"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
            placeholder="John Doe"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Phone Number *
          </label>
          <input
            v-model="formData.customerPhone"
            type="tel"
            required
            class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
            placeholder="+1 (555) 123-4567"
          />
        </div>
      </div>
  
      <!-- Email -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Email Address *
        </label>
        <input
          v-model="formData.customerEmail"
          type="email"
          required
          class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
          placeholder="john@example.com"
        />
      </div>
  
      <!-- Order Notes -->
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Special Instructions
        </label>
        <textarea
          v-model="formData.notes"
          rows="3"
          class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
          placeholder="Allergies, delivery instructions, etc."
        ></textarea>
      </div>
  
      <!-- Submit -->
      <button
        type="submit"
        :disabled="props.isLoading"
        class="w-full btn-primary py-4 text-lg font-bold disabled:opacity-50 disabled:cursor-not-allowed dark:text-white"
      >
        <template v-if="props.isLoading">
          <Loader2 class="w-5 h-5 animate-spin inline mr-2" />
          Processing...
        </template>
        <template v-else>
          {{ props.submitText }}
        </template>
      </button>
    </form>
  </template>
  
  <script setup lang="ts">
  import { Loader2 } from 'lucide-vue-next'
  import { ref } from 'vue'
  
  interface FormData {
    customerName: string
    customerPhone: string
    customerEmail: string
    notes: string
  }
  
  const props = withDefaults(defineProps<{
    isLoading?: boolean
    submitText?: string
  }>(), {
    submitText: 'Continue to Payment'
  })
  
  const emit = defineEmits<{
    submit: [data: FormData]
  }>()
  
  const formData = ref<FormData>({
    customerName: '',
    customerPhone: '',
    customerEmail: '',
    notes: ''
  })
  
  const submitForm = () => {
    emit('submit', formData.value)
  }
  </script>