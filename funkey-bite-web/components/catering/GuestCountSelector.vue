<template>
    <div>
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">
        Number of Guests
        <span class="text-sm font-normal text-gray-500">({{ guestCount }})</span>
      </h3>
      
      <div class="space-y-6">
        
        <div class="px-2">
          <input
            type="range"
            :min="minGuests"
            :max="maxGuests"
            v-model="guestCount"
            class="w-full h-2 bg-gray-200 dark:bg-slate-700 rounded-lg appearance-none cursor-pointer"
            @input="handleSliderInput"
          />
          <div class="flex justify-between text-sm text-gray-600 dark:text-gray-400 mt-2">
            <span>{{ minGuests }} guests</span>
            <span class="font-bold text-brand-500">{{ guestCount }} guests</span>
            <span>{{ maxGuests }}+ guests</span>
          </div>
        </div>
        
        
        <div class="flex flex-wrap gap-2">
          <button
            v-for="count in quickCounts"
            :key="count"
            @click="handleQuickSelect(count)"
            class="px-4 py-2 rounded-lg transition-colors"
            :class="[
              guestCount === count
                ? 'bg-brand-500 text-white'
                : 'bg-gray-100 dark:bg-slate-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-slate-600'
            ]"
          >
            {{ count }} guests
          </button>
        </div>
        
        
        <div class="bg-gray-50 dark:bg-slate-700 p-4 rounded-xl">
          <div class="flex justify-between items-center">
            <div>
              <div class="text-gray-600 dark:text-gray-400">Estimated Cost</div>
              <div class="text-2xl font-bold text-brand-500">
                ${{ estimatedCost.toFixed(0) }}<span class="text-lg text-gray-500">+</span>
              </div>
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400 text-right">
              Based on {{ selectedPackage?.pricePerPerson || 25 }}/person
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, watch } from 'vue'
  
  interface PackageDetails {
    pricePerPerson: number
    minGuests: number
  }
  
  const props = withDefaults(defineProps<{
    minGuests?: number
    maxGuests?: number
    selectedPackage?: PackageDetails
  }>(), {
    minGuests: 10,
    maxGuests: 500
  })
  
  const emit = defineEmits<{
    update: [count: number]
  }>()
  
  const guestCount = ref(props.selectedPackage?.minGuests || 50)
  const quickCounts = [20, 50, 100, 200]
  
  watch(() => props.selectedPackage, (newPackage) => {
    if (newPackage) {
      guestCount.value = newPackage.minGuests
      emitUpdate()
    }
  })
  
  const estimatedCost = computed(() => {
    const pricePerPerson = props.selectedPackage?.pricePerPerson || 25
    return guestCount.value * pricePerPerson
  })
  
  const handleSliderInput = (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target) {
      const value = parseInt(target.value)
      guestCount.value = value
      emitUpdate()
    }
  }
  
  const handleQuickSelect = (count: number) => {
    guestCount.value = count
    emitUpdate()
  }
  
  const emitUpdate = () => {
    emit('update', guestCount.value)
  }
  </script>