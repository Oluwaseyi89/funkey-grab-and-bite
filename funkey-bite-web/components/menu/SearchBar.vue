<template>
    <div class="relative max-w-2xl mx-auto w-full">
      <Search class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search menu items..."
        class="w-full pl-12 pr-4 py-4 bg-white dark:bg-slate-800 border border-gray-300 dark:border-slate-600 rounded-xl focus:ring-2 focus:ring-brand-500 focus:border-transparent outline-none transition-all"
        @input="$emit('search', searchQuery)"
      />
      <button
        v-if="searchQuery"
        @click="clearSearch"
        class="absolute right-4 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
      >
        <X class="w-5 h-5" />
      </button>
    </div>
  </template>
  
  <script setup lang="ts">
  import { Search, X } from 'lucide-vue-next'
  import { ref } from 'vue'
  
  const searchQuery = ref('')
  
  // Define emits first
  const emit = defineEmits<{
    search: [query: string]
  }>()
  
  const clearSearch = () => {
    searchQuery.value = ''
    emit('search', '') // ✅ Now emit is defined
  }
  </script>