import type { OrderStatus, PaymentMethod, PaymentStatus } from '~/types/order'

export const ORDER_PAYMENT_METHOD_LABELS: Record<PaymentMethod, string> = {
  cash: 'Pay on delivery / pickup',
  transfer: 'Bank transfer',
  card: 'Card',
}

export const PAYMENT_STATUS_LABELS: Record<PaymentStatus, string> = {
  not_required: 'No payment needed',
  pending: 'Awaiting payment',
  paid: 'Paid',
  failed: 'Payment failed',
  expired: 'Not paid in time',
  needs_review: 'Payment being checked',
  refund_due: 'Refund on its way',
  refunded: 'Refunded',
}

export const ORDER_STATUS_LABELS: Record<OrderStatus, string> = {
  pending: 'Received',
  confirmed: 'Confirmed',
  preparing: 'Being prepared',
  ready: 'Ready',
  completed: 'Completed',
  cancelled: 'Cancelled',
}

export const paymentStatusClass = (status?: PaymentStatus) => {
  switch (status) {
    case 'paid':
      return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300'
    case 'pending':
      return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300'
    case 'failed':
    case 'needs_review':
    case 'refund_due':
      return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300'
    default:
      return 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
  }
}
