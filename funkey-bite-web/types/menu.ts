export interface MenuCategory {
    id: string
    name: string
    description: string
    displayOrder: number
    isActive: boolean
  }
  
  export interface MenuItem {
    id: string
    categoryId: string
    name: string
    description: string
    price: number
    imageUrl: string
    isAvailable: boolean
    isPreOrder: boolean
    preparationTime: number
    tags: string[]
    nutritionalInfo?: {
      calories: number
      protein: number
      carbs: number
      fat: number
    }
  }
  
  export interface CartItem {
    menuItem: MenuItem
    quantity: number
    specialInstructions?: string
  }
  
  export interface Promotion {
    id: string
    title: string
    description: string
    discountType: 'percentage' | 'fixed' | 'bogo'
    discountValue: number
    validFrom: string
    validUntil: string
    isActive: boolean
    applicableItems?: string[]
  }