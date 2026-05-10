import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import {
  ArrowLeft,
  User,
  Mail,
  Phone,
  Calendar,
  MapPin,
  Package,
  DollarSign,
  TrendingUp,
  Clock,
  CheckCircle,
  XCircle,
  Edit,
  MessageSquare,
  ShoppingBag,
  Award,
  BarChart3,
  MoreVertical,
  Star,
  ExternalLink
} from 'lucide-react';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { getUsers, getOrders } from '../../api/adminApi';
import { isBackendUnavailableError } from '../../api/apiHelpers';
import type { User as Customer, Order } from '../../types';

const CustomerDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<'overview' | 'orders' | 'activity'>('overview');

  const { data: customer, isLoading: customerLoading } = useQuery({
    queryKey: ['customer', id],
    queryFn: async () => {
      if (!id) throw new Error('No customer ID');
      try {
        const data = await getUsers({ page: 1, limit: 1 });
        const customers = data.data || [];
        const foundCustomer = customers.find((c: Customer) => c.id === parseInt(id));
        
        if (!foundCustomer) {
          toast.error('Customer not found');
          navigate('/customers');
          return null;
        }
        return foundCustomer;
      } catch (error) {
        toast.error('Failed to load customer details');
        navigate('/customers');
        return null;
      }
    },
    enabled: !!id,
  });

  const { data: ordersData, isLoading: ordersLoading } = useQuery({
    queryKey: ['customer-orders', id],
    queryFn: async () => {
      if (!id) return { data: [], pagination: { total: 0 } };
      try {
        const data = await getOrders({ page: 1, limit: 10, userId: parseInt(id) });
        return data;
      } catch (error) {
        if (isBackendUnavailableError(error)) {
          toast.error('Backend is unavailable. Showing empty customer orders state.');
          return { data: [], pagination: { total: 0 } };
        }

        toast.error('Failed to load customer orders');
        throw error;
      }
    },
    enabled: !!customer && !!id,
  });

  const orders = ordersData?.data || [];
  const totalOrders = ordersData?.pagination?.total || 0;

  const calculateCustomerStats = () => {
    if (!orders.length) {
      return {
        totalSpent: 0,
        avgOrderValue: 0,
        favoriteCategory: 'N/A',
        lastOrderDate: 'Never',
        orderFrequency: 0,
      };
    }

    const totalSpent = orders.reduce((sum, order) => sum + order.totalAmount, 0);
    const avgOrderValue = totalSpent / orders.length;
    
    const favoriteCategory = 'Shawarma'; // This would come from order items analysis
    
    const lastOrder = orders[0];
    const lastOrderDate = new Date(lastOrder.createdAt).toLocaleDateString();
    
    const firstOrder = orders[orders.length - 1];
    const daysBetween = Math.ceil(
      (new Date(lastOrder.createdAt).getTime() - new Date(firstOrder.createdAt).getTime()) / 
      (1000 * 60 * 60 * 24)
    );
    const orderFrequency = daysBetween > 0 ? orders.length / (daysBetween || 1) : orders.length;

    return {
      totalSpent,
      avgOrderValue,
      favoriteCategory,
      lastOrderDate,
      orderFrequency: orderFrequency.toFixed(2),
    };
  };

  const stats = calculateCustomerStats();

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const getCustomerTier = () => {
    if (stats.totalSpent > 500) return { name: 'Gold', color: 'text-yellow-600', bgColor: 'bg-yellow-100 dark:bg-yellow-900/20' };
    if (stats.totalSpent > 200) return { name: 'Silver', color: 'text-gray-600', bgColor: 'bg-gray-100 dark:bg-gray-700' };
    return { name: 'Bronze', color: 'text-amber-700', bgColor: 'bg-amber-100 dark:bg-amber-900/20' };
  };

  const tier = getCustomerTier();

  if (customerLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center space-x-4">
          <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded w-32 animate-pulse"></div>
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 space-y-4">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="h-32 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
            ))}
          </div>
          <div className="space-y-4">
            {[...Array(2)].map((_, i) => (
              <div key={i} className="h-48 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (!customer) {
    return (
      <div className="text-center py-12">
        <User className="h-16 w-16 text-gray-400 mx-auto mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Customer Not Found
        </h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6">
          The customer you're looking for doesn't exist.
        </p>
        <button
          onClick={() => navigate('/customers')}
          className="inline-flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
        >
          <ArrowLeft className="h-4 w-4" />
          <span>Back to Customers</span>
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center space-x-4">
          <button
            onClick={() => navigate('/customers')}
            className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          >
            <ArrowLeft className="h-5 w-5" />
          </button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              {customer.fullName}
            </h1>
            <p className="text-gray-600 dark:text-gray-400">
              Customer since {formatDate(customer.createdAt)}
            </p>
          </div>
        </div>
        
        <div className="flex items-center space-x-3">
          <button className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
            <MessageSquare className="h-4 w-4" />
            <span>Message</span>
          </button>
          <button className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors">
            <Edit className="h-4 w-4" />
            <span>Edit Customer</span>
          </button>
        </div>
      </div>

      
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          
          <div className="space-y-6">
            <div className="flex items-center space-x-4">
              <div className="h-20 w-20 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                <span className="text-2xl font-bold text-primary-600 dark:text-primary-400">
                  {customer.fullName.charAt(0)}
                </span>
              </div>
              <div>
                <h2 className="text-xl font-bold text-gray-900 dark:text-white">
                  {customer.fullName}
                </h2>
                <div className="flex items-center space-x-2 mt-1">
                  <span className={`px-3 py-1 rounded-full text-sm font-medium ${tier.bgColor} ${tier.color}`}>
                    <div className="flex items-center space-x-1">
                      <Award className="h-3 w-3" />
                      <span>{tier.name} Tier</span>
                    </div>
                  </span>
                  <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                    customer.isActive
                      ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                      : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                  }`}>
                    {customer.isActive ? 'Active' : 'Inactive'}
                  </span>
                </div>
              </div>
            </div>

            
            <div className="space-y-4">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                Contact Information
              </h3>
              
              <div className="space-y-3">
                <div className="flex items-center space-x-3">
                  <Phone className="h-5 w-5 text-gray-400" />
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">{customer.phone}</p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Phone Number</p>
                  </div>
                </div>
                
                {customer.email && (
                  <div className="flex items-center space-x-3">
                    <Mail className="h-5 w-5 text-gray-400" />
                    <div>
                      <p className="font-medium text-gray-900 dark:text-white">{customer.email}</p>
                      <p className="text-sm text-gray-500 dark:text-gray-400">Email Address</p>
                    </div>
                  </div>
                )}
                
                <div className="flex items-center space-x-3">
                  <Calendar className="h-5 w-5 text-gray-400" />
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">
                      {formatDate(customer.createdAt)}
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Joined Date</p>
                  </div>
                </div>
                
                {customer.lastLogin && (
                  <div className="flex items-center space-x-3">
                    <Clock className="h-5 w-5 text-gray-400" />
                    <div>
                      <p className="font-medium text-gray-900 dark:text-white">
                        {formatDate(customer.lastLogin)}
                      </p>
                      <p className="text-sm text-gray-500 dark:text-gray-400">Last Login</p>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>

          
          <div className="space-y-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
              Customer Stats
            </h3>
            
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
                <div className="flex items-center space-x-3">
                  <ShoppingBag className="h-5 w-5 text-primary-500" />
                  <div>
                    <p className="text-2xl font-bold text-gray-900 dark:text-white">
                      {totalOrders}
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Total Orders</p>
                  </div>
                </div>
              </div>
              
              <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
                <div className="flex items-center space-x-3">
                  <DollarSign className="h-5 w-5 text-green-500" />
                  <div>
                    <p className="text-2xl font-bold text-gray-900 dark:text-white">
                      ${stats.totalSpent.toFixed(2)}
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Total Spent</p>
                  </div>
                </div>
              </div>
              
              <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
                <div className="flex items-center space-x-3">
                  <BarChart3 className="h-5 w-5 text-blue-500" />
                  <div>
                    <p className="text-2xl font-bold text-gray-900 dark:text-white">
                      ${stats.avgOrderValue.toFixed(2)}
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Avg. Order Value</p>
                  </div>
                </div>
              </div>
              
              <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
                <div className="flex items-center space-x-3">
                  <TrendingUp className="h-5 w-5 text-purple-500" />
                  <div>
                    <p className="text-2xl font-bold text-gray-900 dark:text-white">
                      {stats.orderFrequency}
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Orders per Week</p>
                  </div>
                </div>
              </div>
            </div>

            
            <div className="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <Star className="h-5 w-5 text-yellow-500" />
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">
                      {stats.favoriteCategory}
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Favorite Category</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="font-medium text-gray-900 dark:text-white">
                    {stats.lastOrderDate}
                  </p>
                  <p className="text-sm text-gray-500 dark:text-gray-400">Last Order</p>
                </div>
              </div>
            </div>
          </div>

          
          <div className="space-y-6">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
              Quick Actions
            </h3>
            
            <div className="space-y-3">
              <button className="w-full flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg transition-colors">
                <div className="flex items-center space-x-3">
                  <MessageSquare className="h-5 w-5 text-gray-400" />
                  <span className="font-medium">Send Message</span>
                </div>
                <ExternalLink className="h-4 w-4 text-gray-400" />
              </button>
              
              <button className="w-full flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg transition-colors">
                <div className="flex items-center space-x-3">
                  <ShoppingBag className="h-5 w-5 text-gray-400" />
                  <span className="font-medium">Create Order</span>
                </div>
                <ExternalLink className="h-4 w-4 text-gray-400" />
              </button>
              
              <button className="w-full flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg transition-colors">
                <div className="flex items-center space-x-3">
                  <Award className="h-5 w-5 text-gray-400" />
                  <span className="font-medium">Update Tier</span>
                </div>
                <ExternalLink className="h-4 w-4 text-gray-400" />
              </button>
              
              <button
                onClick={() => navigate(`/customers/${customer.id}/edit`)}
                className="w-full flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <Edit className="h-5 w-5 text-gray-400" />
                  <span className="font-medium">Edit Profile</span>
                </div>
                <ExternalLink className="h-4 w-4 text-gray-400" />
              </button>
            </div>
          </div>
        </div>
      </div>

      
      <div className="border-b border-gray-200 dark:border-gray-700">
        <nav className="flex space-x-8">
          <button
            onClick={() => setActiveTab('overview')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'overview'
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'
            }`}
          >
            Overview
          </button>
          <button
            onClick={() => setActiveTab('orders')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'orders'
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'
            }`}
          >
            Orders ({totalOrders})
          </button>
          <button
            onClick={() => setActiveTab('activity')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'activity'
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'
            }`}
          >
            Activity
          </button>
        </nav>
      </div>

      
      {activeTab === 'overview' && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                Recent Orders
              </h3>
              <Link
                to="/orders"
                className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
              >
                View all →
              </Link>
            </div>
            
            {ordersLoading ? (
              <div className="text-center py-8">
                <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-primary-500 mx-auto"></div>
              </div>
            ) : orders.length === 0 ? (
              <div className="text-center py-8">
                <ShoppingBag className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-500 dark:text-gray-400">No orders yet</p>
              </div>
            ) : (
              <div className="space-y-4">
                {orders.slice(0, 5).map((order) => (
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
                          #{order.orderNumber}
                        </p>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                          {new Date(order.createdAt).toLocaleDateString()}
                        </p>
                      </div>
                    </div>
                    
                    <div className="text-right">
                      <p className="font-semibold text-gray-900 dark:text-white">
                        ${order.totalAmount.toFixed(2)}
                      </p>
                      <span className={`px-2 py-1 text-xs rounded-full ${
                        order.status === 'completed'
                          ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                          : order.status === 'cancelled'
                          ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                          : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
                      }`}>
                        {order.status}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                Customer Notes
              </h3>
              <button className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
                Add Note
              </button>
            </div>
            
            <div className="space-y-4">
              <div className="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
                <div className="flex items-start justify-between mb-2">
                  <span className="text-sm font-medium text-gray-900 dark:text-white">
                    First-time customer
                  </span>
                  <span className="text-xs text-gray-500 dark:text-gray-400">
                    {formatDate(customer.createdAt)}
                  </span>
                </div>
                <p className="text-sm text-gray-600 dark:text-gray-400">
                  New customer registered through the website.
                </p>
              </div>
              
              <div className="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
                <div className="flex items-start justify-between mb-2">
                  <span className="text-sm font-medium text-gray-900 dark:text-white">
                    Phone verification
                  </span>
                  <span className="text-xs text-gray-500 dark:text-gray-400">
                    {customer.lastLogin ? formatDate(customer.lastLogin) : 'Pending'}
                  </span>
                </div>
                <p className="text-sm text-gray-600 dark:text-gray-400">
                  {customer.isVerified ? 'Phone number verified successfully.' : 'Phone verification pending.'}
                </p>
              </div>
              
              <textarea
                placeholder="Add a note about this customer..."
                rows={3}
                className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
              />
            </div>
          </div>
        </div>
      )}

      {activeTab === 'orders' && (
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
          {ordersLoading ? (
            <div className="p-8 text-center">
              <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 mx-auto"></div>
              <p className="mt-4 text-gray-500 dark:text-gray-400">Loading orders...</p>
            </div>
          ) : orders.length === 0 ? (
            <div className="p-8 text-center">
              <ShoppingBag className="h-16 w-16 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                No Orders Yet
              </h3>
              <p className="text-gray-500 dark:text-gray-400">
                This customer hasn't placed any orders yet.
              </p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 dark:bg-gray-700/50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Order #
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Date
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Type
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Items
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Total
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                  {orders.map((order) => (
                    <tr key={order.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm font-medium text-gray-900 dark:text-white">
                          #{order.orderNumber}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm text-gray-900 dark:text-white">
                          {new Date(order.createdAt).toLocaleDateString()}
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400">
                          {new Date(order.createdAt).toLocaleTimeString()}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="px-2 py-1 text-xs rounded-full bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300">
                          {order.orderType}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <div className="text-sm text-gray-900 dark:text-white">
                          {order.items.length} items
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="text-sm font-semibold text-gray-900 dark:text-white">
                          ${order.totalAmount.toFixed(2)}
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                          order.status === 'completed'
                            ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                            : order.status === 'cancelled'
                            ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                            : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
                        }`}>
                          {order.status}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <Link
                          to={`/orders/${order.id}`}
                          className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
                        >
                          View Details →
                        </Link>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {activeTab === 'activity' && (
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-6">
            Recent Activity
          </h3>
          
          <div className="space-y-6">
            
            <div className="relative">
              
              <div className="absolute left-6 top-0 bottom-0 w-0.5 bg-gray-200 dark:bg-gray-700"></div>
              
              
              <div className="space-y-8">
                
                <div className="relative flex items-start">
                  <div className="absolute left-5 h-3 w-3 rounded-full bg-green-500 mt-2"></div>
                  <div className="ml-12">
                    <div className="flex items-center space-x-2 mb-1">
                      <User className="h-4 w-4 text-gray-400" />
                      <span className="font-medium text-gray-900 dark:text-white">
                        Account Created
                      </span>
                    </div>
                    <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
                      Customer registered on the website
                    </p>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      {formatDate(customer.createdAt)}
                    </p>
                  </div>
                </div>

                
                {orders.length > 0 && (
                  <div className="relative flex items-start">
                    <div className="absolute left-5 h-3 w-3 rounded-full bg-blue-500 mt-2"></div>
                    <div className="ml-12">
                      <div className="flex items-center space-x-2 mb-1">
                        <ShoppingBag className="h-4 w-4 text-gray-400" />
                        <span className="font-medium text-gray-900 dark:text-white">
                          First Order Placed
                        </span>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
                        Placed order #{orders[orders.length - 1].orderNumber} for ${orders[orders.length - 1].totalAmount.toFixed(2)}
                      </p>
                      <p className="text-xs text-gray-500 dark:text-gray-400">
                        {formatDate(orders[orders.length - 1].createdAt)}
                      </p>
                    </div>
                  </div>
                )}

                
                {customer.lastLogin && (
                  <div className="relative flex items-start">
                    <div className="absolute left-5 h-3 w-3 rounded-full bg-purple-500 mt-2"></div>
                    <div className="ml-12">
                      <div className="flex items-center space-x-2 mb-1">
                        <Clock className="h-4 w-4 text-gray-400" />
                        <span className="font-medium text-gray-900 dark:text-white">
                          Last Login
                        </span>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
                        Customer logged into their account
                      </p>
                      <p className="text-xs text-gray-500 dark:text-gray-400">
                        {formatDate(customer.lastLogin)}
                      </p>
                    </div>
                  </div>
                )}

                
                {orders.length > 0 && (
                  <div className="relative flex items-start">
                    <div className="absolute left-5 h-3 w-3 rounded-full bg-primary-500 mt-2"></div>
                    <div className="ml-12">
                      <div className="flex items-center space-x-2 mb-1">
                        <Package className="h-4 w-4 text-gray-400" />
                        <span className="font-medium text-gray-900 dark:text-white">
                          Latest Order
                        </span>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
                        Placed order #{orders[0].orderNumber} for ${orders[0].totalAmount.toFixed(2)}
                      </p>
                      <p className="text-xs text-gray-500 dark:text-gray-400">
                        {formatDate(orders[0].createdAt)}
                      </p>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default CustomerDetails;