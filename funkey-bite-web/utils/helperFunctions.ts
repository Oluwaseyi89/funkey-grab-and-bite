import type { Promotion } from '../types/menu'
 
 
 
  export const formatDiscountType = (type?: string) => {
    return type === 'bogo' ? 'BOGO' : (type?.toUpperCase() || '')
  }
  
  export const formatDiscountValue = (value?: number, type?: string) => {
    if (!value || !type) return ''
    switch(type) {
      case 'percentage':
        return `${value}% OFF`
      case 'fixed':
        return `$${value.toFixed(2)} OFF`
      case 'bogo':
        return 'Buy 1 Get 1'
      default:
        return `${value} OFF`
    }
  }
  
  export const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric',
      year: 'numeric'
    })
  }
  
  export const formatDateFull = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', { 
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }
  
  export const getTimeRemaining = (validUntil: string) => {
    const now = new Date()
    const endDate = new Date(validUntil)
    const diffTime = endDate.getTime() - now.getTime()
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    
    if (diffDays <= 0) return 'Ending today!'
    if (diffDays === 1) return 'Ends tomorrow'
    if (diffDays <= 7) return `${diffDays} days left`
    return `${diffDays} days remaining`
  }
  
  export const getDiscountedPrice = (promo: Promotion, originalPrice: number) => {
    switch(promo.discountType) {
      case 'percentage':
        const discount = originalPrice * (promo.discountValue / 100)
        return `$${(originalPrice - discount).toFixed(2)}`
      case 'fixed':
        const newPrice = Math.max(0, originalPrice - promo.discountValue)
        return `$${newPrice.toFixed(2)}`
      case 'bogo':
        return `$${originalPrice.toFixed(2)} (BOGO)`
      default:
        return `$${originalPrice.toFixed(2)}`
    }
  }