<template>
    <div>
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Order Type</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <button
          v-for="type in orderTypes"
          :key="type.id"
          @click="$emit('select', type.id)"
          class="p-6 border-2 rounded-xl text-center transition-all"
          :class="[
            selectedType === type.id
              ? 'border-brand-500 bg-brand-50 dark:bg-brand-900/20'
              : 'border-gray-200 dark:border-slate-700 hover:border-brand-300'
          ]"
        >
          <component :is="type.icon" class="w-10 h-10 mx-auto mb-3" :class="type.iconClass" />
          <h4 class="font-bold text-gray-900 dark:text-white mb-1">{{ type.name }}</h4>
          <p class="text-sm text-gray-600 dark:text-gray-400">{{ type.description }}</p>
          <div v-if="type.fee" class="mt-2 text-brand-500 font-medium">
            +${{ type.fee }} fee
          </div>
        </button>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { Store, Truck, Users } from 'lucide-vue-next'
  import type { OrderType } from '../../types/order'
  
  const orderTypes = [
    {
      id: 'pickup' as OrderType,
      name: 'Pickup',
      description: 'Order online, pick up in-store',
      icon: Store,
      iconClass: 'text-blue-500',
      fee: 0
    },
    {
      id: 'delivery' as OrderType,
      name: 'Delivery',
      description: 'Get it delivered to your door',
      icon: Truck,
      iconClass: 'text-green-500',
      fee: 3.99
    },
    {
      id: 'catering' as OrderType,
      name: 'Catering',
      description: 'For events & large orders',
      icon: Users,
      iconClass: 'text-purple-500',
      fee: 0
    }
  ]
  
  defineProps<{
    selectedType?: OrderType
  }>()
  
  defineEmits<{
    select: [type: OrderType]
  }>()
  </script>