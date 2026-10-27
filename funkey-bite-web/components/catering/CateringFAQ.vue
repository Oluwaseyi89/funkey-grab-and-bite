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
            <p>
              <template v-for="(part, i) in faq.parts" :key="i">
                <strong v-if="part.bold">{{ part.text }}</strong>
                <template v-else>{{ part.text }}</template>
              </template>
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
        { text: 'Each package shows its minimum notice. Custom requests need at least ' },
        { text: '3 days', bold: true },
        { text: ", and more for large events. The request form only offers dates we can cater." }
      ]
    },
    {
      question: 'What is the difference between a package and a custom request?',
      parts: [
        { text: 'A package has a set price per guest, a guest range and a list of dishes to choose from, so you see an estimate straight away. A custom request lets you describe any menu and service you want, and we quote for it.' }
      ]
    },
    {
      question: 'How is the price decided?',
      parts: [
        { text: "For packages, the estimate is the price per guest times the number of guests. We review every request and send you a quote with the final price before anything is confirmed." }
      ]
    },
    {
      question: 'Can you accommodate dietary needs?',
      parts: [
        { text: 'Yes. Tick vegetarian, vegan, halal, no pork, nut-free, gluten-free or dairy-free on the form, and add any allergies in the notes.' }
      ]
    },
    {
      question: 'Do you serve at the venue or just deliver?',
      parts: [
        { text: 'Both. Choose drop-off (food delivered ready to serve), a buffet with our staff, or full service with servers, setup and cleanup.' }
      ]
    },
    {
      question: 'How do I check or change my request?',
      parts: [
        { text: "Use your reference and phone number on the Check your request page to see its status and quote. To change or cancel it, call us and quote your reference." }
      ]
    }
  ]

  const toggleFAQ = (index: number) => {
    openIndex.value = openIndex.value === index ? null : index
  }
  </script>