<template>
  <div>
    <HeroBanner />
    <CategoriesGrid />
    <FeaturedItems />
    <DealsSection />

    <section class="container-narrow py-6 md:py-10" aria-labelledby="how-title">
      <h2 id="how-title" class="section-title">Ordering takes a minute</h2>
      <ol class="mt-5 grid gap-3 sm:grid-cols-3 md:gap-5">
        <li v-for="(step, index) in steps" :key="step.title" class="card flex gap-4 p-5 sm:flex-col">
          <span class="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl" :class="step.tint">
            <component :is="step.icon" class="h-6 w-6" aria-hidden="true" />
          </span>
          <div>
            <p class="text-xs font-semibold text-ink-400">Step {{ index + 1 }}</p>
            <h3 class="font-display font-semibold text-ink-900 dark:text-white">{{ step.title }}</h3>
            <p class="mt-1 text-sm text-ink-500 dark:text-ink-400">{{ step.text }}</p>
          </div>
        </li>
      </ol>
    </section>

    <CateringCTA />

    <section class="container-narrow pb-10 pt-2 md:pb-16">
      <NuxtLink to="/order/track" class="card flex items-center gap-4 p-5 transition hover:shadow-lift">
        <span class="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-leaf-50 text-leaf-700 dark:bg-leaf-800/30 dark:text-leaf-300">
          <PackageSearch class="h-6 w-6" aria-hidden="true" />
        </span>
        <span class="flex-1">
          <span class="block font-display font-semibold text-ink-900 dark:text-white">Already ordered?</span>
          <span class="block text-sm text-ink-500 dark:text-ink-400">Track your order and payment</span>
        </span>
        <ChevronRight class="h-5 w-5 text-ink-400" aria-hidden="true" />
      </NuxtLink>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { ChevronRight, CreditCard, PackageSearch, ShoppingBag, UtensilsCrossed } from 'lucide-vue-next'
import HeroBanner from '../components/home/HeroBanner.vue'
import CategoriesGrid from '../components/home/CategoriesGrid.vue'
import FeaturedItems from '../components/home/FeaturedItems.vue'
import DealsSection from '../components/home/DealsSection.vue'
import CateringCTA from '../components/home/CateringCTA.vue'
import { useMenuStore } from '../stores/menu'
import { usePageSeo } from '../composables/usePageSeo'

usePageSeo({
  title: 'Fast Food & Catering',
  description: 'Order chicken & chips, shawarma, noodles and lunch packs for pickup or delivery, or book catering for your event.',
  path: '/',
})

const menu = useMenuStore()
onMounted(() => menu.load())

const steps = [
  { icon: UtensilsCrossed, title: 'Pick your food', text: 'Tap Add on anything you fancy. Change quantities right on the menu.', tint: 'bg-brand-50 text-brand-700 dark:bg-brand-900/30 dark:text-brand-300' },
  { icon: CreditCard, title: 'Pay your way', text: 'Card or bank transfer online, or pay when it arrives where available.', tint: 'bg-accent-50 text-accent-700 dark:bg-accent-700/20 dark:text-accent-300' },
  { icon: ShoppingBag, title: 'Pick up or get it delivered', text: 'Follow your order’s progress any time from Orders.', tint: 'bg-leaf-50 text-leaf-700 dark:bg-leaf-800/30 dark:text-leaf-300' },
]
</script>
