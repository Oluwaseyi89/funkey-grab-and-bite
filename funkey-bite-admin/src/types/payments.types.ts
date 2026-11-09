import type { PaymentStatus } from './order.types';

/** An online transfer order on the payments page. */
export interface OnlinePayment {
  orderId: number;
  orderNumber: string;
  customerName: string;
  customerPhone: string;
  totalAmount: number;
  status: string;
  paymentMethod: 'transfer' | 'card';
  paymentStatus: PaymentStatus;
  paymentReference?: string;
  receivedAmount?: number;
  reviewReason?: 'amount_mismatch' | 'paid_after_cancellation' | string;
  referenceCount: number;
  expiresAt?: string;
  paidAt?: string;
  createdAt: string;
}

/** A Paystack payment that matched no order. */
export interface UnmatchedPayment {
  id: number;
  event: string;
  reference: string;
  amount?: number;
  currency?: string;
  receivedAt: string;
}

export type PaymentReviewDecision = 'accept' | 'refund' | 'refunded';

export interface PaymentRecheckResult {
  orderId: number;
  paid: boolean;
  outcome: string;
  message: string;
}

/** A background job; its payload is never sent to the admin. */
export interface BackgroundJob {
  id: number;
  kind: string;
  status: 'pending' | 'running' | 'succeeded' | 'dead';
  attempts: number;
  maxAttempts: number;
  runAt: string;
  lastError?: string;
  createdAt: string;
  updatedAt: string;
  completedAt?: string;
}

export interface AuditEntry {
  id: number;
  occurredAt: string;
  actorType: 'admin' | 'customer' | 'system' | 'paystack';
  actorId?: string;
  action: string;
  entityType: string;
  entityId?: string;
  requestId?: string;
  idempotencyKey?: string;
  ipAddress?: string;
  before?: Record<string, unknown>;
  after?: Record<string, unknown>;
  metadata?: Record<string, unknown>;
}

export interface Page {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}
