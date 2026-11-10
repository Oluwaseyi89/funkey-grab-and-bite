<template>
  <div class="relative overflow-hidden bg-gradient-to-br from-accent-100 via-brand-50 to-leaf-100 dark:from-ink-800 dark:via-ink-800 dark:to-ink-700">
    <img
      v-if="src && !failed"
      :src="src"
      :alt="alt"
      :loading="eager ? 'eager' : 'lazy'"
      decoding="async"
      class="h-full w-full object-cover transition duration-500"
      :class="loaded ? 'opacity-100' : 'opacity-0'"
      @load="loaded = true"
      @error="failed = true"
    />
    <div v-if="!src || failed" class="absolute inset-0 flex items-center justify-center" role="img" :aria-label="alt">
      <UtensilsCrossed class="h-1/3 w-1/3 max-h-10 max-w-10 text-brand-300 dark:text-ink-500" aria-hidden="true" />
    </div>
    <slot />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { UtensilsCrossed } from 'lucide-vue-next'

// A dish photo that fades in, with a warm placeholder when there's no photo
// or it fails to load, so a broken image never shows.
const props = withDefaults(defineProps<{ src?: string; alt: string; eager?: boolean }>(), { src: '', eager: false })

const loaded = ref(false)
const failed = ref(false)
watch(() => props.src, () => { loaded.value = false; failed.value = false })
</script>
