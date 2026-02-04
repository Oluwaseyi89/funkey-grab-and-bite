<template>
    <div>
      <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Event Type</h3>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <button
          v-for="event in eventTypes"
          :key="event.id"
          @click="$emit('select', event.id)"
          class="p-4 border-2 rounded-xl text-center transition-all hover:border-brand-300"
          :class="[
            selectedEvent === event.id
              ? 'border-brand-500 bg-brand-50 dark:bg-brand-900/20'
              : 'border-gray-200 dark:border-slate-700'
          ]"
        >
          <component :is="event.icon" class="w-8 h-8 mx-auto mb-2" :class="event.iconColor" />
          <span class="font-medium text-gray-900 dark:text-white">{{ event.name }}</span>
        </button>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { 
    Briefcase, 
    Users, 
    Heart, 
    GraduationCap,
    Cake,
    Music,
    Trophy,
    Church
  } from 'lucide-vue-next'
  
  const eventTypes = [
    { id: 'corporate', name: 'Corporate', icon: Briefcase, iconColor: 'text-blue-500' },
    { id: 'wedding', name: 'Wedding', icon: Heart, iconColor: 'text-red-500' },
    { id: 'birthday', name: 'Birthday', icon: Cake, iconColor: 'text-amber-500' },
    { id: 'family', name: 'Family', icon: Users, iconColor: 'text-green-500' },
    { id: 'graduation', name: 'Graduation', icon: GraduationCap, iconColor: 'text-purple-500' },
    { id: 'baby-shower', name: 'Baby Shower', icon: Church, iconColor: 'text-pink-500' },
    { id: 'sports', name: 'Sports Event', icon: Trophy, iconColor: 'text-orange-500' },
    { id: 'concert', name: 'Concert', icon: Music, iconColor: 'text-indigo-500' }
  ]
  
  defineProps<{
    selectedEvent?: string
  }>()
  
  defineEmits<{
    select: [eventId: string]
  }>()
  </script>