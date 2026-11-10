<template>
  <section v-if="menu.isLoading || categories.length" class="container-narrow py-6 md:py-10" aria-labelledby="cravings-title">
    <h2 id="cravings-title" class="section-title">What are you craving?</h2>
    <ul class="scroll-row -mx-4 mt-4 px-4 pb-2 sm:-mx-6 sm:px-6 md:mx-0 md:grid md:grid-cols-6 md:gap-4 md:overflow-visible md:px-0">
      <template v-if="menu.isLoading && !categories.length">
        <li v-for="i in 5" :key="i" class="w-24 shrink-0 md:w-auto">
          <div class="skeleton mx-auto aspect-square w-full rounded-full" />
          <div class="skeleton mx-auto mt-2 h-3 w-2/3" />
        </li>
      </template>
      <li v-for="category in categories" :key="category.id" class="w-24 shrink-0 snap-start md:w-auto">
        <NuxtLink :to="`/menu#cat-${category.id}`" class="group block text-center">
          <MenuImage
            :src="category.photo"
            :alt="''"
            class="mx-auto aspect-square w-full rounded-full ring-4 ring-white transition group-hover:scale-105 group-hover:ring-brand-100 dark:ring-ink-900 dark:group-hover:ring-brand-900/50"
          />
          <span class="mt-2 block text-sm font-semibold leading-tight text-ink-800 group-hover:text-brand-600 dark:text-ink-100">
            {{ category.name }}
          </span>
          <span class="block text-xs text-ink-500 dark:text-ink-400">{{ plural(category.count, 'item') }}</span>
        </NuxtLink>
      </li>
    </ul>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import MenuImage from '../menu/MenuImage.vue'
import { useMenuStore } from '../../stores/menu'
import { plural } from '../../utils/format'

// Round category links with a photo from each category, linking straight to
// that part of the menu.
const menu = useMenuStore()

const categories = computed(() =>
  menu.visibleCategories.map((category) => {
    const items = menu.getItemsByCategory(category.id)
    return {
      id: category.id,
      name: category.name,
      count: items.length,
      photo: items.find((item) => item.imageUrl)?.imageUrl ?? '',
    }
  }),
)
</script>
