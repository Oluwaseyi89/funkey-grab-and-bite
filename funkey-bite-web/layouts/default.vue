<template>
  <div class="min-h-screen flex flex-col bg-gradient-to-br from-white via-red-50 to-white dark:from-slate-900 dark:via-slate-800 dark:to-slate-900">

    <a href="#main-content" class="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-[9999] focus:px-4 focus:py-2 focus:bg-brand-500 focus:text-white focus:rounded-lg focus:shadow-lg">
      Skip to main content
    </a>

    <LayoutHeader />

    <main id="main-content" class="flex-grow pt-28 pb-12 md:pt-32">
      <slot />
    </main>
    

    <LayoutFooter />
  </div>
</template>

<script setup lang="ts">
import { useSeoMeta } from 'nuxt/app'; 

useSeoMeta({
  title: 'Funkey Grab & Bite - Fast Food & Catering',
  ogTitle: 'Funkey Grab & Bite - Fast Food & Catering',
  description: 'Fast-food business offering quick take-out meals, lunch packs, and reliable indoor/outdoor catering services',
  ogDescription: 'Fast-food business offering quick take-out meals, lunch packs, and reliable indoor/outdoor catering services',
  ogImage: 'https://funkey-static-assets.s3.amazonaws.com/branding/og-image.jpg',
  twitterCard: 'summary_large_image',
})
</script>