import React from 'react';
import { Link } from 'react-router-dom';
import { Package, Clock, CheckCircle, XCircle, AlertCircle } from 'lucide-react';
import type { Order, OrderStatus } from '../../types';

interface RecentOrdersProps {
  orders?: Order[];
  loading?: boolean;
}

const RecentOrders: React.FC<RecentOrdersProps> = ({ orders = [], loading = false }) => {
  const getStatusIcon = (status: OrderStatus) => {
    switch (status) {
      case 'completed': return <CheckCircle className="h-4 w-4 text-green-500" />;
      case 'cancelled': return <XCircle className="h-4 w-4 text-red-500" />;
      case 'pending': return <Clock className="h-4 w-4 text-yellow-500" />;
      default: return <AlertCircle className="h-4 w-4 text-blue-500" />;
    }
  };

  const getStatusColor = (status: OrderStatus) => {
    switch (status) {
      case 'completed': return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
      case 'cancelled': return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
      case 'pending': return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
      case 'preparing': return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
      case 'ready': return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200';
      default: return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200';
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  if (loading) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            Recent Orders
          </h2>
          <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-24"></div>
        </div>
        {[...Array(5)].map((_, i) => (
          <div key={i} className="flex items-center justify-between py-3 border-b border-gray-200 dark:border-gray-700 last:border-0">
            <div className="flex items-center space-x-4">
              <div className="h-10 w-10 rounded-lg bg-gray-200 dark:bg-gray-700"></div>
              <div className="space-y-1">
                <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-32"></div>
                <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-24"></div>
              </div>
            </div>
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-20"></div>
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
          Recent Orders
        </h2>
        <Link
          to="/orders"
          className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300 font-medium"
        >
          View all →
        </Link>
      </div>

      <div className="space-y-4">
        {orders.length === 0 ? (
          <div className="text-center py-8">
            <Package className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500 dark:text-gray-400">No orders yet</p>
          </div>
        ) : (
          orders.slice(0, 5).map((order) => (
            <div
              key={order.id}
              className="flex items-center justify-between p-4 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors"
            >
              <div className="flex items-center space-x-4">
                <div className="h-10 w-10 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                  <Package className="h-5 w-5 text-primary-600 dark:text-primary-400" />
                </div>
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">
                    {order.customerName}
                  </p>
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    #{order.orderNumber} • {order.orderType}
                  </p>
                </div>
              </div>

              <div className="flex items-center space-x-4">
                <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(order.status)}`}>
                  <span className="flex items-center space-x-1">
                    {getStatusIcon(order.status)}
                    <span>{order.status}</span>
                  </span>
                </span>
                <div className="text-right">
                  <p className="font-semibold text-gray-900 dark:text-white">
                    ${order.totalAmount.toFixed(2)}
                  </p>
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    {formatDate(order.createdAt)}
                  </p>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default RecentOrders;