<template>
    <div>

      <PageHeader
      title-before="Get"
      highlight-text="in Touch"
      variant="gradient"
      alignment="center"
      :narrow="true"
      header-class="pb-12"
    >
      

        <template #subtitle>
            Have questions, feedback, or catering inquiries? We're here to help! Reach out through any channel below.
        </template>

        <div class="mt-8 max-w-2xl mx-auto">

            <div class="flex flex-wrap justify-center gap-4">
            <a href="#contact-form" class="btn-primary">Send Message</a>
            <a href="tel:+15551234567" class="btn-secondary">
                <Phone class="w-4 h-4 inline mr-2" />
                Call Now
            </a>
            </div>
        </div>
    </PageHeader>

  
      <div class="section-padding mx-5 px-3 py-3 md:mx-8 md:py-8 md:px-8">
        <div class="container-narrow">
          <div class="grid lg:grid-cols-3 gap-8">

            <div class="space-y-8">
              <ContactInfo />
              <BusinessHours />
              <ContactMap />
            </div>
  

            <div class="lg:col-span-2">
              <div class="bg-white dark:bg-slate-800 rounded-xl p-6 md:p-8 border border-gray-200 dark:border-slate-700">
                <ContactForm
                  :is-submitting="isSubmitting"
                  @submit="handleSubmit"
                />
              </div>
  

              <div class="mt-8">
                <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-4">Common Questions</h3>
                <div class="space-y-4">
                  <div 
                    v-for="(faq, index) in faqs" 
                    :key="index"
                    class="border border-gray-200 dark:border-slate-700 rounded-lg overflow-hidden"
                  >
                    <button
                      @click="toggleFAQ(index)"
                      class="w-full px-4 py-3 text-left flex justify-between items-center hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
                    >
                      <span class="font-medium text-gray-900 dark:text-white">{{ faq.question }}</span>
                      <ChevronDown 
                        class="w-5 h-5 text-gray-500 transition-transform duration-200" 
                        :class="{ 'rotate-180': openFaqIndex === index }"
                      />
                    </button>
                    
                    <div
                      v-show="openFaqIndex === index"
                      class="px-4 pb-4 text-gray-600 dark:text-gray-400"
                    >
                      {{ faq.answer }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
  

          <div class="mt-16 pt-12 border-t border-gray-200 dark:border-slate-700">
            <h2 class="text-3xl font-bold text-gray-900 dark:text-white text-center mb-12">
              Meet Our <span class="text-gradient">Team</span>
            </h2>
            
            <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div 
                v-for="member in team" 
                :key="member.name"
                class="text-center"
              >
                <div class="w-32 h-32 mx-auto mb-4 bg-gradient-to-br from-brand-100 to-accent-100 dark:from-brand-900/30 dark:to-accent-900/30 rounded-full flex items-center justify-center">
                  <div class="w-28 h-28 bg-white dark:bg-slate-800 rounded-full flex items-center justify-center">
                    <Users class="w-16 h-16 text-brand-500" />
                  </div>
                </div>
                <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-1">{{ member.name }}</h3>
                <p class="text-brand-500 font-medium mb-2">{{ member.role }}</p>
                <p class="text-gray-600 dark:text-gray-400 text-sm">{{ member.bio }}</p>
              </div>
            </div>
          </div>
  

          <div class="mt-16 bg-gradient-to-r from-brand-500 to-accent-500 rounded-2xl p-8 text-center text-white">
            <h2 class="text-3xl font-bold mb-4">Need Immediate Assistance?</h2>
            <p class="text-xl opacity-90 mb-6 max-w-2xl mx-auto">
              Call us now for orders, catering inquiries, or urgent matters.
            </p>
            <a 
              href="tel:+15551234567" 
              class="inline-flex items-center bg-white text-brand-600 hover:bg-gray-100 font-bold py-4 px-8 rounded-xl transition-colors text-lg"
            >
              <Phone class="w-5 h-5 mr-2" />
              Call Now: (555) 123-4567
            </a>
          </div>
        </div>
      </div>
  

      <div 
        v-if="showSuccessModal"
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
      >
        <div class="bg-white dark:bg-slate-800 rounded-2xl max-w-md w-full p-8 text-center">
          <div class="w-20 h-20 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-6">
            <CheckCircle class="w-10 h-10 text-green-600 dark:text-green-400" />
          </div>
          
          <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">Message Sent!</h2>
          <p class="text-gray-600 dark:text-gray-400 mb-6">
            Thank you for reaching out. We'll get back to you within 24 hours.
          </p>
          
          <button @click="showSuccessModal = false" class="w-full btn-primary">
            Continue Browsing
          </button>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue'
  import { 
    Phone, 
    Users, 
    ChevronDown, 
    CheckCircle 
  } from 'lucide-vue-next'
  
  import PageHeader from '../../components/layout/PageHeader.vue'
  import ContactForm from '../../components/contact/ContactForm.vue'
  import ContactInfo from '../../components/contact/ContactInfo.vue'
  import BusinessHours from '../../components/contact/BusinessHours.vue'
  import ContactMap from '../../components/contact/ContactMap.vue'
  
  const isSubmitting = ref(false)
  const showSuccessModal = ref(false)
  const openFaqIndex = ref<number | null>(0)
  
  const faqs = [
    {
      question: 'What are your delivery hours?',
      answer: 'We deliver from 11:00 AM to 9:00 PM daily. Last orders for delivery must be placed by 8:30 PM.'
    },
    {
      question: 'Do you offer catering for small events?',
      answer: 'Yes! We cater events of all sizes, starting from 10 guests. Contact us for a customized quote.'
    },
    {
      question: 'Can I modify or cancel my order?',
      answer: 'You can modify or cancel your order within 15 minutes of placing it by calling us directly.'
    },
    {
      question: 'Do you have vegetarian/vegan options?',
      answer: 'Absolutely! We offer a variety of vegetarian and vegan dishes. Check our menu or ask our staff.'
    },
    {
      question: 'How do I apply for a job?',
      answer: 'Send your resume to careers@funkeygrabandbite.com or visit us in person to fill out an application.'
    }
  ]
  
  const team = [
    {
      name: 'Alex Johnson',
      role: 'General Manager',
      bio: '15+ years in food service. Passionate about creating memorable dining experiences.'
    },
    {
      name: 'Maria Rodriguez',
      role: 'Head Chef',
      bio: 'Trained in culinary arts with a focus on fresh, locally-sourced ingredients.'
    },
    {
      name: 'David Chen',
      role: 'Catering Director',
      bio: 'Specializes in event planning and large-scale catering operations.'
    }
  ]
  
  const handleSubmit = async (formData: any) => {
    isSubmitting.value = true
    
    try {
      await new Promise(resolve => setTimeout(resolve, 1500))
      
      
      showSuccessModal.value = true
      
    } catch (error) {
      console.error('Message submission failed:', error)
    } finally {
      isSubmitting.value = false
    }
  }
  
  const toggleFAQ = (index: number) => {
    openFaqIndex.value = openFaqIndex.value === index ? null : index
  }
  </script>