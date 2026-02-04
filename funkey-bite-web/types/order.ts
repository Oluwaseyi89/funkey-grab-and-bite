export type OrderStatus = 'pending' | 'confirmed' | 'preparing' | 'ready' | 'completed' | 'cancelled'
export type OrderType = 'pickup' | 'delivery' | 'catering'

export interface Order {
  id: string
  orderNumber: string
  customerName: string
  customerPhone: string
  customerEmail?: string
  orderType: OrderType
  status: OrderStatus
  totalAmount: number
  notes?: string
  pickupTime?: string
  items: OrderItem[]
  createdAt: string
  estimatedReadyTime?: string
}

export interface OrderItem {
  menuItemId: string
  name: string
  quantity: number
  unitPrice: number
  specialInstructions?: string
}

export interface CateringRequest {
  id: string
  eventName?: string
  contactName: string
  contactPhone: string
  contactEmail?: string
  eventDate: string
  eventTime?: string
  guestCount: number
  eventType: string
  budget?: number
  specialRequests?: string
  status: 'pending' | 'confirmed' | 'declined' | 'completed'
  createdAt: string
}