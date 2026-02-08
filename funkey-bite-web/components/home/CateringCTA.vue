<template>
    <section class="section-padding mt-8 mb-8 md:mt-12 mx-5">
      <div class="container-narrow">
        <div class="grid lg:grid-cols-2 gap-12 items-center">
          <div>
            <h2 class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white mb-6">
              Plan Your <span class="text-gradient">Perfect</span> Event
            </h2>
            <p class="text-gray-600 dark:text-gray-300 mb-6 text-lg">
              From corporate meetings to family gatherings, our catering service brings delicious food and professional service to your event.
            </p>
            <ul class="space-y-4 mb-8">
              <li v-for="feature in features" :key="feature" class="flex items-center">
                <CheckCircle class="w-6 h-6 text-green-500 mr-3" />
                <span>{{ feature }}</span>
              </li>
            </ul>
            <NuxtLink to="/catering" class="btn-primary text-lg px-8 py-4">
              Request Catering Quote
              <ArrowRight class="w-5 h-5 inline ml-2" />
            </NuxtLink>
          </div>
          <div class="relative">
            <div class="bg-gradient-to-br from-brand-100 to-accent-100 dark:from-brand-900/30 dark:to-accent-900/20 rounded-3xl p-8">
              <img
                src="/images/catering_event"
                alt="Catering Event"
                class="rounded-2xl shadow-2xl w-full h-auto"
                loading="lazy"
              />
            </div>
          </div>
        </div>
      </div>
    </section>
  </template>
  
  <script setup lang="ts">
  import { CheckCircle, ArrowRight } from 'lucide-vue-next'
  
  defineProps<{
    features?: string[]
    imageUrl?: string
  }>()
  
  // Default values
  const defaultFeatures = [
    'Indoor & Outdoor Catering',
    'Customized Menu Options',
    'Professional Staff & Setup'
  ]
  
  const defaultImageUrl = 'https://funkey-static-assets.s3.amazonaws.com/catering/event.jpg'
  </script>