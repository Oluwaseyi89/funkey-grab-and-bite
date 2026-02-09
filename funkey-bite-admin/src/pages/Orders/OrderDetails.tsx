import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  ArrowLeft, 
  Package, 
  Clock, 
  CheckCircle, 
  XCircle, 
  AlertCircle,
  Printer,
  Mail,
  Phone,
  MapPin,
  Calendar,
  DollarSign,
  User,
  ChevronDown
} from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getOrder, updateOrderStatus } from '../../api/adminApi';
import type { Order, OrderStatus } from '../../types';

const OrderDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [isUpdating, setIsUpdating] = useState(false);

  const { data: order, isLoading, refetch } = useQuery({
    queryKey: ['order', id],
    queryFn: async () => {
      try {
        if (!id) throw new Error('No order ID');
        const data = await getOrder(parseInt(id));
        return data;
      } catch (error) {
        toast.error('Failed to load order details');
        navigate('/orders');
        throw error;
      }
    },
    enabled: !!id,
  });

  const handleStatusUpdate = async (newStatus: OrderStatus) => {
    if (!order || !id) return;
    
    setIsUpdating(true);
    try {
      await updateOrderStatus(parseInt(id), { status: newStatus });
      toast.success(`Order status updated to ${newStatus}`);
      refetch();
    } catch (error) {
      toast.error('Failed to update order status');
    } finally {
      setIsUpdating(false);
    }
  };

  const getStatusInfo = (status: OrderStatus) => {
    switch (status) {
      case 'pending':
        return { 
          color: 'text-yellow-600 bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200',
          icon: <Clock className="h-5 w-5" />,
          text: 'Pending - Awaiting confirmation'
        };
      case 'confirmed':
        return { 
          color: 'text-blue-600 bg-blue-50 dark:bg-blue-900/20 border-blue-200',
          icon: <AlertCircle className="h-5 w-5" />,
          text: 'Confirmed - Order accepted'
        };
      case 'preparing':
        return { 
          color: 'text-purple-600 bg-purple-50 dark:bg-purple-900/20 border-purple-200',
          icon: <Package className="h-5 w-5" />,
          text: 'Preparing - In kitchen'
        };
      case 'ready':
        return { 
          color: 'text-green-600 bg-green-50 dark:bg-green-900/20 border-green-200',
          icon: <CheckCircle className="h-5 w-5" />,
          text: 'Ready for pickup/delivery'
        };
      case 'completed':
        return { 
          color: 'text-gray-600 bg-gray-50 dark:bg-gray-900/20 border-gray-200',
          icon: <CheckCircle className="h-5 w-5" />,
          text: 'Completed - Order fulfilled'
        };
      case 'cancelled':
        return { 
          color: 'text-red-600 bg-red-50 dark:bg-red-900/20 border-red-200',
          icon: <XCircle className="h-5 w-5" />,
          text: 'Cancelled - Order cancelled'
        };
      default:
        return { 
          color: 'text-gray-600 bg-gray-50 dark:bg-gray-900/20 border-gray-200',
          icon: <AlertCircle className="h-5 w-5" />,
          text: 'Unknown status'
        };
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (isLoading) {
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

  if (!order) {
    return (
      <div className="text-center py-12">
        <Package className="h-16 w-16 text-gray-400 mx-auto mb-4" />
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
          Order Not Found
        </h2>
        <p className="text-gray-500 dark:text-gray-400 mb-6">
          The order you're looking for doesn't exist.
        </p>
        <button
          onClick={() => navigate('/orders')}
          className="inline-flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
        >
          <ArrowLeft className="h-4 w-4" />
          <span>Back to Orders</span>
        </button>
      </div>
    );
  }

  const statusInfo = getStatusInfo(order.status);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center space-x-4">
          <button
            onClick={() => navigate('/orders')}
            className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          >
            <ArrowLeft className="h-5 w-5" />
          </button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
              Order #{order.orderNumber}
            </h1>
            <p className="text-gray-600 dark:text-gray-400">
              Placed on {formatDate(order.createdAt)}
            </p>
          </div>
        </div>
        
        <div className="flex items-center space-x-3">
          <button className="flex items-center space-x-2 px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors">
            <Printer className="h-4 w-4" />
            <span>Print</span>
          </button>
          <button className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors">
            <Mail className="h-4 w-4" />
            <span>Contact Customer</span>
          </button>
        </div>
      </div>

      {/* Status Banner */}
      <div className={`p-4 rounded-xl border ${statusInfo.color}`}>
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            {statusInfo.icon}
            <div>
              <p className="font-semibold">Current Status: {order.status.toUpperCase()}</p>
              <p className="text-sm opacity-80">{statusInfo.text}</p>
            </div>
          </div>
          
          <div className="relative group">
            <button 
              disabled={isUpdating}
              className="flex items-center space-x-2 px-4 py-2 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 rounded-lg transition-colors disabled:opacity-50"
            >
              <span>Update Status</span>
              <ChevronDown className="h-4 w-4" />
            </button>
            <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 py-1 z-10 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none group-hover:pointer-events-auto">
              <button
                onClick={() => handleStatusUpdate('confirmed')}
                className="w-full text-left px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                Mark as Confirmed
              </button>
              <button
                onClick={() => handleStatusUpdate('preparing')}
                className="w-full text-left px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                Mark as Preparing
              </button>
              <button
                onClick={() => handleStatusUpdate('ready')}
                className="w-full text-left px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                Mark as Ready
              </button>
              <button
                onClick={() => handleStatusUpdate('completed')}
                className="w-full text-left px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700"
              >
                Mark as Completed
              </button>
              <button
                onClick={() => handleStatusUpdate('cancelled')}
                className="w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20"
              >
                Cancel Order
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column - Order Details */}
        <div className="lg:col-span-2 space-y-6">
          {/* Order Items */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              Order Items
            </h2>
            <div className="space-y-4">
              {order.items.map((item, index) => (
                <div key={index} className="flex items-center justify-between p-4 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                  <div className="flex items-center space-x-4">
                    <div className="h-12 w-12 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                      <Package className="h-5 w-5 text-primary-600 dark:text-primary-400" />
                    </div>
                    <div>
                      <p className="font-medium text-gray-900 dark:text-white">
                        {item.name}
                      </p>
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        Quantity: {item.quantity}
                      </p>
                      {item.specialInstructions && (
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                          Note: {item.specialInstructions}
                        </p>
                      )}
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="font-semibold text-gray-900 dark:text-white">
                      ${(item.unitPrice * item.quantity).toFixed(2)}
                    </p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      ${item.unitPrice} each
                    </p>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
              <div className="flex justify-between items-center">
                <span className="text-gray-600 dark:text-gray-400">Subtotal</span>
                <span className="font-medium">${order.totalAmount.toFixed(2)}</span>
              </div>
              <div className="flex justify-between items-center mt-2">
                <span className="text-gray-600 dark:text-gray-400">Tax</span>
                <span className="font-medium">$0.00</span>
              </div>
              <div className="flex justify-between items-center mt-2">
                <span className="text-gray-600 dark:text-gray-400">Delivery Fee</span>
                <span className="font-medium">$0.00</span>
              </div>
              <div className="flex justify-between items-center mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
                <span className="text-lg font-semibold">Total</span>
                <span className="text-2xl font-bold text-primary-600">
                  ${order.totalAmount.toFixed(2)}
                </span>
              </div>
            </div>
          </div>

          {/* Order Notes */}
          {order.notes && (
            <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                Order Notes
              </h2>
              <div className="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
                <p className="text-gray-700 dark:text-gray-300">{order.notes}</p>
              </div>
            </div>
          )}
        </div>

        {/* Right Column - Customer & Order Info */}
        <div className="space-y-6">
          {/* Customer Info */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              Customer Information
            </h2>
            <div className="space-y-4">
              <div className="flex items-center space-x-3">
                <User className="h-5 w-5 text-gray-400" />
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">{order.customerName}</p>
                  <p className="text-sm text-gray-500 dark:text-gray-400">Customer</p>
                </div>
              </div>
              <div className="flex items-center space-x-3">
                <Phone className="h-5 w-5 text-gray-400" />
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">{order.customerPhone}</p>
                  <p className="text-sm text-gray-500 dark:text-gray-400">Phone Number</p>
                </div>
              </div>
              {order.customerEmail && (
                <div className="flex items-center space-x-3">
                  <Mail className="h-5 w-5 text-gray-400" />
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">{order.customerEmail}</p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Email Address</p>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Order Info */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              Order Information
            </h2>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-gray-600 dark:text-gray-400">Order Type</span>
                <span className="font-medium capitalize">{order.orderType}</span>
              </div>
              {order.pickupTime && (
                <div className="flex items-center justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Pickup Time</span>
                  <span className="font-medium">
                    {new Date(order.pickupTime).toLocaleTimeString()}
                  </span>
                </div>
              )}
              {order.estimatedReadyTime && (
                <div className="flex items-center justify-between">
                  <span className="text-gray-600 dark:text-gray-400">Estimated Ready</span>
                  <span className="font-medium">
                    {new Date(order.estimatedReadyTime).toLocaleTimeString()}
                  </span>
                </div>
              )}
              <div className="flex items-center justify-between">
                <span className="text-gray-600 dark:text-gray-400">Payment Method</span>
                <span className="font-medium">Cash on Delivery</span>
              </div>
            </div>
          </div>

          {/* Order Timeline */}
          <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
            <h2 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              Order Timeline
            </h2>
            <div className="space-y-4">
              <div className="flex items-start space-x-3">
                <div className="h-6 w-6 rounded-full bg-green-500 flex items-center justify-center mt-1">
                  <CheckCircle className="h-4 w-4 text-white" />
                </div>
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">Order Placed</p>
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    {formatDate(order.createdAt)}
                  </p>
                </div>
              </div>
              {order.status !== 'pending' && (
                <div className="flex items-start space-x-3">
                  <div className="h-6 w-6 rounded-full bg-blue-500 flex items-center justify-center mt-1">
                    <AlertCircle className="h-4 w-4 text-white" />
                  </div>
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">Order Confirmed</p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">
                      {formatDate(new Date().toISOString())}
                    </p>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default OrderDetails;