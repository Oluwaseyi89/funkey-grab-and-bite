<template>
  <header class="fixed inset-x-0 top-0 z-50 border-b border-ink-200/60 bg-white/95 backdrop-blur-xl dark:border-ink-800/80 dark:bg-ink-950/95">
    <div class="container-narrow flex h-14 items-center gap-3 md:h-16">
      <NuxtLink to="/" class="-ml-1 flex h-9 shrink-0 items-center rounded-lg px-1 md:h-10" aria-label="Funkey Grab & Bite home">
        <BrandLogo alt="" />
      </NuxtLink>

      <nav class="ml-6 hidden items-center gap-1 md:flex" aria-label="Main">
        <NuxtLink
          v-for="link in links"
          :key="link.to"
          :to="link.to"
          class="rounded-full px-3.5 py-2 text-sm font-medium text-ink-600 transition hover:bg-ink-100 hover:text-ink-900 dark:text-ink-300 dark:hover:bg-ink-800 dark:hover:text-white"
          active-class="!bg-ink-900 !text-white dark:!bg-white dark:!text-ink-900"
        >
          {{ link.label }}
        </NuxtLink>
      </nav>

      <div class="ml-auto flex items-center gap-1.5">
        <StoreStatus class="mr-2 hidden lg:inline-flex" />

        <ClientOnly>
          <button type="button" class="icon-btn hidden md:inline-flex" :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'" @click="toggle">
            <Sun v-if="isDark" class="h-5 w-5" aria-hidden="true" />
            <Moon v-else class="h-5 w-5" aria-hidden="true" />
          </button>
        </ClientOnly>

        <NuxtLink
          :to="signedIn ? '/account' : '/account/login'"
          class="icon-btn md:w-auto md:gap-2 md:px-3"
          :aria-label="signedIn ? 'My account' : 'Sign in'"
        >
          <CircleUser class="h-5 w-5" aria-hidden="true" />
          <span class="hidden text-sm font-medium md:inline">{{ signedIn ? 'Account' : 'Sign in' }}</span>
        </NuxtLink>

        <button
          type="button"
          class="relative inline-flex h-11 min-w-11 items-center justify-center gap-2 rounded-full bg-ink-900 px-3.5 text-sm font-semibold text-white transition hover:bg-ink-800 active:scale-95 dark:bg-white dark:text-ink-900"
          :aria-label="cartLabel"
          @click="cart.openCart()"
        >
          <ShoppingBag class="h-5 w-5" aria-hidden="true" />
          <span v-if="count" class="tabular-nums">{{ count }}</span>
          <span v-else class="hidden sm:inline">Cart</span>
          <span v-if="count" class="hidden tabular-nums sm:inline">· {{ naira(cart.subtotal) }}</span>
          <span
            v-if="count"
            :key="cart.bump"
            class="absolute -right-0.5 -top-0.5 h-3 w-3 rounded-full border-2 border-white bg-brand-500 animate-pop dark:border-ink-950"
            aria-hidden="true"
          />
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { CircleUser, Moon, ShoppingBag, Sun } from 'lucide-vue-next'
import BrandLogo from './BrandLogo.vue'
import StoreStatus from './StoreStatus.vue'
import { useCartStore } from '../../stores/cart'
import { useAuthStore } from '../../stores/auth'
import { useTheme } from '../../composables/useTheme'
import { naira, plural } from '../../utils/format'

const cart = useCartStore()
const auth = useAuthStore()
const { isDark, toggle } = useTheme()

// The cart and session are restored in the browser, so they are only shown
// once the page is running there.
const mounted = ref(false)
onMounted(() => { mounted.value = true })
const count = computed(() => (mounted.value ? cart.totalItems : 0))
const signedIn = computed(() => mounted.value && auth.isSignedIn)
const cartLabel = computed(() =>
  count.value ? `Open cart, ${plural(count.value, 'item')}, ${naira(cart.subtotal)}` : 'Open cart, empty',
)

const links = [
  { to: '/menu', label: 'Menu' },
  { to: '/catering', label: 'Catering' },
  { to: '/promotion', label: 'Deals' },
  { to: '/order/track', label: 'Track order' },
  { to: '/contact', label: 'Contact' },
]
</script>
