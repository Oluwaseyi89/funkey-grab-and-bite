import React, { useMemo } from 'react';
import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';
import { useTheme } from '../../contexts/ThemeContext';
import { naira, nairaShort } from '../../utils/format';
import type { SalesPoint } from '../../types';

// Money received each day this week (the last 7 days) against the same
// weekday a week earlier. Built from the API's last 14 days.
const WeekChart: React.FC<{ daily: SalesPoint[] }> = ({ daily }) => {
  const { theme } = useTheme();
  const dark = theme === 'dark';

  const rows = useMemo(() => {
    const last14 = daily.slice(-14);
    const lastWeek = last14.slice(0, 7);
    const thisWeek = last14.slice(7);
    return thisWeek.map((day, i) => {
      const date = new Date(`${day.date}T12:00:00`);
      return {
        label: i === thisWeek.length - 1 ? 'Today' : date.toLocaleDateString('en-NG', { weekday: 'short' }),
        date: date.toLocaleDateString('en-NG', { weekday: 'long', day: 'numeric', month: 'short' }),
        thisWeek: day.revenue,
        lastWeek: lastWeek[i]?.revenue ?? 0,
        orders: day.orders,
      };
    });
  }, [daily]);

  const totals = rows.reduce((acc, r) => ({ thisWeek: acc.thisWeek + r.thisWeek, lastWeek: acc.lastWeek + r.lastWeek }), { thisWeek: 0, lastWeek: 0 });
  const empty = totals.thisWeek === 0 && totals.lastWeek === 0;

  return (
    <figure>
      <div className="h-56 sm:h-64" aria-hidden="true">
        {empty ? (
          <div className="flex h-full items-center justify-center rounded-xl bg-gray-50 text-sm text-gray-500 dark:bg-gray-800/50">
            No payments recorded in the last two weeks yet.
          </div>
        ) : (
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={rows} barGap={2} margin={{ top: 8, right: 4, bottom: 0, left: -12 }}>
              <CartesianGrid vertical={false} stroke={dark ? '#1e293b' : '#e2e8f0'} />
              <XAxis dataKey="label" tickLine={false} axisLine={false} tick={{ fill: '#64748b', fontSize: 12 }} />
              <YAxis tickLine={false} axisLine={false} width={56} tick={{ fill: '#64748b', fontSize: 11 }} tickFormatter={(v) => nairaShort(Number(v))} />
              <Tooltip
                cursor={{ fill: dark ? 'rgba(148,163,184,0.08)' : 'rgba(15,23,42,0.04)' }}
                contentStyle={{ borderRadius: 12, border: 'none', background: dark ? '#1e293b' : '#fff', boxShadow: '0 8px 24px rgba(15,23,42,0.15)' }}
                labelFormatter={(_, payload) => payload?.[0]?.payload?.date ?? ''}
                formatter={(value, name) => [naira(Number(value)), name === 'thisWeek' ? 'This week' : 'A week earlier']}
              />
              <Bar dataKey="lastWeek" fill={dark ? '#334155' : '#cbd5e1'} radius={[6, 6, 0, 0]} maxBarSize={22} />
              <Bar dataKey="thisWeek" fill="#c3241a" radius={[6, 6, 0, 0]} maxBarSize={22} />
            </BarChart>
          </ResponsiveContainer>
        )}
      </div>
      <figcaption className="mt-3 flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-gray-500 dark:text-gray-400">
        <span className="inline-flex items-center gap-1.5"><span className="h-2.5 w-2.5 rounded-sm bg-primary-600" aria-hidden="true" />This week {naira(totals.thisWeek)}</span>
        <span className="inline-flex items-center gap-1.5"><span className="h-2.5 w-2.5 rounded-sm bg-gray-300 dark:bg-gray-700" aria-hidden="true" />A week earlier {naira(totals.lastWeek)}</span>
      </figcaption>
      {/* The same figures as a table, for screen readers. */}
      <div className="sr-only">
      <table>
        <caption>Money received per day, this week and a week earlier</caption>
        <thead><tr><th>Day</th><th>This week</th><th>A week earlier</th></tr></thead>
        <tbody>
          {rows.map((r) => (
            <tr key={r.date}><td>{r.date}</td><td>{naira(r.thisWeek)}</td><td>{naira(r.lastWeek)}</td></tr>
          ))}
        </tbody>
      </table>
      </div>
    </figure>
  );
};

export default WeekChart;
