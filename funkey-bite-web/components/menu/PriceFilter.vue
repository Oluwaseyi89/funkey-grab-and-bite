<template>
    <div class="space-y-4">
      <div class="px-2">
        <input
          type="range"
          :min="minPrice"
          :max="maxPrice"
          v-model="selectedMax"
          class="w-full h-2 bg-gray-200 dark:bg-slate-700 rounded-lg appearance-none cursor-pointer"
          @change="$emit('priceChange', selectedMax)"
        />
        <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mt-2">
          <span>₦{{ minPrice }}</span>
          <span class="font-bold text-brand-500">Up to ₦{{ selectedMax }}</span>
          <span>₦{{ maxPrice }}</span>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
    import { ref } from 'vue';


  defineProps<{
    minPrice: number
    maxPrice: number
  }>()
  
  const selectedMax = ref(30000)
  
  defineEmits<{
    priceChange: [maxPrice: number]
  }>()
  </script>