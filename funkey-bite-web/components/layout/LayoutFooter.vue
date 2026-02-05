<template>
  <footer class="mt-auto bg-gray-50 dark:bg-slate-900 px-8 md:px-12 py-8 md:py-12">
    <!-- Main Footer -->
    <div class="container-narrow px-4">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-12 lg:gap-16">
        
        <!-- Brand & Description -->
        <div class="space-y-6">
          <NuxtLink to="/" class="flex items-center space-x-4 group">
            <div class="w-14 h-14 bg-gradient-to-br from-brand-500 to-accent-500 rounded-xl flex items-center justify-center transform group-hover:scale-110 transition-transform duration-300 shadow-lg">
              <span class="text-white font-bold text-2xl">F</span>
            </div>
            <div>
              <h2 class="text-2xl font-bold text-gray-900 dark:text-white leading-tight">Funkey Grab & Bite</h2>
              <p class="text-brand-500 dark:text-brand-300 text-sm mt-1">Fast Food & Catering</p>
            </div>
          </NuxtLink>
          <p class="text-gray-600 dark:text-gray-300 text-sm leading-relaxed">
            From crispy chicken & chips to authentic shawarma, we serve happiness in every bite. 
            Perfect for quick meals, lunch packs, and unforgettable catering experiences.
          </p>
          <div class="flex space-x-5 pt-2">
            <a href="#" class="text-gray-500 hover:text-brand-500 dark:text-gray-400 dark:hover:text-white transition-colors transform hover:scale-110 duration-200">
              <Facebook class="w-6 h-6" />
            </a>
            <a href="#" class="text-gray-500 hover:text-brand-500 dark:text-gray-400 dark:hover:text-white transition-colors transform hover:scale-110 duration-200">
              <Instagram class="w-6 h-6" />
            </a>
            <a href="#" class="text-gray-500 hover:text-brand-500 dark:text-gray-400 dark:hover:text-white transition-colors transform hover:scale-110 duration-200">
              <Twitter class="w-6 h-6" />
            </a>
          </div>
        </div>

        <!-- Quick Links -->
        <div>
          <h3 class="text-xl font-bold mb-6 text-gray-900 dark:text-white pb-2 border-b border-gray-200 dark:border-slate-800/50">Quick Links</h3>
          <ul class="space-y-4">
            <li v-for="link in quickLinks" :key="link.name">
              <NuxtLink 
                :to="link.href" 
                class="text-gray-700 hover:text-brand-500 dark:text-gray-300 dark:hover:text-brand-300 transition-colors flex items-center group py-2"
              >
                <div class="w-1.5 h-1.5 bg-brand-500 rounded-full mr-3 opacity-0 group-hover:opacity-100 transition-opacity"></div>
                <span class="group-hover:translate-x-1 transition-transform duration-200">
                  {{ link.name }}
                </span>
              </NuxtLink>
            </li>
          </ul>
        </div>

        <!-- Menu Categories -->
        <div>
          <h3 class="text-xl font-bold mb-6 text-gray-900 dark:text-white pb-2 border-b border-gray-200 dark:border-slate-800/50">Our Menu</h3>
          <ul class="space-y-4">
            <li v-for="category in mockCategories" :key="category.id" 
             @click="navigateToMenu(category.id)"
            class="text-gray-700 hover:text-brand-500 cursor-pointer dark:text-gray-300 dark:hover:text-brand-300 transition-colors flex items-center group py-2"
            >
                <div class="w-1.5 h-1.5 bg-brand-500 rounded-full mr-3 opacity-0 group-hover:opacity-100 transition-opacity"></div>
                <span class="group-hover:translate-x-1 transition-transform duration-200">
                  {{ category.name }}
                </span>
            </li>
          </ul>
        </div>

        <!-- Contact & Hours -->
        <div>
          <h3 class="text-xl font-bold mb-6 text-gray-900 dark:text-white pb-2 border-b border-gray-200 dark:border-slate-800/50">Contact Us</h3>
          <div class="space-y-5">
            <div class="flex items-start space-x-4">
              <MapPin class="w-5 h-5 text-brand-500 dark:text-brand-400 mt-1 flex-shrink-0" />
              <p class="text-gray-700 dark:text-gray-300 text-sm">123 Food Street, City Center, Your City</p>
            </div>
            <div class="flex items-center space-x-4">
              <Phone class="w-5 h-5 text-brand-500 dark:text-brand-400 flex-shrink-0" />
              <a href="tel:+1234567890" class="text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white transition-colors text-sm">
                (555) 123-4567
              </a>
            </div>
            <div class="flex items-center space-x-4">
              <Mail class="w-5 h-5 text-brand-500 dark:text-brand-400 flex-shrink-0" />
              <a href="mailto:orders@funkeygrabandbite.com" class="text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white transition-colors text-sm">
                orders@funkeygrabandbite.com
              </a>
            </div>
            
            <!-- Business Hours -->
            <div class="pt-6">
              <h4 class="font-bold mb-4 text-gray-800 dark:text-gray-200 text-lg">Opening Hours</h4>
              <div class="space-y-3 text-sm">
                <div class="flex justify-between items-center py-2 border-b border-gray-200 dark:border-slate-800/30">
                  <span class="text-gray-600 dark:text-gray-400">Mon - Fri:</span>
                  <span class="text-gray-800 dark:text-gray-300 font-medium">8:00 AM - 10:00 PM</span>
                </div>
                <div class="flex justify-between items-center py-2 border-b border-gray-200 dark:border-slate-800/30">
                  <span class="text-gray-600 dark:text-gray-400">Saturday:</span>
                  <span class="text-gray-800 dark:text-gray-300 font-medium">9:00 AM - 11:00 PM</span>
                </div>
                <div class="flex justify-between items-center py-2">
                  <span class="text-gray-600 dark:text-gray-400">Sunday:</span>
                  <span class="text-gray-800 dark:text-gray-300 font-medium">10:00 AM - 9:00 PM</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Divider -->
      <div class="border-t border-gray-200 dark:border-slate-800 my-12"></div>

      <!-- Bottom Bar -->
      <div class="flex flex-col md:flex-row justify-between items-center space-y-6 md:space-y-0">
        <div class="text-gray-600 dark:text-gray-500 text-sm text-center md:text-left">
          © {{ new Date().getFullYear() }} Funkey Grab & Bite. All rights reserved.
        </div>
        <div class="flex flex-wrap items-center justify-center gap-6 text-sm">
          <NuxtLink to="/privacy" class="text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white transition-colors hover:underline">
            Privacy Policy
          </NuxtLink>
          <div class="w-px h-4 bg-gray-300 dark:bg-slate-700"></div>
          <NuxtLink to="/terms" class="text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white transition-colors hover:underline">
            Terms of Service
          </NuxtLink>
          <div class="w-px h-4 bg-gray-300 dark:bg-slate-700"></div>
          <NuxtLink to="/catering" class="text-brand-500 hover:text-brand-600 dark:text-brand-300 dark:hover:text-brand-200 font-medium flex items-center group">
            <span class="mr-2">Book Catering</span>
            <ArrowRight class="w-4 h-4 transform group-hover:translate-x-1 transition-transform" />
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- WhatsApp Float Button -->
    <a 
      href="https://wa.me/1234567890" 
      target="_blank" 
      class="fixed bottom-8 right-8 bg-green-500 text-white p-4 rounded-full shadow-lg hover:bg-green-600 transition-all duration-300 z-40 hover:scale-110 hover:shadow-2xl"
      aria-label="Chat on WhatsApp"
    >
      <MessageCircle class="w-7 h-7" />
    </a>
  </footer>
</template>

<script setup lang="ts">
import { 
  Facebook, 
  Instagram, 
  Twitter, 
  MapPin, 
  Phone, 
  Mail, 
  ArrowRight,
  MessageCircle 
} from 'lucide-vue-next'

import { navigateTo } from 'nuxt/app';

import { mockCategories } from '../../utils/mockData';

const quickLinks = [
  { name: 'Home', href: '/' },
  { name: 'Menu', href: '/menu' },
  { name: 'Order Online', href: '/order' },
  { name: 'Catering Services', href: '/catering' },
  { name: 'About Us', href: '/about' },
  { name: 'Contact', href: '/contact' },
]

const navigateToMenu = (categoryId: string) => {
    navigateTo(`/menu?category=${categoryId}`)
  }
</script>