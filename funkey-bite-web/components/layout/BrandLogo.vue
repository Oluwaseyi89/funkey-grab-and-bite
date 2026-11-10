<template>
  <span class="inline-flex h-full items-center">
    <img src="/brand/logo.png" :alt="alt" class="h-full w-auto dark:hidden" width="512" height="202" decoding="async" />
    <img src="/brand/logo-dark.png" alt="" aria-hidden="true" class="hidden h-full w-auto dark:block" width="512" height="202" decoding="async" />
  </span>
</template>

<script setup lang="ts">
// The Funkey Grab & Bite wordmark, with a white-lettered version for dark mode.
withDefaults(defineProps<{ alt?: string }>(), { alt: 'Funkey Grab & Bite' })
</script>
