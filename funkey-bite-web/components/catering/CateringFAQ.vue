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
            :aria-expanded="openIndex === index"
            :aria-controls="`faq-panel-${index}`"
          >
            <span class="font-medium text-gray-900 dark:text-white flex-1 mr-2">{{ faq.question }}</span>
            <ChevronDown 
              class="w-5 h-5 text-gray-500 transition-transform duration-200" 
              :class="{ 'rotate-180': openIndex === index }"
              aria-hidden="true"
            />
          </button>
          
          <div
            :id="`faq-panel-${index}`"
            v-show="openIndex === index"
            class="px-6 pb-4 text-gray-600 dark:text-gray-400"
          >
            <p v-for="(part, i) in faq.parts" :key="i">
              <strong v-if="part.bold">{{ part.text }}</strong>
              <template v-else>{{ part.text }}</template>
            </p>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ChevronDown } from 'lucide-vue-next'
  import { ref } from 'vue'
  
  interface FaqPart {
    text: string
    bold?: boolean
  }

  interface Faq {
    question: string
    parts: FaqPart[]
  }

  const openIndex = ref<number | null>(0)
  
  const faqs: Faq[] = [
    {
      question: 'How far in advance should I book?',
      parts: [
        { text: 'We recommend booking at least ' },
        { text: '2 weeks in advance', bold: true },
        { text: ' for small events and ' },
        { text: '4-6 weeks', bold: true },
        { text: ' for large events or weddings. Last-minute bookings may be available based on our schedule.' }
      ]
    },
    {
      question: 'What is included in the catering packages?',
      parts: [
        { text: 'All packages include food, serving utensils, basic table setup, and staff. Premium and Executive packages include additional menu options, professional staff, and extended service hours. See package details for specifics.' }
      ]
    },
    {
      question: 'Do you provide alcoholic beverages?',
      parts: [
        { text: "We can arrange bar service through our licensed partners. Additional fees and permits may apply. We'll discuss your beverage needs during the consultation." }
      ]
    },
    {
      question: 'Can you accommodate dietary restrictions?',
      parts: [
        { text: 'Absolutely! We can accommodate vegetarian, vegan, gluten-free, halal, and other dietary needs. Please mention any restrictions in your special requests.' }
      ]
    },
    {
      question: 'What is your cancellation policy?',
      parts: [
        { text: 'Cancellations made ' },
        { text: '7+ days', bold: true },
        { text: ' before the event receive a full refund. ' },
        { text: '3-6 days', bold: true },
        { text: " receive 50% refund. Less than 48 hours is non-refundable." }
      ]
    },
    {
      question: 'Do you provide setup and cleanup?',
      parts: [
        { text: "Yes! Our staff handles setup, food service, and cleanup. We'll coordinate with your venue for setup times and requirements." }
      ]
    }
  ]
  
  const toggleFAQ = (index: number) => {
    openIndex.value = openIndex.value === index ? null : index
  }
  </script>