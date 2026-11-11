// Sales figures from the API. Revenue is money received (paid, not
// cancelled) on the day it was paid; orders are those placed, cancelled ones
// left out. Days are in the business timezone.

export interface PeriodTotals {
  revenue: number;
  orders: number;
  paidOrders: number;
  averageOrder: number;
  cancelled: number;
}

// A day, week (starting Monday) or month; empty ones come back as zero.
export interface SalesPoint {
  date: string;
  revenue: number;
  orders: number;
  paidOrders: number;
}

export interface HourPoint {
  hour: number;
  orders: number;
}

export interface TopItem {
  menuItemId: number;
  name: string;
  quantity: number;
  revenue: number;
}

export interface Share {
  key: string;
  orders: number;
  revenue: number;
}

export interface OrderQueue {
  pending: number;
  confirmed: number;
  preparing: number;
  ready: number;
}

export interface DashboardOverview {
  timezone: string;
  today: PeriodTotals;
  // Yesterday up to the same time of day.
  yesterdaySoFar: PeriodTotals;
  queue: OrderQueue;
  paymentsToReview: number;
  awaitingPayment: number;
  lowStock: number;
  cateringToQuote: number;
  // The last 14 days including today, oldest first.
  daily: SalesPoint[];
  hourly: HourPoint[];
  topItems: TopItem[];
}

export type ReportGroup = 'day' | 'week' | 'month';

export interface SalesSummary {
  from: string;
  to: string;
  group: ReportGroup;
  timezone: string;
  totals: PeriodTotals;
  previous: PeriodTotals;
  previousFrom: string;
  previousTo: string;
  series: SalesPoint[];
  hourly: HourPoint[];
  topItems: TopItem[];
  paymentMethods: Share[];
  orderTypes: Share[];
}
