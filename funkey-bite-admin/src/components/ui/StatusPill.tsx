import React from 'react';

// Order statuses in the words staff use, with one colour each.
export const ORDER_STATUS: Record<string, { label: string; className: string }> = {
  pending: { label: 'New', className: 'badge-pending' },
  confirmed: { label: 'Confirmed', className: 'badge-confirmed' },
  preparing: { label: 'Cooking', className: 'badge-preparing' },
  ready: { label: 'Ready', className: 'badge-ready' },
  completed: { label: 'Completed', className: 'badge-completed' },
  cancelled: { label: 'Cancelled', className: 'badge-cancelled' },
};

export const PAYMENT_STATUS: Record<string, { label: string; className: string }> = {
  not_required: { label: 'No payment', className: 'badge-completed' },
  pending: { label: 'Unpaid', className: 'badge-pending' },
  paid: { label: 'Paid', className: 'badge-ready' },
  failed: { label: 'Failed', className: 'badge-cancelled' },
  expired: { label: 'Expired', className: 'badge-completed' },
  needs_review: { label: 'Check payment', className: 'badge-cancelled' },
  refund_due: { label: 'Refund due', className: 'badge-cancelled' },
  refunded: { label: 'Refunded', className: 'badge-completed' },
};

// What a kitchen does next with an order, as one tap.
export const NEXT_STATUS: Record<string, { status: string; label: string } | undefined> = {
  pending: { status: 'confirmed', label: 'Confirm' },
  confirmed: { status: 'preparing', label: 'Start cooking' },
  preparing: { status: 'ready', label: 'Mark ready' },
  ready: { status: 'completed', label: 'Complete' },
};

interface StatusPillProps {
  status?: string | null;
  kind?: 'order' | 'payment';
  className?: string;
}

const StatusPill: React.FC<StatusPillProps> = ({ status, kind = 'order', className = '' }) => {
  if (!status) return null;
  const map = kind === 'order' ? ORDER_STATUS : PAYMENT_STATUS;
  const entry = map[status] ?? { label: status.replace(/_/g, ' '), className: 'badge-completed' };
  return <span className={`badge ${entry.className} ${className}`}>{entry.label}</span>;
};

export default StatusPill;
