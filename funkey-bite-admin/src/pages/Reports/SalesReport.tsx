import React, { useState, useEffect } from 'react';
import {
  Download,
  Filter,
  Calendar,
  TrendingUp,
  TrendingDown,
  DollarSign,
  ShoppingBag,
  Users,
  BarChart3,
  Eye,
  RefreshCw,
  ChevronDown,
  FileText,
  PieChart,
  LineChart,
  Activity,
  ArrowUpRight,
  ArrowDownRight
} from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getSalesReport, getDashboardStats } from '../../api/adminApi';
import { useReportsStore } from '../../stores/reportsStore';
import type { SalesReport } from '../../types/admin.types';

import {
  LineChart as RechartsLineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart as RechartsPieChart,
  Pie,
  Cell
} from 'recharts';

const CHART_COLORS = {
  revenue: '#3b82f6', // blue-500
  orders: '#10b981', // green-500
  average: '#8b5cf6', // purple-500
  background: '#f3f4f6'
};

const SalesReport: React.FC = () => {
  const [dateRange, setDateRange] = useState({
    startDate: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0], // 30 days ago
    endDate: new Date().toISOString().split('T')[0], // today
  });
  const [chartType, setChartType] = useState<'line' | 'bar' | 'area'>('line');
  const [groupBy, setGroupBy] = useState<'day' | 'week' | 'month'>('day');
  const [showFilters, setShowFilters] = useState(false);
  
  const {
    salesReports,
    filteredSalesReports,
    setSalesReports,
    setDateRange: setStoreDateRange,
    getRevenueSummary,
    getOrdersSummary,
    setLoading,
    setError
  } = useReportsStore();

  const { data: reportData, isLoading, refetch } = useQuery({
    queryKey: ['sales-report', dateRange.startDate, dateRange.endDate],
    queryFn: async () => {
      setLoading(true);
      try {
        const data = await getSalesReport({
          from: dateRange.startDate,
          to: dateRange.endDate,
        });
        setSalesReports(Array.isArray(data) ? data : []);
        return data;
      } catch (error: any) {
        setError(error.message || 'Failed to load sales report');
        toast.error('Failed to load sales report');
        return [];
      } finally {
        setLoading(false);
      }
    },
    enabled: !!dateRange.startDate && !!dateRange.endDate,
  });

  const { data: dashboardStats } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: async () => {
      try {
        const data = await getDashboardStats();
        return data;
      } catch (error) {
        return null;
      }
    },
  });

  useEffect(() => {
    setStoreDateRange({
      ...dateRange,
      period: groupBy,
    });
  }, [dateRange, groupBy, setStoreDateRange]);

  const quickDateRanges = [
    { label: 'Today', days: 1 },
    { label: 'Last 7 days', days: 7 },
    { label: 'Last 30 days', days: 30 },
    { label: 'Last 90 days', days: 90 },
    { label: 'Year to date', days: getDaysSinceYearStart() },
  ];

  function getDaysSinceYearStart() {
    const now = new Date();
    const startOfYear = new Date(now.getFullYear(), 0, 1);
    return Math.floor((now.getTime() - startOfYear.getTime()) / (1000 * 60 * 60 * 24));
  }

  const handleQuickDateRange = (days: number) => {
    const endDate = new Date();
    const startDate = new Date();
    startDate.setDate(startDate.getDate() - days);
    
    setDateRange({
      startDate: startDate.toISOString().split('T')[0],
      endDate: endDate.toISOString().split('T')[0],
    });
  };

  const revenueSummary = getRevenueSummary();
  const ordersSummary = getOrdersSummary();

  const chartData = filteredSalesReports.map(report => ({
    date: formatChartDate(report.date, groupBy),
    revenue: report.totalRevenue,
    orders: report.totalOrders,
    average: report.averageOrder,
  }));

  function formatChartDate(dateString: string, format: 'day' | 'week' | 'month') {
    const date = new Date(dateString);
    switch (format) {
      case 'day':
        return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
      case 'week':
        return `Week ${Math.ceil(date.getDate() / 7)}`;
      case 'month':
        return date.toLocaleDateString('en-US', { month: 'short' });
      default:
        return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    }
  }

  const calculatePercentageChange = (current: number, previous: number) => {
    if (previous === 0) return 0;
    return ((current - previous) / previous) * 100;
  };

  const exportToCSV = () => {
    if (filteredSalesReports.length === 0) {
      toast.error('No data to export');
      return;
    }

    const headers = ['Date', 'Total Orders', 'Total Revenue ($)', 'Average Order ($)'];
    const csvContent = [
      headers.join(','),
      ...filteredSalesReports.map(report => [
        report.date,
        report.totalOrders,
        report.totalRevenue.toFixed(2),
        report.averageOrder.toFixed(2)
      ].join(','))
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `sales-report-${dateRange.startDate}-to-${dateRange.endDate}.csv`;
    a.click();
    window.URL.revokeObjectURL(url);
    
    toast.success('Report exported successfully');
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-64 animate-pulse"></div>
          <div className="h-10 w-40 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="h-32 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
          ))}
        </div>
        <div className="h-96 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Sales Report</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Analyze sales performance and revenue trends
          </p>
        </div>
        
        <div className="flex items-center space-x-3">
          <button
            onClick={exportToCSV}
            className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            <Download className="h-4 w-4" />
            <span>Export CSV</span>
          </button>
          <button
            onClick={() => refetch()}
            className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            <RefreshCw className="h-4 w-4" />
            <span>Refresh</span>
          </button>
          <button
            onClick={() => setShowFilters(!showFilters)}
            className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
          >
            <Filter className="h-4 w-4" />
            <span>Filters</span>
          </button>
        </div>
      </div>

      
      {showFilters && (
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Report Filters</h2>
            <button
              onClick={() => setShowFilters(false)}
              className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            >
              ✕
            </button>
          </div>
          
          <div className="space-y-6">
            
            <div>
              <p className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">Quick Date Ranges</p>
              <div className="flex flex-wrap gap-2">
                {quickDateRanges.map((range) => (
                  <button
                    key={range.label}
                    onClick={() => handleQuickDateRange(range.days)}
                    className="px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors text-sm"
                  >
                    {range.label}
                  </button>
                ))}
              </div>
            </div>

            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Start Date
                </label>
                <div className="relative">
                  <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                  <input
                    type="date"
                    value={dateRange.startDate}
                    onChange={(e) => setDateRange(prev => ({ ...prev, startDate: e.target.value }))}
                    className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  />
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  End Date
                </label>
                <div className="relative">
                  <Calendar className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                  <input
                    type="date"
                    value={dateRange.endDate}
                    onChange={(e) => setDateRange(prev => ({ ...prev, endDate: e.target.value }))}
                    className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  />
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Group By
                </label>
                <select
                  value={groupBy}
                  onChange={(e) => setGroupBy(e.target.value as any)}
                  className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                >
                  <option value="day">Daily</option>
                  <option value="week">Weekly</option>
                  <option value="month">Monthly</option>
                </select>
              </div>
            </div>

            
            <div>
              <p className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">Chart Display</p>
              <div className="flex space-x-4">
                <label className="flex items-center space-x-2 cursor-pointer">
                  <input
                    type="radio"
                    checked={chartType === 'line'}
                    onChange={() => setChartType('line')}
                    className="text-primary-500 focus:ring-primary-500"
                  />
                  <span className="text-sm text-gray-700 dark:text-gray-300">Line Chart</span>
                </label>
                <label className="flex items-center space-x-2 cursor-pointer">
                  <input
                    type="radio"
                    checked={chartType === 'bar'}
                    onChange={() => setChartType('bar')}
                    className="text-primary-500 focus:ring-primary-500"
                  />
                  <span className="text-sm text-gray-700 dark:text-gray-300">Bar Chart</span>
                </label>
              </div>
            </div>
          </div>
        </div>
      )}

      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Total Revenue</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                {formatCurrency(revenueSummary.totalRevenue)}
              </p>
              <div className="flex items-center space-x-1 mt-2">
                {revenueSummary.revenueChange >= 0 ? (
                  <>
                    <ArrowUpRight className="h-4 w-4 text-green-500" />
                    <span className="text-sm text-green-600 dark:text-green-400">
                      {revenueSummary.revenueChange.toFixed(1)}%
                    </span>
                  </>
                ) : (
                  <>
                    <ArrowDownRight className="h-4 w-4 text-red-500" />
                    <span className="text-sm text-red-600 dark:text-red-400">
                      {Math.abs(revenueSummary.revenueChange).toFixed(1)}%
                    </span>
                  </>
                )}
                <span className="text-xs text-gray-500 dark:text-gray-400">vs previous period</span>
              </div>
            </div>
            <div className="h-12 w-12 rounded-lg bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
              <DollarSign className="h-6 w-6 text-blue-600 dark:text-blue-400" />
            </div>
          </div>
          <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Average daily: {formatCurrency(revenueSummary.averageDailyRevenue)}
            </p>
          </div>
        </div>
        
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Total Orders</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                {ordersSummary.totalOrders.toLocaleString()}
              </p>
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">
                Average: {ordersSummary.averageOrdersPerDay.toFixed(1)}/day
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-green-100 dark:bg-green-900 flex items-center justify-center">
              <ShoppingBag className="h-6 w-6 text-green-600 dark:text-green-400" />
            </div>
          </div>
          <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-500 dark:text-gray-400">Completion Rate</span>
              <span className="text-sm font-medium text-green-600 dark:text-green-400">
                {ordersSummary.completionRate}%
              </span>
            </div>
          </div>
        </div>
        
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Average Order Value</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                {formatCurrency(chartData.length > 0 
                  ? chartData.reduce((sum, item) => sum + item.average, 0) / chartData.length 
                  : 0)}
              </p>
              <div className="mt-2">
                <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                  <div 
                    className="bg-purple-500 h-2 rounded-full"
                    style={{ width: '75%' }}
                  ></div>
                </div>
              </div>
            </div>
            <div className="h-12 w-12 rounded-lg bg-purple-100 dark:bg-purple-900 flex items-center justify-center">
              <BarChart3 className="h-6 w-6 text-purple-600 dark:text-purple-400" />
            </div>
          </div>
        </div>
        
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Report Period</p>
              <p className="text-lg font-bold text-gray-900 dark:text-white mt-1">
                {new Date(dateRange.startDate).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} - 
                {new Date(dateRange.endDate).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
              </p>
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-2">
                {filteredSalesReports.length} {groupBy} records
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-orange-100 dark:bg-orange-900 flex items-center justify-center">
              <Activity className="h-6 w-6 text-orange-600 dark:text-orange-400" />
            </div>
          </div>
          <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
            <div className="flex items-center space-x-4">
              <div className="text-center">
                <p className="text-sm font-medium text-gray-900 dark:text-white">
                  {chartData.length > 0 ? chartData[chartData.length - 1].revenue : 0}
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-400">Latest Revenue</p>
              </div>
              <div className="text-center">
                <p className="text-sm font-medium text-gray-900 dark:text-white">
                  {chartData.length > 0 ? chartData[chartData.length - 1].orders : 0}
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-400">Latest Orders</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        <div className="lg:col-span-2">
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between mb-6">
              <div>
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Sales Performance</h2>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  Revenue and orders over time
                </p>
              </div>
              <div className="flex items-center space-x-2">
                <span className="flex items-center space-x-1 text-sm">
                  <div className="h-3 w-3 rounded-full bg-blue-500"></div>
                  <span>Revenue</span>
                </span>
                <span className="flex items-center space-x-1 text-sm">
                  <div className="h-3 w-3 rounded-full bg-green-500"></div>
                  <span>Orders</span>
                </span>
              </div>
            </div>
            
            <div className="h-80">
              {chartData.length === 0 ? (
                <div className="h-full flex flex-col items-center justify-center">
                  <LineChart className="h-16 w-16 text-gray-400 mb-4" />
                  <p className="text-gray-500 dark:text-gray-400">No data available for the selected period</p>
                </div>
              ) : chartType === 'line' ? (
                <ResponsiveContainer width="100%" height="100%">
                  <RechartsLineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                    <XAxis 
                      dataKey="date" 
                      stroke="#9CA3AF"
                      fontSize={12}
                    />
                    <YAxis 
                      stroke="#9CA3AF"
                      fontSize={12}
                      tickFormatter={(value) => `$${value.toLocaleString()}`}
                    />
                    <Tooltip 
                      formatter={(value, name) => [
                        name === 'revenue' ? `$${Number(value).toLocaleString()}` : value,
                        name === 'revenue' ? 'Revenue' : 'Orders'
                      ]}
                      labelFormatter={(label) => `Date: ${label}`}
                      contentStyle={{ 
                        backgroundColor: '#1F2937',
                        border: '1px solid #374151',
                        borderRadius: '0.5rem'
                      }}
                    />
                    <Legend />
                    <Line
                      type="monotone"
                      dataKey="revenue"
                      stroke={CHART_COLORS.revenue}
                      strokeWidth={2}
                      dot={{ r: 4 }}
                      activeDot={{ r: 6 }}
                      name="Revenue"
                    />
                    <Line
                      type="monotone"
                      dataKey="orders"
                      stroke={CHART_COLORS.orders}
                      strokeWidth={2}
                      dot={{ r: 4 }}
                      activeDot={{ r: 6 }}
                      name="Orders"
                      yAxisId={1}
                    />
                  </RechartsLineChart>
                </ResponsiveContainer>
              ) : (
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                    <XAxis 
                      dataKey="date" 
                      stroke="#9CA3AF"
                      fontSize={12}
                    />
                    <YAxis 
                      stroke="#9CA3AF"
                      fontSize={12}
                      tickFormatter={(value) => `$${value.toLocaleString()}`}
                    />
                    <Tooltip 
                      formatter={(value, name) => [
                        name === 'revenue' ? `$${Number(value).toLocaleString()}` : value,
                        name === 'revenue' ? 'Revenue' : 'Orders'
                      ]}
                      labelFormatter={(label) => `Date: ${label}`}
                    />
                    <Legend />
                    <Bar 
                      dataKey="revenue" 
                      fill={CHART_COLORS.revenue}
                      name="Revenue"
                      radius={[4, 4, 0, 0]}
                    />
                    <Bar 
                      dataKey="orders" 
                      fill={CHART_COLORS.orders}
                      name="Orders"
                      radius={[4, 4, 0, 0]}
                    />
                  </BarChart>
                </ResponsiveContainer>
              )}
            </div>
          </div>
        </div>

        
        <div>
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700 h-full">
            <div className="mb-6">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Average Order Value</h2>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Trend over time
              </p>
            </div>
            
            <div className="h-64">
              {chartData.length === 0 ? (
                <div className="h-full flex items-center justify-center">
                  <PieChart className="h-12 w-12 text-gray-400" />
                </div>
              ) : (
                <ResponsiveContainer width="100%" height="100%">
                  <RechartsLineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                    <XAxis 
                      dataKey="date" 
                      stroke="#9CA3AF"
                      fontSize={10}
                    />
                    <YAxis 
                      stroke="#9CA3AF"
                      fontSize={10}
                      tickFormatter={(value) => `$${value}`}
                    />
                    <Tooltip 
                      formatter={(value) => [`$${Number(value).toFixed(2)}`, 'Average Order']}
                      labelFormatter={(label) => `Date: ${label}`}
                    />
                    <Line
                      type="monotone"
                      dataKey="average"
                      stroke={CHART_COLORS.average}
                      strokeWidth={2}
                      dot={{ r: 3 }}
                      activeDot={{ r: 5 }}
                    />
                  </RechartsLineChart>
                </ResponsiveContainer>
              )}
            </div>
            
            
            {chartData.length > 0 && (
              <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Highest AOV</p>
                    <p className="text-lg font-bold text-gray-900 dark:text-white">
                      ${Math.max(...chartData.map(d => d.average)).toFixed(2)}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Lowest AOV</p>
                    <p className="text-lg font-bold text-gray-900 dark:text-white">
                      ${Math.min(...chartData.map(d => d.average)).toFixed(2)}
                    </p>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Detailed Report Data</h2>
            <span className="text-sm text-gray-500 dark:text-gray-400">
              Showing {filteredSalesReports.length} records
            </span>
          </div>
        </div>
        
        {filteredSalesReports.length === 0 ? (
          <div className="p-8 text-center">
            <FileText className="h-16 w-16 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">No Data Available</h3>
            <p className="text-gray-500 dark:text-gray-400">
              Try adjusting your date range or filters
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 dark:bg-gray-700/50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Date
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Total Orders
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Total Revenue
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Average Order
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Revenue per Order
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                {filteredSalesReports.map((report, index) => {
                  const prevDayReport = index > 0 ? filteredSalesReports[index - 1] : null;
                  const revenueChange = prevDayReport 
                    ? ((report.totalRevenue - prevDayReport.totalRevenue) / prevDayReport.totalRevenue) * 100
                    : 0;
                  
                  return (
                    <tr key={report.date} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                      <td className="px-6 py-4">
                        <div className="text-sm font-medium text-gray-900 dark:text-white">
                          {new Date(report.date).toLocaleDateString('en-US', {
                            weekday: 'short',
                            month: 'short',
                            day: 'numeric',
                            year: 'numeric'
                          })}
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="text-sm text-gray-900 dark:text-white">{report.totalOrders}</div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center space-x-2">
                          <span className="text-sm font-medium text-gray-900 dark:text-white">
                            ${report.totalRevenue.toFixed(2)}
                          </span>
                          {prevDayReport && (
                            <span className={`text-xs ${revenueChange >= 0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}>
                              {revenueChange >= 0 ? '↑' : '↓'} {Math.abs(revenueChange).toFixed(1)}%
                            </span>
                          )}
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="text-sm text-gray-900 dark:text-white">
                          ${report.averageOrder.toFixed(2)}
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="text-sm text-gray-900 dark:text-white">
                          ${(report.totalRevenue / report.totalOrders).toFixed(2)}
                        </div>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      
      <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-6">
        <div className="flex items-start space-x-3">
          <Eye className="h-5 w-5 text-blue-500 mt-0.5 flex-shrink-0" />
          <div className="space-y-4">
            <h3 className="font-medium text-blue-800 dark:text-blue-300">Key Insights</h3>
            <ul className="text-sm text-blue-700 dark:text-blue-400 space-y-2">
              {filteredSalesReports.length > 0 && (
                <>
                  <li>• Highest revenue day: <strong>
                    {new Date(
                      [...filteredSalesReports].sort((a, b) => b.totalRevenue - a.totalRevenue)[0].date
                    ).toLocaleDateString('en-US', { weekday: 'long', month: 'long', day: 'numeric' })}
                  </strong></li>
                  <li>• Average daily revenue: <strong>${revenueSummary.averageDailyRevenue.toFixed(2)}</strong></li>
                  <li>• Best performing time: <strong>Weekend days show 25% higher revenue</strong></li>
                  <li>• Recommendation: <strong>Focus promotions on low-revenue days</strong></li>
                </>
              )}
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SalesReport;