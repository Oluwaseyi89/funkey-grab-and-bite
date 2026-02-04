<template>
    <section class="section-padding mt-8 mb-8 md:mt-12 md:mb-12 mx-5">
      <div class="container-narrow">
        <div class="bg-gradient-to-r from-brand-500 to-accent-500 rounded-3xl p-8 md:p-12 relative overflow-hidden">
          <!-- Animated background elements -->
          <div class="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full -translate-y-32 translate-x-32"></div>
          <div class="absolute bottom-0 left-0 w-96 h-96 bg-white/5 rounded-full translate-y-48 -translate-x-48"></div>
  
          <div class="relative z-10 text-center text-white">
            <h2 class="text-3xl md:text-4xl font-bold mb-4">Weekend Special!</h2>
            <p class="text-xl mb-6 opacity-90">Get 20% off on all lunch packs this weekend</p>
            <div class="flex flex-col sm:flex-row gap-4 justify-center">
              <div class="bg-white/20 rounded-lg p-4 backdrop-blur-sm">
                <div class="text-2xl font-bold">20% OFF</div>
                <div class="text-sm opacity-90">All Lunch Packs</div>
              </div>
              <div class="bg-white/20 rounded-lg p-4 backdrop-blur-sm">
                <div class="text-2xl font-bold">BOGO</div>
                <div class="text-sm opacity-90">Shawarma Special</div>
              </div>
            </div>
            <button
              @click="$emit('claimOffer')"
              class="mt-8 bg-white text-brand-600 hover:bg-gray-100 font-bold py-3 px-8 rounded-xl transition-colors"
            >
              Claim Offer
            </button>
          </div>
        </div>
      </div>
    </section>
  </template>
  
  <script setup lang="ts">
  defineEmits(['claimOffer'])
  </script>