<template>
  <div class="relative overflow-hidden">
    <div aria-hidden="true" class="pointer-events-none absolute -left-24 -top-24 h-72 w-72 rounded-full bg-brand-200/40 blur-3xl dark:bg-brand-700/15" />
    <div aria-hidden="true" class="pointer-events-none absolute -right-20 top-32 h-64 w-64 rounded-full bg-leaf-200/40 blur-3xl dark:bg-leaf-600/10" />
    <div class="container-narrow relative max-w-md py-8 md:py-14">
      <h1 class="text-center font-display text-3xl font-bold text-ink-900 dark:text-white">
        {{ titleBefore }} <span class="text-gradient">{{ highlight }}</span>
      </h1>
      <p class="mt-2 text-center text-ink-600 dark:text-ink-300">{{ subtitle }}</p>
      <div class="card mt-6 space-y-5 p-6">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// The frame shared by the sign-in, registration and password reset pages.
defineProps<{ titleBefore: string; highlight: string; subtitle: string }>()
</script>
