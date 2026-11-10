<template>
  <fieldset>
    <legend class="sr-only">How do you want your order?</legend>
    <div class="grid grid-cols-2 gap-1 rounded-2xl bg-ink-100 p-1 dark:bg-ink-800">
      <label
        v-for="type in types"
        :key="type.id"
        class="relative flex cursor-pointer flex-col items-center justify-center gap-0.5 rounded-xl px-3 py-3 text-center transition has-[:focus-visible]:ring-2 has-[:focus-visible]:ring-brand-500"
        :class="modelValue === type.id
          ? 'bg-white text-ink-900 shadow-soft dark:bg-ink-900 dark:text-white'
          : 'text-ink-600 hover:text-ink-900 dark:text-ink-300 dark:hover:text-white'"
      >
        <input
          type="radio"
          name="order-type"
          class="sr-only"
          :value="type.id"
          :checked="modelValue === type.id"
          @change="$emit('update:modelValue', type.id)"
        />
        <span class="inline-flex items-center gap-2 font-semibold">
          <component :is="type.icon" class="h-5 w-5" aria-hidden="true" />
          {{ type.name }}
        </span>
        <span class="text-xs" :class="type.unavailable ? 'text-brand-700 dark:text-brand-300' : 'text-ink-500 dark:text-ink-400'">
          {{ type.hint }}
        </span>
      </label>
    </div>
  </fieldset>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Bike, Store } from 'lucide-vue-next'
import { naira } from '../../utils/format'
import type { OrderType } from '../../types/order'

// Pickup or delivery, chosen first so the fee and availability are clear
// before anything else. Catering has its own page.
const props = defineProps<{
  modelValue: OrderType
  deliveryFee?: number | null
  pickupOpen?: boolean
  deliveryOpen?: boolean
}>()

defineEmits<{ 'update:modelValue': [type: OrderType] }>()

const types = computed(() => [
  {
    id: 'pickup' as OrderType,
    name: 'Pickup',
    icon: Store,
    unavailable: props.pickupOpen === false,
    hint: props.pickupOpen === false ? 'Not available now' : 'Free',
  },
  {
    id: 'delivery' as OrderType,
    name: 'Delivery',
    icon: Bike,
    unavailable: props.deliveryOpen === false,
    hint: props.deliveryOpen === false
      ? 'Not available now'
      : props.deliveryFee ? `${naira(props.deliveryFee)} fee` : props.deliveryFee === 0 ? 'Free' : 'Fee at checkout',
  },
])
</script>
