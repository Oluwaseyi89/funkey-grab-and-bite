<template>
  <section v-if="menu.isLoading || items.length" class="py-6 md:py-10" aria-labelledby="popular-title">
    <div class="container-narrow flex items-end justify-between gap-4">
      <div>
        <p class="eyebrow">Popular now</p>
        <h2 id="popular-title" class="section-title mt-1">Crowd favourites</h2>
      </div>
      <NuxtLink to="/menu" class="link inline-flex min-h-[44px] items-center gap-1 text-sm">
        Full menu <ArrowRight class="h-4 w-4" aria-hidden="true" />
      </NuxtLink>
    </div>

    <!-- A swipeable row on phones, a grid on larger screens. -->
    <div class="container-narrow mt-4">
      <ul class="scroll-row -mx-4 px-4 pb-4 sm:-mx-6 sm:px-6 lg:mx-0 lg:grid lg:grid-cols-4 lg:gap-5 lg:overflow-visible lg:px-0">
        <template v-if="menu.isLoading && !items.length">
          <li v-for="i in 4" :key="i" class="w-[16.5rem] shrink-0 lg:w-auto">
            <div class="card overflow-hidden">
              <div class="skeleton aspect-[4/3] rounded-none" />
              <div class="space-y-2 p-4">
                <div class="skeleton h-4 w-2/3" />
                <div class="skeleton h-3 w-full" />
                <div class="skeleton h-10 w-full rounded-full" />
              </div>
            </div>
          </li>
        </template>
        <li v-for="(item, index) in items" :key="item.id" class="flex w-[16.5rem] shrink-0 snap-start lg:w-auto">
          <MenuItemCard :item="item" layout="tile" :eager="index < 2" @open="selected = item" />
        </li>
      </ul>
    </div>

    <ItemSheet :item="selected" @close="selected = null" />
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { ArrowRight } from 'lucide-vue-next'
import MenuItemCard from '../menu/MenuItemCard.vue'
import ItemSheet from '../menu/ItemSheet.vue'
import { useMenuStore } from '../../stores/menu'
import type { MenuItem } from '../../types/menu'

const menu = useMenuStore()
const items = computed(() => menu.featuredItems)
const selected = ref<MenuItem | null>(null)
</script>
