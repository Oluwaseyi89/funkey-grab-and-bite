import { defineNuxtConfig } from 'nuxt/config'

export default defineNuxtConfig({
  css: ['~/assets/css/main.css'],
  devtools: { enabled: true },  
  modules: [
    '@nuxtjs/tailwindcss',
    '@nuxtjs/color-mode',
    '@nuxt/image',
    '@pinia/nuxt',
    '@vueuse/nuxt',
  ],
  features: {
    inlineStyles: true // Inlines component styles directly into the HTML
  },
    
  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.NUXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
      s3BucketUrl: process.env.NUXT_PUBLIC_S3_URL || 'https://funkey-static-assets.s3.amazonaws.com',
      siteUrl: process.env.NUXT_PUBLIC_SITE_URL || 'https://funkeygrabandbite.com',
      environment: process.env.NODE_ENV || 'development',
    }
  },
  
  app: {
    head: {
      htmlAttrs: { 
        lang: 'en',
        prefix: 'og: https://ogp.me/ns#',
      },
      title: 'Funkey Grab & Bite - Fast Food & Catering',
      titleTemplate: '%s | Funkey Grab & Bite',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1, maximum-scale=5' },
        { 
          name: 'description', 
          content: 'Fast-food business offering quick take-out meals, lunch packs, and reliable indoor/outdoor catering services. Order chips & chicken, noodles, shawarma, drinks, and soup bowls.' 
        },
        { 
          name: 'keywords', 
          content: 'fast food, take-out, catering, lunch packs, chips & chicken, noodles, shawarma, food delivery, restaurant, food bowls, soup, drinks, event catering' 
        },
        { name: 'theme-color', content: '#f44336' },
        { name: 'apple-mobile-web-app-capable', content: 'yes' },
        { name: 'apple-mobile-web-app-status-bar-style', content: 'black-translucent' },
        
        { property: 'og:type', content: 'website' },
        { property: 'og:site_name', content: 'Funkey Grab & Bite' },
        { property: 'og:title', content: 'Funkey Grab & Bite - Fast Food & Catering' },
        { 
          property: 'og:description', 
          content: 'Fast-food business offering quick take-out meals, lunch packs, and reliable indoor/outdoor catering services' 
        },
        { property: 'og:url', content: 'https://funkeygrabandbite.com' },
        { property: 'og:image', content: 'https://funkey-static-assets.s3.amazonaws.com/branding/og-image.jpg' },
        { property: 'og:image:width', content: '1200' },
        { property: 'og:image:height', content: '630' },
        { property: 'og:image:alt', content: 'Funkey Grab & Bite Restaurant' },
        { property: 'og:locale', content: 'en_US' },
        
        { name: 'twitter:card', content: 'summary_large_image' },
        { name: 'twitter:site', content: '@funkeygrab' },
        { name: 'twitter:creator', content: '@funkeygrab' },
        { name: 'twitter:title', content: 'Funkey Grab & Bite - Fast Food & Catering' },
        { 
          name: 'twitter:description', 
          content: 'Fast-food business offering quick take-out meals, lunch packs, and reliable indoor/outdoor catering services' 
        },
        { name: 'twitter:image', content: 'https://funkey-static-assets.s3.amazonaws.com/branding/twitter-image.jpg' },
        { name: 'twitter:image:alt', content: 'Funkey Grab & Bite Restaurant' },
        
        { name: 'robots', content: 'index, follow, max-image-preview:large, max-snippet:-1, max-video-preview:-1' },
        { name: 'googlebot', content: 'index, follow' },
        { name: 'author', content: 'Funkey Grab & Bite' },
        { name: 'publisher', content: 'Funkey Grab & Bite' },
        
        { name: 'format-detection', content: 'telephone=no' },
        { name: 'mobile-web-app-capable', content: 'yes' },
      ],
      link: [
        { rel: 'icon', href: '/favicon.ico', sizes: 'any' },
        { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' },
        { rel: 'apple-touch-icon', href: '/apple-touch-icon.png' },
        { rel: 'manifest', href: '/site.webmanifest' },
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: 'anonymous' },
        { 
          rel: 'stylesheet', 
          href: 'https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600;700&family=Inter:wght@300;400;500;600&display=swap' 
        },
        { rel: 'canonical', href: 'https://funkeygrabandbite.com' },
      ],
      script: [
        {
          type: 'application/ld+json',
          innerHTML: JSON.stringify({
            '@context': 'https://schema.org',
            '@type': 'Restaurant',
            // telephone, address and opening hours come from the API at
            // runtime (layouts/default.vue) so they follow the admin settings.
            name: 'Funkey Grab & Bite',
            image: 'https://funkey-static-assets.s3.amazonaws.com/branding/logo.png',
            '@id': 'https://funkeygrabandbite.com',
            url: 'https://funkeygrabandbite.com',
            priceRange: '$$',
            servesCuisine: ['Fast Food', 'American', 'Middle Eastern', 'Asian'],
            aggregateRating: {
              '@type': 'AggregateRating',
              ratingValue: '4.8',
              ratingCount: '150',
              bestRating: '5',
              worstRating: '1'
            },
            hasMenu: 'https://funkeygrabandbite.com/menu',
            acceptsReservations: 'True',
            menu: 'https://funkeygrabandbite.com/menu'
          })
        },
        {
          type: 'application/ld+json',
          innerHTML: JSON.stringify({
            '@context': 'https://schema.org',
            '@type': 'WebSite',
            name: 'Funkey Grab & Bite',
            url: 'https://funkeygrabandbite.com',
            potentialAction: {
              '@type': 'SearchAction',
              target: 'https://funkeygrabandbite.com/menu?search={search_term_string}',
              'query-input': 'required name=search_term_string'
            }
          })
        }
      ]
    }
  },
  
  typescript: {
    strict: true,
    typeCheck: true,
  },
  
  build: {
    transpile: ['@headlessui/vue', 'gsap']
  },
  
  nitro: {
    compressPublicAssets: true,
    prerender: {
      crawlLinks: true,
      routes: [
        '/',
      ]
    }
  },
  
  srcDir: '.',
  
  imports: {
    dirs: [
      'composables',
      'composables/**',
      'stores',
      'utils',
      'types'
    ]
  },
  
  components: [
    {
      path: '~/components',
      pathPrefix: false,
      extensions: ['.vue'],
    }
  ]
})