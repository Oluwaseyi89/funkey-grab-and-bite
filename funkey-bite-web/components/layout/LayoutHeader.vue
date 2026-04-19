<template>
  <header class="fixed top-0 left-0 right-0 z-50 h-[100px] bg-white/80 dark:bg-slate-900/80 backdrop-blur-xl border-b border-gray-200/50 dark:border-slate-800/50 shadow-lg shadow-gray-200/50 dark:shadow-slate-900/50 px-3 md:px-8 transition-all duration-300">
    <nav class="container-narrow h-full" role="navigation" aria-label="Main navigation">
      <div class="flex items-center justify-between h-full">

        <NuxtLink to="/" class="flex items-center space-x-3 group">
          <div class="w-10 h-10 bg-gradient-to-br from-brand-500 to-accent-500 rounded-xl flex items-center justify-center transform group-hover:scale-110 transition-transform duration-300 shadow-lg shadow-brand-500/20">
            <span class="text-white font-bold text-xl">F</span>
          </div>
          <div>
            <span class="text-xl font-bold text-gradient bg-clip-text md:text-2xl">Funkey Grab & Bite</span>
            <p class="text-xs text-gray-600/80 hidden md:block dark:text-gray-400/80">Fast Food & Catering</p>
          </div>
        </NuxtLink>

        <div class="hidden md:flex items-center space-x-8" role="menubar">
          <NuxtLink 
            v-for="nav in navigation" 
            :key="nav.name"
            :to="nav.href"
            class="nav-link text-gray-700 dark:text-gray-300 hover:text-brand-500 dark:hover:text-brand-400 font-medium transition-all relative group py-2 px-3 rounded-lg focus:outline-none focus-visible:ring-2 focus-visible:ring-brand-500"
            exact-active-class="nav-link--active"
            role="menuitem"
            tabindex="0"
          >
            {{ nav.name }}
            <span class="absolute inset-0 rounded-full pointer-events-none nav-link-highlight"></span>
          </NuxtLink>
        </div>

        <div class="flex items-center space-x-4" role="group" aria-label="Header actions">

          <button
            @click="toggleTheme"
            class="p-2.5 rounded-xl bg-white/50 dark:bg-slate-800/50 hover:bg-white/80 dark:hover:bg-slate-700/80 backdrop-blur-sm border border-gray-200/50 dark:border-slate-700/50 transition-all duration-200 hover:scale-105 hover:shadow-lg"
            aria-label="Toggle theme"
          >
            <ClientOnly>
              <Sun class="w-5 h-5 text-amber-500 dark:text-amber-400 transition-transform duration-300" v-if="colorMode === 'dark'" />
              <Moon class="w-5 h-5 text-amber-500 dark:text-amber-400 transition-transform duration-300" v-else />
            </ClientOnly>
          </button>

           <NuxtLink
             to="/order"
             class="relative p-2.5 rounded-xl bg-gradient-to-br from-brand-500 to-accent-500 hover:from-brand-600 hover:to-accent-600 transition-all duration-200 hover:scale-105 hover:shadow-lg shadow-brand-500/30 inline-flex items-center justify-center"
             aria-label="Shopping cart"
           >
               <ShoppingCart class="w-5 h-5 text-white" />
               <span 
                 v-if="totalItems > 0"
                 class="absolute -top-1 -right-1 bg-white text-brand-500 text-xs font-bold rounded-full w-5 h-5 flex items-center justify-center shadow-sm"
                 aria-live="polite"
                 role="status"
                 :aria-label="`${totalItems} items in cart`"
               >
                 {{ totalItems }}
               </span>
           </NuxtLink>

          <button
            @click="isMobileMenuOpen = !isMobileMenuOpen"
            class="md:hidden p-2.5 rounded-xl bg-white/50 dark:bg-slate-800/50 hover:bg-white/80 dark:hover:bg-slate-700/80 backdrop-blur-sm border border-gray-200/50 dark:border-slate-700/50 transition-all duration-200 hover:scale-105"
            aria-label="Toggle mobile menu"
          >
            <Menu class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </button>
        </div>
      </div>

      <!-- Mobile Drawer & Overlay -->
      <transition name="fade">
        <div
          v-if="isMobileMenuOpen"
          class="fixed inset-0 z-40 md:hidden drawer-overlay"
          @click="isMobileMenuOpen = false"
          aria-label="Close mobile menu overlay"
        ></div>
      </transition>
      <transition name="drawer">
        <aside
          v-if="isMobileMenuOpen"
          class="drawer-panel fixed top-0 right-0 h-dvh w-4/5 max-w-[320px] z-50 md:hidden"
          role="menu"
          aria-label="Mobile navigation"
        >
          <!-- Top bar -->
          <div class="drawer-topbar">
            <div class="flex items-center gap-2.5">
              <div class="drawer-logo-mark">
                <span>F</span>
              </div>
              <span class="font-bold text-sm tracking-tight">Funkey Grab & Bite</span>
            </div>
            <button @click="isMobileMenuOpen = false" aria-label="Close" class="drawer-x">
              <X class="w-[18px] h-[18px]" />
            </button>
          </div>

          <!-- Nav -->
          <nav class="flex flex-col px-3 pt-3 pb-4 gap-0.5">
            <NuxtLink
              v-for="nav in navigation"
              :key="nav.name"
              :to="nav.href"
              @click="isMobileMenuOpen = false"
              class="drawer-item"
              :class="{ 'drawer-item--active': $route.path === nav.href }"
              role="menuitem"
            >
              <component :is="nav.icon" class="drawer-item-icon" />
              <span>{{ nav.name }}</span>
            </NuxtLink>
          </nav>

          <!-- Divider -->
          <div class="drawer-divider"></div>

          <!-- Bottom actions -->
          <div class="flex items-center gap-3 px-5 py-4">
            <button @click="toggleTheme" class="drawer-action-btn" aria-label="Toggle theme">
              <ClientOnly>
                <Sun class="w-4 h-4" v-if="colorMode === 'dark'" />
                <Moon class="w-4 h-4" v-else />
              </ClientOnly>
              <span>{{ colorMode === 'dark' ? 'Light mode' : 'Dark mode' }}</span>
            </button>
          </div>
        </aside>
      </transition>
    </nav>
  </header>
