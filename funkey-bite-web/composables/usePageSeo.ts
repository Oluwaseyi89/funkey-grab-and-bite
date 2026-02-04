import { useSeo } from './useSeo'

// Page-specific SEO configurations
export const usePageSeo = (page: string, customConfig?: any) => {
  const configs: Record<string, any> = {
    home: {
      title: 'Funkey Grab & Bite - Fast Food & Catering',
      description: 'Order delicious fast food, lunch packs, and book catering services. Chips & chicken, noodles, shawarma, and more!',
      keywords: 'fast food delivery, lunch packs, catering service, chicken and chips, shawarma, noodles',
    },
    menu: {
      title: 'Menu | Funkey Grab & Bite',
      description: 'Browse our complete menu featuring chips & chicken, noodles, shawarma, drinks, soup bowls, and lunch packs.',
      keywords: 'food menu, restaurant menu, chicken dishes, noodle dishes, shawarma wraps, drinks menu',
    },
    order: {
      title: 'Order Online | Funkey Grab & Bite',
      description: 'Order your favorite meals online for pickup or delivery. Fast and easy ordering process.',
      keywords: 'online food order, food delivery, pickup order, fast food ordering',
    },
    catering: {
      title: 'Catering Services | Funkey Grab & Bite',
      description: 'Professional indoor and outdoor catering services for events, parties, and corporate gatherings.',
      keywords: 'event catering, party catering, corporate catering, outdoor catering, catering service',
    },
    contact: {
      title: 'Contact Us | Funkey Grab & Bite',
      description: 'Get in touch with Funkey Grab & Bite for inquiries, feedback, or catering quotes.',
      keywords: 'contact restaurant, customer service, feedback, inquiries, catering quotes',
    },
  }
  
  const pageConfig = configs[page] || {}
  const finalConfig = { ...pageConfig, ...customConfig }
  
  useSeo(finalConfig)
}