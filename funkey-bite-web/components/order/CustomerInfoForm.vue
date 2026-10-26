<template>
  <form @submit.prevent="submitForm" class="space-y-6" autocomplete="on" novalidate>
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
          autocomplete="name"
          class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
          placeholder="John Doe"
          @input="validateField('customerName')"
        />
        <p v-if="errors.customerName" class="text-red-500 text-xs mt-1">{{ errors.customerName }}</p>
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Phone Number *
        </label>
        <input
          v-model="formData.customerPhone"
          type="tel"
          required
          autocomplete="tel"
          inputmode="tel"
          class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
          placeholder="0801 234 5678"
          @input="validateField('customerPhone')"
        />
        <p v-if="errors.customerPhone" class="text-red-500 text-xs mt-1">{{ errors.customerPhone }}</p>
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
        autocomplete="email"
        inputmode="email"
        class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
        placeholder="john@example.com"
        @input="validateField('customerEmail')"
      />
      <p v-if="errors.customerEmail" class="text-red-500 text-xs mt-1">{{ errors.customerEmail }}</p>
    </div>

    <!-- Delivery Address -->
    <div v-if="props.orderType === 'delivery'">
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
        Delivery Address *
      </label>
      <textarea
        v-model="formData.deliveryAddress"
        rows="2"
        required
        autocomplete="street-address"
        class="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-transparent text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
        placeholder="House number, street, area and a landmark"
        @input="validateField('deliveryAddress')"
      ></textarea>
      <p v-if="errors.deliveryAddress" class="text-red-500 text-xs mt-1">{{ errors.deliveryAddress }}</p>
    </div>

    <!-- Order Notes -->
    <div>
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
        Special Instructions
      </label>
      <textarea
        v-model="formData.notes"
        rows="3"
        autocomplete="off"
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
import { ref, reactive } from 'vue'

import type { CustomerInfo, OrderType } from '../../types/order'

type FormData = CustomerInfo

const props = withDefaults(defineProps<{
  isLoading?: boolean
  submitText?: string
  orderType?: OrderType
  // Details entered earlier, so going back a step doesn't clear the form.
  initialData?: Partial<CustomerInfo>
}>(), {
  submitText: 'Continue to Payment',
  orderType: 'pickup',
  initialData: () => ({})
})

const emit = defineEmits<{
  submit: [data: FormData]
}>()

const formData = ref<FormData>({
  customerName: props.initialData.customerName ?? '',
  customerPhone: props.initialData.customerPhone ?? '',
  customerEmail: props.initialData.customerEmail ?? '',
  deliveryAddress: props.initialData.deliveryAddress ?? '',
  notes: props.initialData.notes ?? ''
})

// Spaces, dashes and brackets are fine to type; the API strips them too.
const cleanPhone = (value: string) => value.trim().replace(/(?!^\+)[^\d]/g, '')

const errors = reactive<{ [key: string]: string }>({})

function validateField(field: keyof FormData) {
  switch (field) {
    case 'customerName':
      errors.customerName = formData.value.customerName.trim().length < 2
        ? 'Please enter your full name.'
        : ''
      break
    case 'customerPhone':
      errors.customerPhone =
        !/^\+?\d{7,15}$/.test(cleanPhone(formData.value.customerPhone))
          ? 'Enter a valid phone number, e.g. 0801 234 5678.'
          : ''
      break
    case 'deliveryAddress':
      errors.deliveryAddress = props.orderType === 'delivery' && formData.value.deliveryAddress.trim().length < 5
        ? 'Enter the address we should deliver to.'
        : ''
      break
    case 'customerEmail':
      errors.customerEmail =
        !/^\S+@\S+\.\S+$/.test(formData.value.customerEmail.trim())
          ? 'Enter a valid email address.'
          : ''
      break
  }
}

function validateAll() {
  validateField('customerName')
  validateField('customerPhone')
  validateField('customerEmail')
  validateField('deliveryAddress')
  return !errors.customerName && !errors.customerPhone && !errors.customerEmail && !errors.deliveryAddress
}

function submitForm() {
  if (!validateAll()) return
  emit('submit', {
    ...formData.value,
    customerName: formData.value.customerName.trim(),
    customerPhone: formData.value.customerPhone.trim(),
    customerEmail: formData.value.customerEmail.trim(),
    deliveryAddress: props.orderType === 'delivery' ? formData.value.deliveryAddress.trim() : '',
  })
}
</script>