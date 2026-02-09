// src/pages/Reports/Analytics.tsx
import React, { useState } from 'react';
import {
  TrendingUp,
  TrendingDown,
  Users,
  ShoppingBag,
  DollarSign,
  Clock,
  PieChart,
  BarChart3,
  Activity,
  Target,
  Calendar,
  Filter,
  Download,
  RefreshCw,
  Eye,
  ArrowUpRight,
  ArrowDownRight,
  ChevronRight,
  Sparkles,
  Award,
  Zap,
  Star,
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon
} from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getDashboardStats, getTodayStats } from '../../api/adminApi';
import { useReportsStore } from '../../stores/reportsStore';

// Charting library
import {
  ResponsiveContainer,
  LineChart,
  Line,
  BarChart,
  Bar,
  PieChart as RechartsPieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  RadarChart,
  Radar,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  AreaChart,
  Area
} from 'recharts';

// Color palette
const CHART_COLORS = {
  primary: '#3b82f6',
  secondary: '#10b981',
  tertiary: '#8b5cf6',
  quaternary: '#f59e0b',
  danger: '#ef4444',
  warning: '#f59e0b',
  success: '#10b981',
  gray: '#6b7280'
};

const Analytics: React.FC = () => {
  const [timeRange, setTimeRange] = useState<'today' | 'week' | 'month' | 'quarter' | 'year'>('month');
  const [activeTab, setActiveTab] = useState<'overview' | 'performance' | 'customers' | 'products'>('overview');
  const [showFilters, setShowFilters] = useState(false);
  
  const { getRevenueSummary, getOrdersSummary, getCustomersSummary, setLoading, setError } = useReportsStore();

  // Fetch dashboard stats
  const { data: dashboardStats, isLoading: isLoadingStats, refetch: refetchStats } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: async () => {
      setLoading(true);
      try {
        const data = await getDashboardStats();
        return data;
      } catch (error: any) {
        setError(error.message || 'Failed to load dashboard statistics');
        toast.error('Failed to load analytics data');
        return null;
      } finally {
        setLoading(false);
      }
    },
  });

  // Fetch today's stats
  const { data: todayStats, refetch: refetchTodayStats } = useQuery({
    queryKey: ['today-stats'],
    queryFn: async () => {
      try {
        const data = await getTodayStats();
        return data;
      } catch (error) {
        return null;
      }
    },
  });

  // Mock data for charts (replace with actual API data when available)
  const revenueTrendData = [
    { date: 'Jan', revenue: 45000, orders: 320, avgOrder: 140.62 },
    { date: 'Feb', revenue: 52000, orders: 380, avgOrder: 136.84 },
    { date: 'Mar', revenue: 48000, orders: 350, avgOrder: 137.14 },
    { date: 'Apr', revenue: 61000, orders: 420, avgOrder: 145.24 },
    { date: 'May', revenue: 58000, orders: 410, avgOrder: 141.46 },
    { date: 'Jun', revenue: 65000, orders: 450, avgOrder: 144.44 },
    { date: 'Jul', revenue: 72000, orders: 490, avgOrder: 146.94 },
    { date: 'Aug', revenue: 68000, orders: 470, avgOrder: 144.68 },
    { date: 'Sep', revenue: 75000, orders: 510, avgOrder: 147.06 },
    { date: 'Oct', revenue: 82000, orders: 550, avgOrder: 149.09 },
    { date: 'Nov', revenue: 89000, orders: 590, avgOrder: 150.85 },
    { date: 'Dec', revenue: 95000, orders: 620, avgOrder: 153.23 },
  ];

  const categoryDistribution = [
    { name: 'Chips & Chicken', value: 35, color: CHART_COLORS.primary },
    { name: 'Shawarma', value: 25, color: CHART_COLORS.secondary },
    { name: 'Noodles', value: 20, color: CHART_COLORS.tertiary },
    { name: 'Drinks', value: 12, color: CHART_COLORS.quaternary },
    { name: 'Soup & Bowls', value: 8, color: CHART_COLORS.warning },
  ];

  const peakHoursData = [
    { hour: '8 AM', orders: 45, revenue: 1800 },
    { hour: '10 AM', orders: 85, revenue: 3400 },
    { hour: '12 PM', orders: 150, revenue: 7500 },
    { hour: '2 PM', orders: 120, revenue: 6000 },
    { hour: '4 PM', orders: 95, revenue: 4750 },
    { hour: '6 PM', orders: 180, revenue: 10800 },
    { hour: '8 PM', orders: 140, revenue: 8400 },
    { hour: '10 PM', orders: 60, revenue: 3000 },
  ];

  const customerMetrics = [
    { metric: 'New Customers', value: 245, change: +12.5, icon: Users },
    { metric: 'Returning Rate', value: '68%', change: +5.2, icon: TrendingUp },
    { metric: 'Avg Order Value', value: '$147.50', change: +8.3, icon: DollarSign },
    { metric: 'Order Frequency', value: '3.2x/month', change: +2.1, icon: Clock },
  ];

  const performanceMetrics = dashboardStats ? [
    { label: 'Total Orders', value: dashboardStats.totalOrders, change: 15.2, icon: ShoppingBag },
    { label: 'Total Revenue', value: `$${dashboardStats.totalRevenue.toLocaleString()}`, change: 22.8, icon: DollarSign },
    { label: 'Pending Orders', value: dashboardStats.pendingOrders, change: -8.5, icon: Clock },
    { label: 'Active Catering', value: dashboardStats.activeCatering, change: 32.1, icon: Users },
    { label: 'New Customers', value: dashboardStats.newCustomers, change: 18.7, icon: Users },
    { label: 'Avg Order Value', value: `$${(dashboardStats.totalRevenue / dashboardStats.totalOrders || 0).toFixed(2)}`, change: 5.4, icon: BarChart3 },
  ] : [];

  const popularItems = dashboardStats?.popularItems || [];

  // Calculate metrics
  const calculateGrowth = (current: number, previous: number) => {
    if (previous === 0) return 0;
    return ((current - previous) / previous) * 100;
  };

  // Format numbers
  const formatNumber = (num: number) => {
    if (num >= 1000000) return `$${(num / 1000000).toFixed(1)}M`;
    if (num >= 1000) return `$${(num / 1000).toFixed(1)}K`;
    return `$${num}`;
  };

  const handleRefresh = () => {
    refetchStats();
    refetchTodayStats();
    toast.success('Analytics data refreshed');
  };

  const exportAnalytics = () => {
    toast.success('Analytics report exported (PDF generation would be implemented here)');
  };

  if (isLoadingStats) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-64 animate-pulse"></div>
          <div className="flex space-x-2">
            <div className="h-10 w-32 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
            <div className="h-10 w-32 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
          </div>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="h-32 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
          ))}
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="h-96 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
          <div className="h-96 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Business Analytics</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Comprehensive insights and performance metrics
          </p>
        </div>
        
        <div className="flex items-center space-x-3">
          <button
            onClick={exportAnalytics}
            className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            <Download className="h-4 w-4" />
            <span>Export Report</span>
          </button>
          <button
            onClick={handleRefresh}
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
            <span>Time Range</span>
          </button>
        </div>
      </div>

      {/* Time Range Filters */}
      {showFilters && (
        <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
          <div className="flex flex-wrap gap-2">
            {(['today', 'week', 'month', 'quarter', 'year'] as const).map((range) => (
              <button
                key={range}
                onClick={() => {
                  setTimeRange(range);
                  setShowFilters(false);
                }}
                className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                  timeRange === range
                    ? 'bg-primary-500 text-white'
                    : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                }`}
              >
                {range.charAt(0).toUpperCase() + range.slice(1)}
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Tab Navigation */}
      <div className="bg-white dark:bg-gray-800 rounded-xl p-2 border border-gray-200 dark:border-gray-700">
        <div className="flex space-x-2">
          {([
            { id: 'overview', label: 'Overview', icon: Activity },
            { id: 'performance', label: 'Performance', icon: TrendingUpIcon },
            { id: 'customers', label: 'Customers', icon: Users },
            { id: 'products', label: 'Products', icon: PieChart },
          ] as const).map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`flex-1 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                activeTab === tab.id
                  ? 'bg-primary-500 text-white'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'
              }`}
            >
              <div className="flex items-center justify-center space-x-2">
                <tab.icon className="h-4 w-4" />
                <span>{tab.label}</span>
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* Performance Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {performanceMetrics.map((metric, index) => (
          <div key={index} className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">{metric.label}</p>
                <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                  {metric.value}
                </p>
                <div className="flex items-center space-x-1 mt-2">
                  {metric.change >= 0 ? (
                    <>
                      <ArrowUpRight className="h-4 w-4 text-green-500" />
                      <span className="text-sm text-green-600 dark:text-green-400">
                        {metric.change > 0 ? '+' : ''}{metric.change.toFixed(1)}%
                      </span>
                    </>
                  ) : (
                    <>
                      <ArrowDownRight className="h-4 w-4 text-red-500" />
                      <span className="text-sm text-red-600 dark:text-red-400">
                        {metric.change.toFixed(1)}%
                      </span>
                    </>
                  )}
                  <span className="text-xs text-gray-500 dark:text-gray-400">vs last period</span>
                </div>
              </div>
              <div className={`h-12 w-12 rounded-lg flex items-center justify-center ${
                metric.change >= 0 
                  ? 'bg-green-100 dark:bg-green-900'
                  : 'bg-red-100 dark:bg-red-900'
              }`}>
                <metric.icon className={`h-6 w-6 ${
                  metric.change >= 0 
                    ? 'text-green-600 dark:text-green-400'
                    : 'text-red-600 dark:text-red-400'
                }`} />
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Charts Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Revenue Trend Chart */}
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Revenue Trend</h2>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Monthly revenue and order growth
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
          
          <div className="h-72">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={revenueTrendData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                <XAxis dataKey="date" stroke="#9CA3AF" fontSize={12} />
                <YAxis stroke="#9CA3AF" fontSize={12} tickFormatter={(value) => `$${value/1000}K`} />
                <Tooltip 
                  formatter={(value, name) => [
                    name === 'revenue' ? `$${Number(value).toLocaleString()}` : value,
                    name === 'revenue' ? 'Revenue' : 'Orders'
                  ]}
                  contentStyle={{ 
                    backgroundColor: '#1F2937',
                    border: '1px solid #374151',
                    borderRadius: '0.5rem'
                  }}
                />
                <Area 
                  type="monotone" 
                  dataKey="revenue" 
                  stroke={CHART_COLORS.primary}
                  fill={CHART_COLORS.primary}
                  fillOpacity={0.1}
                  strokeWidth={2}
                />
                <Area 
                  type="monotone" 
                  dataKey="orders" 
                  stroke={CHART_COLORS.secondary}
                  fill={CHART_COLORS.secondary}
                  fillOpacity={0.1}
                  strokeWidth={2}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Category Distribution */}
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Category Performance</h2>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Revenue distribution by category
              </p>
            </div>
            <PieChart className="h-5 w-5 text-gray-400" />
          </div>
          
          <div className="h-72">
            <ResponsiveContainer width="100%" height="100%">
              <RechartsPieChart>
                <Pie
                  data={categoryDistribution}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="value"
                >
                  {categoryDistribution.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip 
                  formatter={(value) => [`${value}%`, 'Market Share']}
                  contentStyle={{ 
                    backgroundColor: '#1F2937',
                    border: '1px solid #374151',
                    borderRadius: '0.5rem'
                  }}
                />
                <Legend />
              </RechartsPieChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* Peak Hours & Popular Items */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Peak Hours */}
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Peak Order Hours</h2>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Order volume throughout the day
              </p>
            </div>
            <Clock className="h-5 w-5 text-gray-400" />
          </div>
          
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={peakHoursData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#374151" />
                <XAxis dataKey="hour" stroke="#9CA3AF" fontSize={12} />
                <YAxis stroke="#9CA3AF" fontSize={12} />
                <Tooltip 
                  formatter={(value, name) => [
                    name === 'revenue' ? `$${Number(value).toLocaleString()}` : value,
                    name === 'revenue' ? 'Revenue' : 'Orders'
                  ]}
                  contentStyle={{ 
                    backgroundColor: '#1F2937',
                    border: '1px solid #374151',
                    borderRadius: '0.5rem'
                  }}
                />
                <Bar 
                  dataKey="orders" 
                  fill={CHART_COLORS.primary}
                  radius={[4, 4, 0, 0]}
                  name="Orders"
                />
                <Bar 
                  dataKey="revenue" 
                  fill={CHART_COLORS.secondary}
                  radius={[4, 4, 0, 0]}
                  name="Revenue"
                />
              </BarChart>
            </ResponsiveContainer>
          </div>
          
          <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-900 dark:text-white">Peak Hour: 6 PM</p>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  180 orders • $10,800 revenue
                </p>
              </div>
              <div className="text-right">
                <p className="text-sm font-medium text-green-600 dark:text-green-400">+22%</p>
                <p className="text-xs text-gray-500 dark:text-gray-400">vs avg hour</p>
              </div>
            </div>
          </div>
        </div>

        {/* Popular Items */}
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Top Performing Items</h2>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Best-selling menu items
              </p>
            </div>
            <Award className="h-5 w-5 text-gray-400" />
          </div>
          
          <div className="space-y-4">
            {popularItems.length > 0 ? (
              popularItems.slice(0, 5).map((item, index) => (
                <div key={item.id} className="flex items-center justify-between p-3 hover:bg-gray-50 dark:hover:bg-gray-700/50 rounded-lg">
                  <div className="flex items-center space-x-3">
                    <div className={`h-8 w-8 rounded-lg flex items-center justify-center ${
                      index === 0 ? 'bg-yellow-100 dark:bg-yellow-900' :
                      index === 1 ? 'bg-gray-100 dark:bg-gray-700' :
                      index === 2 ? 'bg-orange-100 dark:bg-orange-900' :
                      'bg-blue-100 dark:bg-blue-900'
                    }`}>
                      <span className={`text-sm font-bold ${
                        index === 0 ? 'text-yellow-600 dark:text-yellow-400' :
                        index === 1 ? 'text-gray-600 dark:text-gray-400' :
                        index === 2 ? 'text-orange-600 dark:text-orange-400' :
                        'text-blue-600 dark:text-blue-400'
                      }`}>
                        #{index + 1}
                      </span>
                    </div>
                    <div>
                      <p className="font-medium text-gray-900 dark:text-white">{item.name}</p>
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        {item.totalSold} sold • ${item.revenue.toFixed(2)}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="font-medium text-gray-900 dark:text-white">
                      ${item.revenue.toFixed(2)}
                    </p>
                    <p className="text-sm text-green-600 dark:text-green-400 flex items-center justify-end">
                      <TrendingUp className="h-3 w-3 mr-1" />
                      {item.totalSold > 100 ? 'High' : 'Medium'} demand
                    </p>
                  </div>
                </div>
              ))
            ) : (
              <div className="text-center py-8">
                <Sparkles className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500 dark:text-gray-400">
                  No popular items data available
                </p>
              </div>
            )}
          </div>
          
          <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
            <button className="w-full flex items-center justify-center space-x-2 text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
              <span>View All Items</span>
              <ChevronRight className="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>

      {/* Customer Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {customerMetrics.map((metric, index) => (
          <div key={index} className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-500 dark:text-gray-400">{metric.metric}</p>
                <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                  {metric.value}
                </p>
                <div className="flex items-center space-x-1 mt-2">
                  {metric.change >= 0 ? (
                    <>
                      <ArrowUpRight className="h-4 w-4 text-green-500" />
                      <span className="text-sm text-green-600 dark:text-green-400">
                        +{metric.change}%
                      </span>
                    </>
                  ) : (
                    <>
                      <ArrowDownRight className="h-4 w-4 text-red-500" />
                      <span className="text-sm text-red-600 dark:text-red-400">
                        {metric.change}%
                      </span>
                    </>
                  )}
                  <span className="text-xs text-gray-500 dark:text-gray-400">from last month</span>
                </div>
              </div>
              <div className="h-12 w-12 rounded-lg bg-purple-100 dark:bg-purple-900 flex items-center justify-center">
                <metric.icon className="h-6 w-6 text-purple-600 dark:text-purple-400" />
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Performance Radar Chart */}
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Business Health Score</h2>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              Overall performance across key metrics
            </p>
          </div>
          <Target className="h-5 w-5 text-gray-400" />
        </div>
        
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <RadarChart cx="50%" cy="50%" outerRadius="80%" data={[
              { metric: 'Revenue Growth', value: 85, fullMark: 100 },
              { metric: 'Order Volume', value: 78, fullMark: 100 },
              { metric: 'Customer Satisfaction', value: 92, fullMark: 100 },
              { metric: 'Operational Efficiency', value: 75, fullMark: 100 },
              { metric: 'Profit Margins', value: 88, fullMark: 100 },
              { metric: 'Growth Potential', value: 82, fullMark: 100 },
            ]}>
              <PolarGrid stroke="#374151" />
              <PolarAngleAxis dataKey="metric" stroke="#9CA3AF" fontSize={12} />
              <PolarRadiusAxis angle={30} domain={[0, 100]} stroke="#9CA3AF" />
              <Radar
                name="Performance"
                dataKey="value"
                stroke={CHART_COLORS.primary}
                fill={CHART_COLORS.primary}
                fillOpacity={0.3}
                strokeWidth={2}
              />
              <Tooltip 
                formatter={(value) => [`${value}/100`, 'Score']}
                contentStyle={{ 
                  backgroundColor: '#1F2937',
                  border: '1px solid #374151',
                  borderRadius: '0.5rem'
                }}
              />
              <Legend />
            </RadarChart>
          </ResponsiveContainer>
        </div>
        
        <div className="mt-6 grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="text-center">
            <div className="text-3xl font-bold text-green-600 dark:text-green-400">8.7</div>
            <div className="text-sm text-gray-500 dark:text-gray-400">Overall Score</div>
          </div>
          <div className="text-center">
            <div className="text-3xl font-bold text-blue-600 dark:text-blue-400">92%</div>
            <div className="text-sm text-gray-500 dark:text-gray-400">Customer Satisfaction</div>
          </div>
          <div className="text-center">
            <div className="text-3xl font-bold text-purple-600 dark:text-purple-400">A-</div>
            <div className="text-sm text-gray-500 dark:text-gray-400">Business Grade</div>
          </div>
        </div>
      </div>

      {/* Insights & Recommendations */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Key Insights */}
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-6">
          <div className="flex items-start space-x-3">
            <Eye className="h-5 w-5 text-blue-500 mt-0.5 flex-shrink-0" />
            <div className="space-y-4">
              <h3 className="font-medium text-blue-800 dark:text-blue-300">Key Insights</h3>
              <ul className="text-sm text-blue-700 dark:text-blue-400 space-y-2">
                <li className="flex items-start space-x-2">
                  <TrendingUpIcon className="h-4 w-4 text-green-500 mt-0.5 flex-shrink-0" />
                  <span><strong>Revenue growing at 22.8%</strong> month-over-month</span>
                </li>
                <li className="flex items-start space-x-2">
                  <Zap className="h-4 w-4 text-yellow-500 mt-0.5 flex-shrink-0" />
                  <span><strong>6 PM is peak hour</strong> with 180 orders ($10,800 revenue)</span>
                </li>
                <li className="flex items-start space-x-2">
                  <Star className="h-4 w-4 text-purple-500 mt-0.5 flex-shrink-0" />
                  <span><strong>Chips & Chicken category</strong> leads with 35% market share</span>
                </li>
                <li className="flex items-start space-x-2">
                  <Users className="h-4 w-4 text-green-500 mt-0.5 flex-shrink-0" />
                  <span><strong>Customer retention rate</strong> improved by 5.2% this month</span>
                </li>
              </ul>
            </div>
          </div>
        </div>

        {/* Recommendations */}
        <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl p-6">
          <div className="flex items-start space-x-3">
            <Target className="h-5 w-5 text-green-500 mt-0.5 flex-shrink-0" />
            <div className="space-y-4">
              <h3 className="font-medium text-green-800 dark:text-green-300">Recommendations</h3>
              <ul className="text-sm text-green-700 dark:text-green-400 space-y-2">
                <li className="flex items-start space-x-2">
                  <span className="font-bold">1.</span>
                  <span><strong>Increase staffing during 4-8 PM</strong> to handle peak order volume</span>
                </li>
                <li className="flex items-start space-x-2">
                  <span className="font-bold">2.</span>
                  <span><strong>Promote slow-moving categories</strong> during off-peak hours</span>
                </li>
                <li className="flex items-start space-x-2">
                  <span className="font-bold">3.</span>
                  <span><strong>Launch loyalty program</strong> to boost returning customer rate</span>
                </li>
                <li className="flex items-start space-x-2">
                  <span className="font-bold">4.</span>
                  <span><strong>Optimize delivery routes</strong> based on order density patterns</span>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Analytics;