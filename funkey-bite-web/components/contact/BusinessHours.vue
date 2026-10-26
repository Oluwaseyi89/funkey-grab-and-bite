<template>
    <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700">
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Business Hours</h3>
      
      <div v-if="status === 'ready'" class="space-y-3">
        <div 
          v-for="group in hoursGroups" 
          :key="group.label"
          class="flex justify-between items-center py-2 border-b border-gray-100 dark:border-slate-700 last:border-0"
          :class="{ 'text-brand-500 font-bold': group.days.includes(today) }"
        >
          <span class="text-gray-700 dark:text-gray-300">{{ group.label }}</span>
          <span class="font-medium text-right">{{ group.hours }}</span>
        </div>
        <p class="pt-2 text-sm" :class="settings?.isOpenNow ? 'text-green-600 dark:text-green-400' : 'text-gray-500 dark:text-gray-400'">
          {{ settings?.isOpenNow ? 'Open now' : 'Closed now' }} · Times shown in {{ timezoneLabel }}
        </p>
      </div>
      <p v-else-if="status === 'unavailable'" class="text-gray-600 dark:text-gray-400">
        Our opening hours are unavailable right now. Please call us to confirm.
      </p>
      <p v-else class="text-gray-500 dark:text-gray-400">Loading opening hours…</p>
      
      <!-- Holiday Notice -->
      <div class="mt-6 p-3 bg-amber-50 dark:bg-amber-900/20 rounded-lg">
        <div class="flex items-center">
          <AlertCircle class="w-5 h-5 text-amber-500 mr-2" />
          <span class="text-sm text-amber-700 dark:text-amber-300">
            Closed on major holidays. Catering available 24/7 by appointment.
          </span>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { AlertCircle } from 'lucide-vue-next'
  import { onMounted } from 'vue'
  import { useBusinessSettings } from '../../composables/useBusinessSettings'
  
  const { settings, status, hoursGroups, today, timezoneLabel, load } = useBusinessSettings()
  
  onMounted(() => load())
  </script>