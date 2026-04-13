<template>
  <div>

    <HeroBanner />

    <CategoriesGrid :categories="categories" />

    <FeaturedItems :items="featuredItems" :is-loading="menuStore.isLoading" />

    <!-- <PromotionsBanner 
    :promotions="activePromotions"
    @claim-offer="handleClaimOffer"
    @view-details="viewPromotionDetails"
  /> -->

  <PromotionsBanner 
  :promotions="mockPromotions"
  :menu-items="mockMenuItems"
  @claim-offer="handleClaimOffer"
/>

    <CateringCTA />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMenuStore } from '../stores/menu'

import { mockPromotions, mockMenuItems } from '../utils/mockData'
import type { Promotion } from '../types/menu'

import HeroBanner from '../components/home/HeroBanner.vue'
import CategoriesGrid from '../components/home/CategoriesGrid.vue'
import FeaturedItems from '../components/home/FeaturedItems.vue'
import PromotionsBanner from '../components/home/PromotionsBanner.vue'
import CateringCTA from '../components/home/CateringCTA.vue'

import gsap from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { navigateTo } from 'nuxt/app'

const menuStore = useMenuStore()
const categories = ref(menuStore.categories)
const featuredItems = computed(() => menuStore.featuredItems)

const activePromotions = mockPromotions.filter(promo => promo.isActive)

onMounted(async () => {
  await menuStore.fetchCategories()
  await menuStore.fetchMenuItems()
  categories.value = menuStore.categories

  if (import.meta.client) {
    gsap.registerPlugin(ScrollTrigger)
    
    gsap.utils.toArray('.section-padding').forEach((section: any) => {
      gsap.from(section, {
        scrollTrigger: {
          trigger: section,
          start: 'top 80%',
          toggleActions: 'play none none reverse',
        },
        opacity: 0,
        y: 50,
        duration: 0.8,
        ease: 'power2.out',
      })
    })
  }
})

const handleClaimOffer = () => {
  navigateTo('/menu?promo=weekend-special')
}

const viewPromotionDetails = (promotion: Promotion) => {
  console.log('Viewing details:', promotion)
}
</script>