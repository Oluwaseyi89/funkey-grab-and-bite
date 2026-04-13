<template>
    <div class="space-y-8">
      <h2 class="text-3xl font-bold text-gray-900 dark:text-white text-center">
        Our <span class="text-gradient">Catering</span> Packages
      </h2>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div
          v-for="pkg in packages"
          :key="pkg.id"
          class="bg-white dark:bg-slate-800 rounded-2xl p-8 border-2 transition-all hover:scale-105 hover:shadow-2xl"
          :class="[
            pkg.popular 
              ? 'border-gray-200 dark:border-slate-700 shadow-xl' 
              : 'border-gray-200 dark:border-slate-700'
          ]"
        >
          
          <div v-if="pkg.popular" class="absolute -top-3 left-1/2 transform -translate-x-1/2">
            <span class="bg-brand-500 text-white px-4 py-1 rounded-full text-sm font-bold">
              Most Popular
            </span>
          </div>
          
          
          <div class="text-center mb-6">
            <component :is="pkg.icon" class="w-12 h-12 mx-auto mb-4" :class="pkg.iconColor" />
            <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">{{ pkg.name }}</h3>
            <div class="text-4xl font-bold text-brand-500 mb-2">
              ${{ pkg.pricePerPerson }}<span class="text-lg text-gray-500">/person</span>
            </div>
            <p class="text-gray-600 dark:text-gray-400">{{ pkg.description }}</p>
          </div>
          
          
          <ul class="space-y-3 mb-8">
            <li v-for="feature in pkg.features" :key="feature" class="flex items-center">
              <CheckCircle class="w-5 h-5 text-green-500 mr-3 flex-shrink-0" />
              <span class="text-gray-700 dark:text-gray-300">{{ feature }}</span>
            </li>
          </ul>
          
          
          <button
            @click="$emit('select', pkg.id)"
            class="w-full py-3 rounded-xl font-bold transition-colors"
            :class="[
              pkg.popular
                ? 'btn-primary'
                : 'bg-gray-100 dark:bg-slate-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-slate-600'
            ]"
          >
            {{ selectedPackage === pkg.id ? '✓ Selected' : 'Select Package' }}
          </button>
          
          
          <div class="text-center text-sm text-gray-500 dark:text-gray-400 mt-4">
            Minimum {{ pkg.minGuests }} guests
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { Users, Award, Sparkles, CheckCircle } from 'lucide-vue-next'
  
  interface CateringPackage {
    id: string
    name: string
    description: string
    pricePerPerson: number
    minGuests: number
    features: string[]
    popular: boolean
    icon: any
    iconColor: string
  }
  
  const packages: CateringPackage[] = [
    {
      id: 'standard',
      name: 'Standard',
      description: 'Perfect for office meetings and small gatherings',
      pricePerPerson: 18.99,
      minGuests: 10,
      features: [
        '3 Main Course Options',
        '2 Side Dishes',
        'Soft Drinks & Water',
        'Basic Table Setup',
        '4 Hours Service'
      ],
      popular: false,
      icon: Users,
      iconColor: 'text-blue-500'
    },
    {
      id: 'premium',
      name: 'Premium',
      description: 'Ideal for birthdays, anniversaries, and family events',
      pricePerPerson: 24.99,
      minGuests: 20,
      features: [
        '5 Main Course Options',
        '3 Side Dishes',
        'Appetizers & Desserts',
        'Premium Table Setup',
        'Professional Staff',
        '6 Hours Service'
      ],
      popular: true,
      icon: Award,
      iconColor: 'text-brand-500'
    },
    {
      id: 'executive',
      name: 'Executive',
      description: 'For corporate events and large celebrations',
      pricePerPerson: 34.99,
      minGuests: 50,
      features: [
        '7 Main Course Options',
        '5 Side Dishes',
        'Full Bar Service Available',
        'Custom Menu Design',
        'Dedicated Event Coordinator',
        '8+ Hours Service'
      ],
      popular: false,
      icon: Sparkles,
      iconColor: 'text-purple-500'
    }
  ]
  
  defineProps<{
    selectedPackage?: string
  }>()
  
  defineEmits<{
    select: [packageId: string]
  }>()
  </script>