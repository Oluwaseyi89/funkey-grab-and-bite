import React, { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { AlertTriangle, CheckCircle, Loader2, RefreshCw, RotateCcw, Search, Undo2, X } from 'lucide-react';
import toast from 'react-hot-toast';
import {
  getOnlinePayments,
  getUnmatchedPayments,
  recheckPayment,
  resolvePaymentReview,
} from '../../api/adminApi';
import { useIdempotencyKey } from '../../hooks/useIdempotencyKey';
import type { OnlinePayment, PaymentReviewDecision } from '../../types';

const formatNaira = (amount: number) =>
  `₦${amount.toLocaleString('en-NG', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;

const formatDate = (value: string) =>
  new Date(value).toLocaleString('en-NG', { dateStyle: 'medium', timeStyle: 'short' });

const TABS: { status: string; label: string }[] = [
  { status: 'needs_review', label: 'Needs review' },
  { status: 'refund_due', label: 'Refund due' },
  { status: 'pending', label: 'Awaiting transfer' },
  { status: 'paid', label: 'Paid' },
  { status: 'expired', label: 'Not paid in time' },
  { status: 'refunded', label: 'Refunded' },
  { status: '', label: 'All' },
];

const STATUS_STYLE: Record<string, string> = {
  needs_review: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
  refund_due: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
  pending: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
  paid: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
};

const REVIEW_REASONS: Record<string, string> = {
  amount_mismatch: 'Wrong amount received',
  paid_after_cancellation: 'Paid after the order was cancelled for non-payment',
};

const DECISIONS: Record<PaymentReviewDecision, { title: string; explain: string; confirm: string }> = {
  accept: {
    title: 'Accept this payment',
    explain: 'The order is marked paid and confirmed, and the customer is told. A cancelled order is reopened and its items taken from stock again.',
    confirm: 'Accept payment',
  },
  refund: {
    title: 'Refund this payment',
    explain: 'The order is cancelled (its stock returned) and the payment marked refund due. Send the refund from your Paystack dashboard, then mark it refunded here.',
    confirm: 'Mark refund due',
  },
  refunded: {
    title: 'Record the refund as made',
    explain: 'Confirm you have sent the money back to the customer.',
    confirm: 'Mark refunded',
  },
};

interface PendingDecision {
  payment: OnlinePayment;
  decision: PaymentReviewDecision;
}

const PaymentList: React.FC = () => {
  const queryClient = useQueryClient();
  const [params, setParams] = useSearchParams();
  const status = params.get('status') ?? 'needs_review';
  const [page, setPage] = useState(1);
  const [searchInput, setSearchInput] = useState('');
  const [search, setSearch] = useState('');
  const [pending, setPending] = useState<PendingDecision | null>(null);
  const [note, setNote] = useState('');
  const [checkingId, setCheckingId] = useState<number | null>(null);
  const decisionKey = useIdempotencyKey();

  useEffect(() => {
    const timer = setTimeout(() => {
      setSearch(searchInput.trim());
      setPage(1);
    }, 300);
    return () => clearTimeout(timer);
  }, [searchInput]);

  const { data, isLoading, isError } = useQuery({
    queryKey: ['online-payments', status, search, page],
    queryFn: () => getOnlinePayments({ status: status || undefined, search: search || undefined, page, limit: 20 }),
  });
  const { data: unmatched } = useQuery({ queryKey: ['unmatched-payments'], queryFn: getUnmatchedPayments });

  const refresh = () => {
    queryClient.invalidateQueries({ queryKey: ['online-payments'] });
    queryClient.invalidateQueries({ queryKey: ['orders'] });
  };

  const decide = useMutation({
    mutationFn: ({ payment, decision }: PendingDecision) =>
      resolvePaymentReview(payment.orderId, decision, note.trim(), decisionKey.current()),
    onSuccess: () => {
      decisionKey.reset();
      toast.success('Decision recorded');
      setPending(null);
      setNote('');
      refresh();
    },
    onError: (error: Error) => toast.error(error.message || 'Failed to record the decision'),
  });

  const recheck = async (payment: OnlinePayment) => {
    setCheckingId(payment.orderId);
    try {
      const result = await recheckPayment(payment.orderId);
      (result.paid ? toast.success : toast)(result.message);
      if (result.paid) refresh();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to check with Paystack');
    } finally {
      setCheckingId(null);
    }
  };

  const openDecision = (payment: OnlinePayment, decision: PaymentReviewDecision) => {
    // A different decision is a different action, so it gets its own key.
    if (pending?.payment.orderId !== payment.orderId || pending?.decision !== decision) decisionKey.reset();
    setNote('');
    setPending({ payment, decision });
  };

  const selectTab = (next: string) => {
    setPage(1);
    setParams(next ? { status: next } : { status: '' });
  };

  const payments = data?.payments ?? [];
  const counts = data?.counts ?? {};
  const totalPages = data?.pagination.totalPages ?? 0;
  const totalAll = Object.values(counts).reduce((sum, n) => sum + n, 0);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Payments</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Online transfers, and payments that need your decision.
        </p>
      </div>

      <div className="flex flex-wrap gap-2" role="tablist">
        {TABS.map((tab) => {
          const count = tab.status ? counts[tab.status] ?? 0 : totalAll;
          const active = tab.status === status;
          return (
            <button
              key={tab.label}
              role="tab"
              aria-selected={active}
              onClick={() => selectTab(tab.status)}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                active
                  ? 'bg-primary-500 text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
              }`}
            >
              {tab.label}
              <span className={`ml-2 ${active ? 'text-white/80' : 'text-gray-500 dark:text-gray-400'}`}>{count}</span>
            </button>
          );
        })}
      </div>

      <div className="relative max-w-md">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 h-5 w-5" />
        <input
          type="search"
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          placeholder="Order number, customer, phone or reference"
          className="w-full pl-10 pr-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg dark:text-white"
        />
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center text-gray-500">Loading payments…</div>
        ) : isError ? (
          <div className="p-8 text-center text-red-600">Failed to load payments.</div>
        ) : payments.length === 0 ? (
          <div className="p-8 text-center text-gray-500 dark:text-gray-400">
            {status === 'needs_review' ? 'Nothing needs your review.' : 'No payments here.'}
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 dark:bg-gray-700/50 text-left text-xs uppercase tracking-wider text-gray-500 dark:text-gray-400">
                <tr>
                  <th className="px-4 py-3">Order</th>
                  <th className="px-4 py-3">Customer</th>
                  <th className="px-4 py-3">Expected</th>
                  <th className="px-4 py-3">Received</th>
                  <th className="px-4 py-3">Payment</th>
                  <th className="px-4 py-3">Placed</th>
                  <th className="px-4 py-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                {payments.map((payment) => (
                  <tr key={payment.orderId} className="align-top">
                    <td className="px-4 py-3">
                      <Link to={`/orders/${payment.orderId}`} className="font-medium text-primary-600 dark:text-primary-400 hover:underline">
                        #{payment.orderNumber}
                      </Link>
                      <div className="text-xs text-gray-500 dark:text-gray-400 capitalize">{payment.status}</div>
                    </td>
                    <td className="px-4 py-3">
                      <div className="text-gray-900 dark:text-white">{payment.customerName}</div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">{payment.customerPhone}</div>
                    </td>
                    <td className="px-4 py-3 text-gray-900 dark:text-white">{formatNaira(payment.totalAmount)}</td>
                    <td className="px-4 py-3 text-gray-900 dark:text-white">
                      {payment.receivedAmount != null ? formatNaira(payment.receivedAmount) : '—'}
                    </td>
                    <td className="px-4 py-3">
                      <span className={`inline-block px-2 py-1 rounded-full text-xs font-medium ${STATUS_STYLE[payment.paymentStatus] ?? 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'}`}>
                        {TABS.find((t) => t.status === payment.paymentStatus)?.label ?? payment.paymentStatus}
                      </span>
                      {payment.reviewReason && (
                        <div className="mt-1 text-xs text-red-700 dark:text-red-300 flex items-start gap-1">
                          <AlertTriangle className="h-3 w-3 mt-0.5 shrink-0" />
                          {REVIEW_REASONS[payment.reviewReason] ?? payment.reviewReason}
                        </div>
                      )}
                    </td>
                    <td className="px-4 py-3 text-gray-500 dark:text-gray-400 whitespace-nowrap">{formatDate(payment.createdAt)}</td>
                    <td className="px-4 py-3">
                      <div className="flex justify-end gap-2 flex-wrap">
                        {payment.paymentStatus === 'needs_review' && (
                          <>
                            <button onClick={() => openDecision(payment, 'accept')} className="inline-flex items-center gap-1 px-3 py-1.5 rounded-lg bg-green-600 hover:bg-green-700 text-white text-xs font-medium">
                              <CheckCircle className="h-3.5 w-3.5" /> Accept
                            </button>
                            <button onClick={() => openDecision(payment, 'refund')} className="inline-flex items-center gap-1 px-3 py-1.5 rounded-lg border border-red-300 text-red-700 dark:border-red-800 dark:text-red-300 text-xs font-medium hover:bg-red-50 dark:hover:bg-red-900/20">
                              <Undo2 className="h-3.5 w-3.5" /> Refund
                            </button>
                          </>
                        )}
                        {payment.paymentStatus === 'refund_due' && (
                          <button onClick={() => openDecision(payment, 'refunded')} className="inline-flex items-center gap-1 px-3 py-1.5 rounded-lg border border-gray-300 dark:border-gray-600 text-xs font-medium hover:bg-gray-50 dark:hover:bg-gray-700">
                            <RotateCcw className="h-3.5 w-3.5" /> Mark refunded
                          </button>
                        )}
                        {payment.paymentStatus === 'pending' && payment.paymentReference && (
                          <button
                            onClick={() => recheck(payment)}
                            disabled={checkingId === payment.orderId}
                            className="inline-flex items-center gap-1 px-3 py-1.5 rounded-lg border border-gray-300 dark:border-gray-600 text-xs font-medium hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                          >
                            {checkingId === payment.orderId ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <RefreshCw className="h-3.5 w-3.5" />}
                            Re-check with Paystack
                          </button>
                        )}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        {totalPages > 1 && (
          <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200 dark:border-gray-700 text-sm">
            <button disabled={page <= 1} onClick={() => setPage((p) => p - 1)} className="px-3 py-1 rounded-lg disabled:opacity-50 hover:bg-gray-100 dark:hover:bg-gray-700">Previous</button>
            <span className="text-gray-500 dark:text-gray-400">Page {page} of {totalPages}</span>
            <button disabled={page >= totalPages} onClick={() => setPage((p) => p + 1)} className="px-3 py-1 rounded-lg disabled:opacity-50 hover:bg-gray-100 dark:hover:bg-gray-700">Next</button>
          </div>
        )}
      </div>

      {(unmatched?.payments.length ?? 0) > 0 && (
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-amber-300 dark:border-amber-800 p-4">
          <h2 className="font-semibold text-gray-900 dark:text-white flex items-center gap-2">
            <AlertTriangle className="h-4 w-4 text-amber-500" /> Payments that match no order
          </h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1 mb-3">
            Paystack reported these but no order has their reference. Look them up in your Paystack dashboard.
          </p>
          <ul className="divide-y divide-gray-200 dark:divide-gray-700 text-sm">
            {unmatched!.payments.map((p) => (
              <li key={p.id} className="py-2 flex flex-wrap justify-between gap-2">
                <span className="font-mono text-gray-900 dark:text-white break-all">{p.reference}</span>
                <span className="text-gray-600 dark:text-gray-400">
                  {p.amount != null ? formatNaira(p.amount) : '—'} · {formatDate(p.receivedAt)}
                </span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {pending && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" role="dialog" aria-modal="true" aria-labelledby="decision-title">
          <div className="w-full max-w-md bg-white dark:bg-gray-800 rounded-xl shadow-xl p-6 space-y-4">
            <div className="flex items-start justify-between gap-4">
              <h2 id="decision-title" className="text-lg font-semibold text-gray-900 dark:text-white">
                {DECISIONS[pending.decision].title}
              </h2>
              <button onClick={() => setPending(null)} aria-label="Close" className="text-gray-400 hover:text-gray-600">
                <X className="h-5 w-5" />
              </button>
            </div>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Order #{pending.payment.orderNumber}: expected {formatNaira(pending.payment.totalAmount)}
              {pending.payment.receivedAmount != null && <>, received {formatNaira(pending.payment.receivedAmount)}</>}.
            </p>
            <p className="text-sm text-gray-700 dark:text-gray-300">{DECISIONS[pending.decision].explain}</p>
            <div>
              <label htmlFor="decision-note" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Note (optional, kept in the activity log)
              </label>
              <textarea
                id="decision-note"
                value={note}
                onChange={(e) => setNote(e.target.value)}
                maxLength={500}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 dark:text-white"
              />
            </div>
            <div className="flex justify-end gap-3">
              <button onClick={() => setPending(null)} className="px-4 py-2 rounded-lg border border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700">
                Cancel
              </button>
              <button
                onClick={() => decide.mutate(pending)}
                disabled={decide.isPending}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-primary-500 hover:bg-primary-600 text-white disabled:opacity-50"
              >
                {decide.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
                {DECISIONS[pending.decision].confirm}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PaymentList;
