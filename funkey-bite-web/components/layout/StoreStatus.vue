<template>
  <span
    v-if="label"
    class="inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-semibold"
    :class="isOpen
      ? 'bg-leaf-50 text-leaf-700 ring-1 ring-inset ring-leaf-200 dark:bg-leaf-800/30 dark:text-leaf-300 dark:ring-leaf-700/50'
      : 'bg-ink-100 text-ink-700 ring-1 ring-inset ring-ink-200 dark:bg-ink-800 dark:text-ink-200 dark:ring-ink-700'"
  >
    <span class="relative flex h-2 w-2" aria-hidden="true">
      <span v-if="isOpen" class="absolute inline-flex h-full w-full animate-ping rounded-full bg-leaf-400 opacity-60 motion-reduce:hidden" />
      <span class="relative inline-flex h-2 w-2 rounded-full" :class="isOpen ? 'bg-leaf-500' : 'bg-ink-400'" />
    </span>
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import { useStoreStatus } from '../../composables/useStoreStatus'

// "Open now · until 10:00 PM" / "Closed · opens 8:00 AM tomorrow".
const { isOpen, label } = useStoreStatus()
</script>
