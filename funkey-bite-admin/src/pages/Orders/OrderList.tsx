import React, { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { CalendarDays, ChevronLeft, ChevronRight, ClipboardList, Download, Search, X } from 'lucide-react';
import toast from 'react-hot-toast';
import { getOrders } from '../../api/adminApi';
import PageHeader from '../../components/ui/PageHeader';
import OrderCard from '../../components/orders/OrderCard';
import { ORDER_STATUS, PAYMENT_STATUS } from '../../components/ui/StatusPill';
import { useDashboardOverview } from '../../hooks/useDashboardOverview';
import { downloadCsv } from '../../utils/csv';
import { formatDateTime, isoDate } from '../../utils/format';

// "active" is everything the kitchen still has to finish.
const TABS = [
  { id: 'active', label: 'In progress' },
  { id: 'pending', label: 'New' },
  { id: 'confirmed', label: 'Confirmed' },
  { id: 'preparing', label: 'Cooking' },
  { id: 'ready', label: 'Ready' },
  { id: 'completed', label: 'Completed' },
  { id: 'cancelled', label: 'Cancelled' },
  { id: '', label: 'All' },
];

const PAGE_SIZE = 24;

const OrderList: React.FC = () => {
  // The tab and filters live in the address, so Back and refresh keep them.
  const [params, setParams] = useSearchParams();
  // No status in the address means "In progress"; "all" means every order.
  const rawStatus = params.get('status');
  const status = rawStatus === null ? 'active' : rawStatus === 'all' ? '' : rawStatus;
  const unpaidOnDelivery = params.get('unpaid') === '1';
  const page = Math.max(1, Number(params.get('page') ?? 1));
  const [search, setSearch] = useState(params.get('q') ?? '');
  const query = params.get('q') ?? '';
  const { data: overview } = useDashboardOverview();

  const update = (changes: Record<string, string | null>) => {
    const next = new URLSearchParams(params);
    Object.entries(changes).forEach(([key, value]) => (value === null || value === '' ? next.delete(key) : next.set(key, value)));
    if (!('page' in changes)) next.delete('page');
    setParams(next, { replace: true });
  };

  // Search after a short pause in typing.
  useEffect(() => {
    const id = setTimeout(() => {
      if (search.trim() !== query) update({ q: search.trim() || null });
    }, 350);
    return () => clearTimeout(id);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [search]);

  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['orders', status, unpaidOnDelivery, page, query],
    queryFn: () =>
      getOrders({
        page,
        limit: PAGE_SIZE,
        status: status || undefined,
        search: query || undefined,
        ...(unpaidOnDelivery ? { paymentMethod: 'cash' as const, paymentStatus: 'pending' as const } : {}),
      }),
    placeholderData: (previous) => previous,
  });

  const orders = data?.data ?? [];
  const totalPages = data?.pagination?.totalPages ?? 0;
  const total = data?.pagination?.total ?? orders.length;

  const tabCount = (id: string) => {
    if (!overview) return 0;
    const q = overview.queue;
    return id === 'active' ? q.pending + q.confirmed + q.preparing + q.ready : id === 'pending' ? q.pending : id === 'confirmed' ? q.confirmed : id === 'preparing' ? q.preparing : id === 'ready' ? q.ready : 0;
  };

  const exportCsv = () => {
    if (!orders.length) {
      toast.error('There are no orders here to export');
      return;
    }
    downloadCsv(
      `orders-${status || 'all'}-${isoDate(new Date())}.csv`,
      ['Order', 'Placed', 'Customer', 'Phone', 'Type', 'Status', 'Payment method', 'Payment', 'Total (NGN)', 'Items'],
      orders.map((o) => [
        o.orderNumber,
        formatDateTime(o.createdAt),
        o.customerName,
        o.customerPhone,
        o.orderType,
        ORDER_STATUS[o.status]?.label ?? o.status,
        o.paymentMethod,
        PAYMENT_STATUS[o.paymentStatus]?.label ?? o.paymentStatus,
        o.totalAmount,
        (o.items ?? []).map((i) => `${i.quantity}x ${i.name}`).join('; '),
      ]),
    );
  };

  return (
    <div>
      <PageHeader
        title="Orders"
        description={isLoading ? 'Loading…' : `${total} ${total === 1 ? 'order' : 'orders'}${status ? ` · ${TABS.find((t) => t.id === status)?.label.toLowerCase()}` : ''}`}
        actions={
          <>
            <Link to="/orders/calendar" className="btn-secondary">
              <CalendarDays className="h-4 w-4" aria-hidden="true" /> Calendar
            </Link>
            <button type="button" className="btn-secondary" onClick={exportCsv}>
              <Download className="h-4 w-4" aria-hidden="true" /> Export this list
            </button>
          </>
        }
      />

      {/* Tabs and search stay in reach while scrolling on phones. */}
      <div className="sticky top-14 z-20 -mx-3 mb-4 space-y-3 border-b border-gray-200 bg-gray-50/95 px-3 py-3 backdrop-blur dark:border-gray-800 dark:bg-gray-950/95 sm:-mx-4 sm:px-4 md:top-16 lg:mx-0 lg:rounded-2xl lg:border lg:px-4">
        <nav className="scroll-row" aria-label="Order status">
          {TABS.map((tab) => {
            const count = tabCount(tab.id);
            const active = status === tab.id;
            return (
              <button
                key={tab.id || 'all'}
                type="button"
                className={`chip ${active ? 'chip-active' : ''}`}
                aria-pressed={active}
                onClick={() => update({ status: tab.id === 'active' ? null : tab.id || 'all' })}
              >
                {tab.label}
                {count > 0 && <span className={`rounded-full px-1.5 text-xs font-bold ${active ? 'bg-white/20' : 'bg-gray-100 dark:bg-gray-800'}`}>{count}</span>}
              </button>
            );
          })}
        </nav>
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
          <div className="relative flex-1">
            <label htmlFor="order-search" className="sr-only">Search orders</label>
            <Search className="pointer-events-none absolute left-3.5 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400" aria-hidden="true" />
            <input
              id="order-search"
              type="search"
              className="input-primary pl-11 pr-11"
              placeholder="Order number, name or phone"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
            {search && (
              <button type="button" className="icon-btn absolute right-0.5 top-1/2 -translate-y-1/2" aria-label="Clear search" onClick={() => setSearch('')}>
                <X className="h-5 w-5" aria-hidden="true" />
              </button>
            )}
          </div>
          <label className="flex min-h-[44px] items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
            <input type="checkbox" className="h-5 w-5 accent-primary-600" checked={unpaidOnDelivery} onChange={(e) => update({ unpaid: e.target.checked ? '1' : null })} />
            Pay on delivery, not collected yet
          </label>
        </div>
      </div>

      {isLoading ? (
        <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3" aria-busy="true" aria-label="Loading orders">
          {Array.from({ length: 6 }).map((_, i) => <div key={i} className="skeleton h-44 rounded-2xl" />)}
        </div>
      ) : isError ? (
        <div className="card p-8 text-center" role="alert">
          <p className="font-medium text-gray-900 dark:text-white">The orders didn’t load</p>
          <button type="button" className="btn-primary mt-4" onClick={() => refetch()}>Try again</button>
        </div>
      ) : orders.length === 0 ? (
        <div className="card p-10 text-center">
          <ClipboardList className="mx-auto h-10 w-10 text-gray-300" aria-hidden="true" />
          <p className="mt-3 font-medium text-gray-900 dark:text-white">
            {query ? `No orders match “${query}”` : status === 'active' ? 'All caught up' : 'No orders here'}
          </p>
          <p className="mt-1 text-sm text-gray-500">{status === 'active' && !query ? 'New orders will appear here the moment they come in.' : 'Try another tab or search.'}</p>
        </div>
      ) : (
        <ul className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
          {orders.map((order) => (
            <li key={order.id}>
              <OrderCard order={order} />
            </li>
          ))}
        </ul>
      )}

      {totalPages > 1 && (
        <nav className="mt-5 flex items-center justify-between" aria-label="Pages">
          <button type="button" className="btn-secondary" disabled={page <= 1} onClick={() => update({ page: String(page - 1) })}>
            <ChevronLeft className="h-4 w-4" aria-hidden="true" /> Newer
          </button>
          <span className="text-sm text-gray-500">Page {page} of {totalPages}</span>
          <button type="button" className="btn-secondary" disabled={page >= totalPages} onClick={() => update({ page: String(page + 1) })}>
            Older <ChevronRight className="h-4 w-4" aria-hidden="true" />
          </button>
        </nav>
      )}
    </div>
  );
};

export default OrderList;
