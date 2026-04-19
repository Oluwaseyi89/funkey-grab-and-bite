import type { MenuCategory, MenuItem, Promotion } from '~/types/menu'
import type { Order } from '~/types/order'
import { useRuntimeConfig } from 'nuxt/app'

export const mockCategories: MenuCategory[] = [
  {
    id: '1',
    name: 'Chips & Chicken',
    description: 'Crispy golden fries with our signature chicken',
    displayOrder: 1,
    isActive: true,
  },
  {
    id: '2',
    name: 'Noodles',
    description: 'Delicious noodle dishes with fresh ingredients',
    displayOrder: 2,
    isActive: true,
  },
  {
    id: '3',
    name: 'Shawarma',
    description: 'Authentic shawarma wraps and plates',
    displayOrder: 3,
    isActive: true,
  },
  {
    id: '4',
    name: 'Drinks',
    description: 'Refreshing beverages',
    displayOrder: 4,
    isActive: true,
  },
  {
    id: '5',
    name: 'Soup & Food Bowls',
    description: 'Hearty soups and nutritious bowls (Pre-order)',
    displayOrder: 5,
    isActive: true,
  },
  {
    id: '6',
    name: 'Lunch Packs',
    description: 'Complete meal deals for lunch',
    displayOrder: 6,
    isActive: true,
  },
]

export const mockMenuItems: MenuItem[] = [
  {
    id: '1',
    categoryId: '1',
    name: 'Classic Chicken & Chips',
    description: 'Crispy chicken with golden fries and special sauce',
    price: 18186,
    imageUrl: '/images/chicken_and_chips.jpg',
    isAvailable: true,
    isPreOrder: false,
    preparationTime: 15,
    tags: ['best seller', 'spicy'],
    nutritionalInfo: {
      calories: 850,
      protein: 45,
      carbs: 65,
      fat: 32,
    },
  },
  {
    id: '2',
    categoryId: '1',
    name: 'Spicy Chicken Wings',
    description: 'Crispy wings tossed in our signature spicy sauce',
    price: 15386,
    imageUrl: '/images/spicy_chicken_wings.jpg',
    isAvailable: true,
    isPreOrder: false,
    preparationTime: 20,
    tags: ['spicy', 'shareable'],
  },
  {
    id: '3',
    categoryId: '2',
    name: 'Stir Fry Noodles',
    description: 'Fresh noodles with vegetables and choice of protein',
    price: 16786,
    imageUrl: '/images/stir_fry_noodles.jpeg',
    isAvailable: true,
    isPreOrder: false,
    preparationTime: 12,
    tags: ['vegetarian option'],
  },
  {
    id: '4',
    categoryId: '3',
    name: 'Chicken Shawarma Wrap',
    description: 'Marinated chicken with fresh veggies in warm pita',
    price: 12586,
    imageUrl: '/images/chicken_shawarma.avif',
    isAvailable: true,
    isPreOrder: false,
    preparationTime: 10,
    tags: ['quick', 'popular'],
  },
  {
    id: '5',
    categoryId: '4',
    name: 'Fresh Lemonade',
    description: 'Homemade lemonade with mint',
    price: 5586,
    imageUrl: '/images/fresh_lemonade.jpg',
    isAvailable: true,
    isPreOrder: false,
    preparationTime: 5,
    tags: ['refreshing'],
  },
  {
    id: '6',
    categoryId: '5',
    name: 'Hearty Beef Stew Bowl',
    description: 'Slow-cooked beef stew with vegetables and rice',
    price: 20986,
    imageUrl: '/images/hearty_beef_and_stew.jpg',
    isAvailable: true,
    isPreOrder: true,
    preparationTime: 30,
    tags: ['pre-order', 'hearty'],
  },
]

export const mockPromotions: Promotion[] = [
  {
    id: '1',
    title: 'Weekend Special',
    description: '20% off all lunch packs this weekend',
    discountType: 'percentage',
    discountValue: 20,
    validFrom: new Date().toISOString(),
    validUntil: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
    isActive: true,
  },
  {
    id: '2',
    title: 'Buy One Get One',
    description: 'Buy any shawarma, get one free',
    discountType: 'bogo',
    discountValue: 1,
    validFrom: new Date().toISOString(),
    validUntil: new Date(Date.now() + 3 * 24 * 60 * 60 * 1000).toISOString(),
    isActive: true,
    applicableItems: ['4'],
  },
]

export const mockOrders: Order[] = [
  {
    id: '1',
    orderNumber: 'FG-2024-001',
    customerName: 'John Doe',
    customerPhone: '+1234567890',
    orderType: 'pickup',
    status: 'completed',
    totalAmount: 36372,
    items: [
      {
        menuItemId: '1',
        name: 'Classic Chicken & Chips',
        quantity: 2,
          unitPrice: 254440,
      },
    ],
    createdAt: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
  },
]

export const heroBannerBackgroundImages: string[] = [
  '/images/hero_background/hero_one.avif', 
  '/images/hero_background/hero_two.jpg', 
  '/images/hero_background/hero_three.jpg'
]


export const getS3ImageUrl = (path: string): string => {
  const config = useRuntimeConfig()
  return `${config.public.s3BucketUrl}/${path}`
}