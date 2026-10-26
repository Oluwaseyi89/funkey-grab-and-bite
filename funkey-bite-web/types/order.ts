export type OrderStatus = 'pending' | 'confirmed' | 'preparing' | 'ready' | 'completed' | 'cancelled'
export type OrderType = 'pickup' | 'delivery' | 'catering'
// 'cash' is pay on delivery (or at pickup); 'transfer' is online payment.
export type PaymentMethod = 'cash' | 'transfer'
export type PaymentCollectedVia = 'cash' | 'pos' | 'transfer'
export type PaymentStatus = 'not_required' | 'pending' | 'paid' | 'failed'

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
  paymentMethod?: PaymentMethod
  paymentStatus?: PaymentStatus
  paymentReference?: string
  paymentAccountNumber?: string
  paymentAccountName?: string
  paymentBankName?: string
  paymentExpiresAt?: string
  paymentPaidAt?: string
  paymentCollectedVia?: PaymentCollectedVia
  deliveryAddress?: string
  invoiceNumber?: string
  subtotalAmount?: number
  deliveryFee?: number
  taxAmount?: number
  discountAmount?: number
}

// Contact details captured at checkout, before an order can be placed.
export interface CustomerInfo {
  customerName: string
  customerPhone: string
  customerEmail: string
  deliveryAddress: string
  notes: string
}

// The checkout request: the API identifies menu items by their numeric ID.
export interface OrderCreatePayload {
  customerName: string
  customerPhone: string
  customerEmail: string
  orderType: OrderType
  paymentMethod: PaymentMethod
  deliveryAddress?: string
  notes?: string
  items: {
    menuItemId: number
    name: string
    quantity: number
    unitPrice: number
    specialInstructions?: string
  }[]
}

export interface InvoiceParty {
  name: string
  phone: string
  email?: string
  address?: string
}

export interface InvoiceLine {
  name: string
  quantity: number
  unitPrice: number
  total: number
}

// Raised for pay-on-delivery orders.
export interface Invoice {
  invoiceNumber: string
  orderNumber: string
  issuedAt: string
  business: InvoiceParty
  customer: InvoiceParty
  orderType: OrderType
  items: InvoiceLine[]
  subtotal: number
  deliveryFee: number
  tax: number
  discount: number
  total: number
  paymentMethod: PaymentMethod
  paymentStatus: PaymentStatus
  paymentCollectedVia?: PaymentCollectedVia
  paidAt?: string
}

// What the invoice page needs to fetch an order's invoice.
export interface InvoiceLookup {
  orderNumber: string
  phone: string
}

export interface OrderLookupResult {
  order: Order
  invoice: Invoice | null
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