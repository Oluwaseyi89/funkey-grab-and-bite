// src/stores/reportsStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { 
  SalesReport, 
  AnalyticsReport, 
  DateRangeFilter, 
  ReportFilters,
  RevenueDataPoint,
  OrderDataPoint,
  CustomerDataPoint 
} from '../types';

interface ReportsState {
  // Sales Reports
  salesReports: SalesReport[];
  filteredSalesReports: SalesReport[];
  currentSalesReport: SalesReport | null;
  
  // Analytics Data
  analyticsData: AnalyticsReport | null;
  revenueData: RevenueDataPoint[];
  orderData: OrderDataPoint[];
  customerData: CustomerDataPoint[];
  
  // Filters
  filters: ReportFilters;
  dateRange: DateRangeFilter;
  isLoading: boolean;
  error: string | null;
  
  // Actions
  setSalesReports: (reports: SalesReport[]) => void;
  setAnalyticsData: (data: AnalyticsReport) => void;
  setRevenueData: (data: RevenueDataPoint[]) => void;
  setOrderData: (data: OrderDataPoint[]) => void;
  setCustomerData: (data: CustomerDataPoint[]) => void;
  
  // Filtering
  setDateRange: (range: DateRangeFilter) => void;
  setFilters: (filters: Partial<ReportFilters>) => void;
  applyFilters: () => SalesReport[];
  clearFilters: () => void;
  
  // Data Processing
  getRevenueSummary: () => {
    totalRevenue: number;
    averageDailyRevenue: number;
    revenueChange: number;
  };
  
  getOrdersSummary: () => {
    totalOrders: number;
    averageOrdersPerDay: number;
    completionRate: number;
  };
  
  getCustomersSummary: () => {
    newCustomers: number;
    returningCustomers: number;
    customerGrowth: number;
  };
  
  // Loading & Error
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useReportsStore = create<ReportsState>()(
  persist(
    (set, get) => ({
      // Initial state
      salesReports: [],
      filteredSalesReports: [],
      currentSalesReport: null,
      analyticsData: null,
      revenueData: [],
      orderData: [],
      customerData: [],
      isLoading: false,
      error: null,
      
      // Default filters
      filters: {
        dateRange: {
          startDate: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0], // 30 days ago
          endDate: new Date().toISOString().split('T')[0], // today
          period: 'month',
        },
      },
      dateRange: {
        startDate: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        endDate: new Date().toISOString().split('T')[0],
        period: 'month',
      },
      
      // Actions
      setSalesReports: (reports) => set({ 
        salesReports: reports,
        filteredSalesReports: reports,
      }),
      
      setAnalyticsData: (data) => set({ analyticsData: data }),
      
      setRevenueData: (data) => set({ revenueData: data }),
      
      setOrderData: (data) => set({ orderData: data }),
      
      setCustomerData: (data) => set({ customerData: data }),
      
      // Filtering
      setDateRange: (range) => {
        set({ 
          dateRange: range,
          filters: { ...get().filters, dateRange: range }
        });
        get().applyFilters();
      },
      
      setFilters: (filters) => {
        set({ filters: { ...get().filters, ...filters } });
        get().applyFilters();
      },
      
      applyFilters: () => {
        const { salesReports, filters } = get();
        
        let filtered = [...salesReports];
        
        // Filter by date range
        if (filters.dateRange) {
          const startDate = new Date(filters.dateRange.startDate);
          const endDate = new Date(filters.dateRange.endDate);
          
          filtered = filtered.filter(report => {
            const reportDate = new Date(report.date);
            return reportDate >= startDate && reportDate <= endDate;
          });
        }
        
        // Filter by category, orderType, paymentMethod would go here
        // when those filters are implemented
        
        set({ filteredSalesReports: filtered });
        return filtered;
      },
      
      clearFilters: () => {
        const defaultFilters = {
          dateRange: {
            startDate: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
            endDate: new Date().toISOString().split('T')[0],
            period: 'month',
          },
        };
        set({ 
          filters: defaultFilters,
          dateRange: defaultFilters.dateRange,
          filteredSalesReports: get().salesReports,
        });
      },
      
      // Data Processing
      getRevenueSummary: () => {
        const { filteredSalesReports } = get();
        
        if (filteredSalesReports.length === 0) {
          return { totalRevenue: 0, averageDailyRevenue: 0, revenueChange: 0 };
        }
        
        const totalRevenue = filteredSalesReports.reduce((sum, report) => sum + report.totalRevenue, 0);
        const averageDailyRevenue = totalRevenue / filteredSalesReports.length;
        
        // Calculate revenue change compared to previous period
        let revenueChange = 0;
        if (filteredSalesReports.length >= 2) {
          const firstHalf = filteredSalesReports.slice(0, Math.floor(filteredSalesReports.length / 2));
          const secondHalf = filteredSalesReports.slice(Math.floor(filteredSalesReports.length / 2));
          
          const firstHalfRevenue = firstHalf.reduce((sum, report) => sum + report.totalRevenue, 0);
          const secondHalfRevenue = secondHalf.reduce((sum, report) => sum + report.totalRevenue, 0);
          
          revenueChange = ((secondHalfRevenue - firstHalfRevenue) / firstHalfRevenue) * 100;
        }
        
        return {
          totalRevenue,
          averageDailyRevenue,
          revenueChange,
        };
      },
      
      getOrdersSummary: () => {
        const { filteredSalesReports } = get();
        
        if (filteredSalesReports.length === 0) {
          return { totalOrders: 0, averageOrdersPerDay: 0, completionRate: 0 };
        }
        
        const totalOrders = filteredSalesReports.reduce((sum, report) => sum + report.totalOrders, 0);
        const averageOrdersPerDay = totalOrders / filteredSalesReports.length;
        const completionRate = 95; // This would need actual completion data from backend
        
        return {
          totalOrders,
          averageOrdersPerDay,
          completionRate,
        };
      },
      
      getCustomersSummary: () => {
        const { customerData } = get();
        
        if (customerData.length === 0) {
          return { newCustomers: 0, returningCustomers: 0, customerGrowth: 0 };
        }
        
        const newCustomers = customerData.reduce((sum, data) => sum + data.newCustomers, 0);
        const returningCustomers = customerData.reduce((sum, data) => sum + data.returningCustomers, 0);
        
        // Calculate customer growth
        let customerGrowth = 0;
        if (customerData.length >= 2) {
          const firstHalf = customerData.slice(0, Math.floor(customerData.length / 2));
          const secondHalf = customerData.slice(Math.floor(customerData.length / 2));
          
          const firstHalfCustomers = firstHalf.reduce((sum, data) => sum + data.newCustomers, 0);
          const secondHalfCustomers = secondHalf.reduce((sum, data) => sum + data.newCustomers, 0);
          
          customerGrowth = ((secondHalfCustomers - firstHalfCustomers) / firstHalfCustomers) * 100;
        }
        
        return {
          newCustomers,
          returningCustomers,
          customerGrowth,
        };
      },
      
      // Loading & Error
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),
      clearError: () => set({ error: null }),
    }),
    {
      name: 'reports-storage',
      partialize: (state) => ({
        filters: state.filters,
        dateRange: state.dateRange,
      }),
    }
  )
);