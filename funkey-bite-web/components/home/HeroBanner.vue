<template>
    <div ref="heroRef" class="relative overflow-hidden bg-gradient-to-br px-5 py-5 from-brand-50 to-white dark:from-slate-900 dark:to-slate-800 rounded-3xl mx-4 md:mx-8 lg:mx-16 mb-8 md:mb-12">
      <div class="absolute inset-0">
        <div 
    class="absolute inset-0 bg-cover bg-center transition-all duration-1000 ease-in-out"
    :style="{ backgroundImage: `url(${heroBannerBackgroundImages[currentImageIndex]})` }"
  ></div>
  <div class="absolute inset-0 bg-black/20 dark:bg-black/40"></div>

        <div class="absolute top-0 left-0 w-64 h-64 bg-brand-100 dark:bg-brand-900/20 rounded-full -translate-x-32 -translate-y-32"></div>
        <div class="absolute bottom-0 right-0 w-96 h-96 bg-accent-100 dark:bg-accent-900/20 rounded-full translate-x-48 translate-y-48"></div>
      </div>
      
      <div class="relative z-10 section-padding">
        <div class="grid lg:grid-cols-2 gap-12 items-center">

          <div>
            <h1 ref="titleRef" class="text-4xl md:text-5xl lg:text-6xl font-bold text-gray-900 dark:text-white mb-6">
              Taste the<br>
              <span class="text-gradient">Funkey</span> Difference
            </h1>
            
            <p ref="descRef" class="text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-2xl">
              From crispy chicken & chips to authentic shawarma, we serve happiness in every bite.
              Perfect for quick meals, lunch packs, and unforgettable catering experiences.
            </p>
            
            <div ref="buttonsRef" class="flex flex-wrap gap-4">
              <NuxtLink to="/menu" class="btn-primary text-lg px-8 py-4">
                <ShoppingBag class="w-5 h-5 inline mr-2" />
                Order Now
              </NuxtLink>
              <NuxtLink to="/catering" class="btn-secondary text-lg px-8 py-4">
                <Calendar class="w-5 h-5 inline mr-2" />
                Book Catering
              </NuxtLink>
            </div>
            

            <div ref="statsRef" class="mt-12 grid grid-cols-3 gap-6">
              <div class="text-center">
                <div class="text-3xl font-bold text-brand-500 dark:text-brand-400">500+</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Happy Customers</div>
              </div>
              <div class="text-center">
                <div class="text-3xl font-bold text-brand-500 dark:text-brand-400">50+</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Menu Items</div>
              </div>
              <div class="text-center">
                <div class="text-3xl font-bold text-brand-500 dark:text-brand-400">24/7</div>
                <div class="text-sm text-gray-600 dark:text-gray-400">Online Orders</div>
              </div>
            </div>
          </div>
          

          <div ref="imageRef" class="relative">
            <div class="relative w-full h-64 md:h-96 lg:h-[500px]">

              <div class="absolute top-0 left-1/4 w-32 h-32 bg-white dark:bg-slate-800 rounded-2xl shadow-2xl transform -rotate-12 animate-float">
                <div class="p-4">
                  <Utensils class="w-16 h-16 mx-auto text-brand-500" />
                  <div class="text-center mt-2 font-semibold">Fresh Food</div>
                </div>
              </div>
              
              <div class="absolute top-1/3 right-0 w-28 h-28 bg-white dark:bg-slate-800 rounded-2xl shadow-2xl transform rotate-12 animate-float" style="animation-delay: 0.5s">
                <div class="p-3">
                  <Clock class="w-12 h-12 mx-auto text-green-500" />
                  <div class="text-center mt-1 text-sm font-semibold">Quick Service</div>
                </div>
              </div>
              
              <div class="absolute bottom-0 left-0 w-36 h-36 bg-white dark:bg-slate-800 rounded-2xl shadow-2xl transform -rotate-6 animate-float" style="animation-delay: 1s">
                <div class="p-4">
                  <Users class="w-16 h-16 mx-auto text-purple-500" />
                  <div class="text-center mt-2 font-semibold">Catering</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ShoppingBag, Calendar, Utensils, Clock, Users } from 'lucide-vue-next'
  import { useNuxtApp } from 'nuxt/app'
  import { ref } from 'vue'
  import { onMounted } from 'vue'
  import { heroBannerBackgroundImages } from '../../utils/mockData';
  
  const { $gsap, $animations } = useNuxtApp()
  const gsap = $gsap as typeof import('gsap')

  const heroRef = ref<HTMLElement>()
  const titleRef = ref<HTMLElement>()
  const descRef = ref<HTMLElement>()
  const buttonsRef = ref<HTMLElement>()
  const statsRef = ref<HTMLElement>()
  const imageRef = ref<HTMLElement>()
  const currentImageIndex = ref(0)

  
  onMounted(() => {
    if (import.meta.client) {

      const tl = gsap.timeline()
      tl.from(titleRef.value!, {
        opacity: 0,
        y: 50,
        duration: 0.8,
        ease: 'power2.out'
      })
      .from(descRef.value!, {
        opacity: 0,
        y: 30,
        duration: 0.6,
        ease: 'power2.out'
      }, '-=0.4')
      .from(buttonsRef.value?.children!, {
        opacity: 0,
        y: 20,
        duration: 0.5,
        stagger: 0.1,
        ease: 'power2.out'
      }, '-=0.3')
      .from(statsRef.value?.children!, {
        opacity: 0,
        y: 20,
        duration: 0.4,
        stagger: 0.1,
        ease: 'power2.out'
      }, '-=0.2')
      .from(imageRef.value?.children!, {
        opacity: 0,
        scale: 0.8,
        duration: 0.8,
        stagger: 0.2,
        ease: 'back.out(1.7)'
      }, '-=0.4')
  
      gsap.from(heroRef.value!, {
        scrollTrigger: {
          trigger: heroRef.value,
          start: 'top top',
          end: 'bottom top',
          scrub: 1,
        },
        y: 50,
      })
    }

  const interval = setInterval(() => {
    currentImageIndex.value = (currentImageIndex.value + 1) % heroBannerBackgroundImages.length
  }, 10000)
  
  return () => clearInterval(interval)
  })
  </script>

<style scoped>
    .bg-transition {
      transition: background-image 1s ease-in-out;
    }
</style>