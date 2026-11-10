<template>
  <div class="flex min-h-dvh flex-col">
    <a
      href="#main-content"
      class="sr-only focus:not-sr-only focus:fixed focus:left-4 focus:top-4 focus:z-[100] focus:rounded-xl focus:bg-brand-600 focus:px-4 focus:py-2 focus:text-white"
    >
      Skip to main content
    </a>

    <LayoutHeader />

    <main id="main-content" class="flex-1 pt-14 md:pt-16" tabindex="-1">
      <slot />
    </main>

    <!-- Room below the footer for the tab bar and the cart bar on phones. -->
    <LayoutFooter class="pb-[calc(9rem+env(safe-area-inset-bottom))] md:pb-0" />

    <CartBar />
    <BottomNav />
    <CartSheet />

    <!-- Cart changes, read out by screen readers. -->
    <p class="sr-only" aria-live="polite" aria-atomic="true">{{ cart.announcement }}</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useHead, useSeoMeta } from 'nuxt/app'
import LayoutHeader from '../components/layout/LayoutHeader.vue'
import LayoutFooter from '../components/layout/LayoutFooter.vue'
import BottomNav from '../components/layout/BottomNav.vue'
import CartBar from '../components/cart/CartBar.vue'
import CartSheet from '../components/cart/CartSheet.vue'
import { useBusinessSettings } from '../composables/useBusinessSettings'
import { useCartStore } from '../stores/cart'

const { settings, load } = useBusinessSettings()
const cart = useCartStore()

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
  title: 'Fast Food & Catering',
  ogTitle: 'Funkey Grab & Bite - Fast Food & Catering',
  description: 'Order chicken & chips, shawarma, noodles and lunch packs for pickup or delivery, or book catering for your event.',
  ogDescription: 'Order chicken & chips, shawarma, noodles and lunch packs for pickup or delivery, or book catering for your event.',
  ogImage: 'https://funkey-static-assets.s3.amazonaws.com/branding/og-image.jpg',
  twitterCard: 'summary_large_image',
})
</script>