</template>

<style scoped>
/* transitions */
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.drawer-enter-active, .drawer-leave-active { transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
.drawer-enter-from, .drawer-leave-to { transform: translateX(100%); }

/* overlay */
.drawer-overlay {
  background: rgba(2, 6, 23, 0.6);
}

/* panel */
.drawer-panel {
  background: #ffffff;
  display: flex;
  flex-direction: column;
  box-shadow: -8px 0 40px rgba(2, 6, 23, 0.18);
}
.dark .drawer-panel { background: #111827; }

/* topbar */
.drawer-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1rem;
  height: 64px;
  border-bottom: 1px solid #f1f5f9;
  flex-shrink: 0;
  color: #111827;
}
.dark .drawer-topbar {
  border-bottom-color: #1f2937;
  color: #f9fafb;
}

/* logo mark */
.drawer-logo-mark {
  width: 2rem;
  height: 2rem;
  border-radius: 8px;
  background: linear-gradient(135deg, #f59e0b, #ef4444);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.drawer-logo-mark span {
  color: #fff;
  font-weight: 800;
  font-size: 0.875rem;
  line-height: 1;
}

/* close button */
.drawer-x {
  width: 2rem;
  height: 2rem;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
  transition: background 0.15s, color 0.15s;
  outline: none;
}
.drawer-x:hover { background: #f3f4f6; color: #111827; }
.dark .drawer-x { color: #9ca3af; }
.dark .drawer-x:hover { background: #1f2937; color: #f9fafb; }
.drawer-x:focus-visible { outline: 2px solid #f59e0b; outline-offset: 2px; }

/* nav items */
.drawer-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.8rem 0.875rem;
  border-radius: 10px;
  font-size: 0.9375rem;
  font-weight: 500;
  color: #374151;
  text-decoration: none;
  transition: background 0.15s, color 0.15s;
  outline: none;
  position: relative;
}
.dark .drawer-item { color: #9ca3af; }

.drawer-item:hover {
  background: #fef3c7;
  color: #b45309;
}
.dark .drawer-item:hover {
  background: #1f2937;
  color: #fbbf24;
}

.drawer-item--active {
  background: #fffbeb;
  color: #92400e;
  font-weight: 700;
}
.drawer-item--active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 25%;
  height: 50%;
  width: 3px;
  border-radius: 0 3px 3px 0;
  background: linear-gradient(180deg, #f59e0b, #ef4444);
}
.dark .drawer-item--active {
  background: #1c1917;
  color: #fbbf24;
}

.drawer-item-icon {
  width: 1.125rem;
  height: 1.125rem;
  flex-shrink: 0;
  opacity: 0.7;
}
.drawer-item--active .drawer-item-icon { opacity: 1; }

/* divider */
.drawer-divider {
  height: 1px;
  background: #f1f5f9;
  margin: 0 1rem;
}
.dark .drawer-divider { background: #1f2937; }

/* action button (theme toggle) */
.drawer-action-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.8125rem;
  font-weight: 500;
  color: #6b7280;
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  transition: background 0.15s, color 0.15s;
  outline: none;
}
.drawer-action-btn:hover { background: #f3f4f6; color: #374151; }
.dark .drawer-action-btn { color: #6b7280; }
.dark .drawer-action-btn:hover { background: #1f2937; color: #9ca3af; }
</style>

<script setup lang="ts">
import { Menu, ShoppingCart, Sun, Moon, Home, UtensilsCrossed, ShoppingBag, CalendarRange, Phone, X } from 'lucide-vue-next'
import { useColorMode } from '@vueuse/core'
import { useCartStore } from '../../stores/cart'
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { watch } from 'vue'
import { storeToRefs } from 'pinia'

const colorMode = useColorMode()
const cart = useCartStore()
const { totalItems } = storeToRefs(cart)

const isMobileMenuOpen = ref(false)

const navigation = [
  { name: 'Home', href: '/', icon: Home },
  { name: 'Menu', href: '/menu', icon: UtensilsCrossed },
  { name: 'Order', href: '/order', icon: ShoppingBag },
  { name: 'Catering', href: '/catering', icon: CalendarRange },
  { name: 'Contact', href: '/contact', icon: Phone },
]

const toggleTheme = () => {
  colorMode.value = colorMode.value === 'dark' ? 'light' : 'dark'
}

const route = useRoute()
watch(() => route.path, () => {
  isMobileMenuOpen.value = false
})
</script>

<style scoped>
.text-gradient {
  background: linear-gradient(135deg, #f59e0b 0%, #ef4444 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

header {
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
}

.glass-effect {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.7) 0%, rgba(255, 255, 255, 0.4) 100%);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
}

.dark .glass-effect {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.8) 0%, rgba(15, 23, 42, 0.6) 100%);
}
/* Enhanced nav link styles for active state */
/* Drool-worthy nav link highlight: animated border glow and soft shadow */


.nav-link--active {
  position: relative;
  background: transparent;
  color: #f44336 !important;
  font-weight: 700;
  z-index: 1;
  border-radius: 9999px;
  padding-left: 1.25rem !important;
  padding-right: 1.25rem !important;
  letter-spacing: 0.01em;
  transition: color 0.3s;
  overflow: visible;
}

/* Only show the animated highlight for the active link */
.nav-link--active .nav-link-highlight {
  z-index: 0;
  border-radius: 9999px;
  border: 2.5px solid;
  border-image: linear-gradient(90deg, #f44336 0%, #f59e0b 100%) 1;
  background: none !important;
  box-shadow: 0 0 0 6px #f4433622, 0 8px 32px 0 #f59e0b22;
  animation: nav-glow 2s infinite ease-in-out;
  display: block;
}
.nav-link-highlight {
  display: none;
}

@keyframes nav-glow {
  0% {
    box-shadow: 0 0 0 4px #f4433633, 0 8px 32px 0 #f59e0b26, 0 2px 8px 0 #f59e0b1a;
    border: 2px solid #f44336;
  }
  50% {
    box-shadow: 0 0 0 8px #f59e0b33, 0 12px 48px 0 #f4433626, 0 2px 8px 0 #f59e0b1a;
    border: 2px solid #f59e0b;
  }
  100% {
    box-shadow: 0 0 0 4px #f4433633, 0 8px 32px 0 #f59e0b26, 0 2px 8px 0 #f59e0b1a;
    border: 2px solid #f44336;
  }
}
.animate-nav-glow {
  z-index: 0;
  border-radius: 9999px;
  border: 2.5px solid;
  border-image: linear-gradient(90deg, #f44336 0%, #f59e0b 100%) 1;
  background: none !important;
  box-shadow: 0 0 0 6px #f4433622, 0 8px 32px 0 #f59e0b22;
  animation: nav-glow 2s infinite ease-in-out;
}

.nav-link--active .text-gradient {
  -webkit-text-fill-color: #fff;
}
</style>