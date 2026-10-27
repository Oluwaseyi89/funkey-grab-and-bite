<template>
    <div>
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Order Type</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
        <button
          v-for="type in orderTypes"
          :key="type.id"
          @click="$emit('select', type.id)"
          class="p-4 sm:p-6 border-2 rounded-xl text-center transition-all"
          :class="[
            selectedType === type.id
              ? 'border-brand-500 bg-brand-50 dark:bg-brand-900/20'
              : 'border-gray-200 dark:border-slate-700 hover:border-brand-300'
          ]"
        >
          <component :is="type.icon" class="w-10 h-10 mx-auto mb-3" :class="type.iconClass" />
          <h4 class="font-bold text-gray-900 dark:text-white mb-1">{{ type.name }}</h4>
          <p class="text-sm text-gray-600 dark:text-gray-400">{{ type.description }}</p>
          <div v-if="type.id === 'delivery' && deliveryFee" class="mt-2 text-brand-500 font-medium">
            +₦{{ deliveryFee.toLocaleString('en-NG') }} fee
          </div>
        </button>
      </div>
      <p class="mt-3 text-sm text-gray-600 dark:text-gray-400">
        Planning an event or a large order?
        <NuxtLink to="/catering" class="text-brand-500 hover:text-brand-600 font-medium">Request catering</NuxtLink>
      </p>
    </div>
  </template>
  
  <script setup lang="ts">
  import { Store, Truck } from 'lucide-vue-next'
  import type { OrderType } from '../../types/order'
  
  const orderTypes = [
    {
      id: 'pickup' as OrderType,
      name: 'Pickup',
      description: 'Order online, pick up in-store',
      icon: Store,
      iconClass: 'text-blue-500',
    },
    {
      id: 'delivery' as OrderType,
      name: 'Delivery',
      description: 'Get it delivered to your door',
      icon: Truck,
      iconClass: 'text-green-500',
    },
  ]
  
  // Catering has its own request flow on the catering page, so the cart only
  // offers pickup and delivery. The fee comes from the business settings.
  defineProps<{
    selectedType?: OrderType
    deliveryFee?: number
  }>()
  
  defineEmits<{
    select: [type: OrderType]
  }>()
  </script>