// One way to show money and dates across the admin.

const whole = new Intl.NumberFormat('en-NG', { maximumFractionDigits: 0 });
const exact = new Intl.NumberFormat('en-NG', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
const compact = new Intl.NumberFormat('en-NG', { notation: 'compact', maximumFractionDigits: 1 });

// ₦4,200 for whole amounts, ₦4,200.50 when there are kobo.
export const naira = (amount: number | null | undefined): string => {
  const value = Number(amount ?? 0);
  return `₦${Math.round(value * 100) % 100 === 0 ? whole.format(value) : exact.format(value)}`;
};

// ₦1.2M, ₦45K: for chart axes.
export const nairaShort = (amount: number): string => `₦${compact.format(amount)}`;

// 10 Nov 2026, 1:14 PM
export const formatDateTime = (value?: string | null): string =>
  value ? new Date(value).toLocaleString('en-NG', { dateStyle: 'medium', timeStyle: 'short' }) : '';

// 10 Nov 2026
export const formatDay = (value?: string | null): string =>
  value ? new Date(value).toLocaleDateString('en-NG', { dateStyle: 'medium' }) : '';

// 1:14 PM
export const formatTime = (value?: string | null): string =>
  value ? new Date(value).toLocaleTimeString('en-NG', { hour: 'numeric', minute: '2-digit' }) : '';

// "just now", "5 min ago", "2 h ago", then the date.
export const timeAgo = (value?: string | null, now = Date.now()): string => {
  if (!value) return '';
  const seconds = Math.max(0, Math.round((now - new Date(value).getTime()) / 1000));
  if (seconds < 45) return 'just now';
  const minutes = Math.round(seconds / 60);
  if (minutes < 60) return `${minutes} min ago`;
  const hours = Math.round(minutes / 60);
  if (hours < 24) return `${hours} h ago`;
  return formatDay(value);
};

export const plural = (count: number, one: string, many = `${one}s`): string => `${count} ${count === 1 ? one : many}`;

// The % change from previous to current, or null when there's nothing to compare with.
export const percentChange = (current: number, previous: number): number | null =>
  previous > 0 ? ((current - previous) / previous) * 100 : null;

// A local calendar date as YYYY-MM-DD (not UTC, which can be the day before).
export const isoDate = (date: Date): string =>
  `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`;

// 0 -> "12 am", 13 -> "1 pm".
export const hourLabel = (hour: number): string => `${hour % 12 === 0 ? 12 : hour % 12} ${hour < 12 ? 'am' : 'pm'}`;
