import React from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Area,
  AreaChart,
} from 'recharts';
import { TrendingUp, Calendar } from 'lucide-react';

interface SalesChartProps {
  data?: any[];
  loading?: boolean;
}

const SalesChart: React.FC<SalesChartProps> = ({ data = [], loading = false }) => {
  const mockData = [
    { date: 'Mon', sales: 1200, orders: 24 },
    { date: 'Tue', sales: 1800, orders: 32 },
    { date: 'Wed', sales: 1500, orders: 28 },
    { date: 'Thu', sales: 2200, orders: 42 },
    { date: 'Fri', sales: 3000, orders: 55 },
    { date: 'Sat', sales: 2800, orders: 48 },
    { date: 'Sun', sales: 1900, orders: 35 },
  ];

  const chartData = data.length > 0 ? data : mockData;

  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-white dark:bg-gray-800 p-4 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700">
          <p className="text-sm font-medium text-gray-900 dark:text-white mb-2">
            {label}
          </p>
          <div className="space-y-1">
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Sales: <span className="font-semibold text-primary-600">${payload[0].value}</span>
            </p>
            <p className="text-sm text-gray-600 dark:text-gray-400">
              Orders: <span className="font-semibold text-blue-600">{payload[1].value}</span>
            </p>
          </div>
        </div>
      );
    }
    return null;
  };

  if (loading) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700 animate-pulse">
        <div className="flex items-center justify-between mb-6">
          <div>
            <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded w-32 mb-2"></div>
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-48"></div>
          </div>
          <div className="h-10 w-24 bg-gray-200 dark:bg-gray-700 rounded"></div>
        </div>
        <div className="h-64 bg-gray-200 dark:bg-gray-700 rounded"></div>
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
      <div className="flex flex-col md:flex-row md:items-center justify-between mb-6">
        <div>
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            Sales Overview
          </h2>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            Weekly revenue and order trends
          </p>
        </div>
        
        <div className="flex items-center space-x-4 mt-4 md:mt-0">
          <div className="flex items-center space-x-2">
            <div className="h-3 w-3 rounded-full bg-primary-500"></div>
            <span className="text-sm text-gray-600 dark:text-gray-400">Revenue</span>
          </div>
          <div className="flex items-center space-x-2">
            <div className="h-3 w-3 rounded-full bg-blue-500"></div>
            <span className="text-sm text-gray-600 dark:text-gray-400">Orders</span>
          </div>
          <button className="flex items-center space-x-2 px-3 py-2 text-sm text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors">
            <Calendar className="h-4 w-4" />
            <span>This Week</span>
          </button>
        </div>
      </div>

      <div className="h-72">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" stroke="#374151" vertical={false} />
            <XAxis 
              dataKey="date" 
              axisLine={false}
              tickLine={false}
              tick={{ fill: '#6B7280' }}
            />
            <YAxis 
              axisLine={false}
              tickLine={false}
              tick={{ fill: '#6B7280' }}
              tickFormatter={(value) => `$${value}`}
            />
            <Tooltip content={<CustomTooltip />} />
            <Area
              type="monotone"
              dataKey="sales"
              stroke="#E40A2D"
              fill="#E40A2D"
              fillOpacity={0.1}
              strokeWidth={2}
            />
            <Line
              type="monotone"
              dataKey="orders"
              stroke="#3B82F6"
              strokeWidth={2}
              dot={{ r: 4 }}
              activeDot={{ r: 6 }}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>

      <div className="flex items-center justify-between mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
        <div>
          <p className="text-sm text-gray-500 dark:text-gray-400">Total This Week</p>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">
            ${chartData.reduce((sum, day) => sum + day.sales, 0).toLocaleString()}
          </p>
        </div>
        <div className="flex items-center text-green-600">
          <TrendingUp className="h-5 w-5 mr-2" />
          <span className="font-semibold">+12.5%</span>
          <span className="text-gray-500 dark:text-gray-400 ml-2">from last week</span>
        </div>
      </div>
    </div>
  );
};

export default SalesChart;