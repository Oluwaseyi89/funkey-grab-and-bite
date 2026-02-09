// src/types/reports.types.ts
// export interface SalesReport {
//     date: string;
//     totalOrders: number;
//     totalRevenue: number;
//     averageOrder: number;
//   }
  
  export interface AnalyticsReport {
    period: 'daily' | 'weekly' | 'monthly' | 'yearly';
    revenueData: RevenueDataPoint[];
    orderData: OrderDataPoint[];
    customerData: CustomerDataPoint[];
    popularItems: PopularItem[];
  }
  
  export interface RevenueDataPoint {
    date: string;
    revenue: number;
    orders: number;
    averageOrder: number;
  }
  
  export interface OrderDataPoint {
    date: string;
    count: number;
    completed: number;
    cancelled: number;
  }
  
  export interface CustomerDataPoint {
    date: string;
    newCustomers: number;
    returningCustomers: number;
    totalCustomers: number;
  }
  
  export interface PopularItem {
    id: number;
    name: string;
    category: string;
    quantitySold: number;
    revenue: number;
    percentageChange: number;
  }
  
  export interface DateRangeFilter {
    startDate: string;
    endDate: string;
    period: 'day' | 'week' | 'month' | 'year';
  }
  
  export interface ReportFilters {
    dateRange: DateRangeFilter;
    category?: string;
    orderType?: 'pickup' | 'delivery' | 'catering';
    paymentMethod?: string;
  }