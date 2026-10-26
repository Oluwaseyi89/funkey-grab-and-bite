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
import { onMounted } from 'vue'
import { useHead, useSeoMeta } from 'nuxt/app'; 
import { useBusinessSettings } from '../composables/useBusinessSettings'

const { settings, load } = useBusinessSettings()

onMounted(() => load())

// Live contact details and hours for the Restaurant structured data; the
// shared @id merges this with the static entry in nuxt.config.ts.
useHead(() => {
  const s = settings.value
  if (!s) return {}
  return {
    script: [
      {
        key: 'business-ld',
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'Restaurant',
          '@id': 'https://funkeygrabandbite.com',
          name: s.businessName,
          telephone: s.phoneNumber,
          email: s.email,
          address: {
            '@type': 'PostalAddress',
            streetAddress: s.address,
            addressLocality: s.city || undefined,
            addressRegion: s.state || undefined,
            postalCode: s.postalCode || undefined,
            addressCountry: s.country || undefined,
          },
          openingHoursSpecification: s.openingHours
            .filter((h) => h.isOpen)
            .map((h) => ({
              '@type': 'OpeningHoursSpecification',
              dayOfWeek: h.day,
              opens: h.openTime,
              closes: h.closeTime,
            })),
        }),
      },
    ],
  }
})

useSeoMeta({
  title: 'Funkey Grab & Bite - Fast Food & Catering',
  ogTitle: 'Funkey Grab & Bite - Fast Food & Catering',
  description: 'Fast-food business offering quick take-out meals, lunch packs, and reliable indoor/outdoor catering services',
  ogDescription: 'Fast-food business offering quick take-out meals, lunch packs, and reliable indoor/outdoor catering services',
  ogImage: 'https://funkey-static-assets.s3.amazonaws.com/branding/og-image.jpg',
  twitterCard: 'summary_large_image',
})
</script>