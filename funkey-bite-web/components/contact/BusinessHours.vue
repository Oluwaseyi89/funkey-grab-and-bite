<template>
    <div class="bg-white dark:bg-slate-800 rounded-xl p-6 border border-gray-200 dark:border-slate-700">
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Business Hours</h3>
      
      <div class="space-y-3">
        <div 
          v-for="day in businessHours" 
          :key="day.name"
          class="flex justify-between items-center py-2 border-b border-gray-100 dark:border-slate-700 last:border-0"
          :class="{ 'text-brand-500 font-bold': isToday(day.name) }"
        >
          <span class="text-gray-700 dark:text-gray-300">{{ day.name }}</span>
          <span class="font-medium">{{ day.hours }}</span>
        </div>
      </div>
      
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
  import { computed } from 'vue'
  
  const businessHours = [
    { name: 'Monday - Friday', hours: '8:00 AM - 10:00 PM' },
    { name: 'Saturday', hours: '9:00 AM - 11:00 PM' },
    { name: 'Sunday', hours: '10:00 AM - 9:00 PM' },
    { name: 'Holiday Hours', hours: '12:00 PM - 6:00 PM' }
  ]
  
  const isToday = (dayName: string) => {
    const today = new Date().getDay()
    const dayMap: Record<number, string> = {
      1: 'Monday - Friday',
      2: 'Monday - Friday',
      3: 'Monday - Friday',
      4: 'Monday - Friday',
      5: 'Monday - Friday',
      6: 'Saturday',
      0: 'Sunday'
    }
    return dayMap[today] === dayName
  }
  </script>