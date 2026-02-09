import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { 
  Search, 
  Filter, 
  User, 
  Mail, 
  Phone, 
  Calendar,
  Eye,
  Edit,
  CheckCircle,
  XCircle,
  ChevronLeft,
  ChevronRight,
  Download,
  UserPlus
} from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getUsers, updateUserStatus } from '../../api/adminApi';
import type { User as CustomUser } from '../../types';

const CustomerList: React.FC = () => {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [searchQuery, setSearchQuery] = useState('');

  const { data: customersData, isLoading, refetch } = useQuery({
    queryKey: ['customers', page, limit, statusFilter, searchQuery],
    queryFn: async () => {
      try {
        const params: any = { page, limit };
        if (statusFilter) params.isActive = statusFilter === 'active';
        if (searchQuery) params.search = searchQuery;
        
        const data = await getUsers(params);
        return data;
      } catch (error) {
        toast.error('Failed to load customers');
        throw error;
      }
    },
  });

  const customers = customersData?.data || [];
  const pagination = customersData?.pagination || { total: 0, totalPages: 0 };

  const handleStatusUpdate = async (userId: number, isActive: boolean) => {
    try {
      await updateUserStatus(userId, { isActive });
      toast.success(`User ${isActive ? 'activated' : 'deactivated'} successfully`);
      refetch();
    } catch (error) {
      toast.error('Failed to update user status');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const getOrderCount = (customer: CustomUser) => {
    // This would come from API in a real app
    return Math.floor(Math.random() * 20); // Mock data
  };

  const getTotalSpent = (customer: CustomUser) => {
    // This would come from API in a real app
    return (Math.random() * 1000).toFixed(2); // Mock data
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Customers</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage and view your customer database
          </p>
        </div>
        
        <div className="flex items-center space-x-3">
          <button className="flex items-center space-x-2 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
            <Download className="h-4 w-4" />
            <span>Export</span>
          </button>
          <button className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors">
            <UserPlus className="h-4 w-4" />
            <span>Add Customer</span>
          </button>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
        <div className="flex flex-col lg:flex-row gap-4">
          {/* Search */}
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
              <input
                type="search"
                placeholder="Search by name, email, or phone..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-10 pr-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
              />
            </div>
          </div>

          {/* Status Filter */}
          <div className="flex items-center space-x-4">
            <Filter className="h-5 w-5 text-gray-400" />
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="px-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
            >
              <option value="">All Status</option>
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
              <option value="verified">Verified</option>
              <option value="unverified">Unverified</option>
            </select>
          </div>
        </div>
      </div>

      {/* Customers Table */}
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary-500 mx-auto"></div>
            <p className="mt-4 text-gray-500 dark:text-gray-400">Loading customers...</p>
          </div>
        ) : customers.length === 0 ? (
          <div className="p-8 text-center">
            <User className="h-16 w-16 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-500 dark:text-gray-400">No customers found</p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 dark:bg-gray-700/50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Customer
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Contact
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Orders
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Total Spent
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Joined
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                  {customers.map((customer) => (
                    <tr key={customer.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                      <td className="px-6 py-4">
                        <div className="flex items-center space-x-3">
                          <div className="h-10 w-10 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                            <span className="text-primary-600 dark:text-primary-400 font-semibold">
                              {customer.fullName.charAt(0)}
                            </span>
                          </div>
                          <div>
                            <div className="font-medium text-gray-900 dark:text-white">
                              {customer.fullName}
                            </div>
                            <div className="text-sm text-gray-500 dark:text-gray-400">
                              ID: {customer.id}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="space-y-1">
                          <div className="flex items-center space-x-2">
                            <Phone className="h-3 w-3 text-gray-400" />
                            <span className="text-sm text-gray-900 dark:text-white">
                              {customer.phone}
                            </span>
                          </div>
                          {customer.email && (
                            <div className="flex items-center space-x-2">
                              <Mail className="h-3 w-3 text-gray-400" />
                              <span className="text-sm text-gray-500 dark:text-gray-400 truncate max-w-xs">
                                {customer.email}
                              </span>
                            </div>
                          )}
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="text-center">
                          <span className="font-semibold text-gray-900 dark:text-white">
                            {getOrderCount(customer)}
                          </span>
                          <p className="text-xs text-gray-500 dark:text-gray-400">orders</p>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="font-semibold text-gray-900 dark:text-white">
                          ${getTotalSpent(customer)}
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center space-x-3">
                          <div className="flex flex-col space-y-1">
                            <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                              customer.isActive
                                ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                                : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                            }`}>
                              {customer.isActive ? 'Active' : 'Inactive'}
                            </span>
                            <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                              customer.isVerified
                                ? 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
                                : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200'
                            }`}>
                              {customer.isVerified ? 'Verified' : 'Unverified'}
                            </span>
                          </div>
                          <button
                            onClick={() => handleStatusUpdate(customer.id, !customer.isActive)}
                            className="p-1 text-gray-400 hover:text-primary-500 transition-colors"
                            title={customer.isActive ? 'Deactivate' : 'Activate'}
                          >
                            {customer.isActive ? (
                              <XCircle className="h-4 w-4" />
                            ) : (
                              <CheckCircle className="h-4 w-4" />
                            )}
                          </button>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center space-x-2 text-sm text-gray-500 dark:text-gray-400">
                          <Calendar className="h-4 w-4" />
                          <span>{formatDate(customer.createdAt)}</span>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div className="flex items-center space-x-2">
                          <Link
                            to={`/customers/${customer.id}`}
                            className="p-1 text-gray-400 hover:text-primary-500 transition-colors"
                            title="View Details"
                          >
                            <Eye className="h-5 w-5" />
                          </Link>
                          <button className="p-1 text-gray-400 hover:text-primary-500 transition-colors" title="Edit">
                            <Edit className="h-5 w-5" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {pagination.totalPages > 1 && (
              <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700">
                <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                  <div className="text-sm text-gray-500 dark:text-gray-400">
                    Showing <span className="font-medium">{((page - 1) * limit) + 1}</span> to{' '}
                    <span className="font-medium">{Math.min(page * limit, pagination.total)}</span> of{' '}
                    <span className="font-medium">{pagination.total}</span> customers
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

      {/* Customer Stats */}
      {customers.length > 0 && (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
            <div className="text-center">
              <p className="text-2xl font-bold text-gray-900 dark:text-white">{pagination.total}</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">Total Customers</p>
            </div>
          </div>
          <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
            <div className="text-center">
              <p className="text-2xl font-bold text-green-600">{customers.filter(c => c.isActive).length}</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">Active Customers</p>
            </div>
          </div>
          <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
            <div className="text-center">
              <p className="text-2xl font-bold text-blue-600">{customers.filter(c => c.isVerified).length}</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">Verified Accounts</p>
            </div>
          </div>
          <div className="bg-white dark:bg-gray-800 rounded-xl p-4 border border-gray-200 dark:border-gray-700">
            <div className="text-center">
              <p className="text-2xl font-bold text-purple-600">
                {Math.round(customers.reduce((acc, c) => acc + getOrderCount(c), 0) / customers.length) || 0}
              </p>
              <p className="text-sm text-gray-500 dark:text-gray-400">Avg Orders/Customer</p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default CustomerList;