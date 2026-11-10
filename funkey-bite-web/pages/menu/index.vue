<template>
  <div>
    <div class="container-narrow pt-5 md:pt-10">
      <div class="flex flex-wrap items-end justify-between gap-3">
        <div>
          <p class="eyebrow">Pickup &amp; delivery</p>
          <h1 class="font-display text-3xl font-bold text-ink-900 dark:text-white md:text-5xl">Our menu</h1>
        </div>
        <StoreStatus />
      </div>
    </div>

    <!-- Search and category tabs stay under the header while scrolling. -->
    <div class="sticky top-14 z-30 mt-4 border-b border-ink-200/60 bg-ink-50/95 backdrop-blur-xl dark:border-ink-800 dark:bg-ink-950/95 md:top-16">
      <div class="container-narrow space-y-3 py-3">
        <div class="relative" role="search">
          <label for="menu-search" class="sr-only">Search the menu</label>
          <Search class="pointer-events-none absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-ink-400" aria-hidden="true" />
          <input
            id="menu-search"
            v-model="search"
            type="search"
            enterkeyhint="search"
            autocomplete="off"
            class="form-input rounded-full pl-12 pr-12"
            placeholder="Search chicken, shawarma, drinks…"
          />
          <button v-if="search" type="button" class="icon-btn absolute right-1 top-1/2 -translate-y-1/2" aria-label="Clear search" @click="search = ''">
            <X class="h-5 w-5" aria-hidden="true" />
          </button>
        </div>
        <nav v-if="!search && sections.length > 1" ref="tabsEl" class="scroll-row -mx-4 px-4 sm:-mx-6 sm:px-6 lg:mx-0 lg:px-0" aria-label="Menu sections">
          <a
            v-for="section in sections"
            :key="section.id"
            :href="`#${section.id}`"
            :data-section="section.id"
            class="chip snap-start"
            :class="activeSection === section.id ? 'chip-active' : ''"
            :aria-current="activeSection === section.id ? 'true' : undefined"
            @click.prevent="goTo(section.id)"
          >
            {{ section.title }}
            <span class="text-xs opacity-60">{{ section.items.length }}</span>
          </a>
        </nav>
      </div>
    </div>

    <div class="container-narrow py-4 md:py-8">
      <!-- Loading -->
      <div v-if="menu.isLoading && !menu.loaded" class="space-y-6" aria-busy="true" aria-label="Loading the menu">
        <div v-for="i in 4" :key="i" class="flex gap-4">
          <div class="flex-1 space-y-2 py-2">
            <div class="skeleton h-4 w-2/3" />
            <div class="skeleton h-3 w-full" />
            <div class="skeleton h-3 w-1/3" />
          </div>
          <div class="skeleton h-28 w-28 rounded-2xl" />
        </div>
      </div>

      <!-- Couldn't load -->
      <div v-else-if="menu.error" class="card mx-auto max-w-md p-8 text-center" role="alert">
        <WifiOff class="mx-auto h-10 w-10 text-ink-400" aria-hidden="true" />
        <p class="mt-3 font-semibold text-ink-900 dark:text-white">The menu didn’t load</p>
        <p class="mt-1 text-sm text-ink-500 dark:text-ink-400">{{ menu.error }}</p>
        <button type="button" class="btn-primary mt-5" @click="menu.load(true)">Try again</button>
      </div>

      <!-- Search results -->
      <section v-else-if="search" aria-live="polite">
        <h2 class="text-sm font-medium text-ink-500 dark:text-ink-400">
          {{ results.length ? `${plural(results.length, 'dish', 'dishes')} for “${search}”` : `Nothing matches “${search}”` }}
        </h2>
        <div v-if="results.length" :class="listClass">
          <MenuItemCard v-for="item in results" :key="item.id" :item="item" :layout="layout" @open="open" />
        </div>
        <div v-else class="py-12 text-center">
          <p class="text-ink-600 dark:text-ink-300">Try another word, or browse the full menu.</p>
          <button type="button" class="btn-secondary mt-4" @click="search = ''">Show everything</button>
        </div>
      </section>

      <!-- Nothing on the menu yet -->
      <div v-else-if="!sections.length" class="py-16 text-center">
        <p class="font-semibold text-ink-900 dark:text-white">The menu is being updated</p>
        <p class="mt-1 text-sm text-ink-500 dark:text-ink-400">Please check back soon.</p>
      </div>

      <!-- Full menu by section -->
      <template v-else>
        <section v-for="section in sections" :id="section.id" :key="section.id" :data-section="section.id" class="scroll-mt-[12.5rem] pb-6 pt-4 md:scroll-mt-[13rem] md:pb-10">
          <div class="mb-1 flex items-baseline justify-between gap-3">
            <h2 class="font-display text-xl font-bold text-ink-900 dark:text-white md:text-2xl">{{ section.title }}</h2>
            <span class="text-sm text-ink-500 dark:text-ink-400">{{ plural(section.items.length, 'item') }}</span>
          </div>
          <p v-if="section.description" class="text-sm text-ink-500 dark:text-ink-400">{{ section.description }}</p>
          <div :class="listClass">
            <MenuItemCard
              v-for="(item, index) in section.items"
              :key="`${section.id}-${item.id}`"
              :item="item"
              :layout="layout"
              :eager="section === sections[0] && index < 3"
              @open="open"
            />
          </div>
        </section>
      </template>
    </div>

    <ItemSheet :item="selected" @close="selected = null" />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMediaQuery } from '@vueuse/core'
