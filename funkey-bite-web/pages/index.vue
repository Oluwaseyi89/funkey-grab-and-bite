<template>
  <div>
    <!-- Hero Banner -->
    <HeroBanner />

    <!-- Featured Categories -->
    <CategoriesGrid :categories="categories" />

    <!-- Featured Items -->
    <FeaturedItems :items="featuredItems" :is-loading="menuStore.isLoading" />

    <!-- Promotions Banner -->
    <!-- <PromotionsBanner @claim-offer="handleClaimOffer" /> -->
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

    <!-- Catering CTA -->
    <CateringCTA />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMenuStore } from '../stores/menu'

import { mockPromotions, mockMenuItems } from '../utils/mockData'
import type { Promotion } from '../types/menu'

// Components
import HeroBanner from '../components/home/HeroBanner.vue'
import CategoriesGrid from '../components/home/CategoriesGrid.vue'
import FeaturedItems from '../components/home/FeaturedItems.vue'
import PromotionsBanner from '../components/home/PromotionsBanner.vue'
import CateringCTA from '../components/home/CateringCTA.vue'

// GSAP
import gsap from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { navigateTo } from 'nuxt/app'

const menuStore = useMenuStore()
const categories = ref(menuStore.categories)
const featuredItems = computed(() => menuStore.featuredItems)

const activePromotions = mockPromotions.filter(promo => promo.isActive)


onMounted(async () => {
  // Fetch menu data
  await menuStore.fetchCategories()
  await menuStore.fetchMenuItems()
  categories.value = menuStore.categories

  // GSAP animations (client-side only)
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
  // Navigate to menu or show modal
  navigateTo('/menu?promo=weekend-special')
}


const viewPromotionDetails = (promotion: Promotion) => {
  console.log('Viewing details:', promotion)
  // Navigate to promotion details page or open modal
  // Example: navigateTo(`/promotions/${promotion.id}`)
}
</script>