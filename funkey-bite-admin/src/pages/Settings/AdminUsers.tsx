// src/pages/Settings/AdminUsers.tsx
import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  Search,
  Filter,
  Plus,
  Edit,
  Trash2,
  Eye,
  Shield,
  Mail,
  User,
  CheckCircle,
  XCircle,
  MoreVertical,
  AlertCircle,
  ChevronLeft,
  ChevronRight,
  Lock,
  UserPlus,
  UserCheck,
  UserX,
  Save
} from 'lucide-react';
import { useQuery, useMutation } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { 
  getAdminUsers, 
  createAdminUser, 
  updateAdminUser, 
  deleteAdminUser 
} from '../../api/adminApi';
import type { AdminUser } from '../../types';
import { useAuthStore } from '../../stores/authStore';

// Validation schemas
const adminUserSchema = z.object({
  username: z.string().min(3, 'Username must be at least 3 characters'),
  email: z.string().email('Valid email is required'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  role: z.enum(['admin', 'manager', 'staff']),
  isActive: z.boolean().default(true),
});

const updateAdminSchema = adminUserSchema.omit({ password: true }).partial();

type AdminUserFormData = z.infer<typeof adminUserSchema>;
type AdminUserFormInput = z.input<typeof adminUserSchema>;
type AdminUserFormOutput = z.output<typeof adminUserSchema>;
type UpdateAdminFormData = z.infer<typeof updateAdminSchema>;

const AdminUsers: React.FC = () => {
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [roleFilter, setRoleFilter] = useState<string>('all');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [selectedAdmin, setSelectedAdmin] = useState<AdminUser | null>(null);
  
  const { user: currentUser } = useAuthStore();

  // Fetch admin users
  const { data: adminsData, isLoading, refetch } = useQuery({
    queryKey: ['admin-users', page, limit, statusFilter, roleFilter, searchQuery],
    queryFn: async () => {
      try {
        const params: any = { page, limit };
        if (statusFilter !== 'all') params.status = statusFilter;
        if (roleFilter !== 'all') params.role = roleFilter;
        if (searchQuery) params.search = searchQuery;
        
        const data = await getAdminUsers(params);
        return data;
      } catch (error) {
        toast.error('Failed to load admin users');
        return { data: [], pagination: { total: 0, totalPages: 0 } };
      }
    },
  });

  const admins = adminsData?.data || [];
  const pagination = adminsData?.pagination || { total: 0, totalPages: 0 };

  // Create admin mutation
  const createMutation = useMutation({
    mutationFn: createAdminUser,
    onSuccess: () => {
      toast.success('Admin user created successfully!');
      setShowCreateModal(false);
      refetch();
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to create admin user');
    },
  });

  // Update admin mutation
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateAdminFormData }) =>
      updateAdminUser(id, data),
    onSuccess: () => {
      toast.success('Admin user updated successfully!');
      setShowEditModal(false);
      setSelectedAdmin(null);
      refetch();
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to update admin user');
    },
  });

  // Delete admin mutation
  const deleteMutation = useMutation({
    mutationFn: deleteAdminUser,
    onSuccess: () => {
      toast.success('Admin user deleted successfully!');
      refetch();
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to delete admin user');
    },
  });

  // Create form
  const {
    register: registerCreate,
    handleSubmit: handleSubmitCreate,
    formState: { errors: createErrors, isSubmitting: isCreateSubmitting },
    reset: resetCreateForm,
  } = useForm<AdminUserFormInput, unknown, AdminUserFormOutput>({
    resolver: zodResolver(adminUserSchema),
    defaultValues: {
      role: 'staff',
      isActive: true,
    },
  });

  // Edit form
  const {
    register: registerEdit,
    handleSubmit: handleSubmitEdit,
    formState: { errors: editErrors, isSubmitting: isEditSubmitting },
    reset: resetEditForm,
  } = useForm<UpdateAdminFormData>({
    resolver: zodResolver(updateAdminSchema),
  });

  const handleCreate = (data: AdminUserFormOutput) => {
    createMutation.mutate(data);
  };

  const handleEdit = (data: UpdateAdminFormData) => {
    if (!selectedAdmin) return;
    updateMutation.mutate({ id: selectedAdmin.id, data });
  };

  const handleDelete = (admin: AdminUser) => {
    if (!confirm(`Are you sure you want to delete admin user "${admin.username}"? This action cannot be undone.`)) return;
    
    // Prevent self-deletion
    if (admin.id === currentUser?.id) {
      toast.error('You cannot delete your own account');
      return;
    }
    
    deleteMutation.mutate(admin.id);
  };

  const openEditModal = (admin: AdminUser) => {
    // Prevent self-modification
    if (admin.id === currentUser?.id) {
      toast.error('You cannot modify your own account from here. Use profile settings.');
      return;
    }
    
    setSelectedAdmin(admin);
    resetEditForm({
      username: admin.username,
      email: admin.email,
      role: admin.role as any,
      isActive: admin.isActive,
    });
    setShowEditModal(true);
  };

  const getRoleColor = (role: string) => {
    switch (role) {
      case 'admin':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200';
      case 'manager':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
      case 'staff':
        return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200';
    }
  };

  const getRolePermissions = (role: string) => {
    switch (role) {
      case 'admin':
        return 'Full access to all features and user management';
      case 'manager':
        return 'Can manage orders, menu, and reports';
      case 'staff':
        return 'Can view orders and update basic information';
      default:
        return 'Limited access';
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return 'Never';
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-48 animate-pulse"></div>
          <div className="h-10 w-32 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
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
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Admin User Management</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-1">
            Manage administrative accounts and permissions
          </p>
        </div>
        
        <div className="flex items-center space-x-3">
          <button
            onClick={() => setShowCreateModal(true)}
            className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
          >
            <UserPlus className="h-4 w-4" />
            <span>Add Admin User</span>
          </button>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Total Admins</p>
              <p className="text-2xl font-bold text-gray-900 dark:text-white mt-1">
                {pagination.total}
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
              <Shield className="h-6 w-6 text-primary-600 dark:text-primary-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Active Admins</p>
              <p className="text-2xl font-bold text-green-600 dark:text-green-400 mt-1">
                {admins.filter(a => a.isActive).length}
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-green-100 dark:bg-green-900 flex items-center justify-center">
              <UserCheck className="h-6 w-6 text-green-600 dark:text-green-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Inactive Admins</p>
              <p className="text-2xl font-bold text-red-600 dark:text-red-400 mt-1">
                {admins.filter(a => !a.isActive).length}
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-red-100 dark:bg-red-900 flex items-center justify-center">
              <UserX className="h-6 w-6 text-red-600 dark:text-red-400" />
            </div>
          </div>
        </div>
        
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-500 dark:text-gray-400">Last Login</p>
              <p className="text-lg font-bold text-gray-900 dark:text-white mt-1">
                {admins.length > 0 && admins[0].lastLogin 
                  ? formatDate(admins[0].lastLogin)
                  : 'N/A'}
              </p>
            </div>
            <div className="h-12 w-12 rounded-lg bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
              <Lock className="h-6 w-6 text-blue-600 dark:text-blue-400" />
            </div>
          </div>
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
                placeholder="Search by username or email..."
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
              </select>
            </div>

            <select
              value={roleFilter}
              onChange={(e) => setRoleFilter(e.target.value)}
              className="px-4 py-2 bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:text-white"
            >
              <option value="all">All Roles</option>
              <option value="admin">Admin</option>
              <option value="manager">Manager</option>
              <option value="staff">Staff</option>
            </select>
          </div>
        </div>
      </div>

      {/* Admin Users Table */}
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        {admins.length === 0 ? (
          <div className="p-8 text-center">
            <Shield className="h-16 w-16 text-gray-400 mx-auto mb-4" />
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">
              No Admin Users Found
            </h2>
            <p className="text-gray-500 dark:text-gray-400 mb-6">
              Add your first admin user to help manage the business.
            </p>
            <button
              onClick={() => setShowCreateModal(true)}
              className="inline-flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
            >
              <UserPlus className="h-4 w-4" />
              <span>Add Admin User</span>
            </button>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 dark:bg-gray-700/50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Admin User
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Role & Permissions
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Last Activity
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
                  {admins.map((admin) => {
                    const isCurrentUser = admin.id === currentUser?.id;
                    
                    return (
                      <tr key={admin.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                        <td className="px-6 py-4">
                          <div className="flex items-center space-x-3">
                            <div className="h-10 w-10 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
                              <User className="h-5 w-5 text-primary-600 dark:text-primary-400" />
                            </div>
                            <div>
                              <div className="font-medium text-gray-900 dark:text-white flex items-center space-x-2">
                                <span>{admin.username}</span>
                                {isCurrentUser && (
                                  <span className="px-2 py-0.5 text-xs bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200 rounded">
                                    You
                                  </span>
                                )}
                              </div>
                              <div className="text-sm text-gray-500 dark:text-gray-400">
                                {admin.email}
                              </div>
                              <div className="text-xs text-gray-400 dark:text-gray-500">
                                ID: {admin.id} • Created {formatDate(admin.createdAt)}
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="space-y-2">
                            <span className={`px-3 py-1 rounded-full text-xs font-medium ${getRoleColor(admin.role)}`}>
                              {admin.role.charAt(0).toUpperCase() + admin.role.slice(1)}
                            </span>
                            <p className="text-xs text-gray-500 dark:text-gray-400 max-w-xs">
                              {getRolePermissions(admin.role)}
                            </p>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex items-center space-x-2">
                            {admin.isActive ? (
                              <CheckCircle className="h-5 w-5 text-green-500" />
                            ) : (
                              <XCircle className="h-5 w-5 text-red-500" />
                            )}
                            <span className={admin.isActive ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}>
                              {admin.isActive ? 'Active' : 'Inactive'}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="text-sm text-gray-900 dark:text-white">
                            {admin.lastLogin ? formatDate(admin.lastLogin) : 'Never'}
                          </div>
                          <div className="text-xs text-gray-500 dark:text-gray-400">
                            {admin.lastLogin ? 'Last login' : 'No login history'}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex items-center space-x-2">
                            <button
                              onClick={() => openEditModal(admin)}
                              disabled={isCurrentUser}
                              className="p-1 text-gray-400 hover:text-primary-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                              title={isCurrentUser ? "Use profile settings to edit your account" : "Edit"}
                            >
                              <Edit className="h-4 w-4" />
                            </button>
                            <button
                              onClick={() => handleDelete(admin)}
                              disabled={isCurrentUser}
                              className="p-1 text-gray-400 hover:text-red-500 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                              title={isCurrentUser ? "Cannot delete your own account" : "Delete"}
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
                    <span className="font-medium">{Math.min(page * limit, pagination.total)}</span> of{' '}
                    <span className="font-medium">{pagination.total}</span> admin users
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

      {/* Create Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white dark:bg-gray-800 rounded-xl w-full max-w-md">
            <div className="p-6">
              <div className="flex items-center justify-between mb-6">
                <div className="flex items-center space-x-3">
                  <UserPlus className="h-5 w-5 text-primary-500" />
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                    Add New Admin User
                  </h2>
                </div>
                <button
                  onClick={() => setShowCreateModal(false)}
                  className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                >
                  ✕
                </button>
              </div>
              
              <form onSubmit={handleSubmitCreate(handleCreate)} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Username *
                  </label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type="text"
                      {...registerCreate('username')}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                      placeholder="e.g., johndoe"
                    />
                  </div>
                  {createErrors.username && (
                    <p className="mt-1 text-sm text-red-600">{createErrors.username.message}</p>
                  )}
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Email Address *
                  </label>
                  <div className="relative">
                    <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type="email"
                      {...registerCreate('email')}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                      placeholder="admin@example.com"
                    />
                  </div>
                  {createErrors.email && (
                    <p className="mt-1 text-sm text-red-600">{createErrors.email.message}</p>
                  )}
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Password *
                  </label>
                  <input
                    type="password"
                    {...registerCreate('password')}
                    className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    placeholder="••••••••"
                  />
                  {createErrors.password && (
                    <p className="mt-1 text-sm text-red-600">{createErrors.password.message}</p>
                  )}
                </div>
                
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Role *
                    </label>
                    <select
                      {...registerCreate('role')}
                      className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    >
                      <option value="staff">Staff</option>
                      <option value="manager">Manager</option>
                      <option value="admin">Admin</option>
                    </select>
                    {createErrors.role && (
                      <p className="mt-1 text-sm text-red-600">{createErrors.role.message}</p>
                    )}
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Status
                    </label>
                    <div className="flex items-center space-x-3">
                      <label className="flex items-center space-x-2 cursor-pointer">
                        <input
                          type="checkbox"
                          {...registerCreate('isActive')}
                          className="rounded border-gray-300 dark:border-gray-600 text-primary-500 focus:ring-primary-500"
                        />
                        <span className="text-sm text-gray-700 dark:text-gray-300">Active</span>
                      </label>
                    </div>
                  </div>
                </div>
                
                <div className="pt-4 border-t border-gray-200 dark:border-gray-700">
                  <div className="flex items-center justify-end space-x-3">
                    <button
                      type="button"
                      onClick={() => setShowCreateModal(false)}
                      className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                    >
                      Cancel
                    </button>
                    <button
                      type="submit"
                      disabled={isCreateSubmitting}
                      className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors disabled:opacity-50"
                    >
                      <UserPlus className="h-4 w-4" />
                      <span>{isCreateSubmitting ? 'Creating...' : 'Create Admin User'}</span>
                    </button>
                  </div>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}

      {/* Edit Modal */}
      {showEditModal && selectedAdmin && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white dark:bg-gray-800 rounded-xl w-full max-w-md">
            <div className="p-6">
              <div className="flex items-center justify-between mb-6">
                <div className="flex items-center space-x-3">
                  <Edit className="h-5 w-5 text-primary-500" />
                  <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                    Edit Admin User
                  </h2>
                </div>
                <button
                  onClick={() => {
                    setShowEditModal(false);
                    setSelectedAdmin(null);
                  }}
                  className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                >
                  ✕
                </button>
              </div>
              
              <form onSubmit={handleSubmitEdit(handleEdit)} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Username *
                  </label>
                  <div className="relative">
                    <User className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type="text"
                      {...registerEdit('username')}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    />
                  </div>
                  {editErrors.username && (
                    <p className="mt-1 text-sm text-red-600">{editErrors.username.message}</p>
                  )}
                </div>
                
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Email Address *
                  </label>
                  <div className="relative">
                    <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-5 w-5" />
                    <input
                      type="email"
                      {...registerEdit('email')}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    />
                  </div>
                  {editErrors.email && (
                    <p className="mt-1 text-sm text-red-600">{editErrors.email.message}</p>
                  )}
                </div>
                
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Role *
                    </label>
                    <select
                      {...registerEdit('role')}
                      className="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    >
                      <option value="staff">Staff</option>
                      <option value="manager">Manager</option>
                      <option value="admin">Admin</option>
                    </select>
                    {editErrors.role && (
                      <p className="mt-1 text-sm text-red-600">{editErrors.role.message}</p>
                    )}
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Status
                    </label>
                    <div className="flex items-center space-x-3">
                      <label className="flex items-center space-x-2 cursor-pointer">
                        <input
                          type="checkbox"
                          {...registerEdit('isActive')}
                          className="rounded border-gray-300 dark:border-gray-600 text-primary-500 focus:ring-primary-500"
                        />
                        <span className="text-sm text-gray-700 dark:text-gray-300">Active</span>
                      </label>
                    </div>
                  </div>
                </div>
                
                <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-4">
                  <div className="flex items-start space-x-3">
                    <AlertCircle className="h-5 w-5 text-yellow-500 mt-0.5 flex-shrink-0" />
                    <div className="text-sm text-yellow-700 dark:text-yellow-400">
                      <p className="font-medium">Note:</p>
                      <p>Password changes must be done by the user in their profile settings.</p>
                    </div>
                  </div>
                </div>
                
                <div className="pt-4 border-t border-gray-200 dark:border-gray-700">
                  <div className="flex items-center justify-end space-x-3">
                    <button
                      type="button"
                      onClick={() => {
                        setShowEditModal(false);
                        setSelectedAdmin(null);
                      }}
                      className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                    >
                      Cancel
                    </button>
                    <button
                      type="submit"
                      disabled={isEditSubmitting}
                      className="flex items-center space-x-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors disabled:opacity-50"
                    >
                      <Save className="h-4 w-4" />
                      <span>{isEditSubmitting ? 'Updating...' : 'Update Admin User'}</span>
                    </button>
                  </div>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}

      {/* Security Notice */}
      <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-xl p-6">
        <div className="flex items-start space-x-3">
          <Shield className="h-5 w-5 text-blue-500 mt-0.5 flex-shrink-0" />
          <div className="space-y-3">
            <h3 className="font-medium text-blue-800 dark:text-blue-300">
              Admin User Security Guidelines
            </h3>
            <ul className="text-sm text-blue-700 dark:text-blue-400 space-y-2">
              <li>• Only grant admin access to trusted team members</li>
              <li>• Use the principle of least privilege - give only necessary permissions</li>
              <li>• Regularly review and update admin user access</li>
              <li>• Deactivate admin users who no longer need access</li>
              <li>• Never share admin credentials via email or chat</li>
              <li>• All admin activities are logged for security audits</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AdminUsers;