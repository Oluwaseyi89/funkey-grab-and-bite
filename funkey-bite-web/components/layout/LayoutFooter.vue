<template>
  <footer class="border-t border-ink-200/70 bg-white dark:border-ink-800 dark:bg-ink-900">
    <div class="container-narrow grid gap-8 py-10 md:grid-cols-[1.4fr_1fr_1fr_1.2fr] md:py-14">
      <div class="space-y-3">
        <NuxtLink to="/" class="inline-flex h-11" aria-label="Funkey Grab & Bite home">
          <BrandLogo alt="" />
        </NuxtLink>
        <p class="max-w-xs text-sm text-ink-600 dark:text-ink-300">
          Chicken &amp; chips, shawarma, noodles and lunch packs for pickup or delivery, plus catering for your events.
        </p>
        <StoreStatus />
      </div>

      <nav aria-label="Food" class="hidden md:block">
        <h2 class="mb-3 text-sm font-semibold text-ink-900 dark:text-white">Food</h2>
        <ul class="space-y-2 text-sm">
          <li><NuxtLink to="/menu" class="footer-link">Full menu</NuxtLink></li>
          <li v-for="category in menu.visibleCategories.slice(0, 5)" :key="category.id">
            <NuxtLink :to="`/menu#cat-${category.id}`" class="footer-link">{{ category.name }}</NuxtLink>
          </li>
        </ul>
      </nav>

      <nav aria-label="Help">
        <h2 class="mb-3 text-sm font-semibold text-ink-900 dark:text-white">Help</h2>
        <ul class="grid grid-cols-2 gap-2 text-sm md:grid-cols-1">
          <li><NuxtLink to="/order/track" class="footer-link">Track an order</NuxtLink></li>
          <li><NuxtLink to="/catering" class="footer-link">Catering</NuxtLink></li>
          <li><NuxtLink to="/promotion" class="footer-link">Deals</NuxtLink></li>
          <li><NuxtLink to="/order/invoice" class="footer-link">Find an invoice</NuxtLink></li>
          <li><NuxtLink to="/contact" class="footer-link">Contact us</NuxtLink></li>
          <li><NuxtLink to="/account" class="footer-link">My account</NuxtLink></li>
        </ul>
      </nav>

      <div class="space-y-3 text-sm">
        <h2 class="text-sm font-semibold text-ink-900 dark:text-white">Visit or call</h2>
        <p v-if="fullAddress" class="flex gap-2 text-ink-600 dark:text-ink-300">
          <MapPin class="mt-0.5 h-4 w-4 shrink-0 text-brand-500" aria-hidden="true" />
          <a v-if="mapsUrl" :href="mapsUrl" target="_blank" rel="noopener noreferrer" class="hover:underline">{{ fullAddress }}</a>
        </p>
        <p v-if="phoneHref" class="flex gap-2">
          <Phone class="mt-0.5 h-4 w-4 shrink-0 text-brand-500" aria-hidden="true" />
          <a :href="phoneHref" class="text-ink-700 hover:underline dark:text-ink-200">{{ settings?.phoneNumber }}</a>
        </p>
        <p v-if="settings?.email" class="flex gap-2">
          <Mail class="mt-0.5 h-4 w-4 shrink-0 text-brand-500" aria-hidden="true" />
          <a :href="`mailto:${settings.email}`" class="break-all text-ink-700 hover:underline dark:text-ink-200">{{ settings.email }}</a>
        </p>
        <details v-if="status === 'ready'" class="group rounded-xl bg-ink-50 px-4 py-2 dark:bg-ink-800/60">
          <summary class="flex min-h-[40px] cursor-pointer list-none items-center justify-between font-medium text-ink-800 dark:text-ink-100">
            Opening hours
            <ChevronDown class="h-4 w-4 transition group-open:rotate-180" aria-hidden="true" />
          </summary>
          <dl class="space-y-1.5 pb-2 pt-1">
            <div v-for="group in hoursGroups" :key="group.label" class="flex justify-between gap-4">
              <dt class="text-ink-500 dark:text-ink-400">{{ group.label }}</dt>
              <dd class="text-right font-medium text-ink-800 dark:text-ink-100">{{ group.hours }}</dd>
            </div>
          </dl>
        </details>
      </div>
    </div>

    <div class="border-t border-ink-100 dark:border-ink-800">
      <div class="container-narrow flex flex-col items-center justify-between gap-2 py-5 text-xs text-ink-500 dark:text-ink-400 md:flex-row">
        <p>© {{ year }} Funkey Grab &amp; Bite</p>
        <p class="flex gap-4">
          <NuxtLink to="/legal/privacy" class="hover:underline">Privacy</NuxtLink>
          <NuxtLink to="/legal/terms" class="hover:underline">Terms</NuxtLink>
        </p>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { ChevronDown, Mail, MapPin, Phone } from 'lucide-vue-next'
import BrandLogo from './BrandLogo.vue'
import StoreStatus from './StoreStatus.vue'
import { useBusinessSettings } from '../../composables/useBusinessSettings'
import { useMenuStore } from '../../stores/menu'

const { settings, status, fullAddress, mapsUrl, hoursGroups, phoneHref, load } = useBusinessSettings()
const menu = useMenuStore()
const year = new Date().getFullYear()

onMounted(() => {
  load()
  menu.load()
})
</script>

<style scoped>
.footer-link {
  @apply inline-flex min-h-[32px] items-center text-ink-600 transition hover:text-brand-600 dark:text-ink-300 dark:hover:text-brand-400;
}
</style>
