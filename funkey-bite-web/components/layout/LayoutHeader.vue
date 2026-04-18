<template>
  <header class="fixed top-0 left-0 right-0 z-50 h-[100px] bg-white/80 dark:bg-slate-900/80 backdrop-blur-xl border-b border-gray-200/50 dark:border-slate-800/50 shadow-lg shadow-gray-200/50 dark:shadow-slate-900/50 px-3 md:px-8 transition-all duration-300">
    <nav class="container-narrow h-full" role="navigation" aria-label="Main navigation">
      <div class="flex items-center justify-between h-full">

        <NuxtLink to="/" class="flex items-center space-x-3 group">
          <div class="w-10 h-10 bg-gradient-to-br from-brand-500 to-accent-500 rounded-xl flex items-center justify-center transform group-hover:scale-110 transition-transform duration-300 shadow-lg shadow-brand-500/20">
            <span class="text-white font-bold text-xl">F</span>
          </div>
          <div>
            <h1 class="text-xl font-bold text-gradient bg-clip-text md:text-2xl">Funkey Grab & Bite</h1>
            <p class="text-xs text-gray-600/80 dark:text-gray-400/80">Fast Food & Catering</p>
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

           <NuxtLink to="/order" aria-label="Go to order page">
             <button
               @click="cart.toggleCart()"
               class="relative p-2.5 rounded-xl bg-gradient-to-br from-brand-500 to-accent-500 hover:from-brand-600 hover:to-accent-600 transition-all duration-200 hover:scale-105 hover:shadow-lg shadow-brand-500/30"
               aria-label="Shopping cart"
             >
               <ShoppingCart class="w-5 h-5 text-white" />
               <span 
                 v-if="cart.totalItems > 0"
                 class="absolute -top-1 -right-1 bg-white text-brand-500 text-xs font-bold rounded-full w-5 h-5 flex items-center justify-center shadow-sm"
                 aria-label="Cart item count"
                 role="status"
               >
                 {{ cart.totalItems }}
               </span>
             </button>
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

      <Transition
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="transform -translate-y-4 opacity-0"
        enter-to-class="transform translate-y-0 opacity-100"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="transform translate-y-0 opacity-100"
        leave-to-class="transform -translate-y-4 opacity-0"
      >
        <div 
          v-if="isMobileMenuOpen"
          class="md:hidden mt-4 rounded-2xl bg-white/80 dark:bg-slate-900/80 backdrop-blur-xl border border-gray-200/50 dark:border-slate-800/50 shadow-xl shadow-gray-200/20 dark:shadow-slate-900/50 py-6 px-4"
          role="menu"
          aria-label="Mobile navigation menu"
        >
          <div class="flex flex-col space-y-2">
            <NuxtLink 
              v-for="nav in navigation" 
              :key="nav.name"
              :to="nav.href"
              @click="isMobileMenuOpen = false"
              class="nav-link text-gray-700 dark:text-gray-300 hover:text-brand-500 dark:hover:text-brand-400 font-medium py-3 px-6 rounded-xl hover:bg-white/50 dark:hover:bg-slate-800/50 transition-all duration-200 focus:outline-none focus-visible:ring-2 focus-visible:ring-brand-500"
              exact-active-class="nav-link--active"
              role="menuitem"
              tabindex="0"
            >
              {{ nav.name }}
              <span class="absolute inset-0 rounded-full pointer-events-none nav-link-highlight"></span>
            </NuxtLink>
          </div>
        </div>
      </Transition>
    </nav>
  </header>
</template>

<script setup lang="ts">
import { Menu, ShoppingCart, Sun, Moon } from 'lucide-vue-next'
import { useColorMode } from '@vueuse/core'
import { useCartStore } from '../../stores/cart'
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { watch } from 'vue'

const colorMode = useColorMode()
const cart = useCartStore()

const isMobileMenuOpen = ref(false)

const navigation = [
  { name: 'Home', href: '/' },
  { name: 'Menu', href: '/menu' },
  { name: 'Order', href: '/order' },
  { name: 'Catering', href: '/catering' },
  { name: 'Contact', href: '/contact' },
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