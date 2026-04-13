import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import {
  AlertTriangle,
  Filter,
  Search,
  CheckCircle,
  XCircle,
  Clock,
  Bell,
  BellOff,
  RefreshCw,
  Download,
  ExternalLink,
  Package,
  TrendingUp,
  TrendingDown,
  MoreVertical,
  Eye,
  User,
  Calendar,
  AlertCircle,
  CheckCheck,
  Archive,
  Send,
  Mail,
  ChevronLeft,
  ChevronRight,
  BarChart3
} from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getInventoryAlerts, updateCateringStatus } from '../../api/adminApi';
import type { InventoryAlert } from '../../types';
import { useInventoryStore } from '../../stores/inventoryStore';

const InventoryAlerts: React.FC = () => {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [typeFilter, setTypeFilter] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedAlerts, setSelectedAlerts] = useState<number[]>([]);
  const [showResolved, setShowResolved] = useState(false);

  const { markAlertAsRead, markAllAlertsAsRead, updateAlert } = useInventoryStore();

  const { data: alertsData, isLoading, refetch } = useQuery({
    queryKey: ['inventory-alerts-admin', page, limit, statusFilter, typeFilter, searchQuery],
    queryFn: async () => {
      try {
        const data = await getInventoryAlerts({ 
          resolved: showResolved ? undefined : false 
        });
        if (Array.isArray(data)) {
          return data;
        }
        return [];
      } catch (error) {
        toast.error('Failed to load inventory alerts');
        return [];
      }
    },
  });

  const alerts = alertsData || [];
  const filteredAlerts = alerts.filter(alert => {
    if (statusFilter === 'unread' && alert.readAt) return false;
    if (statusFilter === 'read' && !alert.readAt) return false;
    if (statusFilter === 'resolved' && !alert.isResolved) return false;
    if (statusFilter === 'unresolved' && alert.isResolved) return false;
    
    if (typeFilter !== 'all' && alert.alertType !== typeFilter) return false;
    
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase().trim();
      return alert.message.toLowerCase().includes(query);
    }
    
    return true;
  });

  const pagedAlerts = filteredAlerts.slice((page - 1) * limit, page * limit);
  const totalPages = Math.ceil(filteredAlerts.length / limit);

  const unreadAlerts = alerts.filter(a => !a.readAt).length;
  const unresolvedAlerts = alerts.filter(a => !a.isResolved).length;
  const criticalAlerts = alerts.filter(a => a.alertType === 'critical').length;
  const todayAlerts = alerts.filter(a => {
    const alertDate = new Date(a.createdAt);
    const today = new Date();
    return alertDate.toDateString() === today.toDateString();
  }).length;

  const handleSelectAll = () => {
    if (selectedAlerts.length === pagedAlerts.length) {
      setSelectedAlerts([]);
    } else {
      setSelectedAlerts(pagedAlerts.map(alert => alert.id));
    }
  };

  const handleSelectAlert = (alertId: number) => {
    setSelectedAlerts(prev => 
      prev.includes(alertId) 
        ? prev.filter(id => id !== alertId)
        : [...prev, alertId]
    );
  };

  const handleMarkAsRead = (alertId: number) => {
    markAlertAsRead(alertId);
    toast.success('Alert marked as read');
  };

  const handleMarkAsResolved = (alertId: number) => {
    updateAlert(alertId, { 
      isResolved: true, 
      resolvedAt: new Date().toISOString() 
    });
    toast.success('Alert marked as resolved');
  };

  const handleBulkResolve = () => {
    if (selectedAlerts.length === 0) {
      toast.error('No alerts selected');
      return;
    }
    
    selectedAlerts.forEach(alertId => {
      updateAlert(alertId, { 
        isResolved: true, 
        resolvedAt: new Date().toISOString() 
      });
    });
    
    toast.success(`${selectedAlerts.length} alerts resolved`);
    setSelectedAlerts([]);
  };

  const handleBulkMarkAsRead = () => {
    if (selectedAlerts.length === 0) {
      toast.error('No alerts selected');
      return;
    }
    
    selectedAlerts.forEach(alertId => {
      markAlertAsRead(alertId);
    });
    
    toast.success(`${selectedAlerts.length} alerts marked as read`);
    setSelectedAlerts([]);
  };

  const handleNotifyStaff = () => {
    if (selectedAlerts.length === 0) {
      toast.error('No alerts selected');
      return;
    }
    
    toast.success(`Notifications sent for ${selectedAlerts.length} alerts`);
  };

  const handleExportAlerts = () => {
    toast.success('Export started');
  };

  const getAlertTypeInfo = (type: string) => {
    switch (type) {
      case 'critical':
        return {
          color: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
          icon: <AlertTriangle className="h-4 w-4" />,
          text: 'Critical',
          priority: 1
        };
      case 'low_stock':
        return {
          color: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
          icon: <TrendingDown className="h-4 w-4" />,
          text: 'Low Stock',
          priority: 2
        };
      case 'expiring':
        return {
          color: 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200',
          icon: <Clock className="h-4 w-4" />,
          text: 'Expiring Soon',
          priority: 3
        };
      case 'restocked':
        return {
          color: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
          icon: <TrendingUp className="h-4 w-4" />,
          text: 'Restocked',
          priority: 4
        };
      default:
        return {
          color: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200',
          icon: <Bell className="h-4 w-4" />,
          text: type,
          priority: 5
        };
    }
  };

  const getAlertStatusInfo = (alert: InventoryAlert) => {
    if (alert.isResolved) {
      return {
        color: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
        icon: <CheckCircle className="h-4 w-4" />,
        text: 'Resolved',
        badge: 'bg-green-500'
      };
    }
    
    if (!alert.readAt) {
      return {
        color: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
        icon: <Bell className="h-4 w-4" />,
        text: 'Unread',
        badge: 'bg-red-500'
      };
    }
    
    return {
      color: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
      icon: <Eye className="h-4 w-4" />,
      text: 'Read',
      badge: 'bg-blue-500'
    };
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    
    if (diffHours < 1) {
      return 'Just now';
    } else if (diffHours < 24) {
      return `${diffHours}h ago`;
    } else if (diffHours < 168) { // 7 days
      return `${Math.floor(diffHours / 24)}d ago`;
    }
    
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-48 animate-pulse"></div>
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
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Inventory Alerts Management
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Monitor and manage inventory alerts and notifications
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
          <button
            onClick={handleExportAlerts}
            className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            <Download className="h-4 w-4" />
            <span>Export</span>
          </button>
        </div>
      </div>

      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Active Alerts</p>
              <p className="text-2xl font-bold text-red-600 dark:text-red-400 mt-1">
                {unresolvedAlerts}
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                {todayAlerts} new today
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
              <p className="text-sm text-gray-500 dark:text-gray-400">Unread Alerts</p>
              <p className="text-2xl font-bold text-yellow-600 dark:text-yellow-400 mt-1">
                {unreadAlerts}
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Requires attention
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-yellow-100 dark:bg-yellow-900 flex items-center justify-center">
              <Bell className="h-6 w-6 text-yellow-600 dark:text-yellow-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Critical Alerts</p>
              <p className="text-2xl font-bold text-purple-600 dark:text-purple-400 mt-1">
                {criticalAlerts}
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Immediate action needed
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-purple-100 dark:bg-purple-900 flex items-center justify-center">
              <AlertCircle className="h-6 w-6 text-purple-600 dark:text-purple-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Total Alerts</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                {alerts.length}
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                All time
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
              <BarChart3 className="h-6 w-6 text-primary-600 dark:text-primary-400" />
            </div>
          </div>
        </div>
      </div>

      
      {selectedAlerts.length > 0 && (
        <div className="bg-primary-50 dark:bg-primary-900/20 border border-primary-200 dark:border-primary-800 rounded-xl p-4">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div className="flex items-center space-x-3">
              <div className="h-8 w-8 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                <span className="text-primary-600 dark:text-primary-400 font-semibold">
                  {selectedAlerts.length}
                </span>
              </div>
              <div>
                <p className="font-medium text-primary-800 dark:text-primary-300">
                  {selectedAlerts.length} alert{selectedAlerts.length !== 1 ? 's' : ''} selected
                </p>
                <p className="text-sm text-primary-600 dark:text-primary-400">
                  Choose an action below
                </p>
              </div>
            </div>
            
            <div className="flex flex-wrap gap-2">
              <button
                onClick={handleBulkMarkAsRead}
                className="px-4 py-2 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg transition-colors flex items-center space-x-2"
              >
                <CheckCheck className="h-4 w-4" />
                <span>Mark as Read</span>
              </button>
              <button
                onClick={handleBulkResolve}
                className="px-4 py-2 bg-green-500 hover:bg-green-600 text-white rounded-lg transition-colors flex items-center space-x-2"
              >
                <CheckCircle className="h-4 w-4" />
                <span>Resolve Alerts</span>
              </button>
              <button
                onClick={handleNotifyStaff}
                className="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg transition-colors flex items-center space-x-2"
              >
                <Send className="h-4 w-4" />
                <span>Notify Staff</span>
              </button>
              <button
                onClick={() => setSelectedAlerts([])}
                className="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors flex items-center space-x-2"
              >
                <XCircle className="h-4 w-4" />
                <span>Clear Selection</span>
              </button>
            </div>
          </div>
        </div>
      )}

      
      <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
        <div className="flex flex-col lg:flex-row gap-4">
          
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
              <input
                type="search"
                placeholder="Search alert messages..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
              />
            </div>
          </div>

          
          <div className="flex flex-wrap items-center gap-4">
            <div className="flex items-center space-x-2">
              <Filter className="h-5 w-5 text-gray-400" />
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="px-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
              >
                <option value="all">All Status</option>
                <option value="unread">Unread Only</option>
                <option value="read">Read Only</option>
                <option value="unresolved">Unresolved</option>
                <option value="resolved">Resolved</option>
              </select>
            </div>

            <div className="flex items-center space-x-2">
              <select
                value={typeFilter}
                onChange={(e) => setTypeFilter(e.target.value)}
                className="px-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
              >
                <option value="all">All Types</option>
                <option value="critical">Critical</option>
                <option value="low_stock">Low Stock</option>
                <option value="expiring">Expiring Soon</option>
                <option value="restocked">Restocked</option>
              </select>
            </div>

            <button
              onClick={() => setShowResolved(!showResolved)}
              className={`flex items-center space-x-2 px-4 py-2 rounded-lg transition-colors ${
                showResolved
                  ? 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300'
                  : 'bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 text-gray-600 dark:text-gray-400'
              }`}
            >
              {showResolved ? (
                <BellOff className="h-4 w-4" />
              ) : (
                <Bell className="h-4 w-4" />
              )}
              <span>{showResolved ? 'Hide Resolved' : 'Show Resolved'}</span>
            </button>
          </div>
        </div>
      </div>

      
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        {filteredAlerts.length === 0 ? (
          <div className="p-8 text-center">
            <CheckCircle className="h-16 w-16 text-green-400 mx-auto mb-4" />
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
              No Alerts Found
            </h2>
            <p className="text-gray-500 dark:text-gray-400">
              {searchQuery || statusFilter !== 'all' || typeFilter !== 'all'
                ? 'Try changing your filters'
                : 'All inventory alerts are currently resolved.'}
            </p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 dark:bg-gray-700/50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider w-12">
                      <input
                        type="checkbox"
                        checked={selectedAlerts.length === pagedAlerts.length && pagedAlerts.length > 0}
                        onChange={handleSelectAll}
                        className="rounded border-gray-300 dark:border-gray-600 text-primary-500 focus:ring-primary-500"
                      />
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Alert Details
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Type
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Time
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                  {pagedAlerts.map((alert) => {
                    const typeInfo = getAlertTypeInfo(alert.alertType);
                    const statusInfo = getAlertStatusInfo(alert);
                    const isSelected = selectedAlerts.includes(alert.id);
                    
                    return (
                      <tr 
                        key={alert.id} 
                        className={`hover:bg-gray-50 dark:hover:bg-gray-700/30 ${isSelected ? 'bg-primary-50 dark:bg-primary-900/10' : ''}`}
                      >
                        <td className="px-6 py-4">
                          <input
                            type="checkbox"
                            checked={isSelected}
                            onChange={() => handleSelectAlert(alert.id)}
                            className="rounded border-gray-300 dark:border-gray-600 text-primary-500 focus:ring-primary-500"
                          />
                        </td>
                        <td className="px-6 py-4">
                          <div className="space-y-1">
                            <div className="flex items-start space-x-3">
                              {!alert.readAt && (
                                <div className="h-2 w-2 rounded-full bg-red-500 mt-2"></div>
                              )}
                              <div className="flex-1">
                                <p className="font-medium text-gray-900 dark:text-white">
                                  {alert.message}
                                </p>
                                <div className="flex items-center space-x-2 text-sm text-gray-500 dark:text-gray-400 mt-1">
                                  <Package className="h-3 w-3" />
                                  <span>Item ID: {alert.inventoryItemId}</span>
                                </div>
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <span className={`px-3 py-1 rounded-full text-xs font-medium inline-flex items-center space-x-1 ${typeInfo.color}`}>
                            {typeInfo.icon}
                            <span>{typeInfo.text}</span>
                          </span>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex items-center space-x-2">
                            <span className={`px-3 py-1 rounded-full text-xs font-medium inline-flex items-center space-x-1 ${statusInfo.color}`}>
                              {statusInfo.icon}
                              <span>{statusInfo.text}</span>
                            </span>
                            {alert.isResolved && alert.resolvedAt && (
                              <div className="text-xs text-gray-500">
                                {formatDate(alert.resolvedAt)}
                              </div>
                            )}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="space-y-1">
                            <div className="text-sm text-gray-900 dark:text-white">
                              {formatDate(alert.createdAt)}
                            </div>
                            <div className="text-xs text-gray-500 dark:text-gray-400">
                              Created
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex items-center space-x-2">
                            {!alert.readAt && (
                              <button
                                onClick={() => handleMarkAsRead(alert.id)}
                                className="p-1 text-gray-400 hover:text-blue-500 transition-colors"
                                title="Mark as Read"
                              >
                                <Eye className="h-4 w-4" />
                              </button>
                            )}
                            
                            {!alert.isResolved && (
                              <button
                                onClick={() => handleMarkAsResolved(alert.id)}
                                className="p-1 text-gray-400 hover:text-green-500 transition-colors"
                                title="Resolve Alert"
                              >
                                <CheckCircle className="h-4 w-4" />
                              </button>
                            )}
                            
                            <Link
                              to={`/inventory/${alert.inventoryItemId}`}
                              className="p-1 text-gray-400 hover:text-primary-500 transition-colors"
                              title="View Item"
                            >
                              <ExternalLink className="h-4 w-4" />
                            </Link>
                            
                            <button className="p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
                              <MoreVertical className="h-4 w-4" />
                            </button>
                          </div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>

            
            {totalPages > 1 && (
              <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700">
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                  <div className="text-sm text-gray-500 dark:text-gray-400">
                    Showing <span className="font-medium">{((page - 1) * limit) + 1}</span> to{' '}
                    <span className="font-medium">{Math.min(page * limit, filteredAlerts.length)}</span> of{' '}
                    <span className="font-medium">{filteredAlerts.length}</span> alerts
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
                      {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                        let pageNum = i + 1;
                        if (totalPages > 5) {
                          if (page > 3) pageNum = page - 2 + i;
                          if (page > totalPages - 2) pageNum = totalPages - 4 + i;
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
                      onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                      disabled={page === totalPages}
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

      
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            Alert Configuration
          </h2>
          <button className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
            Save Settings
          </button>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          
          <div className="space-y-4">
            <h3 className="font-medium text-gray-900 dark:text-white">Notification Settings</h3>
            <div className="space-y-3">
              <label className="flex items-center space-x-3 cursor-pointer">
                <div className="relative">
                  <input type="checkbox" className="sr-only" defaultChecked />
                  <div className="w-10 h-6 rounded-full bg-primary-500">
                    <div className="absolute top-1 left-5 w-4 h-4 rounded-full bg-white"></div>
                  </div>
                </div>
                <span className="text-sm text-gray-700 dark:text-gray-300">Email Notifications</span>
              </label>
              <label className="flex items-center space-x-3 cursor-pointer">
                <div className="relative">
                  <input type="checkbox" className="sr-only" defaultChecked />
                  <div className="w-10 h-6 rounded-full bg-primary-500">
                    <div className="absolute top-1 left-5 w-4 h-4 rounded-full bg-white"></div>
                  </div>
                </div>
                <span className="text-sm text-gray-700 dark:text-gray-300">SMS Alerts</span>
              </label>
              <label className="flex items-center space-x-3 cursor-pointer">
                <div className="relative">
                  <input type="checkbox" className="sr-only" />
                  <div className="w-10 h-6 rounded-full bg-gray-300 dark:bg-gray-600">
                    <div className="absolute top-1 left-1 w-4 h-4 rounded-full bg-white"></div>
                  </div>
                </div>
                <span className="text-sm text-gray-700 dark:text-gray-300">Push Notifications</span>
              </label>
            </div>
          </div>

          
          <div className="space-y-4">
            <h3 className="font-medium text-gray-900 dark:text-white">Alert Thresholds</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm text-gray-600 dark:text-gray-400 mb-1">
                  Low Stock Threshold (%)
                </label>
                <input
                  type="number"
                  defaultValue="30"
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white"
                />
              </div>
              <div>
                <label className="block text-sm text-gray-600 dark:text-gray-400 mb-1">
                  Critical Stock Level
                </label>
                <input
                  type="number"
                  defaultValue="10"
                  className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white"
                />
              </div>
            </div>
          </div>

          
          <div className="space-y-4">
            <h3 className="font-medium text-gray-900 dark:text-white">Staff Notifications</h3>
            <div className="space-y-3">
              <button className="w-full flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg transition-colors">
                <div className="flex items-center space-x-3">
                  <User className="h-5 w-5 text-gray-400" />
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">Kitchen Manager</p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">manager@funkey.com</p>
                  </div>
                </div>
                <Mail className="h-4 w-4 text-gray-400" />
              </button>
              <button className="w-full flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 rounded-lg transition-colors">
                <div className="flex items-center space-x-3">
                  <User className="h-5 w-5 text-gray-400" />
                  <div>
                    <p className="font-medium text-gray-900 dark:text-white">Head Chef</p>
                    <p className="text-sm text-gray-500 dark:text-gray-400">chef@funkey.com</p>
                  </div>
                </div>
                <Mail className="h-4 w-4 text-gray-400" />
              </button>
              <button className="w-full text-center py-2 text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
                + Add Staff Member
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default InventoryAlerts;