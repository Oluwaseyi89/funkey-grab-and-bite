import React from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { AlertTriangle, ArrowRight, Banknote, ChefHat, CircleDollarSign, ClipboardList, Package2, PartyPopper, ShoppingBag, Timer } from 'lucide-react';
import { getOrders } from '../api/adminApi';
import StatCard from '../components/dashboard/StatCard';
import WeekChart from '../components/dashboard/WeekChart';
import OrderCard from '../components/orders/OrderCard';
import { useDashboardOverview } from '../hooks/useDashboardOverview';
import { useAuthStore } from '../stores/authStore';
import { hourLabel, naira, percentChange, plural } from '../utils/format';

const greeting = () => {
  const hour = new Date().getHours();
  return hour < 12 ? 'Good morning' : hour < 17 ? 'Good afternoon' : 'Good evening';
};

// The first screen: what needs doing now, today's real takings, and the
// orders still in progress. Every number comes from the API.
const Dashboard: React.FC = () => {
  const user = useAuthStore((state) => state.user);
  const { data: overview, isLoading, isError, refetch } = useDashboardOverview();
  const { data: active, isLoading: activeLoading } = useQuery({
    queryKey: ['orders', 'dashboard-in-progress'],
    queryFn: () => getOrders({ page: 1, limit: 6, status: 'active' }),
  });

  const queue = overview?.queue;
  const inProgress = queue ? queue.pending + queue.confirmed + queue.preparing + queue.ready : 0;
  const today = overview?.today;
  const yesterday = overview?.yesterdaySoFar;

  // Tiles for things waiting on someone; only the ones with something in them.
  const attention = overview
    ? [
        { show: overview.queue.pending > 0, to: '/orders?status=pending', icon: ClipboardList, label: 'New orders to confirm', count: overview.queue.pending, tone: 'primary' },
        { show: overview.queue.confirmed + overview.queue.preparing > 0, to: '/orders?status=preparing', icon: ChefHat, label: 'In the kitchen', count: overview.queue.confirmed + overview.queue.preparing, tone: 'neutral' },
        { show: overview.queue.ready > 0, to: '/orders?status=ready', icon: ShoppingBag, label: 'Ready to hand over', count: overview.queue.ready, tone: 'leaf' },
        { show: overview.paymentsToReview > 0, to: '/payments', icon: AlertTriangle, label: 'Payments to check', count: overview.paymentsToReview, tone: 'danger' },
        { show: overview.lowStock > 0, to: '/inventory/alerts', icon: Package2, label: 'Items low on stock', count: overview.lowStock, tone: 'amber' },
        { show: overview.cateringToQuote > 0, to: '/catering', icon: PartyPopper, label: 'Catering to quote', count: overview.cateringToQuote, tone: 'neutral' },
      ].filter((tile) => tile.show)
    : [];

  const toneClass: Record<string, string> = {
    primary: 'bg-primary-600 text-white',
    danger: 'bg-red-600 text-white',
    leaf: 'bg-leaf-600 text-white',
    amber: 'bg-amber-500 text-white',
    neutral: 'bg-gray-900 text-white dark:bg-white dark:text-gray-900',
  };

  const busiest = overview?.hourly.reduce((best, h) => (h.orders > best.orders ? h : best), { hour: -1, orders: 0 });

  return (
    <div className="space-y-6">
      <div>
        <p className="text-sm text-gray-500 dark:text-gray-400">{new Date().toLocaleDateString('en-NG', { weekday: 'long', day: 'numeric', month: 'long' })}</p>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white md:text-3xl">
          {greeting()}{user?.username ? `, ${user.username}` : ''}
        </h1>
      </div>

      {isError && (
        <div className="card flex items-center justify-between gap-3 p-4" role="alert">
          <p className="text-sm text-gray-700 dark:text-gray-200">Today’s figures didn’t load.</p>
          <button type="button" className="btn-secondary" onClick={() => refetch()}>Try again</button>
        </div>
      )}

      {/* Needs attention */}
      {overview && (
        <section aria-labelledby="attention-title">
          <h2 id="attention-title" className="sr-only">Needs attention</h2>
          {attention.length === 0 ? (
            <p className="card flex items-center gap-3 p-4 text-sm text-gray-700 dark:text-gray-200">
              <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-leaf-50 text-leaf-600 dark:bg-leaf-500/10" aria-hidden="true">✓</span>
              Nothing waiting. New orders will appear here and in the bell.
            </p>
          ) : (
            <ul className="-mx-3 flex snap-x gap-3 overflow-x-auto px-3 pb-1 sm:mx-0 sm:grid sm:grid-cols-3 sm:overflow-visible sm:px-0 lg:grid-cols-6">
              {attention.map((tile) => {
                const Icon = tile.icon;
                return (
                  <li key={tile.label} className="w-40 shrink-0 snap-start sm:w-auto">
                    <Link to={tile.to} className={`flex h-full flex-col justify-between gap-3 rounded-2xl p-4 shadow-soft transition active:scale-[0.98] ${toneClass[tile.tone]}`}>
                      <Icon className="h-5 w-5 opacity-90" aria-hidden="true" />
                      <span>
                        <span className="block font-display text-3xl font-bold leading-none">{tile.count}</span>
                        <span className="mt-1 block text-sm font-medium opacity-90">{tile.label}</span>
                      </span>
                    </Link>
                  </li>
                );
              })}
            </ul>
          )}
        </section>
      )}

      {/* Today */}
      <section aria-labelledby="today-title">
        <h2 id="today-title" className="mb-3 text-lg font-semibold text-gray-900 dark:text-white">Today so far</h2>
        <div className="grid grid-cols-2 gap-3 lg:grid-cols-4">
          <StatCard
            label="Money received"
            value={naira(today?.revenue)}
            icon={Banknote}
            change={today && yesterday ? percentChange(today.revenue, yesterday.revenue) : null}
            comparison="vs yesterday by now"
            loading={isLoading}
          />
          <StatCard
            label="Orders placed"
            value={String(today?.orders ?? 0)}
            icon={ClipboardList}
            change={today && yesterday ? percentChange(today.orders, yesterday.orders) : null}
            comparison="vs yesterday by now"
            loading={isLoading}
          />
          <StatCard label="Average paid order" value={naira(today?.averageOrder)} icon={CircleDollarSign} loading={isLoading} />
          <StatCard
            label="Unpaid online orders"
            value={String(overview?.awaitingPayment ?? 0)}
            icon={Timer}
            loading={isLoading}
          />
        </div>
        <p className="mt-2 text-xs text-gray-500 dark:text-gray-400">
          Money received counts paid orders that weren’t cancelled. Pay-on-delivery orders count once their payment is recorded.
        </p>
      </section>

      <div className="grid gap-6 lg:grid-cols-5">
        {/* Orders in progress */}
        <section className="lg:col-span-3" aria-labelledby="progress-title">
          <div className="mb-3 flex items-center justify-between">
            <h2 id="progress-title" className="text-lg font-semibold text-gray-900 dark:text-white">
              In progress{inProgress ? ` (${inProgress})` : ''}
            </h2>
            <Link to="/orders" className="inline-flex min-h-[44px] items-center gap-1 text-sm font-semibold text-primary-600 dark:text-primary-400">
              All orders <ArrowRight className="h-4 w-4" aria-hidden="true" />
            </Link>
          </div>
          {activeLoading ? (
            <div className="grid gap-3 sm:grid-cols-2">{[0, 1].map((i) => <div key={i} className="skeleton h-44 rounded-2xl" />)}</div>
          ) : (active?.data ?? []).length === 0 ? (
            <p className="card p-6 text-center text-sm text-gray-500">No orders in progress right now.</p>
          ) : (
            <ul className="grid gap-3 sm:grid-cols-2">
              {(active?.data ?? []).map((order) => (
                <li key={order.id}><OrderCard order={order} /></li>
              ))}
            </ul>
          )}
        </section>

        {/* Week and best sellers */}
        <div className="space-y-6 lg:col-span-2">
          <section className="card p-4" aria-labelledby="week-title">
            <h2 id="week-title" className="font-semibold text-gray-900 dark:text-white">Money received, this week</h2>
            <p className="mb-3 text-xs text-gray-500 dark:text-gray-400">Each day compared with the same day last week</p>
            {overview ? <WeekChart daily={overview.daily} /> : <div className="skeleton h-56" />}
          </section>

          <section className="card p-4" aria-labelledby="top-title">
            <div className="flex items-baseline justify-between">
              <h2 id="top-title" className="font-semibold text-gray-900 dark:text-white">Best sellers</h2>
              <span className="text-xs text-gray-500">Last 7 days</span>
            </div>
            {overview && overview.topItems.length > 0 ? (
              <ol className="mt-3 space-y-2">
                {overview.topItems.map((item, index) => {
                  const max = overview.topItems[0].quantity || 1;
                  return (
                    <li key={item.menuItemId} className="text-sm">
                      <div className="flex justify-between gap-3">
                        <span className="truncate text-gray-800 dark:text-gray-200"><span className="mr-2 text-gray-400">{index + 1}</span>{item.name}</span>
                        <span className="shrink-0 tabular-nums text-gray-500">{plural(item.quantity, 'sold', 'sold')}</span>
                      </div>
                      <div className="mt-1 h-1.5 rounded-full bg-gray-100 dark:bg-gray-800" aria-hidden="true">
                        <div className="h-1.5 rounded-full bg-primary-500" style={{ width: `${(item.quantity / max) * 100}%` }} />
                      </div>
                    </li>
                  );
                })}
              </ol>
            ) : (
              <p className="mt-3 text-sm text-gray-500">{overview ? 'No orders in the last 7 days.' : 'Loading…'}</p>
            )}
            {busiest && busiest.hour >= 0 && (
              <p className="mt-4 rounded-xl bg-gray-50 px-3 py-2 text-xs text-gray-600 dark:bg-gray-800/60 dark:text-gray-300">
                Busiest hour today: <strong>{hourLabel(busiest.hour)}</strong> ({plural(busiest.orders, 'order')})
              </p>
            )}
          </section>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
