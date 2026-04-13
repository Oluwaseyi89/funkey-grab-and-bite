<template>
    <div class="mx-5 px-3 py-3 md:mx-8 md:py-8 md:px-8">
     
  <PageHeader
    title-before="Catering for"
    highlight-text="Every"
    title-after="Occasion"
    subtitle="From corporate meetings to wedding celebrations, we bring delicious food and professional service to your event. Serving 10 to 500+ guests."
    variant="solid"
    solid-color="brand-50"
    solid-color-dark="brand-600"
    :narrow="true"
    alignment="center"
    header-class="dark:from-brand-500 dark:to-accent-500"
    title-class="text-gray-900 dark:text-white"
    subtitle-class="text-gray-700 dark:text-white/90"
  >

      <div class="flex flex-col sm:flex-row gap-4 justify-center mt-8">
        <button 
          @click="scrollToForm"
          class="text-gray-900 dark:text-white font-bold py-3 px-8 rounded-xl transition-colors hover:scale-105 shadow-lg hover:shadow-xl"
        >
          Request Quote
        </button>
        <a 
          href="tel:+1234567890" 
          class="bg-white/80 dark:bg-white/20 hover:bg-white dark:hover:bg-white/30 backdrop-blur-sm font-medium py-3 px-8 rounded-xl transition-colors hover:scale-105 flex items-center justify-center gap-2 border border-brand-200 dark:border-white/30 text-brand-700 dark:text-white"
        >
          <Phone class="w-5 h-5" />
          (555) 123-4567
        </a>
      </div>
    </PageHeader>

      <section class="section-padding px-3 py-3 md:py-8 md:px-8 my-3 md:my-8 bg-transparent dark:bg-slate-900">
        <div class="container-narrow">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white text-center mb-12">
            Why Choose <span class="text-gradient">Funkey</span> Catering
          </h2>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div v-for="feature in features" :key="feature.title" class="text-center">
              <div class="w-16 h-16 bg-brand-100 dark:bg-brand-900/30 rounded-full flex items-center justify-center mx-auto mb-6">
                <component :is="feature.icon" class="w-8 h-8 text-brand-500" />
              </div>
              <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-3">{{ feature.title }}</h3>
              <p class="text-gray-600 dark:text-gray-400">{{ feature.description }}</p>
            </div>
          </div>
        </div>
      </section>
  

      <section class="section-padding mx-5 md:mx-8 my-5 md:my-8 px-3 md:px-8 py-3 md:py-8 bg-transparent">
        <div class="container-narrow">
          <CateringPackages
            :selected-package="selectedPackage"
            @select="selectedPackage = $event"
          />
        </div>
      </section>
  

      <section ref="formSection" class="section-padding my-5 mx-5 md:my-8 md:mx-8 px-3 py-3 md:py-8 md:px-8 bg-white dark:bg-slate-900">
        <div class="container-narrow">
          <div class="max-w-4xl mx-auto">
            <div class="grid lg:grid-cols-2 gap-12">

              <div class="space-y-8">
                <div>
                  <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-4">Event Details</h2>
                  <EventTypeSelector
                    :selected-event="selectedEvent"
                    @select="selectedEvent = $event"
                  />
                </div>
                
                <GuestCountSelector
                  :selected-package="getPackageDetails(selectedPackage)"
                  @update="guestCount = $event"
                />
                

                <div class="bg-brand-50 dark:brand-900/20 p-6 rounded-xl">
                  <h4 class="font-bold text-brand-700 dark:text-brand-300 mb-3">📍 Service Areas</h4>
                  <div class="grid grid-cols-2 gap-2 text-sm">
                    <div class="flex items-center">
                      <CheckCircle class="w-4 h-4 text-green-500 mr-2" />
                      <span>Downtown</span>
                    </div>
                    <div class="flex items-center">
                      <CheckCircle class="w-4 h-4 text-green-500 mr-2" />
                      <span>North Side</span>
                    </div>
                    <div class="flex items-center">
                      <CheckCircle class="w-4 h-4 text-green-500 mr-2" />
                      <span>South Side</span>
                    </div>
                    <div class="flex items-center">
                      <CheckCircle class="w-4 h-4 text-green-500 mr-2" />
                      <span>West End</span>
                    </div>
                  </div>
                  <p class="text-sm text-gray-600 dark:text-gray-400 mt-3">
                    Additional travel fees may apply for locations outside our primary service area.
                  </p>
                </div>
              </div>
  

              <div>
                <CateringRequestForm
                  :guest-count="guestCount"
                  :selected-package="selectedPackage"
                  :selected-event="selectedEvent"
                  :is-submitting="isSubmitting"
                  @submit="handleSubmit"
                />
              </div>
            </div>
          </div>
        </div>
      </section>
  

      <section class="section-padding mx-3 my-3 md:mx-8 md:my-8 py-3 px-3 md:py-8 md:px-8 bg-gray-50 dark:bg-slate-800">
        <div class="container-narrow">
          <CateringFAQ />
        </div>
      </section>
  

      <section class="section-padding px-5 py-5 md:py-8 md:px-8 mx-3 my-3 md:mx-8 md:my-8 bg-white dark:bg-slate-900">
        <div class="container-narrow">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white text-center mb-12">
            What Our <span class="text-gradient">Clients</span> Say
          </h2>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div 
              v-for="testimonial in testimonials" 
              :key="testimonial.name"
              class="bg-gray-50 dark:bg-slate-800 p-6 rounded-xl"
            >
              <div class="flex items-center mb-4">
                <div class="w-12 h-12 bg-brand-500 rounded-full flex items-center justify-center text-white font-bold mr-4">
                  {{ testimonial.initials }}
                </div>
                <div>
                  <h4 class="font-bold text-gray-900 dark:text-white">{{ testimonial.name }}</h4>
                  <p class="text-sm text-gray-600 dark:text-gray-400">{{ testimonial.event }}</p>
                </div>
              </div>
              <p class="text-gray-700 dark:text-gray-300 italic">"{{ testimonial.quote }}"</p>
              <div class="flex text-amber-500 mt-3">
                <Star v-for="i in 5" :key="i" class="w-4 h-4 fill-current" />
              </div>
            </div>
          </div>
        </div>
      </section>
  

      <section class="bg-gradient-to-r from-brand-600 to-accent-600 px-5 py-5 md:px-8 md:py-8 mx-3 my-3 md:my-8 md:mx-8 text-white">
        <div class="container-narrow section-padding text-center">
          <h2 class="text-3xl font-bold mb-4">Ready to Plan Your Event?</h2>
          <p class="text-xl opacity-90 mb-8 max-w-2xl mx-auto">
            Contact us today for a personalized quote and let's create an unforgettable experience.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <button 
              @click="scrollToForm"
              class="bg-white text-brand-600 hover:bg-gray-100 font-bold py-3 px-8 rounded-xl transition-colors"
            >
              Get Started
            </button>
            <a href="tel:+1234567890" class="bg-white/20 hover:bg-white/30 backdrop-blur-sm font-medium py-3 px-8 rounded-xl transition-colors">
              Call Now: (555) 123-4567
            </a>
          </div>
        </div>
      </section>
  

      <div 
        v-if="showSuccessModal"
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white dark:bg-slate-800 rounded-2xl max-w-md w-full p-8 text-center">
          <div class="w-20 h-20 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-6">
            <CheckCircle class="w-10 h-10 text-green-600 dark:text-green-400" />
          </div>
          
          <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Request Sent!</h2>
          <p class="text-gray-600 dark:text-gray-400 mb-6">
            Thank you for your catering request. We'll contact you within 24 hours to discuss your event details.
          </p>
          
          <div class="space-y-3">
            <button @click="showSuccessModal = false" class="w-full btn-primary">
              Continue Browsing
            </button>
            <button @click="resetForm" class="w-full btn-secondary">
              Submit Another Request
            </button>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue'
  import type { CateringRequest } from '../../types/order'
  import { 
    CheckCircle, 
    Star,
    Clock,
    Award,
    Users,
    Shield,
    Truck
  } from 'lucide-vue-next'
  
  import CateringPackages from '../../components/catering/CateringPackages.vue'
  import EventTypeSelector from '../../components/catering/EventTypeSelector.vue'
  import GuestCountSelector from '../../components/catering/GuestCountSelector.vue'
  import CateringRequestForm from '../../components/catering/CateringRequestForm.vue'
  import CateringFAQ from '../../components/catering/CateringFAQ.vue'
  import { Phone } from 'lucide-vue-next'

  
  const formSection = ref<HTMLElement>()
  const selectedPackage = ref('premium')
  const selectedEvent = ref('corporate')
  const guestCount = ref(50)
  const isSubmitting = ref(false)
  const showSuccessModal = ref(false)
  
  const features = [
    {
      icon: Award,
      title: 'Award-Winning Food',
      description: 'Our chefs create delicious, fresh meals that impress every time.'
    },
    {
      icon: Clock,
      title: 'On-Time Delivery',
      description: 'We guarantee timely setup and service for your event schedule.'
    },
    {
      icon: Shield,
      title: 'Fully Insured',
      description: 'Professional liability insurance for your peace of mind.'
    },
    {
      icon: Users,
      title: 'Professional Staff',
      description: 'Trained, uniformed servers and coordinators.'
    },
    {
      icon: Truck,
      title: 'Full Setup & Cleanup',
      description: 'We handle everything from setup to post-event cleanup.'
    }
  ]
  
  const testimonials = [
    {
      name: 'Sarah Johnson',
      initials: 'SJ',
      event: 'Wedding Reception',
      quote: 'The food was incredible! Our guests are still talking about it. Professional service from start to finish.'
    },
    {
      name: 'Michael Chen',
      initials: 'MC',
      event: 'Corporate Conference',
      quote: 'Fed 200 attendees seamlessly. The team was efficient and the food was a hit with everyone.'
    },
    {
      name: 'Jessica Williams',
      initials: 'JW',
      event: 'Birthday Party',
      quote: 'Exceptional service! They accommodated all our dietary needs and made the party stress-free.'
    }
  ]
  
  const scrollToForm = () => {
    formSection.value?.scrollIntoView({ behavior: 'smooth' })
  }
  
  const getPackageDetails = (pkgId: string) => {
    const packages = {
      standard: { pricePerPerson: 18.99, minGuests: 10 },
      premium: { pricePerPerson: 24.99, minGuests: 20 },
      executive: { pricePerPerson: 34.99, minGuests: 50 }
    }
    return packages[pkgId as keyof typeof packages] || packages.premium
  }
  
  const handleSubmit = async (data: Partial<CateringRequest>) => {
    isSubmitting.value = true
    
    try {
      await new Promise(resolve => setTimeout(resolve, 1500))
      
      
      showSuccessModal.value = true
      
      resetForm()
      
    } catch (error) {
      console.error('Submission failed:', error)
    } finally {
      isSubmitting.value = false
    }
  }
  
  const resetForm = () => {
    selectedPackage.value = 'premium'
    selectedEvent.value = 'corporate'
    guestCount.value = 50
    showSuccessModal.value = false
  }
  </script>