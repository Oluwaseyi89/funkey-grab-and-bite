<template>
  <section class="relative overflow-hidden" aria-labelledby="hero-title">
    <!-- Warm glow behind the page, in the logo's colours. -->
    <div aria-hidden="true" class="pointer-events-none absolute inset-0">
      <div class="absolute -left-24 -top-32 h-80 w-80 rounded-full bg-brand-300/40 blur-3xl dark:bg-brand-600/20" />
      <div class="absolute -right-20 top-10 h-72 w-72 rounded-full bg-accent-300/40 blur-3xl dark:bg-accent-500/10" />
      <div class="absolute bottom-0 left-1/3 h-64 w-64 rounded-full bg-leaf-200/50 blur-3xl dark:bg-leaf-500/10" />
    </div>

    <div class="container-narrow relative grid items-center gap-8 pb-8 pt-6 md:grid-cols-[1.1fr_1fr] md:gap-12 md:pb-16 md:pt-14">
      <div class="animate-rise-in">
        <StoreStatus />
        <h1 id="hero-title" class="mt-4 font-display text-[2.6rem] font-extrabold leading-[1.02] tracking-tight text-ink-900 dark:text-white sm:text-6xl lg:text-7xl">
          Hot, fresh &amp;<br />
          <span class="text-gradient">seriously funkey.</span>
        </h1>
        <p class="mt-4 max-w-md text-base text-ink-600 dark:text-ink-300 md:text-lg">
          Crispy chicken &amp; chips, juicy shawarma and loaded noodles. Order for pickup or delivery, or let us cater your event.
        </p>
        <div class="mt-6 grid grid-cols-[1.3fr_1fr] gap-3 sm:flex">
          <NuxtLink to="/menu" class="btn-primary btn-lg">
            <UtensilsCrossed class="h-5 w-5" aria-hidden="true" />
            Order now
          </NuxtLink>
          <NuxtLink to="/catering" class="btn-secondary btn-lg">
            <CalendarHeart class="h-5 w-5" aria-hidden="true" />
            Catering
          </NuxtLink>
        </div>
        <ul class="mt-6 flex flex-wrap gap-x-5 gap-y-2 text-sm text-ink-600 dark:text-ink-300">
          <li v-for="point in points" :key="point.text" class="inline-flex items-center gap-1.5">
            <component :is="point.icon" class="h-4 w-4 text-leaf-600 dark:text-leaf-400" aria-hidden="true" />
            {{ point.text }}
          </li>
        </ul>
      </div>

      <!-- Dishes from the real menu; a brand composition until photos load. -->
      <div class="relative mx-auto aspect-square w-full max-w-[22rem] md:max-w-md" aria-hidden="true">
        <div class="absolute inset-[8%] rounded-full bg-gradient-to-br from-brand-500 via-brand-400 to-accent-400 shadow-glow" />
        <div class="absolute inset-[8%] rounded-full bg-[radial-gradient(circle_at_30%_25%,rgba(255,255,255,0.35),transparent_55%)]" />
        <template v-if="photos.length">
          <div
            v-for="(photo, index) in photos"
            :key="photo"
            class="absolute overflow-hidden rounded-full border-4 border-white shadow-lift dark:border-ink-900 motion-safe:animate-drift"
            :class="slots[index]"
            :style="{ animationDelay: `${index * -2}s` }"
          >
            <img :src="photo" alt="" class="h-full w-full object-cover" :loading="index === 0 ? 'eager' : 'lazy'" decoding="async" />
          </div>
        </template>
        <div v-else class="absolute inset-0 flex items-center justify-center">
          <div class="flex h-1/2 w-1/2 items-center justify-center rounded-full bg-white/95 p-6 shadow-lift dark:bg-ink-900">
            <BrandLogo alt="" />
          </div>
        </div>
        <div v-if="readyIn" class="absolute bottom-[6%] left-0 rounded-2xl bg-white px-4 py-2.5 shadow-lift dark:bg-ink-900 motion-safe:animate-drift" style="animation-delay: -3s">
          <p class="text-xs text-ink-500 dark:text-ink-400">Popular picks ready in</p>
          <p class="font-display text-sm font-bold text-ink-900 dark:text-white">about {{ readyIn }} minutes</p>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Bike, CalendarHeart, CreditCard, ShoppingBag, UtensilsCrossed } from 'lucide-vue-next'
import BrandLogo from '../layout/BrandLogo.vue'
import StoreStatus from '../layout/StoreStatus.vue'
import { useMenuStore } from '../../stores/menu'

const menu = useMenuStore()

const points = [
  { icon: ShoppingBag, text: 'Pickup' },
  { icon: Bike, text: 'Delivery' },
  { icon: CreditCard, text: 'Card, transfer or pay on delivery' },
]

const photos = computed(() =>
  menu.featuredItems.map((item) => item.imageUrl).filter(Boolean).slice(0, 3),
)

// The quickest of the popular dishes, from the menu's preparation times.
const readyIn = computed(() => {
  const times = menu.featuredItems.map((item) => item.preparationTime).filter((time) => time > 0)
  return times.length ? Math.min(...times) : 0
})

// Where each photo sits on the circle: big centre-left, two smaller ones.
const slots = [
  'left-[6%] top-[14%] h-[58%] w-[58%]',
  'right-[2%] top-[4%] h-[38%] w-[38%]',
  'bottom-[10%] right-[8%] h-[44%] w-[44%]',
]
</script>
