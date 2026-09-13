export type OrderStatus = 'pending' | 'confirmed' | 'preparing' | 'ready' | 'completed' | 'cancelled';
export type OrderType = 'pickup' | 'delivery' | 'catering';
export type PaymentMethod = 'cash' | 'transfer';
export type PaymentStatus = 'not_required' | 'pending' | 'paid' | 'failed';

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
  items: OrderItem[];
}

export interface OrderItemRequest {
  menuItemId: number;
  name: string;
  quantity: number;
  unitPrice: number;
  specialInstructions?: string;
}