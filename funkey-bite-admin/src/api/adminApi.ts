// src/api/adminApi.ts
import { apiGet, apiPost, apiPut, apiPatch, apiDelete } from './apiHelpers';
import type {
  AdminUser,
  AdminLoginResponse,
  Order,
  MenuItem,
  MenuCategory,
  InventoryItem,
  SalesReport,
  BusinessSettings,
  CateringRequest,
  Promotion,
  User,
  AdminStats,
  PaginatedResponse,
  InventoryAlert
} from '../types';

// Admin Authentication
export const adminLogin = (credentials: { email: string; password: string }) =>
  apiPost<AdminLoginResponse>('/api/v1/admin/auth/login', credentials);

export const adminLogout = () => apiPost<void>('/api/v1/admin/auth/logout');

export const updateAdminPassword = (data: { currentPassword: string; newPassword: string }) =>
  apiPatch<void>('/api/v1/admin/auth/password', data);

// Admin User Management
export const getAdminUsers = (params?: { page?: number; limit?: number }) =>
  apiGet<PaginatedResponse<AdminUser>>('/api/v1/admin/users/admins', params);

export const createAdminUser = (data: {
  username: string;
  email: string;
  password: string;
  role: string;
  isActive: boolean;
}) => apiPost<AdminUser>('/api/v1/admin/users/admins', data);

export const updateAdminUser = (id: number, data: Partial<AdminUser>) =>
  apiPut<AdminUser>(`/api/v1/admin/users/admins/${id}`, data);

export const deleteAdminUser = (id: number) =>
  apiDelete<void>(`/api/v1/admin/users/admins/${id}`);

// Dashboard
export const getDashboardStats = () =>
  apiGet<AdminStats>('/api/v1/admin/dashboard/stats');

export const getTodayStats = () =>
  apiGet<AdminStats>('/api/v1/admin/dashboard/stats/today');

// Reports
export const getSalesReport = (params: { from: string; to: string }) =>
  apiGet<SalesReport[]>('/api/v1/admin/reports/sales', params);

// Orders Management
export const getOrders = (params?: {
  page?: number;
  limit?: number;
  status?: string;
}) => apiGet<PaginatedResponse<Order>>('/api/v1/admin/orders', params);

export const getOrder = (id: number) =>
  apiGet<Order>(`/api/v1/admin/orders/${id}`);

export const updateOrderStatus = (id: number, data: { status: string }) =>
  apiPatch<{ message: string; orderId: number; status: string }>(
    `/api/v1/admin/orders/${id}/status`,
    data
  );

// User Management
export const getUsers = (params?: { page?: number; limit?: number }) =>
  apiGet<PaginatedResponse<User>>('/api/v1/admin/users', params);

export const updateUserStatus = (id: number, data: { isActive: boolean }) =>
  apiPatch<{ message: string; userId: number; isActive: boolean }>(
    `/api/v1/admin/users/${id}/status`,
    data
  );

// Menu Management
export const getMenuItems = (params?: {
  page?: number;
  limit?: number;
  categoryId?: number;
  query?: string;
}) => apiGet<MenuItem[]>('/api/v1/admin/menu/items', params);

export const createMenuItem = (data: MenuItem) =>
  apiPost<MenuItem>('/api/v1/admin/menu/items', data);

export const updateMenuItem = (id: number, data: Partial<MenuItem>) =>
  apiPut<MenuItem>(`/api/v1/admin/menu/items/${id}`, data);

export const deleteMenuItem = (id: number) =>
  apiDelete<{ message: string; itemId: number }>(`/api/v1/admin/menu/items/${id}`);

// Categories
export const getCategories = () =>
  apiGet<MenuCategory[]>('/api/v1/menu/categories');

export const createCategory = (data: MenuCategory) =>
  apiPost<MenuCategory>('/api/v1/admin/menu/categories', data);

export const updateCategory = (id: number, data: Partial<MenuCategory>) =>
  apiPut<MenuCategory>(`/api/v1/admin/menu/categories/${id}`, data);

// Inventory Management
export const getInventory = () =>
  apiGet<InventoryItem[]>('/api/v1/admin/inventory');

export const getLowStock = () =>
  apiGet<InventoryItem[]>('/api/v1/admin/inventory/low-stock');

export const updateStock = (data: {
  menuItemId: number;
  quantity: number;
  operation: string;
  reason: string;
}) => apiPatch<InventoryItem>('/api/v1/admin/inventory/stock', data);

export const restockItem = (data: { menuItemId: number; quantity: number; reason: string }) =>
  apiPost<InventoryItem>('/api/v1/admin/inventory/restock', data);

export const getInventoryAlerts = (params?: { resolved?: boolean }) =>
  apiGet<InventoryAlert[]>('/api/v1/admin/inventory/alerts', params);

// Catering Management
export const getCateringRequests = (params?: {
  page?: number;
  limit?: number;
  status?: string;
}) => apiGet<PaginatedResponse<CateringRequest>>('/api/v1/admin/catering/requests', params);

export const updateCateringStatus = (id: number, data: { status: string }) =>
  apiPatch<{ message: string; requestId: number; status: string }>(
    `/api/v1/admin/catering/requests/${id}/status`,
    data
  );

// Promotions Management
export const getPromotions = (params?: {
  page?: number;
  limit?: number;
  status?: string;
}) => apiGet<PaginatedResponse<Promotion>>('/api/v1/admin/promotions', params);

// Add this to your existing adminApi.ts file
export const getPromotionByID = (id: number) =>
  apiGet<Promotion>(`/api/v1/admin/promotions/${id}`);

export const createPromotion = (data: Promotion) =>
  apiPost<Promotion>('/api/v1/admin/promotions', data);

export const updatePromotion = (id: number, data: Partial<Promotion>) =>
  apiPut<Promotion>(`/api/v1/admin/promotions/${id}`, data);

export const deletePromotion = (id: number) =>
  apiDelete<{ message: string; id: number }>(`/api/v1/admin/promotions/${id}`);

// Settings
export const getSettings = () =>
  apiGet<BusinessSettings>('/api/v1/admin/settings');

export const updateSettings = (data: Partial<BusinessSettings>) =>
  apiPut<BusinessSettings>('/api/v1/admin/settings', data);



