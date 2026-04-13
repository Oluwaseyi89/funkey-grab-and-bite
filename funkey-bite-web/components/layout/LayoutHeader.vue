<template>
  <header class="fixed top-0 left-0 right-0 z-50 h-[100px] bg-white/80 dark:bg-slate-900/80 backdrop-blur-xl border-b border-gray-200/50 dark:border-slate-800/50 shadow-lg shadow-gray-200/50 dark:shadow-slate-900/50 px-3 md:px-8 transition-all duration-300">
    <nav class="container-narrow h-full">
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

        <div class="hidden md:flex items-center space-x-8">
          <NuxtLink 
            v-for="nav in navigation" 
            :key="nav.name"
            :to="nav.href"
            class="text-gray-700 dark:text-gray-300 hover:text-brand-500 dark:hover:text-brand-400 font-medium transition-all relative group py-2"
            active-class="text-brand-500 dark:text-brand-400"
          >
            {{ nav.name }}
            <span class="absolute -bottom-1 left-0 w-0 h-0.5 bg-gradient-to-r from-brand-500 to-accent-500 group-hover:w-full transition-all duration-300"></span>
          </NuxtLink>
        </div>

        <div class="flex items-center space-x-4">

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

           <NuxtLink to="/order">
             <button
               @click="cart.toggleCart()"
               class="relative p-2.5 rounded-xl bg-gradient-to-br from-brand-500 to-accent-500 hover:from-brand-600 hover:to-accent-600 transition-all duration-200 hover:scale-105 hover:shadow-lg shadow-brand-500/30"
               aria-label="Shopping cart"
             >
               <ShoppingCart class="w-5 h-5 text-white" />
               <span 
                 v-if="cart.totalItems > 0"
                 class="absolute -top-1 -right-1 bg-white text-brand-500 text-xs font-bold rounded-full w-5 h-5 flex items-center justify-center shadow-sm"
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
        >
          <div class="flex flex-col space-y-2">
            <NuxtLink 
              v-for="nav in navigation" 
              :key="nav.name"
              :to="nav.href"
              @click="isMobileMenuOpen = false"
              class="text-gray-700 dark:text-gray-300 hover:text-brand-500 dark:hover:text-brand-400 font-medium py-3 px-6 rounded-xl hover:bg-white/50 dark:hover:bg-slate-800/50 transition-all duration-200"
              active-class="text-brand-500 dark:text-brand-400 bg-brand-50/50 dark:bg-brand-900/20"
            >
              {{ nav.name }}
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

watch(() => useRoute().path, () => {
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
</style>