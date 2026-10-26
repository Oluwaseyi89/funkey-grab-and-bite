<template>
    <div class="space-y-6">
      <h3 class="text-2xl font-bold text-gray-900 dark:text-white">Contact Information</h3>
      
      <!-- Address -->
      <div class="flex items-start space-x-4 p-4 bg-white dark:bg-slate-800 rounded-xl border border-gray-200 dark:border-slate-700">
        <MapPin class="w-6 h-6 text-brand-500 mt-1 flex-shrink-0" />
        <div>
          <h4 class="font-bold text-gray-900 dark:text-white mb-1">Visit Our Location</h4>
          <p v-if="status === 'ready'" class="text-gray-600 dark:text-gray-400">
            <template v-for="(line, index) in addressLines" :key="index">
              {{ line }}<br v-if="index < addressLines.length - 1">
            </template>
          </p>
          <p v-else class="text-gray-600 dark:text-gray-400">
            {{ status === 'unavailable' ? 'Address unavailable right now' : 'Loading address…' }}
          </p>
          <a 
            v-if="mapsUrl"
            :href="mapsUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-block mt-2 text-brand-500 hover:text-brand-600 font-medium"
          >
            Get Directions →
          </a>
        </div>
      </div>
  
      <!-- Phone -->
      <div class="flex items-center space-x-4 p-4 bg-white dark:bg-slate-800 rounded-xl border border-gray-200 dark:border-slate-700">
        <Phone class="w-6 h-6 text-brand-500 flex-shrink-0" />
        <div>
          <h4 class="font-bold text-gray-900 dark:text-white mb-1">Call Us</h4>
          <a v-if="phoneHref" :href="phoneHref" class="text-gray-600 dark:text-gray-400 hover:text-brand-500">
            {{ settings?.phoneNumber }}
          </a>
          <p v-else class="text-gray-600 dark:text-gray-400">
            {{ status === 'unavailable' ? 'Phone unavailable right now' : 'Loading…' }}
          </p>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
            Available during our opening hours
          </p>
        </div>
      </div>
  
      <!-- Email -->
      <div class="flex items-center space-x-4 p-4 bg-white dark:bg-slate-800 rounded-xl border border-gray-200 dark:border-slate-700">
        <Mail class="w-6 h-6 text-brand-500 flex-shrink-0" />
        <div>
          <h4 class="font-bold text-gray-900 dark:text-white mb-1">Email Us</h4>
          <a v-if="settings?.email" :href="`mailto:${settings.email}`" class="text-gray-600 dark:text-gray-400 hover:text-brand-500">
            {{ settings.email }}
          </a>
          <p v-else class="text-gray-600 dark:text-gray-400">
            {{ status === 'unavailable' ? 'Email unavailable right now' : 'Loading…' }}
          </p>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
            For inquiries, orders and catering requests
          </p>
        </div>
      </div>
  
      <!-- Social Media -->
      <div class="p-4 bg-white dark:bg-slate-800 rounded-xl border border-gray-200 dark:border-slate-700">
        <h4 class="font-bold text-gray-900 dark:text-white mb-3">Follow Us</h4>
        <div class="flex space-x-4">
          <a 
            href="#" 
            class="w-10 h-10 bg-gray-100 dark:bg-slate-700 rounded-full flex items-center justify-center hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
            aria-label="Facebook"
          >
            <Facebook class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </a>
          <a 
            href="#" 
            class="w-10 h-10 bg-gray-100 dark:bg-slate-700 rounded-full flex items-center justify-center hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
            aria-label="Instagram"
          >
            <Instagram class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </a>
          <a 
            href="#" 
            class="w-10 h-10 bg-gray-100 dark:bg-slate-700 rounded-full flex items-center justify-center hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
            aria-label="Twitter"
          >
            <Twitter class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </a>
          <a 
            href="#" 
            class="w-10 h-10 bg-gray-100 dark:bg-slate-700 rounded-full flex items-center justify-center hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
            aria-label="TikTok"
          >
            <MessageCircle class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </a>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { 
    MapPin, 
    Phone, 
    Mail, 
    Facebook, 
    Instagram, 
    Twitter,
    MessageCircle
  } from 'lucide-vue-next'
  import { onMounted } from 'vue'
  import { useBusinessSettings } from '../../composables/useBusinessSettings'
  
  const { settings, status, addressLines, mapsUrl, phoneHref, load } = useBusinessSettings()
  
  onMounted(() => load())
  </script>