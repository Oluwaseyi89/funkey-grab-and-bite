<template>
  <header class="fixed top-0 left-0 right-0 z-50 h-[100px] bg-white/90 dark:bg-slate-900/90 backdrop-blur-lg border-b border-gray-200 dark:border-slate-800">
    <nav class="container-narrow h-full"> <!-- Make nav full height -->
      <div class="flex items-center justify-between h-full"> <!-- This is key -->
        <!-- Logo -->
        <NuxtLink to="/" class="flex items-center space-x-3 group">
          <div class="w-10 h-10 bg-gradient-to-br from-brand-500 to-accent-500 rounded-xl flex items-center justify-center transform group-hover:scale-110 transition-transform duration-300">
            <span class="text-white font-bold text-xl">F</span>
          </div>
          <div>
            <h1 class="text-2xl font-bold text-gradient">Funkey Grab & Bite</h1>
            <p class="text-xs text-gray-600 dark:text-gray-400">Fast Food & Catering</p>
          </div>
        </NuxtLink>

        <!-- Desktop Navigation -->
        <div class="hidden md:flex items-center space-x-8">
          <NuxtLink 
            v-for="nav in navigation" 
            :key="nav.name"
            :to="nav.href"
            class="text-gray-700 dark:text-gray-300 hover:text-brand-500 dark:hover:text-brand-400 font-medium transition-colors relative group py-2"
          >
            {{ nav.name }}
            <span class="absolute -bottom-1 left-0 w-0 h-0.5 bg-brand-500 group-hover:w-full transition-all duration-300"></span>
          </NuxtLink>
        </div>

        <!-- Right Side Actions -->
        <div class="flex items-center space-x-4">
          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="p-2 rounded-lg bg-gray-100 dark:bg-slate-800 hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors"
            aria-label="Toggle theme"
          >
            <ClientOnly>
              <Sun class="w-5 h-5 text-gray-700 dark:text-gray-300" v-if="colorMode === 'dark'" />
              <Moon class="w-5 h-5 text-gray-700 dark:text-gray-300" v-else />
            </ClientOnly>
          </button>

          <!-- Cart Button -->
          <button
            @click="cart.toggleCart()"
            class="relative p-2 rounded-lg bg-brand-500 hover:bg-brand-600 transition-colors"
            aria-label="Shopping cart"
          >
            <ShoppingCart class="w-5 h-5 text-white" />
            <span 
              v-if="cart.totalItems > 0"
              class="absolute -top-1 -right-1 bg-white text-brand-500 text-xs font-bold rounded-full w-5 h-5 flex items-center justify-center"
            >
              {{ cart.totalItems }}
            </span>
          </button>

          <!-- Mobile Menu Button -->
          <button
            @click="isMobileMenuOpen = !isMobileMenuOpen"
            class="md:hidden p-2 rounded-lg bg-gray-100 dark:bg-slate-800 hover:bg-gray-200 dark:hover:bg-slate-700 transition-colors"
            aria-label="Toggle mobile menu"
          >
            <Menu class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </button>
        </div>
      </div>

      <!-- Mobile Menu -->
      <Transition
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="transform opacity-0 scale-95"
        enter-to-class="transform opacity-100 scale-100"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="transform opacity-100 scale-100"
        leave-to-class="transform opacity-0 scale-95"
      >
        <div 
          v-if="isMobileMenuOpen"
          class="md:hidden border-t border-gray-200 dark:border-slate-800 py-4"
        >
          <div class="flex flex-col space-y-3">
            <NuxtLink 
              v-for="nav in navigation" 
              :key="nav.name"
              :to="nav.href"
              @click="isMobileMenuOpen = false"
              class="text-gray-700 dark:text-gray-300 hover:text-brand-500 dark:hover:text-brand-400 font-medium py-2 px-4 rounded-lg hover:bg-gray-100 dark:hover:bg-slate-800 transition-colors"
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
    // colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'
    colorMode.value = colorMode.value === 'dark' ? 'light' : 'dark'

  }
  
  // Close mobile menu on route change
  watch(() => useRoute().path, () => {
    isMobileMenuOpen.value = false
  })
  </script>