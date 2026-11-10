export type OrderStatus = 'pending' | 'confirmed' | 'preparing' | 'ready' | 'completed' | 'cancelled'
export type OrderType = 'pickup' | 'delivery' | 'catering'
// 'cash' is pay on delivery (or at pickup); 'transfer' and 'card' are paid
// online through Paystack.
export type PaymentMethod = 'cash' | 'transfer' | 'card'
export type PaymentCollectedVia = 'cash' | 'pos' | 'transfer'
// expired: an online payment not made in time (the order was cancelled).
// needs_review / refund_due / refunded: a payment we're checking or returning.
export type PaymentStatus =
  | 'not_required'
  | 'pending'
  | 'paid'
  | 'failed'
  | 'expired'
  | 'needs_review'
  | 'refund_due'
  | 'refunded'

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
  // Opens Paystack's card popup; only present just after a card payment starts.
  paymentAccessCode?: string
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
  // Checked and applied by the server; the discount shows on the order.
  promotionCode?: string
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
