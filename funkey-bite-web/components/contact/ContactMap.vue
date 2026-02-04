<template>
    <div class="bg-white dark:bg-slate-800 rounded-xl overflow-hidden border border-gray-200 dark:border-slate-700">
      <!-- Map Container -->
      <div class="relative h-64 md:h-80 bg-gray-200 dark:bg-slate-700">
        <!-- Fallback map image -->
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="text-center">
            <MapPin class="w-12 h-12 text-brand-500 mx-auto mb-4" />
            <h3 class="font-bold text-gray-900 dark:text-white mb-2">Our Location</h3>
            <p class="text-gray-600 dark:text-gray-400">
              123 Food Street, City Center
            </p>
            <a 
              href="https://maps.google.com/?q=123+Food+Street+City+Center"
              target="_blank"
              class="inline-block mt-4 btn-primary text-sm"
            >
              Open in Google Maps
            </a>
          </div>
        </div>
        
        <!-- Map overlay for interactivity -->
        <div class="absolute inset-0 cursor-pointer" @click="openGoogleMaps"></div>
      </div>
      
      <!-- Location Details -->
      <div class="p-6">
        <h3 class="font-bold text-gray-900 dark:text-white mb-2">Getting Here</h3>
        <div class="space-y-3">
          <div class="flex items-start">
            <Train class="w-5 h-5 text-gray-500 mr-3 mt-1" />
            <div>
              <div class="font-medium">Public Transit</div>
              <div class="text-sm text-gray-600 dark:text-gray-400">
                5 min walk from Central Station (Red Line)
              </div>
            </div>
          </div>
          <div class="flex items-start">
            <Car class="w-5 h-5 text-gray-500 mr-3 mt-1" />
            <div>
              <div class="font-medium">Parking</div>
              <div class="text-sm text-gray-600 dark:text-gray-400">
                Street parking & paid lot available (2-hour limit)
              </div>
            </div>
          </div>
          <div class="flex items-start">
            <Accessibility class="w-5 h-5 text-gray-500 mr-3 mt-1" />
            <div>
              <div class="font-medium">Accessibility</div>
              <div class="text-sm text-gray-600 dark:text-gray-400">
                Wheelchair accessible entrance & restrooms
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { MapPin, Train, Car, Accessibility } from 'lucide-vue-next'
  
  const openGoogleMaps = () => {
    window.open('https://maps.google.com/?q=123+Food+Street+City+Center', '_blank')
  }
  </script>