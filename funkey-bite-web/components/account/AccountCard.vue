<template>
  <div class="mt-[-30px]">
    <PageHeader :title-before="titleBefore" :highlight-text="highlight" :subtitle="subtitle" :narrow="true" :alignment="'center'" variant="gradient" />
    <div class="container-narrow px-4 sm:px-6 lg:px-8 py-8 max-w-md">
      <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700 space-y-4">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import PageHeader from '../layout/PageHeader.vue'

// The frame shared by the sign-in, registration and password reset pages.
defineProps<{ titleBefore: string; highlight: string; subtitle: string }>()
</script>
