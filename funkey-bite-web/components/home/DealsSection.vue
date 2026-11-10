<template>
  <section v-if="promotions.length" class="py-6 md:py-10" aria-labelledby="deals-title">
    <div class="container-narrow flex items-end justify-between gap-4">
      <div>
        <p class="eyebrow">Deals</p>
        <h2 id="deals-title" class="section-title mt-1">Save on your next order</h2>
      </div>
      <NuxtLink v-if="promotions.length > 2" to="/promotion" class="link inline-flex min-h-[44px] items-center gap-1 text-sm">
        All deals <ArrowRight class="h-4 w-4" aria-hidden="true" />
      </NuxtLink>
    </div>
    <div class="container-narrow mt-4">
      <ul class="scroll-row -mx-4 px-4 pb-2 sm:-mx-6 sm:px-6 md:mx-0 md:grid md:grid-cols-2 md:gap-5 md:overflow-visible md:px-0 lg:grid-cols-3">
        <li v-for="promotion in promotions.slice(0, 3)" :key="promotion.id" class="w-[85%] shrink-0 snap-center sm:w-80 md:w-auto">
          <PromoCard :promotion="promotion" />
        </li>
      </ul>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ArrowRight } from 'lucide-vue-next'
import PromoCard from '../promo/PromoCard.vue'
import { useApi } from '../../utils/api'
import type { Promotion } from '../../types/menu'

// Live deals from the admin. Hidden entirely when there are none.
const promotions = ref<Promotion[]>([])
onMounted(async () => {
  promotions.value = await useApi().getActivePromotions()
})
</script>
