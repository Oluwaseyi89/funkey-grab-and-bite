<template>
  <div
    class="inline-flex items-center rounded-full"
    :class="size === 'sm' ? 'bg-ink-900 text-white dark:bg-white dark:text-ink-900' : 'border border-ink-200 bg-white dark:border-ink-700 dark:bg-ink-900'"
    role="group"
    :aria-label="`Quantity of ${name}`"
  >
    <button
      type="button"
      class="inline-flex items-center justify-center rounded-full transition active:scale-90 disabled:opacity-40"
      :class="size === 'sm' ? 'h-10 w-10' : 'h-11 w-11'"
      :aria-label="min === 0 && quantity <= 1 ? `Remove ${name}` : `One less ${name}`"
      :disabled="quantity <= min || disabled"
      @click="$emit('decrease')"
    >
      <Trash2 v-if="quantity <= 1 && min === 0" class="h-4 w-4" aria-hidden="true" />
      <Minus v-else class="h-4 w-4" aria-hidden="true" />
    </button>
    <span
      class="min-w-[1.75rem] text-center text-sm font-bold tabular-nums"
      aria-live="polite"
      aria-atomic="true"
    >
      <span class="sr-only">{{ name }}: </span>{{ quantity }}
    </span>
    <button
      type="button"
      class="inline-flex items-center justify-center rounded-full transition active:scale-90 disabled:opacity-40"
      :class="size === 'sm' ? 'h-10 w-10' : 'h-11 w-11'"
      :aria-label="`One more ${name}`"
      :disabled="atMax || disabled"
      @click="$emit('increase')"
    >
      <Plus class="h-4 w-4" aria-hidden="true" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { Minus, Plus, Trash2 } from 'lucide-vue-next'

// - [n] + controls. With min 0 the minus becomes a bin at 1, so the same
// button removes the item.
withDefaults(defineProps<{
  name: string
  quantity: number
  min?: number
  atMax?: boolean
  disabled?: boolean
  size?: 'sm' | 'md'
}>(), { min: 0, atMax: false, disabled: false, size: 'md' })

defineEmits<{ increase: []; decrease: [] }>()
</script>
