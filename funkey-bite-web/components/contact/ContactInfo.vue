<template>
    <div class="space-y-6">
      <h3 class="text-2xl font-bold text-gray-900 dark:text-white">Contact Information</h3>
      
      <!-- Address -->
      <div class="flex items-start space-x-4 p-4 bg-white dark:bg-slate-800 rounded-xl border border-gray-200 dark:border-slate-700">
        <MapPin class="w-6 h-6 text-brand-500 mt-1 flex-shrink-0" />
        <div>
          <h4 class="font-bold text-gray-900 dark:text-white mb-1">Visit Our Location</h4>
          <p class="text-gray-600 dark:text-gray-400">
            123 Food Street<br>
            City Center, Your City<br>
            State 12345
          </p>
          <a 
            href="https://maps.google.com/?q=123+Food+Street+City+Center"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-block mt-2 text-brand-500 hover:text-brand-600 font-medium"
          >
            Get Directions →
          </a>
        </div>
      </div>
  
      <!-- Phone -->
      <div class="flex items-center space-x-4 p-4 bg-white dark:bg-slate-800 rounded-xl border border-gray-200 dark:border-slate-700">
        <Phone class="w-6 h-6 text-brand-500 flex-shrink-0" />
        <div>
          <h4 class="font-bold text-gray-900 dark:text-white mb-1">Call Us</h4>
          <a href="tel:+15551234567" class="text-gray-600 dark:text-gray-400 hover:text-brand-500">
            (555) 123-4567
          </a>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
            Available 8AM-10PM daily
          </p>
        </div>
      </div>
  
      <!-- Email -->
      <div class="flex items-center space-x-4 p-4 bg-white dark:bg-slate-800 rounded-xl border border-gray-200 dark:border-slate-700">
        <Mail class="w-6 h-6 text-brand-500 flex-shrink-0" />
        <div>
          <h4 class="font-bold text-gray-900 dark:text-white mb-1">Email Us</h4>
          <a href="mailto:info@funkeygrabandbite.com" class="text-gray-600 dark:text-gray-400 hover:text-brand-500">
            info@funkeygrabandbite.com
          </a>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
            For general inquiries
          </p>
        </div>
      </div>
  
      <!-- Orders Email -->
      <div class="flex items-center space-x-4 p-4 bg-white dark:bg-slate-800 rounded-xl border border-gray-200 dark:border-slate-700">
        <ShoppingBag class="w-6 h-6 text-brand-500 flex-shrink-0" />
        <div>
          <h4 class="font-bold text-gray-900 dark:text-white mb-1">Orders & Catering</h4>
          <a href="mailto:orders@funkeygrabandbite.com" class="text-gray-600 dark:text-gray-400 hover:text-brand-500">
            orders@funkeygrabandbite.com
          </a>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
            For orders and catering requests
          </p>
        </div>
      </div>
  
      <!-- Social Media -->
      <div class="p-4 bg-white dark:bg-slate-800 rounded-xl border border-gray-200 dark:border-slate-700">
        <h4 class="font-bold text-gray-900 dark:text-white mb-3">Follow Us</h4>
        <div class="flex space-x-4">
          <a 
            href="#" 
            class="w-10 h-10 bg-gray-100 dark:bg-slate-700 rounded-full flex items-center justify-center hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
            aria-label="Facebook"
          >
            <Facebook class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </a>
          <a 
            href="#" 
            class="w-10 h-10 bg-gray-100 dark:bg-slate-700 rounded-full flex items-center justify-center hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
            aria-label="Instagram"
          >
            <Instagram class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </a>
          <a 
            href="#" 
            class="w-10 h-10 bg-gray-100 dark:bg-slate-700 rounded-full flex items-center justify-center hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
            aria-label="Twitter"
          >
            <Twitter class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </a>
          <a 
            href="#" 
            class="w-10 h-10 bg-gray-100 dark:bg-slate-700 rounded-full flex items-center justify-center hover:bg-gray-200 dark:hover:bg-slate-600 transition-colors"
            aria-label="TikTok"
          >
            <MessageCircle class="w-5 h-5 text-gray-700 dark:text-gray-300" />
          </a>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { 
    MapPin, 
    Phone, 
    Mail, 
    ShoppingBag,
    Facebook, 
    Instagram, 
    Twitter,
    MessageCircle
  } from 'lucide-vue-next'
  </script>