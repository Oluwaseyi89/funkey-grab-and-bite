<template>
  <div class="container-narrow max-w-xl py-6 md:py-12">
    <NuxtLink to="/promotion" class="link inline-flex min-h-[44px] items-center gap-1 text-sm">
      <ArrowLeft class="h-4 w-4" aria-hidden="true" /> All deals
    </NuxtLink>

    <div v-if="loading" class="skeleton mt-3 h-64 rounded-3xl" aria-busy="true" />

    <template v-else-if="promotion">
      <h1 class="sr-only">{{ promotion.title }}</h1>
      <div class="mt-3"><PromoCard :promotion="promotion" /></div>
      <section class="card mt-4 p-5" aria-labelledby="deal-terms">
        <h2 id="deal-terms" class="font-display text-lg font-semibold text-ink-900 dark:text-white">How it works</h2>
        <ul class="mt-3 space-y-2 text-sm text-ink-600 dark:text-ink-300">
          <li class="flex gap-2"><Check class="mt-0.5 h-4 w-4 shrink-0 text-leaf-600" aria-hidden="true" />Valid {{ formatDay(promotion.validFrom) }} to {{ formatDay(promotion.validUntil) }}.</li>
          <li v-if="promotion.minOrderAmount" class="flex gap-2"><Check class="mt-0.5 h-4 w-4 shrink-0 text-leaf-600" aria-hidden="true" />For orders of {{ naira(promotion.minOrderAmount) }} or more.</li>
          <li v-if="promotion.maxDiscount" class="flex gap-2"><Check class="mt-0.5 h-4 w-4 shrink-0 text-leaf-600" aria-hidden="true" />Saves up to {{ naira(promotion.maxDiscount) }}.</li>
          <li class="flex gap-2"><Check class="mt-0.5 h-4 w-4 shrink-0 text-leaf-600" aria-hidden="true" />Enter the code at checkout; we confirm the discount before you pay.</li>
        </ul>
      </section>
    </template>

    <div v-else class="card mt-3 p-8 text-center">
      <p class="font-semibold text-ink-900 dark:text-white">This deal has ended</p>
      <p class="mt-1 text-sm text-ink-500 dark:text-ink-400">Have a look at what’s running now.</p>
      <NuxtLink to="/promotion" class="btn-primary mt-5">See current deals</NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeft, Check } from 'lucide-vue-next'
import PromoCard from '../../components/promo/PromoCard.vue'
import { useApi } from '../../utils/api'
import { formatDay, naira } from '../../utils/format'
import { usePageSeo } from '../../composables/usePageSeo'
import type { Promotion } from '../../types/menu'

const route = useRoute()
usePageSeo({ title: 'Deal', description: 'A Funkey Grab & Bite promo code.', path: route.path })

const promotion = ref<Promotion | null>(null)
const loading = ref(true)

onMounted(async () => {
  const all = await useApi().getActivePromotions()
  promotion.value = all.find((promo) => String(promo.id) === String(route.params.id)) ?? null
  loading.value = false
})
</script>
