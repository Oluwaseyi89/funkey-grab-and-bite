import type { UseHeadInput } from '@unhead/vue'
import { useRoute } from 'vue-router'
import { useHead, useRuntimeConfig } from 'nuxt/app'

export interface SeoConfig {
  title?: string
  description?: string
  image?: string
  url?: string
  type?: string
  robots?: string
  keywords?: string
  author?: string
}

export const useSeo = (seoConfig: SeoConfig = {}): void => {
  const route = useRoute()
  const config = useRuntimeConfig()
  
  const siteUrl = config.public.siteUrl || 'https://funkeygrabandbite.com'
  const defaultImage = 'https://funkey-static-assets.s3.amazonaws.com/branding/og-image.jpg'
  
  const title = seoConfig.title || 'Funkey Grab & Bite - Fast Food & Catering'
  const description = seoConfig.description || 
    'Fast-food business offering quick take-out meals, lunch packs, and reliable indoor/outdoor catering services'
  const image = seoConfig.image || defaultImage
  const url = seoConfig.url || `${siteUrl}${route.path}`
  const type = seoConfig.type || 'website'
  const robots = seoConfig.robots || 'index, follow'
  const keywords = seoConfig.keywords || 
    'fast food, take-out, catering, lunch packs, chips & chicken, noodles, shawarma'
  const author = seoConfig.author || 'Funkey Grab & Bite'
  
  const headConfig: UseHeadInput = {
    title,
    meta: [
      { name: 'description', content: description },
      { name: 'keywords', content: keywords },
      { name: 'author', content: author },
      { name: 'robots', content: robots },
      
      { property: 'og:title', content: title },
      { property: 'og:description', content: description },
      { property: 'og:image', content: image },
      { property: 'og:url', content: url },
      { property: 'og:type', content: type },
      { property: 'og:site_name', content: 'Funkey Grab & Bite' },
      { property: 'og:locale', content: 'en_US' },
      
      { name: 'twitter:title', content: title },
      { name: 'twitter:description', content: description },
      { name: 'twitter:image', content: image },
      { name: 'twitter:card', content: 'summary_large_image' },
      { name: 'twitter:site', content: '@funkeygrab' },
      { name: 'twitter:creator', content: '@funkeygrab' },
    ],
    link: [
      { rel: 'canonical', href: url }
    ],
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': type === 'article' ? 'Article' : 'WebPage',
          headline: title,
          description: description,
          image: image,
          url: url,
          author: {
            '@type': 'Organization',
            name: author
          },
          publisher: {
            '@type': 'Organization',
            name: 'Funkey Grab & Bite',
            logo: {
              '@type': 'ImageObject',
              url: 'https://funkey-static-assets.s3.amazonaws.com/branding/logo.png'
            }
          },
          datePublished: new Date().toISOString(),
          dateModified: new Date().toISOString()
        })
      }
    ]
  }
  
  useHead(headConfig)
}