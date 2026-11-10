<template>
  <div class="space-y-4">
    <div v-for="field in visibleFields" :key="field.key">
      <label :for="`co-${field.key}`" class="form-label">
        {{ field.label }}<span v-if="field.required" class="text-brand-600" aria-hidden="true"> *</span>
        <span v-else class="font-normal text-ink-500"> (optional)</span>
      </label>
      <component
        :is="field.multiline ? 'textarea' : 'input'"
        :id="`co-${field.key}`"
        :value="modelValue[field.key]"
        :type="field.multiline ? undefined : field.type"
        :rows="field.multiline ? 2 : undefined"
        :autocomplete="field.autocomplete"
        :inputmode="field.inputmode"
        :placeholder="field.placeholder"
        :required="field.required || undefined"
        :maxlength="field.maxlength"
        :aria-invalid="errors[field.key] ? 'true' : undefined"
        :aria-describedby="describedBy(field)"
        class="form-input"
        :class="field.multiline ? 'resize-none' : ''"
        @input="onInput(field.key, ($event.target as HTMLInputElement).value)"
        @blur="touch(field.key)"
      />
      <p v-if="field.hint && !errors[field.key]" :id="`co-${field.key}-hint`" class="form-hint">{{ field.hint }}</p>
      <p v-if="errors[field.key]" :id="`co-${field.key}-error`" class="form-error">{{ errors[field.key] }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, reactive } from 'vue'
import type { CustomerInfo, OrderType } from '../../types/order'

type Key = keyof CustomerInfo

interface Field {
  key: Key
  label: string
  type?: string
  autocomplete: string
  inputmode?: 'tel' | 'email' | 'text'
  placeholder?: string
  hint?: string
  required: boolean
  multiline?: boolean
  maxlength?: number
}

// The customer's contact details. Errors show once a field has been left (or
// on placing the order), never while the customer is still typing.
const props = defineProps<{ modelValue: CustomerInfo; orderType: OrderType }>()
const emit = defineEmits<{ 'update:modelValue': [value: CustomerInfo] }>()

const fields: Field[] = [
  { key: 'customerName', label: 'Full name', type: 'text', autocomplete: 'name', placeholder: 'e.g. Ada Okafor', required: true, maxlength: 100 },
  { key: 'customerPhone', label: 'Phone number', type: 'tel', autocomplete: 'tel', inputmode: 'tel', placeholder: '0801 234 5678', hint: 'We’ll call if there’s a question about your order.', required: true },
  { key: 'customerEmail', label: 'Email', type: 'email', autocomplete: 'email', inputmode: 'email', placeholder: 'you@example.com', hint: 'For your receipt and updates.', required: true, maxlength: 200 },
  { key: 'deliveryAddress', label: 'Delivery address', autocomplete: 'street-address', placeholder: 'House number, street, area and a landmark', required: true, multiline: true, maxlength: 500 },
  { key: 'notes', label: 'Notes for us', autocomplete: 'off', placeholder: 'Allergies, gate code, where to find you…', required: false, multiline: true, maxlength: 500 },
]

const visibleFields = computed(() => fields.filter((f) => f.key !== 'deliveryAddress' || props.orderType === 'delivery'))

const errors = reactive<Partial<Record<Key, string>>>({})

// Spaces, dashes and brackets are fine to type; the API strips them too.
const cleanPhone = (value: string) => value.trim().replace(/(?!^\+)[^\d]/g, '')

const check = (key: Key): string => {
  const value = (props.modelValue[key] ?? '').trim()
  switch (key) {
    case 'customerName': return value.length < 2 ? 'Enter your name.' : ''
    case 'customerPhone': return /^\+?\d{7,15}$/.test(cleanPhone(value)) ? '' : 'Enter a phone number we can reach, e.g. 0801 234 5678.'
    case 'customerEmail': return /^\S+@\S+\.\S+$/.test(value) ? '' : 'Enter an email address like you@example.com.'
    case 'deliveryAddress': return props.orderType === 'delivery' && value.length < 5 ? 'Enter the address to deliver to.' : ''
    default: return ''
  }
}

const describedBy = (field: Field) =>
  errors[field.key] ? `co-${field.key}-error` : field.hint ? `co-${field.key}-hint` : undefined

const onInput = (key: Key, value: string) => {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
  // Once an error is showing, clear it as soon as the field is fixed.
  if (errors[key]) nextTick(() => { errors[key] = check(key) })
}

const touch = (key: Key) => {
  errors[key] = check(key)
}

// Checks every field, shows what's wrong and moves to the first problem.
const validate = (): boolean => {
  let first: Key | null = null
  for (const field of visibleFields.value) {
    errors[field.key] = check(field.key)
    if (errors[field.key] && !first) first = field.key
  }
  if (first) {
    const el = document.getElementById(`co-${first}`)
    el?.focus()
    el?.scrollIntoView({ block: 'center', behavior: 'smooth' })
  }
  return !first
}

defineExpose({ validate })
</script>
