<template>
    <div class="space-y-6">
      <h3 class="text-2xl font-bold text-gray-900 dark:text-white">Frequently Asked Questions</h3>
      
      <div class="space-y-4">
        <div
          v-for="(faq, index) in faqs"
          :key="index"
          class="border border-gray-200 dark:border-slate-700 rounded-xl overflow-hidden"
        >
          <button
            @click="toggleFAQ(index)"
            class="w-full px-6 py-4 text-left flex justify-between items-center hover:bg-gray-50 dark:hover:bg-slate-700 transition-colors"
          >
            <span class="font-medium text-gray-900 dark:text-white">{{ faq.question }}</span>
            <ChevronDown 
              class="w-5 h-5 text-gray-500 transition-transform duration-200" 
              :class="{ 'rotate-180': openIndex === index }"
            />
          </button>
          
          <div
            v-show="openIndex === index"
            class="px-6 pb-4 text-gray-600 dark:text-gray-400"
            v-html="faq.answer"
          ></div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ChevronDown } from 'lucide-vue-next'
  import { ref } from 'vue'
  
  const openIndex = ref<number | null>(0)
  
  const faqs = [
    {
      question: 'How far in advance should I book?',
      answer: 'We recommend booking at least <strong>2 weeks in advance</strong> for small events and <strong>4-6 weeks</strong> for large events or weddings. Last-minute bookings may be available based on our schedule.'
    },
    {
      question: 'What is included in the catering packages?',
      answer: 'All packages include food, serving utensils, basic table setup, and staff. Premium and Executive packages include additional menu options, professional staff, and extended service hours. See package details for specifics.'
    },
    {
      question: 'Do you provide alcoholic beverages?',
      answer: 'We can arrange bar service through our licensed partners. Additional fees and permits may apply. We\'ll discuss your beverage needs during the consultation.'
    },
    {
      question: 'Can you accommodate dietary restrictions?',
      answer: 'Absolutely! We can accommodate vegetarian, vegan, gluten-free, halal, and other dietary needs. Please mention any restrictions in your special requests.'
    },
    {
      question: 'What is your cancellation policy?',
      answer: 'Cancellations made <strong>7+ days</strong> before the event receive a full refund. <strong>3-6 days</strong> receive 50% refund. Less than 48 hours is non-refundable.'
    },
    {
      question: 'Do you provide setup and cleanup?',
      answer: 'Yes! Our staff handles setup, food service, and cleanup. We\'ll coordinate with your venue for setup times and requirements.'
    }
  ]
  
  const toggleFAQ = (index: number) => {
    openIndex.value = openIndex.value === index ? null : index
  }
  </script>