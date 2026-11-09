export type OrderStatus = 'pending' | 'confirmed' | 'preparing' | 'ready' | 'completed' | 'cancelled';
export type OrderType = 'pickup' | 'delivery' | 'catering';
// 'cash' is pay on delivery (or at pickup); 'transfer' is online payment.
export type PaymentMethod = 'cash' | 'transfer';
// 'not_required' only appears on cash orders from before payments on delivery were tracked.
// expired: an online transfer never made in time (the order was cancelled).
// needs_review: money arrived that can't simply be accepted, e.g. the wrong
// amount; refund_due / refunded follow the admin's decision on it.
export type PaymentStatus =
  | 'not_required'
  | 'pending'
  | 'paid'
  | 'failed'
  | 'expired'
  | 'needs_review'
  | 'refund_due'
  | 'refunded';
export type PaymentCollectedVia = 'cash' | 'pos' | 'transfer';

export interface OrderItem {
  menuItemId: number;
  name: string;
  quantity: number;
  unitPrice: number;
  specialInstructions?: string;
}

export interface Order {
  id: number;
  orderNumber: string;
  userId?: number;
  customerId?: string;
  customerName: string;
  customerPhone: string;
  customerEmail?: string;
  orderType: OrderType;
  status: OrderStatus;
  totalAmount: number;
  notes?: string;
  pickupTime?: string;
  estimatedReadyTime?: string;
  createdAt: string;
  paymentMethod: PaymentMethod;
  paymentStatus: PaymentStatus;
  paymentReference?: string;
  paymentAccountNumber?: string;
  paymentAccountName?: string;
  paymentBankName?: string;
  paymentExpiresAt?: string;
  paymentPaidAt?: string;
  paymentCollectedVia?: PaymentCollectedVia;
  deliveryAddress?: string;
  invoiceNumber?: string;
  subtotalAmount?: number;
  deliveryFee?: number;
  taxAmount?: number;
  discountAmount?: number;
  items: OrderItem[];
}

export interface InvoiceParty {
  name: string;
  phone: string;
  email?: string;
  address?: string;
}

export interface InvoiceLine {
  name: string;
  quantity: number;
  unitPrice: number;
  total: number;
}

// Raised for pay-on-delivery orders.
export interface Invoice {
  invoiceNumber: string;
  orderNumber: string;
  issuedAt: string;
  business: InvoiceParty;
  customer: InvoiceParty;
  orderType: OrderType;
  items: InvoiceLine[];
  subtotal: number;
  deliveryFee: number;
  tax: number;
  discount: number;
  total: number;
  paymentMethod: PaymentMethod;
  paymentStatus: PaymentStatus;
  paymentCollectedVia?: PaymentCollectedVia;
  paidAt?: string;
}

export interface OrderItemRequest {
  menuItemId: number;
  name: string;
  quantity: number;
  unitPrice: number;
  specialInstructions?: string;
}