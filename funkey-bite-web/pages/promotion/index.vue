<template>
  <div>
    <PageHeader eyebrow="Deals" title="Save on your next order" subtitle="Tap “Use code” and it’s applied at checkout. We check every code when you order." />

    <div class="container-narrow py-6 md:py-10">
      <div v-if="loading" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3" aria-busy="true" aria-label="Loading deals">
        <div v-for="i in 3" :key="i" class="skeleton h-56 rounded-3xl" />
      </div>

      <ul v-else-if="promotions.length" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <li v-for="promotion in promotions" :key="promotion.id" :id="`deal-${promotion.id}`">
          <PromoCard :promotion="promotion" />
        </li>
      </ul>

      <div v-else class="card mx-auto max-w-md p-8 text-center">
        <TicketPercent class="mx-auto h-10 w-10 text-ink-300" aria-hidden="true" />
        <p class="mt-3 font-semibold text-ink-900 dark:text-white">No deals running right now</p>
        <p class="mt-1 text-sm text-ink-500 dark:text-ink-400">Check back soon. In the meantime, the full menu is a tap away.</p>
        <NuxtLink to="/menu" class="btn-primary mt-5">Browse the menu</NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { TicketPercent } from 'lucide-vue-next'
import PageHeader from '../../components/layout/PageHeader.vue'
import PromoCard from '../../components/promo/PromoCard.vue'
import { useApi } from '../../utils/api'
import { usePageSeo } from '../../composables/usePageSeo'
import type { Promotion } from '../../types/menu'

usePageSeo({ title: 'Deals', description: 'Current promo codes for Funkey Grab & Bite orders.', path: '/promotion' })

const promotions = ref<Promotion[]>([])
const loading = ref(true)

onMounted(async () => {
  promotions.value = await useApi().getActivePromotions()
  loading.value = false
})
</script>
