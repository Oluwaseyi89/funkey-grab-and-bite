import React, { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Bar, BarChart, CartesianGrid, ComposedChart, Line, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';
import { Banknote, CircleDollarSign, ClipboardList, Download, XCircle } from 'lucide-react';
import toast from 'react-hot-toast';
import { getSalesSummary } from '../../api/adminApi';
import PageHeader from '../../components/ui/PageHeader';
import StatCard from '../../components/dashboard/StatCard';
import { useTheme } from '../../contexts/ThemeContext';
import { downloadCsv } from '../../utils/csv';
import { formatDay, hourLabel, isoDate, naira, nairaShort, percentChange } from '../../utils/format';
import type { ReportGroup, SalesPoint } from '../../types';

type Preset = 'today' | '7d' | '30d' | 'month' | 'lastMonth' | 'year' | 'custom';

const PRESETS: { id: Preset; label: string }[] = [
  { id: 'today', label: 'Today' },
  { id: '7d', label: 'Last 7 days' },
  { id: '30d', label: 'Last 30 days' },
  { id: 'month', label: 'This month' },
  { id: 'lastMonth', label: 'Last month' },
  { id: 'year', label: 'This year' },
  { id: 'custom', label: 'Custom' },
];

const METHOD_LABELS: Record<string, string> = { cash: 'Pay on delivery / pickup', transfer: 'Bank transfer', card: 'Card' };

// The dates a preset covers, both included, in the device's calendar.
const presetRange = (preset: Preset): { from: string; to: string } => {
  const today = new Date();
  const days = (n: number) => { const d = new Date(today); d.setDate(d.getDate() - n); return d; };
  switch (preset) {
    case 'today': return { from: isoDate(today), to: isoDate(today) };
    case '7d': return { from: isoDate(days(6)), to: isoDate(today) };
    case 'month': return { from: isoDate(new Date(today.getFullYear(), today.getMonth(), 1)), to: isoDate(today) };
    case 'lastMonth': return { from: isoDate(new Date(today.getFullYear(), today.getMonth() - 1, 1)), to: isoDate(new Date(today.getFullYear(), today.getMonth(), 0)) };
    case 'year': return { from: isoDate(new Date(today.getFullYear(), 0, 1)), to: isoDate(today) };
    default: return { from: isoDate(days(29)), to: isoDate(today) };
  }
};

const pointLabel = (point: SalesPoint, group: ReportGroup) => {
  const date = new Date(`${point.date}T12:00:00`);
  if (group === 'month') return date.toLocaleDateString('en-NG', { month: 'short', year: '2-digit' });
  if (group === 'week') return `w/c ${date.toLocaleDateString('en-NG', { day: 'numeric', month: 'short' })}`;
  return date.toLocaleDateString('en-NG', { day: 'numeric', month: 'short' });
};

// Sales for a chosen period, from the API: money received (paid, not
// cancelled), orders placed, and a comparison with the period before.
const SalesReport: React.FC = () => {
  const { theme } = useTheme();
  const dark = theme === 'dark';
  const [preset, setPreset] = useState<Preset>('30d');
  const [custom, setCustom] = useState(presetRange('30d'));
  const range = preset === 'custom' ? custom : presetRange(preset);
  const spanDays = Math.round((new Date(range.to).getTime() - new Date(range.from).getTime()) / 86_400_000) + 1;
  const [groupChoice, setGroupChoice] = useState<ReportGroup | 'auto'>('auto');
  const group: ReportGroup = groupChoice !== 'auto' ? groupChoice : spanDays > 120 ? 'month' : spanDays > 45 ? 'week' : 'day';
  const validRange = range.from && range.to && range.from <= range.to;

  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['sales-summary', range.from, range.to, group],
    queryFn: () => getSalesSummary({ from: range.from, to: range.to, group }),
    enabled: Boolean(validRange),
    placeholderData: (previous) => previous,
  });

  const series = useMemo(() => (data?.series ?? []).map((p) => ({ ...p, label: pointLabel(p, group) })), [data, group]);
  const totals = data?.totals;
  const previous = data?.previous;
  const comparison = data ? `vs ${formatDay(data.previousFrom)} – ${formatDay(data.previousTo)}` : '';
  const busyHours = (data?.hourly ?? []).filter((h) => h.orders > 0);
  const totalMethodRevenue = (data?.paymentMethods ?? []).reduce((s, m) => s + m.revenue, 0) || 1;

  const exportCsv = () => {
    if (!data) return;
    downloadCsv(
      `sales-${data.from}-to-${data.to}-by-${group}.csv`,
      [group === 'day' ? 'Date' : group === 'week' ? 'Week starting' : 'Month starting', 'Money received (NGN)', 'Paid orders', 'Orders placed'],
      data.series.map((p) => [p.date, p.revenue, p.paidOrders, p.orders]),
    );
    toast.success('Report downloaded');
  };

  const gridStroke = dark ? '#1e293b' : '#e2e8f0';
  const tooltipStyle = { borderRadius: 12, border: 'none', background: dark ? '#1e293b' : '#fff', boxShadow: '0 8px 24px rgba(15,23,42,0.15)' };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Reports"
        description="Money received counts paid orders that weren’t cancelled, on the day they were paid."
        actions={
          <button type="button" className="btn-secondary" onClick={exportCsv} disabled={!data}>
            <Download className="h-4 w-4" aria-hidden="true" /> Download CSV
          </button>
        }
      />

      {/* Period */}
      <section className="card space-y-3 p-4" aria-label="Report period">
        <div className="scroll-row" role="group" aria-label="Period">
          {PRESETS.map((p) => (
            <button key={p.id} type="button" className={`chip ${preset === p.id ? 'chip-active' : ''}`} aria-pressed={preset === p.id} onClick={() => setPreset(p.id)}>
              {p.label}
            </button>
          ))}
        </div>
        <div className="grid gap-3 sm:grid-cols-3">
          {preset === 'custom' && (
            <>
              <div>
                <label htmlFor="report-from" className="form-label">From</label>
                <input id="report-from" type="date" className="input-primary" value={custom.from} max={custom.to} onChange={(e) => setCustom({ ...custom, from: e.target.value })} />
              </div>
              <div>
                <label htmlFor="report-to" className="form-label">To</label>
                <input id="report-to" type="date" className="input-primary" value={custom.to} min={custom.from} max={isoDate(new Date())} onChange={(e) => setCustom({ ...custom, to: e.target.value })} />
              </div>
            </>
          )}
          <div>
            <label htmlFor="report-group" className="form-label">Show by</label>
            <select id="report-group" className="input-primary" value={groupChoice} onChange={(e) => setGroupChoice(e.target.value as ReportGroup | 'auto')}>
              <option value="auto">Automatic ({group})</option>
              <option value="day">Day</option>
              <option value="week">Week</option>
              <option value="month">Month</option>
            </select>
          </div>
        </div>
        {!validRange && <p className="text-sm text-red-600" role="alert">The start date must be on or before the end date.</p>}
      </section>

      {isError && (
        <div className="card flex items-center justify-between gap-3 p-4" role="alert">
          <p className="text-sm">The report didn’t load.</p>
          <button type="button" className="btn-secondary" onClick={() => refetch()}>Try again</button>
        </div>
      )}

      {/* Totals */}
      <section aria-label="Totals" className="grid grid-cols-2 gap-3 lg:grid-cols-4">
        <StatCard label="Money received" value={naira(totals?.revenue)} icon={Banknote} change={totals && previous ? percentChange(totals.revenue, previous.revenue) : null} comparison={comparison} loading={isLoading} />
        <StatCard label="Orders placed" value={String(totals?.orders ?? 0)} icon={ClipboardList} change={totals && previous ? percentChange(totals.orders, previous.orders) : null} comparison={comparison} loading={isLoading} />
        <StatCard label="Average paid order" value={naira(totals?.averageOrder)} icon={CircleDollarSign} change={totals && previous ? percentChange(totals.averageOrder, previous.averageOrder) : null} comparison={comparison} loading={isLoading} />
        <StatCard label="Cancelled" value={String(totals?.cancelled ?? 0)} icon={XCircle} comparison={totals && totals.orders + totals.cancelled > 0 ? `${((totals.cancelled / (totals.orders + totals.cancelled)) * 100).toFixed(0)}% of orders` : ''} loading={isLoading} />
      </section>

      {/* Over time */}
      <section className="card p-4" aria-labelledby="trend-title">
        <h2 id="trend-title" className="font-semibold text-gray-900 dark:text-white">Money received and orders, by {group}</h2>
        <p className="mb-3 text-xs text-gray-500">Bars are money received; the line is orders placed. Days with no sales show as zero.</p>
        <div className="h-64 sm:h-80" aria-hidden="true">
          {series.length === 0 ? (
            <div className="skeleton h-full" />
          ) : (
            <ResponsiveContainer width="100%" height="100%">
              <ComposedChart data={series} margin={{ top: 8, right: 0, bottom: 0, left: -8 }}>
                <CartesianGrid vertical={false} stroke={gridStroke} />
                <XAxis dataKey="label" tickLine={false} axisLine={false} tick={{ fill: '#64748b', fontSize: 11 }} minTickGap={16} />
                <YAxis yAxisId="money" tickLine={false} axisLine={false} width={56} tick={{ fill: '#64748b', fontSize: 11 }} tickFormatter={(v) => nairaShort(Number(v))} />
                <YAxis yAxisId="orders" orientation="right" tickLine={false} axisLine={false} width={32} allowDecimals={false} tick={{ fill: '#64748b', fontSize: 11 }} />
                <Tooltip
                  contentStyle={tooltipStyle}
                  cursor={{ fill: dark ? 'rgba(148,163,184,0.08)' : 'rgba(15,23,42,0.04)' }}
                  formatter={(value, name) => (name === 'revenue' ? [naira(Number(value)), 'Money received'] : [String(value), 'Orders placed'])}
                />
                <Bar yAxisId="money" dataKey="revenue" fill="#c3241a" radius={[6, 6, 0, 0]} maxBarSize={36} />
                <Line yAxisId="orders" type="monotone" dataKey="orders" stroke={dark ? '#cbd5e1' : '#334155'} strokeWidth={2} dot={series.length <= 31} />
              </ComposedChart>
            </ResponsiveContainer>
          )}
        </div>
        {/* The figures behind the chart, readable and on any screen size. */}
        {series.length > 0 && (
          <details className="mt-4 rounded-xl border border-gray-200 dark:border-gray-800">
            <summary className="flex min-h-[44px] cursor-pointer items-center px-4 text-sm font-medium">Show the figures ({series.length} rows)</summary>
            <div className="max-h-80 overflow-auto">
              <table className="w-full text-sm">
                <thead className="sticky top-0 bg-gray-50 text-left text-xs uppercase text-gray-500 dark:bg-gray-800">
                  <tr><th className="px-4 py-2">{group === 'day' ? 'Date' : group === 'week' ? 'Week of' : 'Month'}</th><th className="px-4 py-2 text-right">Received</th><th className="px-4 py-2 text-right">Orders</th></tr>
                </thead>
                <tbody className="divide-y divide-gray-100 dark:divide-gray-800">
                  {series.map((p) => (
                    <tr key={p.date}><td className="px-4 py-2">{p.label}</td><td className="px-4 py-2 text-right tabular-nums">{naira(p.revenue)}</td><td className="px-4 py-2 text-right tabular-nums">{p.orders}</td></tr>
                  ))}
                </tbody>
              </table>
            </div>
          </details>
        )}
      </section>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Best sellers */}
        <section className="card p-4" aria-labelledby="best-title">
          <h2 id="best-title" className="font-semibold text-gray-900 dark:text-white">Best sellers</h2>
          {(data?.topItems ?? []).length === 0 ? (
            <p className="mt-3 text-sm text-gray-500">No sales in this period.</p>
          ) : (
            <ol className="mt-3 divide-y divide-gray-100 dark:divide-gray-800">
              {data!.topItems.map((item, i) => (
                <li key={item.menuItemId} className="flex items-center justify-between gap-3 py-2 text-sm">
                  <span className="truncate"><span className="mr-2 text-gray-400">{i + 1}</span>{item.name}</span>
                  <span className="shrink-0 text-right tabular-nums">
                    <span className="font-semibold">{item.quantity} sold</span>
                    <span className="block text-xs text-gray-500">{naira(item.revenue)}</span>
                  </span>
                </li>
              ))}
            </ol>
          )}
        </section>

        {/* Busy hours */}
        <section className="card p-4" aria-labelledby="hours-title">
          <h2 id="hours-title" className="font-semibold text-gray-900 dark:text-white">Busy hours</h2>
          <p className="text-xs text-gray-500">Orders placed in each hour of the day, across the period</p>
          {busyHours.length === 0 ? (
            <p className="mt-3 text-sm text-gray-500">No orders in this period.</p>
          ) : (
            <div className="mt-3 h-48" aria-hidden="true">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={data!.hourly.filter((h) => h.hour >= Math.min(...busyHours.map((b) => b.hour)) && h.hour <= Math.max(...busyHours.map((b) => b.hour)))} margin={{ top: 4, right: 0, bottom: 0, left: -20 }}>
                  <CartesianGrid vertical={false} stroke={gridStroke} />
                  <XAxis dataKey="hour" tickLine={false} axisLine={false} tick={{ fill: '#64748b', fontSize: 11 }} tickFormatter={(h) => hourLabel(Number(h))} />
                  <YAxis allowDecimals={false} tickLine={false} axisLine={false} width={40} tick={{ fill: '#64748b', fontSize: 11 }} />
                  <Tooltip contentStyle={tooltipStyle} labelFormatter={(h) => `${hourLabel(Number(h))}–${hourLabel((Number(h) + 1) % 24)}`} formatter={(v) => [String(v), 'Orders']} />
                  <Bar dataKey="orders" fill={dark ? '#94a3b8' : '#475569'} radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          )}
          {busyHours.length > 0 && (
            <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">
              Busiest: <strong>{hourLabel(busyHours.reduce((a, b) => (b.orders > a.orders ? b : a)).hour)}</strong>
            </p>
          )}
        </section>

        {/* How customers paid */}
        <section className="card p-4" aria-labelledby="methods-title">
          <h2 id="methods-title" className="font-semibold text-gray-900 dark:text-white">How customers paid</h2>
          {(data?.paymentMethods ?? []).length === 0 ? (
            <p className="mt-3 text-sm text-gray-500">No payments in this period.</p>
          ) : (
            <ul className="mt-3 space-y-3">
              {data!.paymentMethods.map((m) => (
                <li key={m.key} className="text-sm">
                  <div className="flex justify-between gap-3">
                    <span>{METHOD_LABELS[m.key] ?? m.key}</span>
                    <span className="tabular-nums">{naira(m.revenue)} <span className="text-gray-500">· {m.orders}</span></span>
                  </div>
                  <div className="mt-1 h-2 rounded-full bg-gray-100 dark:bg-gray-800" aria-hidden="true">
                    <div className="h-2 rounded-full bg-primary-500" style={{ width: `${(m.revenue / totalMethodRevenue) * 100}%` }} />
                  </div>
                </li>
              ))}
            </ul>
          )}
        </section>

        {/* Pickup and delivery */}
        <section className="card p-4" aria-labelledby="types-title">
          <h2 id="types-title" className="font-semibold text-gray-900 dark:text-white">Pickup and delivery</h2>
          {(data?.orderTypes ?? []).length === 0 ? (
            <p className="mt-3 text-sm text-gray-500">No orders in this period.</p>
          ) : (
            <ul className="mt-3 grid grid-cols-2 gap-3">
              {data!.orderTypes.map((t) => (
                <li key={t.key} className="rounded-xl bg-gray-50 p-3 dark:bg-gray-800/60">
                  <p className="text-xs capitalize text-gray-500">{t.key}</p>
                  <p className="font-display text-xl font-bold">{t.orders}</p>
                  <p className="text-xs text-gray-500">{naira(t.revenue)} ordered</p>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>
    </div>
  );
};

export default SalesReport;
