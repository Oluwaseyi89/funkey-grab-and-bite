import React, { useEffect, useState } from 'react';
import { TrendingUp, Users, Package, AlertCircle, Clock, DollarSign } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { useSocket } from '../contexts/SocketContext';
import toast from 'react-hot-toast';
import StatCard from '../components/dashboard/StatCard';
import RecentOrders from '../components/dashboard/RecentOrders';
import SalesChart from '../components/dashboard/SalesChart';
import { getDashboardStats, getTodayStats, getOrders } from '../api/adminApi';
import type { AdminStats, Order } from '../types';

const Dashboard: React.FC = () => {
  const { socket } = useSocket();
  const [realtimeOrders, setRealtimeOrders] = useState<Order[]>([]);

  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ['dashboard-stats'],
    queryFn: async () => {
      try {
        const data = await getTodayStats();
        return data;
      } catch (error) {
        toast.error('Failed to load dashboard statistics');
        throw error;
      }
    },
    refetchInterval: 300000, // Refetch every 5 minutes
  });

  const { data: ordersData, refetch: refetchOrders } = useQuery({
    queryKey: ['recent-orders'],
    queryFn: async () => {
      try {
        const data = await getOrders({ page: 1, limit: 10 });
        return data.data || [];
      } catch (error) {
        toast.error('Failed to load recent orders');
        return [];
      }
    },
  });

  useEffect(() => {
    if (!socket) return;

    const handleNewOrder = (order: Order) => {
      toast.success(`New order received: #${order.orderNumber}`);
      setRealtimeOrders((prev) => [order, ...prev.slice(0, 4)]);
      refetchOrders();
    };

    const handleOrderUpdated = (order: Partial<Order> & { id?: number; status?: string }) => {
      const orderLabel = order.orderNumber || order.id || 'Unknown';
      const statusLabel = order.status || 'updated';
      toast(`Order #${orderLabel} updated: ${statusLabel}`);
      refetchOrders();
    };

    socket.on('new_order', handleNewOrder);
    socket.on('order_updated', handleOrderUpdated);

    return () => {
      socket.off('new_order', handleNewOrder);
      socket.off('order_updated', handleOrderUpdated);
    };
  }, [socket, refetchOrders]);

  const statCards = [
    {
      title: "Today's Revenue",
      value: `$${stats?.totalRevenue?.toFixed(2) || '0.00'}`,
      icon: DollarSign,
      change: 12.5,
      color: 'primary' as const,
    },
    {
      title: "Today's Orders",
      value: stats?.totalOrders || 0,
      icon: Package,
      change: 8.2,
      color: 'blue' as const,
    },
    {
      title: 'Pending Orders',
      value: stats?.pendingOrders || 0,
      icon: Clock,
      change: 0,
      color: 'yellow' as const,
    },
    {
      title: 'New Customers',
      value: stats?.newCustomers || 0,
      icon: Users,
      change: 15.3,
      color: 'green' as const,
    },
  ];

  const allOrders = [...realtimeOrders, ...(ordersData || [])].slice(0, 10);

  return (
    <div className="space-y-6">
      
      <div>
        <h1 className="text-2xl md:text-3xl font-bold text-gray-900 dark:text-white">
          Dashboard
        </h1>
        <p className="text-gray-600 dark:text-gray-400 mt-2">
          Welcome back! Here's what's happening with your business today.
        </p>
      </div>

      
      {realtimeOrders.length > 0 && (
        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
          <div className="flex items-center space-x-2 mb-2">
            <TrendingUp className="h-5 w-5 text-blue-500" />
            <span className="font-medium text-blue-800 dark:text-blue-300">
              Real-time Updates
            </span>
          </div>
          <div className="text-sm text-blue-700 dark:text-blue-400">
            {realtimeOrders.length} new order(s) received in real-time
          </div>
        </div>
      )}

      
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
        {statCards.map((card, index) => (
          <StatCard
            key={index}
            {...card}
            loading={statsLoading}
          />
        ))}
      </div>

      
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        
        <div className="lg:col-span-2">
          <SalesChart loading={statsLoading} />
        </div>

        
        <div>
          <RecentOrders 
            orders={allOrders} 
            loading={statsLoading} 
          />
        </div>
      </div>

      
      {stats?.popularItems && stats.popularItems.length > 0 && (
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
            Popular Items
          </h2>
          <div className="space-y-4">
            {stats.popularItems.map((item) => (
              <div
                key={item.id}
                className="flex items-center justify-between p-4 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors"
              >
                <div className="flex items-center space-x-4">
                  <div className="h-12 w-12 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                    <span className="text-primary-600 dark:text-primary-400 font-bold">
                      {item.name.charAt(0)}
                    </span>
                  </div>
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">
                      {item.name}
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {item.totalSold} sold
                    </p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="font-semibold text-gray-900 dark:text-white">
                    ${item.revenue.toFixed(2)}
                  </p>
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    Revenue
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <button className="bg-white dark:bg-gray-800 p-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:border-primary-500 dark:hover:border-primary-500 transition-colors text-left group">
          <div className="flex items-center space-x-3">
            <div className="h-10 w-10 rounded-lg bg-green-100 dark:bg-green-900 flex items-center justify-center group-hover:bg-green-200 dark:group-hover:bg-green-800 transition-colors">
              <Package className="h-5 w-5 text-green-600 dark:text-green-400" />
            </div>
            <div>
              <p className="font-medium text-gray-900 dark:text-white">Process Orders</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {stats?.pendingOrders || 0} pending orders
              </p>
            </div>
          </div>
        </button>

        <button className="bg-white dark:bg-gray-800 p-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:border-primary-500 dark:hover:border-primary-500 transition-colors text-left group">
          <div className="flex items-center space-x-3">
            <div className="h-10 w-10 rounded-lg bg-blue-100 dark:bg-blue-900 flex items-center justify-center group-hover:bg-blue-200 dark:group-hover:bg-blue-800 transition-colors">
              <AlertCircle className="h-5 w-5 text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <p className="font-medium text-gray-900 dark:text-white">Check Inventory</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                View low stock items
              </p>
            </div>
          </div>
        </button>

        <button className="bg-white dark:bg-gray-800 p-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:border-primary-500 dark:hover:border-primary-500 transition-colors text-left group">
          <div className="flex items-center space-x-3">
            <div className="h-10 w-10 rounded-lg bg-purple-100 dark:bg-purple-900 flex items-center justify-center group-hover:bg-purple-200 dark:group-hover:bg-purple-800 transition-colors">
              <TrendingUp className="h-5 w-5 text-purple-600 dark:text-purple-400" />
            </div>
            <div>
              <p className="font-medium text-gray-900 dark:text-white">View Reports</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Sales analytics & insights
              </p>
            </div>
          </div>
        </button>
      </div>
    </div>
  );
};

export default Dashboard;


































