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
    // Current stock; missing when the item's stock isn't tracked.
    stockQuantity?: number
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
  
  // An active promotion from GET /promotions/active. Customers enter its code
  // at checkout, where the server works out the discount.
  export interface Promotion {
    id: number
    code: string
    title: string
    description: string
    promotionType: 'percentage' | 'fixed' | 'bogo'
    discountValue: number
    maxDiscount?: number
    minOrderAmount?: number
    validFrom: string
    validUntil: string
    isActive: boolean
  }

  // The server's answer to "would this code work on this amount?".
  export interface PromotionCheck {
    isValid: boolean
    message?: string
    discount?: number
  }
