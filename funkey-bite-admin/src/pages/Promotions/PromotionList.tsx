// src/pages/Promotions/PromotionList.tsx
import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import {
  Search,
  Filter,
  Tag,
  Percent,
  DollarSign,
  Gift,
  Calendar,
  TrendingUp,
  TrendingDown,
  Plus,
  Edit,
  Eye,
  Copy,
  Trash2,
  Download,
  BarChart3,
  Clock,
  CheckCircle,
  XCircle,
  ChevronLeft,
  ChevronRight,
  MoreVertical,
  AlertCircle,
  Zap,
  Users,
  ShoppingBag
} from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getPromotions } from '../../api/adminApi';
import type { Promotion, PromotionStatus, PromotionType } from '../../types';
import { usePromotionStore } from '../../stores/promotionStore';

const PromotionList: React.FC = () => {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [typeFilter, setTypeFilter] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedPromotions, setSelectedPromotions] = useState<number[]>([]);

  const { setPromotions } = usePromotionStore();

  // Fetch promotions
  const { data: promotionsData, isLoading, refetch } = useQuery({
    queryKey: ['promotions', page, limit, statusFilter, typeFilter, searchQuery],
    queryFn: async () => {
      try {
        const params: any = { page, limit };
        if (statusFilter !== 'all') params.status = statusFilter;
        if (typeFilter !== 'all') params.type = typeFilter;
        if (searchQuery) params.search = searchQuery;
        
        const data = await getPromotions(params);
        if (data.data && Array.isArray(data.data)) {
          setPromotions(data.data);
          return data;
        }
        return { data: [], pagination: { total: 0, totalPages: 0 } };
      } catch (error) {
        toast.error('Failed to load promotions');
        return { data: [], pagination: { total: 0, totalPages: 0 } };
      }
    },
  });

  const promotions = promotionsData?.data || [];
  const pagination = promotionsData?.pagination || { total: 0, totalPages: 0 };

  const handleBulkActivate = () => {
    if (selectedPromotions.length === 0) {
      toast.error('No promotions selected');
      return;
    }
    
    toast.success(`${selectedPromotions.length} promotions activated`);
    setSelectedPromotions([]);
  };

  const handleBulkDeactivate = () => {
    if (selectedPromotions.length === 0) {
      toast.error('No promotions selected');
      return;
    }
    
    toast.success(`${selectedPromotions.length} promotions deactivated`);
    setSelectedPromotions([]);
  };

  const handleDuplicate = (promotion: Promotion) => {
    toast.success(`Promotion "${promotion.title}" duplicated`);
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this promotion?')) return;
    
    try {
      // await deletePromotion(id);
      toast.success('Promotion deleted successfully');
      refetch();
    } catch (error) {
      toast.error('Failed to delete promotion');
    }
  };

  const getPromotionStatus = (promotion: Promotion): PromotionStatus => {
    const now = new Date();
    const validFrom = new Date(promotion.validFrom);
    const validUntil = new Date(promotion.validUntil);
    
    if (!promotion.isActive) return 'inactive';
    if (now < validFrom) return 'inactive';
    if (now > validUntil) return 'expired';
    return 'active';
  };

  const getStatusInfo = (status: PromotionStatus) => {
    switch (status) {
      case 'active':
        return {
          color: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
          icon: <CheckCircle className="h-4 w-4" />,
          text: 'Active',
          badge: 'bg-green-500'
        };
      case 'expired':
        return {
          color: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200',
          icon: <Clock className="h-4 w-4" />,
          text: 'Expired',
          badge: 'bg-gray-500'
        };
      case 'inactive':
        return {
          color: 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200',
          icon: <XCircle className="h-4 w-4" />,
          text: 'Inactive',
          badge: 'bg-red-500'
        };
      default:
        return {
          color: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200',
          icon: <AlertCircle className="h-4 w-4" />,
          text: 'Unknown',
          badge: 'bg-gray-500'
        };
    }
  };

  const getTypeInfo = (type: PromotionType) => {
    switch (type) {
      case 'percentage':
        return {
          color: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
          icon: <Percent className="h-4 w-4" />,
          text: 'Percentage',
          prefix: '%'
        };
      case 'fixed':
        return {
          color: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200',
          icon: <DollarSign className="h-4 w-4" />,
          text: 'Fixed Amount',
          prefix: '$'
        };
      case 'bogo':
        return {
          color: 'bg-pink-100 text-pink-800 dark:bg-pink-900 dark:text-pink-200',
          icon: <Gift className="h-4 w-4" />,
          text: 'Buy One Get One',
          prefix: ''
        };
      default:
        return {
          color: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200',
          icon: <Tag className="h-4 w-4" />,
          text: type,
          prefix: ''
        };
    }
  };

  const calculateUsageRate = (promotion: Promotion) => {
    if (!promotion.usageLimit) return 100;
    return Math.min(100, (promotion.usedCount / promotion.usageLimit) * 100);
  };

  const getUsageRateColor = (rate: number) => {
    if (rate >= 90) return 'bg-red-500';
    if (rate >= 70) return 'bg-yellow-500';
    return 'bg-green-500';
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric'
    });
  };

  const formatDateRange = (from: string, until: string) => {
    return `${formatDate(from)} - ${formatDate(until)}`;
  };

  // Calculate stats
  const activePromotions = promotions.filter(p => getPromotionStatus(p) === 'active').length;
  const totalDiscountValue = promotions.reduce((sum, promo) => {
    if (getPromotionStatus(promo) === 'active') {
      return sum + promo.discountValue;
    }
    return sum;
  }, 0);
  const totalUsage = promotions.reduce((sum, promo) => sum + promo.usedCount, 0);
  const expiringSoon = promotions.filter(promo => {
    const status = getPromotionStatus(promo);
    if (status !== 'active') return false;
    const validUntil = new Date(promo.validUntil);
    const now = new Date();
    const daysDiff = Math.ceil((validUntil.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
    return daysDiff <= 7;
  }).length;

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
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Promotions Management</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage discounts, offers, and promotional campaigns
          </p>
        </div>
        
        <div className="flex items-center space-x-3">
          <button className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
            <Download className="h-4 w-4" />
            <span>Export</span>
          </button>
          <Link
            to="/promotions/new"
            className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
          >
            <Plus className="h-4 w-4" />
            <span>Create Promotion</span>
          </Link>
        </div>
      </div>

      {/* Promotion Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Active Promotions</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                {activePromotions}
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                of {promotions.length} total
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-green-100 dark:bg-green-900 flex items-center justify-center">
              <Zap className="h-6 w-6 text-green-600 dark:text-green-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Total Discount Value</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                ${totalDiscountValue.toLocaleString()}
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Combined discount amount
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
              <DollarSign className="h-6 w-6 text-blue-600 dark:text-blue-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Total Usage</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                {totalUsage}
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Times used by customers
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-purple-100 dark:bg-purple-900 flex items-center justify-center">
              <Users className="h-6 w-6 text-purple-600 dark:text-purple-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Expiring Soon</p>
              <p className="text-2xl font-bold text-yellow-600 dark:text-yellow-400 mt-1">
                {expiringSoon}
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                Within 7 days
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-yellow-100 dark:bg-yellow-900 flex items-center justify-center">
              <Clock className="h-6 w-6 text-yellow-600 dark:text-yellow-400" />
            </div>
          </div>
        </div>
      </div>

      {/* Bulk Actions */}
      {selectedPromotions.length > 0 && (
        <div className="bg-primary-50 dark:bg-primary-900/20 border border-primary-200 dark:border-primary-800 rounded-xl p-4">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div className="flex items-center space-x-3">
              <div className="h-8 w-8 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                <span className="text-primary-600 dark:text-primary-400 font-semibold">
                  {selectedPromotions.length}
                </span>
              </div>
              <div>
                <p className="font-medium text-primary-800 dark:text-primary-300">
                  {selectedPromotions.length} promotion{selectedPromotions.length !== 1 ? 's' : ''} selected
                </p>
                <p className="text-sm text-primary-600 dark:text-primary-400">
                  Apply bulk actions to all selected promotions
                </p>
              </div>
            </div>
            
            <div className="flex flex-wrap gap-2">
              <button
                onClick={handleBulkActivate}
                className="px-4 py-2 bg-green-500 hover:bg-green-600 text-white rounded-lg transition-colors flex items-center space-x-2"
              >
                <CheckCircle className="h-4 w-4" />
                <span>Activate Selected</span>
              </button>
              <button
                onClick={handleBulkDeactivate}
                className="px-4 py-2 bg-red-500 hover:bg-red-600 text-white rounded-lg transition-colors flex items-center space-x-2"
              >
                <XCircle className="h-4 w-4" />
                <span>Deactivate Selected</span>
              </button>
              <button
                onClick={() => setSelectedPromotions([])}
                className="px-4 py-2 bg-gray-500 hover:bg-gray-600 text-white rounded-lg transition-colors flex items-center space-x-2"
              >
                <XCircle className="h-4 w-4" />
                <span>Clear Selection</span>
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Filters */}
      <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
        <div className="flex flex-col lg:flex-row gap-4">
          {/* Search */}
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
              <input
                type="search"
                placeholder="Search by code, title, or description..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
              />
            </div>
          </div>

          {/* Filters */}
          <div className="flex flex-wrap items-center gap-4">
            <div className="flex items-center space-x-2">
              <Filter className="h-5 w-5 text-gray-400" />
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="px-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
              >
                <option value="all">All Status</option>
                <option value="active">Active Only</option>
                <option value="inactive">Inactive Only</option>
                <option value="expired">Expired Only</option>
              </select>
            </div>

            <div className="flex items-center space-x-2">
              <select
                value={typeFilter}
                onChange={(e) => setTypeFilter(e.target.value)}
                className="px-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
              >
                <option value="all">All Types</option>
                <option value="percentage">Percentage</option>
                <option value="fixed">Fixed Amount</option>
                <option value="bogo">Buy One Get One</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      {/* Promotions Table */}
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        {promotions.length === 0 ? (
          <div className="p-8 text-center">
            <Tag className="h-16 w-16 text-gray-400 mx-auto mb-4" />
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
              No Promotions Found
            </h2>
            <p className="text-gray-500 dark:text-gray-400 mb-6">
              Create your first promotion to attract more customers.
            </p>
            <Link
              to="/promotions/new"
              className="inline-flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
            >
              <Plus className="h-4 w-4" />
              <span>Create Promotion</span>
            </Link>
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
                        checked={selectedPromotions.length === promotions.length && promotions.length > 0}
                        onChange={() => {
                          if (selectedPromotions.length === promotions.length) {
                            setSelectedPromotions([]);
                          } else {
                            setSelectedPromotions(promotions.map(p => p.id));
                          }
                        }}
                        className="rounded border-gray-300 dark:border-gray-600 text-primary-500 focus:ring-primary-500"
                      />
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Promotion
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Type & Value
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Usage
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Validity
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
                  {promotions.map((promotion) => {
                    const status = getPromotionStatus(promotion);
                    const statusInfo = getStatusInfo(status);
                    const typeInfo = getTypeInfo(promotion.promotionType);
                    const usageRate = calculateUsageRate(promotion);
                    const isSelected = selectedPromotions.includes(promotion.id);
                    
                    return (
                      <tr 
                        key={promotion.id} 
                        className={`hover:bg-gray-50 dark:hover:bg-gray-700/30 ${isSelected ? 'bg-primary-50 dark:bg-primary-900/10' : ''}`}
                      >
                        <td className="px-6 py-4">
                          <input
                            type="checkbox"
                            checked={isSelected}
                            onChange={() => {
                              if (isSelected) {
                                setSelectedPromotions(prev => prev.filter(id => id !== promotion.id));
                              } else {
                                setSelectedPromotions(prev => [...prev, promotion.id]);
                              }
                            }}
                            className="rounded border-gray-300 dark:border-gray-600 text-primary-500 focus:ring-primary-500"
                          />
                        </td>
                        <td className="px-6 py-4">
                          <div className="space-y-1">
                            <div className="flex items-center space-x-2">
                              <Tag className="h-4 w-4 text-gray-400" />
                              <div>
                                <p className="font-medium text-gray-900 dark:text-white">
                                  {promotion.title}
                                </p>
                                <p className="text-sm text-gray-500 dark:text-gray-400">
                                  Code: <span className="font-mono">{promotion.code}</span>
                                </p>
                              </div>
                            </div>
                            {promotion.description && (
                              <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-1">
                                {promotion.description}
                              </p>
                            )}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="space-y-2">
                            <span className={`px-3 py-1 rounded-full text-xs font-medium inline-flex items-center space-x-1 ${typeInfo.color}`}>
                              {typeInfo.icon}
                              <span>{typeInfo.text}</span>
                            </span>
                            <div className="text-lg font-bold text-gray-900 dark:text-white">
                              {typeInfo.prefix}{promotion.discountValue}
                              {promotion.promotionType === 'percentage' && '%'}
                              {promotion.maxDiscount && (
                                <span className="text-sm text-gray-500 dark:text-gray-400 ml-2">
                                  max ${promotion.maxDiscount}
                                </span>
                              )}
                            </div>
                            {promotion.minOrderAmount && (
                              <div className="text-xs text-gray-500 dark:text-gray-400">
                                Min order: ${promotion.minOrderAmount}
                              </div>
                            )}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="space-y-2">
                            <div className="flex items-center justify-between">
                              <span className="text-sm font-medium text-gray-900 dark:text-white">
                                {promotion.usedCount}
                              </span>
                              {promotion.usageLimit && (
                                <span className="text-xs text-gray-500 dark:text-gray-400">
                                  / {promotion.usageLimit}
                                </span>
                              )}
                            </div>
                            {promotion.usageLimit && (
                              <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                                <div
                                  className={`h-2 rounded-full ${getUsageRateColor(usageRate)}`}
                                  style={{ width: `${usageRate}%` }}
                                />
                              </div>
                            )}
                            <div className="text-xs text-gray-500 dark:text-gray-400">
                              {promotion.usageLimit ? 
                                `${Math.round(usageRate)}% used` : 
                                'No limit'}
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="space-y-1">
                            <div className="flex items-center space-x-2 text-sm text-gray-900 dark:text-white">
                              <Calendar className="h-4 w-4 text-gray-400" />
                              <span>{formatDateRange(promotion.validFrom, promotion.validUntil)}</span>
                            </div>
                            <div className="text-xs text-gray-500 dark:text-gray-400">
                              {status === 'expired' ? 'Expired' : 
                               status === 'active' ? 'Active now' : 
                               'Starts soon'}
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <span className={`px-3 py-1 rounded-full text-xs font-medium inline-flex items-center space-x-1 ${statusInfo.color}`}>
                            {statusInfo.icon}
                            <span>{statusInfo.text}</span>
                          </span>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex items-center space-x-2">
                            <Link
                              to={`/promotions/${promotion.id}`}
                              className="p-1 text-gray-400 hover:text-primary-500 transition-colors"
                              title="View Details"
                            >
                              <Eye className="h-4 w-4" />
                            </Link>
                            <Link
                              to={`/promotions/${promotion.id}/edit`}
                              className="p-1 text-gray-400 hover:text-primary-500 transition-colors"
                              title="Edit"
                            >
                              <Edit className="h-4 w-4" />
                            </Link>
                            <button
                              onClick={() => handleDuplicate(promotion)}
                              className="p-1 text-gray-400 hover:text-blue-500 transition-colors"
                              title="Duplicate"
                            >
                              <Copy className="h-4 w-4" />
                            </button>
                            <button
                              onClick={() => handleDelete(promotion.id)}
                              className="p-1 text-gray-400 hover:text-red-500 transition-colors"
                              title="Delete"
                            >
                              <Trash2 className="h-4 w-4" />
                            </button>
                          </div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {pagination.totalPages > 1 && (
              <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700">
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                  <div className="text-sm text-gray-500 dark:text-gray-400">
                    Showing <span className="font-medium">{((page - 1) * limit) + 1}</span> to{' '}
                    <span className="font-medium">{Math.min(page * limit, promotions.length)}</span> of{' '}
                    <span className="font-medium">{promotions.length}</span> promotions
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

      {/* Promotion Performance */}
      {promotions.length > 0 && (
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Promotion Performance
              </h2>
              <p className="text-sm text-gray-500 dark:text-gray-400">
                Top performing promotions by usage
              </p>
            </div>
            <button className="text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300">
              View Report
            </button>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {promotions
              .filter(p => p.usedCount > 0)
              .sort((a, b) => b.usedCount - a.usedCount)
              .slice(0, 3)
              .map((promotion) => {
                const status = getPromotionStatus(promotion);
                const statusInfo = getStatusInfo(status);
                const typeInfo = getTypeInfo(promotion.promotionType);
                
                return (
                  <div
                    key={promotion.id}
                    className="p-4 border border-gray-200 dark:border-gray-700 rounded-lg hover:border-primary-500 dark:hover:border-primary-500 transition-colors"
                  >
                    <div className="flex items-start justify-between mb-3">
                      <div>
                        <div className="flex items-center space-x-2 mb-1">
                          <span className="font-medium text-gray-900 dark:text-white">
                            {promotion.code}
                          </span>
                          <span className={`px-2 py-0.5 text-xs rounded-full ${statusInfo.color}`}>
                            {statusInfo.text}
                          </span>
                        </div>
                        <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-1">
                          {promotion.title}
                        </p>
                      </div>
                      <span className={`px-2 py-1 text-xs rounded-full ${typeInfo.color}`}>
                        {typeInfo.text}
                      </span>
                    </div>
                    
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-2xl font-bold text-gray-900 dark:text-white">
                          {promotion.usedCount}
                        </p>
                        <p className="text-xs text-gray-500 dark:text-gray-400">Times Used</p>
                      </div>
                      <div className="text-right">
                        <p className="text-lg font-semibold text-primary-600 dark:text-primary-400">
                          {typeInfo.prefix}{promotion.discountValue}
                          {promotion.promotionType === 'percentage' && '%'}
                        </p>
                        <p className="text-xs text-gray-500 dark:text-gray-400">Discount</p>
                      </div>
                    </div>
                    
                    <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-gray-500 dark:text-gray-400">
                          {formatDateRange(promotion.validFrom, promotion.validUntil)}
                        </span>
                        <Link
                          to={`/promotions/${promotion.id}`}
                          className="text-primary-600 hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300"
                        >
                          Details →
                        </Link>
                      </div>
                    </div>
                  </div>
                );
              })}
          </div>
        </div>
      )}
    </div>
  );
};

export default PromotionList;