<template>
  <article
    class="group relative flex transition"
    :class="[
      layout === 'tile'
        ? 'card w-full flex-col overflow-hidden hover:-translate-y-0.5 hover:shadow-lift'
        : 'gap-4 py-4',
      stock.orderable ? '' : 'opacity-70',
    ]"
  >
    <!-- Tile: photo on top. Row: details left, photo right (as in food apps). -->
    <MenuImage
      v-if="layout === 'tile'"
      :src="item.imageUrl"
      :alt="item.name"
      :eager="eager"
      class="aspect-[4/3] w-full"
    >
      <div class="absolute left-3 top-3 flex flex-wrap gap-1.5">
        <span v-for="badge in badges" :key="badge.text" class="badge shadow-sm" :class="badge.class">{{ badge.text }}</span>
      </div>
    </MenuImage>

    <div class="flex min-w-0 flex-1 flex-col" :class="layout === 'tile' ? 'p-4' : ''">
      <h3 class="font-display text-base font-semibold leading-snug text-ink-900 dark:text-white">
        <!-- The whole card opens the item; this button gives it a name and keyboard access. -->
        <button type="button" class="text-left after:absolute after:inset-0 after:content-[''] focus-visible:outline-none" @click="$emit('open', item)">
          {{ item.name }}
        </button>
      </h3>
      <p v-if="item.description" class="mt-1 line-clamp-2 text-sm text-ink-500 dark:text-ink-400">{{ item.description }}</p>
      <div v-if="layout === 'row' && badges.length" class="mt-2 flex flex-wrap gap-1.5">
        <span v-for="badge in badges" :key="badge.text" class="badge" :class="badge.class">{{ badge.text }}</span>
      </div>
      <div class="mt-auto flex items-center justify-between gap-2 pt-3">
        <p class="font-bold tabular-nums text-ink-900 dark:text-white">{{ naira(item.price) }}</p>
        <p v-if="item.preparationTime" class="inline-flex items-center gap-1 text-xs text-ink-500 dark:text-ink-400">
          <Clock class="h-3.5 w-3.5" aria-hidden="true" />{{ item.preparationTime }} min
        </p>
      </div>
      <div v-if="layout === 'tile'" class="relative z-10 mt-3">
        <AddControl :item="item" :stock="stock" block />
      </div>
    </div>

    <div v-if="layout === 'row'" class="relative w-28 shrink-0 sm:w-32">
      <MenuImage :src="item.imageUrl" :alt="item.name" :eager="eager" class="aspect-square w-full rounded-2xl" />
      <div class="relative z-10 -mt-5 flex justify-center">
        <AddControl :item="item" :stock="stock" />
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Clock } from 'lucide-vue-next'
import MenuImage from './MenuImage.vue'
import AddControl from './AddControl.vue'
import { menuStock } from '../../utils/menuStock'
import { naira } from '../../utils/format'
import type { MenuItem } from '../../types/menu'

const props = withDefaults(defineProps<{
  item: MenuItem
  layout?: 'row' | 'tile'
  eager?: boolean
}>(), { layout: 'row', eager: false })

defineEmits<{ open: [item: MenuItem] }>()

const stock = computed(() => menuStock(props.item))

const badges = computed(() => {
  const list: { text: string; class: string }[] = []
  if (!stock.value.orderable) list.push({ text: stock.value.label, class: 'bg-ink-900 text-white dark:bg-white dark:text-ink-900' })
  else if (stock.value.tone === 'low') list.push({ text: stock.value.label, class: 'bg-accent-100 text-accent-700 dark:bg-accent-700/30 dark:text-accent-200' })
  if (props.item.isPreOrder) list.push({ text: 'Order ahead', class: 'bg-ink-100 text-ink-700 dark:bg-ink-800 dark:text-ink-200' })
  const tags = (props.item.tags ?? []).map((tag) => tag.toLowerCase())
  if (tags.includes('best seller') || tags.includes('popular')) list.push({ text: 'Popular', class: 'bg-brand-50 text-brand-700 dark:bg-brand-900/40 dark:text-brand-200' })
  if (tags.includes('spicy')) list.push({ text: 'Spicy', class: 'bg-brand-600 text-white' })
  return list
})
</script>
