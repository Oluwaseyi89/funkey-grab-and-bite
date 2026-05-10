import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import {
  Search,
  Filter,
  Package,
  AlertTriangle,
  TrendingUp,
  TrendingDown,
  Plus,
  Edit,
  Eye,
  BarChart3,
  Download,
  RefreshCw,
  ChevronLeft,
  ChevronRight,
  MoreVertical,
  CheckCircle,
  XCircle,
  Package2,
  ChevronDown,
  Calendar,
  Box,
  Scale
} from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getInventory, getLowStock, getInventoryAlerts } from '../../api/adminApi';
import { isBackendUnavailableError } from '../../api/apiHelpers';
import type { InventoryItem, InventoryAlert } from '../../types';
import { useInventoryStore } from '../../stores/inventoryStore';

const InventoryList: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'all' | 'low-stock' | 'alerts'>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  
  const { setItems, setAlerts } = useInventoryStore();

  const { data: inventoryData, isLoading, refetch } = useQuery({
    queryKey: ['inventory', page, limit, statusFilter, searchQuery],
    queryFn: async () => {
      try {
        const data = await getInventory();
        if (Array.isArray(data)) {
          setItems(data);
          return { data, pagination: { total: data.length, totalPages: Math.ceil(data.length / limit) } };
        }
        return { data: [], pagination: { total: 0, totalPages: 0 } };
      } catch (error) {
        if (isBackendUnavailableError(error)) {
          toast.error('Backend is unavailable. Showing empty inventory state.');
          return { data: [], pagination: { total: 0, totalPages: 0 } };
        }

        toast.error('Failed to load inventory');
        throw error;
      }
    },
  });

  const { data: lowStockData } = useQuery({
    queryKey: ['low-stock'],
    queryFn: async () => {
      try {
        const data = await getLowStock();
        return Array.isArray(data) ? data : [];
      } catch (error) {
        if (isBackendUnavailableError(error)) {
          return [];
        }

        throw error;
      }
    },
  });

  const { data: alertsData } = useQuery({
    queryKey: ['inventory-alerts'],
    queryFn: async () => {
      try {
        const data = await getInventoryAlerts();
        if (Array.isArray(data)) {
          setAlerts(data);
          return data;
        }
        return [];
      } catch (error) {
        if (isBackendUnavailableError(error)) {
          return [];
        }

        throw error;
      }
    },
  });

  const inventoryItems = inventoryData?.data || [];
  const pagination = inventoryData?.pagination || { total: 0, totalPages: 0 };
  const lowStockItems = lowStockData || [];
  const alerts = alertsData || [];
  const unreadAlerts = alerts.filter(a => !a.isResolved && !a.readAt).length;

  const getFilteredItems = () => {
    let items = [...inventoryItems];
    
    if (activeTab === 'low-stock') {
      items = lowStockItems;
    }
    
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase().trim();
      items = items.filter(item =>
        item.name.toLowerCase().includes(query) ||
        item.unit.toLowerCase().includes(query)
      );
    }
    
    if (statusFilter === 'active') {
      items = items.filter(item => item.isActive);
    } else if (statusFilter === 'inactive') {
      items = items.filter(item => !item.isActive);
    }
    
    return items;
  };

  const filteredItems = getFilteredItems();
  const pagedItems = filteredItems.slice((page - 1) * limit, page * limit);

  const handleRestock = (item: InventoryItem) => {
    toast.success(`Restock order placed for ${item.name}`);
  };

  const handleMarkAlertAsResolved = (alertId: number) => {
    toast.success('Alert marked as resolved');
  };

  const getStockStatus = (item: InventoryItem) => {
    const percentage = (item.currentStock / item.minimumStock) * 100;
    
    if (item.currentStock <= item.reorderPoint) {
      return {
        color: 'text-red-600',
        bgColor: 'bg-red-100 dark:bg-red-900/20',
        text: 'Critical',
        icon: <AlertTriangle className="h-4 w-4" />
      };
    } else if (percentage <= 150) {
      return {
        color: 'text-yellow-600',
        bgColor: 'bg-yellow-100 dark:bg-yellow-900/20',
        text: 'Low',
        icon: <TrendingDown className="h-4 w-4" />
      };
    } else {
      return {
        color: 'text-green-600',
        bgColor: 'bg-green-100 dark:bg-green-900/20',
        text: 'Good',
        icon: <TrendingUp className="h-4 w-4" />
      };
    }
  };

  const getAlertStatus = (alert: InventoryAlert) => {
    if (alert.isResolved) {
      return {
        color: 'text-gray-600',
        bgColor: 'bg-gray-100 dark:bg-gray-700',
        text: 'Resolved',
        icon: <CheckCircle className="h-4 w-4" />
      };
    }
    
    if (!alert.readAt) {
      return {
        color: 'text-red-600',
        bgColor: 'bg-red-100 dark:bg-red-900/20',
        text: 'Unread',
        icon: <AlertTriangle className="h-4 w-4" />
      };
    }
    
    return {
      color: 'text-yellow-600',
      bgColor: 'bg-yellow-100 dark:bg-yellow-900/20',
      text: 'Read',
      icon: <AlertTriangle className="h-4 w-4" />
    };
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const calculateTotalValue = () => {
    return inventoryItems.reduce((sum, item) => sum + (item.currentStock * 10), 0);
  };

  const calculateRestockCost = () => {
    return lowStockItems.reduce((sum, item) => {
      const needed = item.minimumStock - item.currentStock;
      return sum + (Math.max(0, needed) * 10); // Estimate $10 per unit
    }, 0);
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-48 animate-pulse"></div>
          <div className="h-10 w-32 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
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
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Inventory Management</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Track and manage your restaurant inventory
          </p>
        </div>
        
        <div className="flex items-center space-x-3">
          <button
            onClick={() => refetch()}
            className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            <RefreshCw className="h-4 w-4" />
            <span>Refresh</span>
          </button>
          <button className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
            <Download className="h-4 w-4" />
            <span>Export</span>
          </button>
          <Link
            to="/inventory/new"
            className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
          >
            <Plus className="h-4 w-4" />
            <span>Add Item</span>
          </Link>
        </div>
      </div>

      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Total Items</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                {inventoryItems.length}
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
              <Package className="h-6 w-6 text-blue-600 dark:text-blue-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Low Stock Items</p>
              <p className="text-2xl font-bold text-red-600 dark:text-red-400 mt-1">
                {lowStockItems.length}
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-red-100 dark:bg-red-900 flex items-center justify-center">
              <AlertTriangle className="h-6 w-6 text-red-600 dark:text-red-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Inventory Value</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                ${calculateTotalValue().toLocaleString()}
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-green-100 dark:bg-green-900 flex items-center justify-center">
              <BarChart3 className="h-6 w-6 text-green-600 dark:text-green-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Active Alerts</p>
              <p className="text-2xl font-bold text-yellow-600 dark:text-yellow-400 mt-1">
                {unreadAlerts}
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-yellow-100 dark:bg-yellow-900 flex items-center justify-center">
              <AlertTriangle className="h-6 w-6 text-yellow-600 dark:text-yellow-400" />
            </div>
          </div>
        </div>
      </div>

      
      <div className="border-b border-gray-200 dark:border-gray-700">
        <nav className="flex space-x-8">
          <button
            onClick={() => setActiveTab('all')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'all'
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'
            }`}
          >
            All Items ({inventoryItems.length})
          </button>
          <button
            onClick={() => setActiveTab('low-stock')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'low-stock'
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'
            }`}
          >
            Low Stock ({lowStockItems.length})
          </button>
          <button
            onClick={() => setActiveTab('alerts')}
            className={`py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'alerts'
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'
            }`}
          >
            Alerts ({alerts.length})
            {unreadAlerts > 0 && (
              <span className="ml-2 px-2 py-0.5 text-xs bg-red-500 text-white rounded-full">
                {unreadAlerts}
              </span>
            )}
          </button>
        </nav>
      </div>

      
      <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
        <div className="flex flex-col lg:flex-row gap-4">
          
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
              <input
                type="search"
                placeholder="Search inventory items..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
              />
            </div>
          </div>

          
          <div className="flex items-center space-x-4">
            <Filter className="h-5 w-5 text-gray-400" />
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="px-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
            >
              <option value="">All Status</option>
              <option value="active">Active Only</option>
              <option value="inactive">Inactive Only</option>
            </select>
          </div>
        </div>
      </div>

      
      {activeTab === 'alerts' ? (
        
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
          {alerts.length === 0 ? (
            <div className="p-8 text-center">
              <CheckCircle className="h-16 w-16 text-green-400 mx-auto mb-4" />
              <p className="text-gray-500 dark:text-gray-400">No active alerts</p>
            </div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-gray-50 dark:bg-gray-700/50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Alert
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Item
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Type
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Created
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                    {alerts.map((alert) => {
                      const statusInfo = getAlertStatus(alert);
                      const relatedItem = inventoryItems.find(item => item.id === alert.inventoryItemId);
                      
                      return (
                        <tr key={alert.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                          <td className="px-6 py-4">
                            <div>
                              <p className="font-medium text-gray-900 dark:text-white">
                                {alert.message}
                              </p>
                            </div>
                          </td>
                          <td className="px-6 py-4">
                            {relatedItem ? (
                              <div className="flex items-center space-x-2">
                                <Package2 className="h-4 w-4 text-gray-400" />
                                <span className="font-medium">{relatedItem.name}</span>
                              </div>
                            ) : (
                              <span className="text-gray-400">Item not found</span>
                            )}
                          </td>
                          <td className="px-6 py-4">
                            <span className="px-2 py-1 text-xs rounded-full bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300 capitalize">
                              {alert.alertType}
                            </span>
                          </td>
                          <td className="px-6 py-4">
                            <span className={`px-3 py-1 rounded-full text-xs font-medium inline-flex items-center space-x-1 ${statusInfo.bgColor} ${statusInfo.color}`}>
                              {statusInfo.icon}
                              <span>{statusInfo.text}</span>
                            </span>
                          </td>
                          <td className="px-6 py-4">
                            <div className="text-sm text-gray-500 dark:text-gray-400">
                              {formatDate(alert.createdAt)}
                            </div>
                          </td>
                          <td className="px-6 py-4">
                            <div className="flex items-center space-x-2">
                              {!alert.isResolved && (
                                <button
                                  onClick={() => handleMarkAlertAsResolved(alert.id)}
                                  className="px-3 py-1 text-sm bg-green-100 hover:bg-green-200 dark:bg-green-900/20 dark:hover:bg-green-900/40 text-green-700 dark:text-green-400 rounded-lg transition-colors"
                                >
                                  Resolve
                                </button>
                              )}
                              {relatedItem && (
                                <Link
                                  to={`/inventory/${relatedItem.id}`}
                                  className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
                                >
                                  View Item
                                </Link>
                              )}
                            </div>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            </>
          )}
        </div>
      ) : (
        
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
          {filteredItems.length === 0 ? (
            <div className="p-8 text-center">
              <Package className="h-16 w-16 text-gray-400 mx-auto mb-4" />
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
                No Items Found
              </h2>
              <p className="text-gray-500 dark:text-gray-400 mb-6">
                {searchQuery || statusFilter
                  ? 'Try changing your filters'
                  : 'Get started by adding your first inventory item.'}
              </p>
              <Link
                to="/inventory/new"
                className="inline-flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
              >
                <Plus className="h-4 w-4" />
                <span>Add Item</span>
              </Link>
            </div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-gray-50 dark:bg-gray-700/50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Item
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Current Stock
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Min Stock
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Status
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Last Restocked
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                    {pagedItems.map((item) => {
                      const statusInfo = getStockStatus(item);
                      const needsRestock = item.currentStock <= item.reorderPoint;
                      
                      return (
                        <tr key={item.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                          <td className="px-6 py-4">
                            <div className="flex items-center space-x-3">
                              <div className="h-10 w-10 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                                <Package2 className="h-5 w-5 text-primary-600 dark:text-primary-400" />
                              </div>
                              <div>
                                <div className="font-medium text-gray-900 dark:text-white">
                                  {item.name}
                                </div>
                                <div className="text-sm text-gray-500 dark:text-gray-400">
                                  {item.unit} • ID: {item.menuItemId}
                                </div>
                              </div>
                            </div>
                          </td>
                          <td className="px-6 py-4">
                            <div className="flex items-baseline space-x-2">
                              <span className={`text-lg font-bold ${needsRestock ? 'text-red-600' : 'text-gray-900 dark:text-white'}`}>
                                {item.currentStock}
                              </span>
                              <span className="text-sm text-gray-500 dark:text-gray-400">
                                / {item.minimumStock}
                              </span>
                            </div>
                            {needsRestock && (
                              <div className="text-xs text-red-600 mt-1">
                                Reorder at: {item.reorderPoint}
                              </div>
                            )}
                          </td>
                          <td className="px-6 py-4">
                            <div className="text-gray-900 dark:text-white">
                              {item.minimumStock} {item.unit}
                            </div>
                          </td>
                          <td className="px-6 py-4">
                            <div className="flex items-center space-x-2">
                              <span className={`px-3 py-1 rounded-full text-xs font-medium inline-flex items-center space-x-1 ${statusInfo.bgColor} ${statusInfo.color}`}>
                                {statusInfo.icon}
                                <span>{statusInfo.text}</span>
                              </span>
                              <span className={`px-2 py-1 text-xs rounded-full ${
                                item.isActive
                                  ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                                  : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                              }`}>
                                {item.isActive ? 'Active' : 'Inactive'}
                              </span>
                            </div>
                          </td>
                          <td className="px-6 py-4">
                            <div className="text-sm text-gray-500 dark:text-gray-400">
                              {item.lastRestocked ? formatDate(item.lastRestocked) : 'Never'}
                            </div>
                          </td>
                          <td className="px-6 py-4">
                            <div className="flex items-center space-x-2">
                              <Link
                                to={`/inventory/${item.id}`}
                                className="p-1 text-gray-400 hover:text-primary-500 transition-colors"
                                title="View Details"
                              >
                                <Eye className="h-4 w-4" />
                              </Link>
                              <Link
                                to={`/inventory/${item.id}/edit`}
                                className="p-1 text-gray-400 hover:text-primary-500 transition-colors"
                                title="Edit"
                              >
                                <Edit className="h-4 w-4" />
                              </Link>
                              {needsRestock && (
                                <button
                                  onClick={() => handleRestock(item)}
                                  className="px-3 py-1 text-sm bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
                                >
                                  Restock
                                </button>
                              )}
                              <div className="relative group">
                                <button className="p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
                                  <MoreVertical className="h-4 w-4" />
                                </button>
                                <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 py-1 z-10 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none group-hover:pointer-events-auto">
                                  <button
                                    onClick={() => handleRestock(item)}
                                    className="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                                  >
                                    Quick Restock
                                  </button>
                                  <Link
                                    to={`/inventory/${item.id}/history`}
                                    className="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 block"
                                  >
                                    View History
                                  </Link>
                                  <button className="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700">
                                    Adjust Stock
                                  </button>
                                  <div className="border-t border-gray-200 dark:border-gray-700 my-1"></div>
                                  <button className="w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20">
                                    Deactivate
                                  </button>
                                </div>
                              </div>
                            </div>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>

              
              {pagination.totalPages > 1 && (
                <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700">
                  <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                    <div className="text-sm text-gray-500 dark:text-gray-400">
                      Showing <span className="font-medium">{((page - 1) * limit) + 1}</span> to{' '}
                      <span className="font-medium">{Math.min(page * limit, filteredItems.length)}</span> of{' '}
                      <span className="font-medium">{filteredItems.length}</span> items
                    </div>
                    <div className="flex items-center space-x-2">
                      <button
                        onClick={() => setPage(p => Math.max(1, p - 1))}
                        disabled={page === 1}
                        className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <ChevronLeft className="h-5 w-5" />
                      </button>
                      
                      <div className="flex items-center space-x-1">
                        {Array.from({ length: Math.min(5, pagination.totalPages) }, (_, i) => {
                          let pageNum = i + 1;
                          if (pagination.totalPages > 5) {
                            if (page > 3) pageNum = page - 2 + i;
                            if (page > pagination.totalPages - 2) pageNum = pagination.totalPages - 4 + i;
                          }
                          return (
                            <button
                              key={pageNum}
                              onClick={() => setPage(pageNum)}
                              className={`px-3 py-1 rounded-lg ${
                                page === pageNum
                                  ? 'bg-primary-500 text-white'
                                  : 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300'
                              }`}
                            >
                              {pageNum}
                            </button>
                          );
                        })}
                      </div>
                      
                      <button
                        onClick={() => setPage(p => Math.min(pagination.totalPages, p + 1))}
                        disabled={page === pagination.totalPages}
                        className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <ChevronRight className="h-5 w-5" />
                      </button>
                    </div>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      )}

      
      {lowStockItems.length > 0 && activeTab !== 'alerts' && (
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Quick Restock Summary
              </h2>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {lowStockItems.length} items need restocking • Estimated cost: ${calculateRestockCost()}
              </p>
            </div>
            <button className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors">
              <RefreshCw className="h-4 w-4" />
              <span>Order All</span>
            </button>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {lowStockItems.slice(0, 6).map((item) => (
              <div
                key={item.id}
                className="p-4 border border-red-200 dark:border-red-800 rounded-lg bg-red-50 dark:bg-red-900/20"
              >
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center space-x-2">
                    <AlertTriangle className="h-4 w-4 text-red-500" />
                    <span className="font-medium text-gray-900 dark:text-white">{item.name}</span>
                  </div>
                  <span className="text-sm font-bold text-red-600">
                    {item.currentStock}/{item.minimumStock}
                  </span>
                </div>
                <div className="text-sm text-gray-600 dark:text-gray-400 mb-3">
                  Needs {Math.max(0, item.minimumStock - item.currentStock)} {item.unit}
                </div>
                <button
                  onClick={() => handleRestock(item)}
                  className="w-full px-3 py-2 text-sm bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 border border-red-300 dark:border-red-700 text-red-600 dark:text-red-400 rounded-lg transition-colors"
                >
                  Order Now
                </button>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default InventoryList;