import { Search, WifiOff, X } from 'lucide-vue-next'
import StoreStatus from '../../components/layout/StoreStatus.vue'
import MenuItemCard from '../../components/menu/MenuItemCard.vue'
import ItemSheet from '../../components/menu/ItemSheet.vue'
import { useMenuStore } from '../../stores/menu'
import { plural } from '../../utils/format'
import { usePageSeo } from '../../composables/usePageSeo'
import type { MenuItem } from '../../types/menu'

usePageSeo({
  title: 'Menu',
  description: 'Chicken & chips, shawarma, noodles, soups, drinks and lunch packs. Order for pickup or delivery.',
  path: '/menu',
})

const route = useRoute()
const router = useRouter()
const menu = useMenuStore()

const search = ref(typeof route.query.search === 'string' ? route.query.search : '')
const selected = ref<MenuItem | null>(null)
const activeSection = ref('')

// Rows on phones (photo on the right, like food apps), tiles from tablets up.
const wide = useMediaQuery('(min-width: 768px)')
const layout = computed(() => (wide.value ? 'tile' : 'row'))
const listClass = computed(() =>
  wide.value ? 'mt-4 grid gap-5 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4' : 'divide-y divide-ink-200/70 dark:divide-ink-800',
)

interface Section { id: string; title: string; description?: string; items: MenuItem[] }

const sections = computed<Section[]>(() => {
  const list: Section[] = []
  const popular = menu.featuredItems
  if (popular.length >= 3 && menu.visibleCategories.length > 1) {
    list.push({ id: 'popular', title: 'Popular', description: 'What everyone’s ordering', items: popular.slice(0, 6) })
  }
  for (const category of menu.visibleCategories) {
    list.push({
      id: `cat-${category.id}`,
      title: category.name,
      description: category.description,
      items: menu.getItemsByCategory(category.id),
    })
  }
  return list
})

const results = computed(() => {
  const query = search.value.trim().toLowerCase()
  if (!query) return []
  return menu.items.filter((item) =>
    [item.name, item.description, ...(item.tags ?? [])].some((text) => text?.toLowerCase().includes(query)),
  )
})

const open = (item: MenuItem) => { selected.value = item }

// Keep the search in the address so Back and shared links keep it.
watch(search, (value) => {
  router.replace({ query: { ...route.query, search: value || undefined }, hash: route.hash })
})

// Highlights the tab of the section being read.
const tabsEl = ref<HTMLElement>()
let observer: IntersectionObserver | undefined
const visible = new Map<string, number>()

// While a tapped tab scrolls the page, the tapped tab stays highlighted.
let lockedUntil = 0

const goTo = (id: string) => {
  activeSection.value = id
  lockedUntil = Date.now() + 1200
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  history.replaceState(history.state, '', `#${id}`)
}

watch(activeSection, async (id) => {
  await nextTick()
  tabsEl.value?.querySelector<HTMLElement>(`[data-section="${id}"]`)?.scrollIntoView({ inline: 'center', block: 'nearest', behavior: 'smooth' })
})

onMounted(async () => {
  observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        const id = (entry.target as HTMLElement).dataset.section ?? ''
        if (entry.isIntersecting) visible.set(id, entry.boundingClientRect.top)
        else visible.delete(id)
      }
      if (Date.now() < lockedUntil) return
      const first = sections.value.find((section) => visible.has(section.id))
      if (first) activeSection.value = first.id
    },
    // The band just below the sticky search and tabs.
    { rootMargin: '-210px 0px -50% 0px' },
  )

  await menu.load()
  await nextTick()
  document.querySelectorAll<HTMLElement>('section[data-section]').forEach((el) => observer?.observe(el))
  activeSection.value ||= sections.value[0]?.id ?? ''

  // Old links: ?category=3 and ?item=12; new ones use #cat-3.
  const category = route.query.category
  const target = typeof category === 'string' ? `cat-${category}` : route.hash.slice(1)
  if (target && document.getElementById(target)) goTo(target)
  const itemId = route.query.item
  if (typeof itemId === 'string') selected.value = menu.items.find((item) => item.id === itemId) ?? null
})

onBeforeUnmount(() => observer?.disconnect())
</script>
