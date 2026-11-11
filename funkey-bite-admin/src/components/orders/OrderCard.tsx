import React from 'react';
import { Link } from 'react-router-dom';
import { Bike, ChevronRight, Loader2, Store } from 'lucide-react';
import StatusPill, { NEXT_STATUS } from '../ui/StatusPill';
import { useAdvanceOrder } from '../../hooks/useAdvanceOrder';
import { naira, timeAgo } from '../../utils/format';
import type { Order } from '../../types';

const paymentLabel = (order: Order) => {
  const method = order.paymentMethod === 'cash' ? (order.orderType === 'delivery' ? 'On delivery' : 'At pickup') : order.paymentMethod === 'card' ? 'Card' : 'Transfer';
  return method;
};

// One order as a card: who, what, how much, where it is, and the next step
// as a single tap. The whole card opens the order.
const OrderCard: React.FC<{ order: Order; compact?: boolean }> = ({ order, compact = false }) => {
  const advance = useAdvanceOrder();
  const next = NEXT_STATUS[order.status];
  // Don't start cooking an online order nobody has paid for yet.
  const waitingForOnlinePayment = order.paymentMethod !== 'cash' && order.paymentStatus === 'pending' && order.status === 'pending';
  const busy = advance.isPending && advance.variables?.id === order.id;
  const items = order.items ?? [];
  const summary = items.slice(0, 3).map((item) => `${item.quantity}× ${item.name}`).join(', ');
  const TypeIcon = order.orderType === 'delivery' ? Bike : Store;

  return (
    <article className="card relative p-4 transition hover:shadow-lift">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <h3 className="font-mono text-sm font-bold text-gray-900 dark:text-white">
            <Link to={`/orders/${order.id}`} className="after:absolute after:inset-0 after:content-['']">
              #{order.orderNumber}
            </Link>
          </h3>
          <p className="mt-0.5 flex items-center gap-1.5 text-xs text-gray-500 dark:text-gray-400">
            <TypeIcon className="h-3.5 w-3.5" aria-hidden="true" />
            <span className="capitalize">{order.orderType}</span> · {timeAgo(order.createdAt)}
          </p>
        </div>
        <p className="shrink-0 text-base font-bold tabular-nums text-gray-900 dark:text-white">{naira(order.totalAmount)}</p>
      </div>

      <p className="mt-2 truncate text-sm font-medium text-gray-800 dark:text-gray-200">{order.customerName}</p>
      {!compact && summary && (
        <p className="mt-0.5 line-clamp-2 text-sm text-gray-500 dark:text-gray-400">
          {summary}
          {items.length > 3 ? ` +${items.length - 3} more` : ''}
        </p>
      )}

      <div className="mt-3 flex flex-wrap items-center gap-1.5">
        <StatusPill status={order.status} />
        <span className="badge bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300">{paymentLabel(order)}</span>
        <StatusPill status={order.paymentStatus} kind="payment" />
      </div>

      {(next || order.status === 'completed' || order.status === 'cancelled') && (
        <div className="relative z-10 mt-3 flex items-center gap-2 border-t border-gray-100 pt-3 dark:border-gray-800">
          {next && !waitingForOnlinePayment ? (
            <button
              type="button"
              className="btn-primary flex-1"
              disabled={busy}
              onClick={() => advance.mutate({ id: order.id, status: next.status, orderNumber: order.orderNumber })}
            >
              {busy && <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />}
              {next.label}
              <span className="sr-only"> order {order.orderNumber}</span>
            </button>
          ) : (
            <p className="flex-1 text-sm text-gray-500 dark:text-gray-400">
              {waitingForOnlinePayment ? 'Waiting for the customer to pay' : 'Nothing left to do'}
            </p>
          )}
          <Link to={`/orders/${order.id}`} className="btn-secondary px-3" aria-label={`Open order ${order.orderNumber}`}>
            <ChevronRight className="h-5 w-5" aria-hidden="true" />
          </Link>
        </div>
      )}
    </article>
  );
};

export default OrderCard;
