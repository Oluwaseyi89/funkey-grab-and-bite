<template>
  <nav
    class="fixed inset-x-0 bottom-0 z-50 border-t border-ink-200/70 bg-white/95 pb-safe backdrop-blur-xl dark:border-ink-800 dark:bg-ink-950/95 md:hidden"
    aria-label="Main"
  >
    <ul class="grid h-16 grid-cols-5">
      <li v-for="tab in tabs" :key="tab.label" class="flex">
        <component
          :is="tab.to ? NuxtLink : 'button'"
          v-bind="tab.to ? { to: tab.to } : { type: 'button', 'aria-haspopup': 'dialog' }"
          class="relative flex flex-1 flex-col items-center justify-center gap-0.5 text-[11px] font-medium transition active:scale-95"
          :class="isActive(tab) ? 'text-brand-600 dark:text-brand-400' : 'text-ink-500 dark:text-ink-400'"
          :aria-current="isActive(tab) ? 'page' : undefined"
          @click="tab.action?.()"
        >
          <span
            class="flex h-8 w-14 items-center justify-center rounded-full transition"
            :class="isActive(tab) ? 'bg-brand-50 dark:bg-brand-900/30' : ''"
          >
            <component :is="tab.icon" class="h-[22px] w-[22px]" :stroke-width="isActive(tab) ? 2.4 : 1.9" aria-hidden="true" />
          </span>
          {{ tab.label }}
          <span
            v-if="tab.badge"
            :key="cart.bump"
            class="absolute left-1/2 top-1 ml-2 inline-flex h-[18px] min-w-[18px] items-center justify-center rounded-full bg-brand-600 px-1 text-[10px] font-bold text-white animate-pop"
            aria-hidden="true"
          >{{ tab.badge }}</span>
          <span v-if="tab.badge" class="sr-only">, {{ tab.badge }} in cart</span>
        </component>
      </li>
    </ul>
  </nav>

  <BottomSheet :open="moreOpen" title="More" @close="moreOpen = false">
    <ul class="grid grid-cols-2 gap-3">
      <li v-for="link in moreLinks" :key="link.to">
        <NuxtLink :to="link.to" class="card flex h-full flex-col gap-2 p-4 transition active:scale-[0.98]" @click="moreOpen = false">
          <component :is="link.icon" class="h-6 w-6 text-brand-600 dark:text-brand-400" aria-hidden="true" />
          <span class="font-semibold text-ink-900 dark:text-white">{{ link.label }}</span>
          <span class="text-xs text-ink-500 dark:text-ink-400">{{ link.hint }}</span>
        </NuxtLink>
      </li>
    </ul>
    <div class="mt-4 flex items-center justify-between rounded-2xl bg-ink-100 px-4 py-3 dark:bg-ink-800">
      <span class="text-sm font-medium text-ink-800 dark:text-ink-100">Dark mode</span>
      <button
        type="button"
        role="switch"
        :aria-checked="isDark"
        class="relative h-7 w-12 rounded-full transition"
        :class="isDark ? 'bg-brand-600' : 'bg-ink-300'"
        @click="toggle"
      >
        <span class="sr-only">Dark mode</span>
        <span class="absolute top-1 h-5 w-5 rounded-full bg-white shadow transition-all" :class="isDark ? 'left-6' : 'left-1'" />
      </button>
    </div>
    <p class="mt-4 flex justify-center gap-4 text-sm text-ink-500 dark:text-ink-400">
      <NuxtLink to="/legal/privacy" class="hover:underline" @click="moreOpen = false">Privacy</NuxtLink>
      <NuxtLink to="/legal/terms" class="hover:underline" @click="moreOpen = false">Terms</NuxtLink>
    </p>
  </BottomSheet>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { NuxtLink } from '#components'
import {
  CalendarHeart, CircleUser, House, LayoutGrid, MessageCircle, Receipt, ShoppingBag, Tag, UtensilsCrossed,
} from 'lucide-vue-next'
import BottomSheet from '../ui/BottomSheet.vue'
import { useCartStore } from '../../stores/cart'
import { useAuthStore } from '../../stores/auth'
import { useTheme } from '../../composables/useTheme'

interface Tab {
  label: string
  icon: unknown
  to?: string
  match?: (path: string) => boolean
  action?: () => void
  badge?: number
}

const route = useRoute()
const cart = useCartStore()
const auth = useAuthStore()
const { isDark, toggle } = useTheme()
const moreOpen = ref(false)
const mounted = ref(false)
onMounted(() => { mounted.value = true })

watch(() => route.fullPath, () => { moreOpen.value = false })

const signedIn = computed(() => mounted.value && auth.isSignedIn)

const tabs = computed<Tab[]>(() => [
  { label: 'Home', icon: House, to: '/', match: (p) => p === '/' },
  { label: 'Menu', icon: UtensilsCrossed, to: '/menu', match: (p) => p.startsWith('/menu') },
  { label: 'Cart', icon: ShoppingBag, action: () => cart.openCart(), badge: mounted.value ? cart.totalItems : 0 },
  {
    label: 'Orders',
    icon: Receipt,
    to: signedIn.value ? '/account' : '/order/track',
    match: (p) => p.startsWith('/order/track') || p.startsWith('/account/orders') || p === '/account',
  },
  { label: 'More', icon: LayoutGrid, action: () => { moreOpen.value = true } },
])

const isActive = (tab: Tab) => (tab.label === 'Cart' ? cart.isOpen : tab.label === 'More' ? moreOpen.value : Boolean(tab.match?.(route.path)))

const moreLinks = computed(() => [
  { to: '/catering', label: 'Catering', hint: 'Events, parties & offices', icon: CalendarHeart },
  { to: '/promotion', label: 'Deals', hint: 'Codes you can use today', icon: Tag },
  { to: '/contact', label: 'Contact', hint: 'Call, email or visit', icon: MessageCircle },
  signedIn.value
    ? { to: '/account', label: 'My account', hint: 'Orders & details', icon: CircleUser }
    : { to: '/account/login', label: 'Sign in', hint: 'Optional, for faster checkout', icon: CircleUser },
])
</script>
