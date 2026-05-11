import { apiGet, apiGetResponse, apiPost, apiPut, apiPatch, apiDelete } from './apiHelpers';
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

type PaginationShape = {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
};

type KeyedListPayload<T, K extends string> = Record<K, T[]> & {
  pagination?: PaginationShape;
  count?: number;
};

const emptyPagination = (page = 1, limit = 20): PaginationShape => ({
  page,
  limit,
  total: 0,
  totalPages: 0,
});

const normalizeKeyedList = <T, K extends string>(
  payload: KeyedListPayload<T, K> | null | undefined,
  key: K,
  fallbackPage = 1,
  fallbackLimit = 20,
): PaginatedResponse<T> => {
  const items = payload?.[key] || [];
  const pagination = payload?.pagination;

  return {
    data: items,
    pagination: pagination || {
      page: fallbackPage,
      limit: fallbackLimit,
      total: payload?.count ?? items.length,
      totalPages: payload?.count ? Math.ceil(payload.count / fallbackLimit) : (items.length > 0 ? 1 : 0),
    },
  };
};

const paginateClientSide = <T>(items: T[], page = 1, limit = 20): PaginatedResponse<T> => {
  const safePage = page > 0 ? page : 1;
  const safeLimit = limit > 0 ? limit : 20;
  const start = (safePage - 1) * safeLimit;
  const pagedItems = items.slice(start, start + safeLimit);

  return {
    data: pagedItems,
    pagination: {
      page: safePage,
      limit: safeLimit,
      total: items.length,
      totalPages: items.length === 0 ? 0 : Math.ceil(items.length / safeLimit),
    },
  };
};

const normalizeMetaPaginated = <T>(
  items: T[],
  meta?: { pagination?: PaginationShape },
  fallbackPage = 1,
  fallbackLimit = 20,
): PaginatedResponse<T> => ({
  data: items,
  pagination: meta?.pagination || emptyPagination(fallbackPage, fallbackLimit),
});

export const adminLogin = (credentials: { email: string; password: string }) =>
  apiPost<AdminLoginResponse>('/api/v1/admin/auth/login', credentials);

export const adminLogout = () => apiPost<void>('/api/v1/admin/auth/logout');

export const updateAdminPassword = (data: { currentPassword: string; newPassword: string }) =>
  apiPatch<void>('/api/v1/admin/auth/password', data);

export const getAdminUsers = (params?: { page?: number; limit?: number }) =>
  apiGet<KeyedListPayload<AdminUser, 'admins'>>('/api/v1/admin/users/admins', params).then((payload) =>
    normalizeKeyedList(payload, 'admins', params?.page, params?.limit),
  );

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

export const getDashboardStats = () =>
  apiGet<AdminStats>('/api/v1/admin/dashboard/stats');

export const getTodayStats = () =>
  apiGet<AdminStats>('/api/v1/admin/dashboard/stats/today');

export const getSalesReport = (params: { from: string; to: string }) =>
  apiGet<SalesReport[]>('/api/v1/admin/reports/sales', params);

export const getOrders = (params?: {
  page?: number;
  limit?: number;
  status?: string;
  search?: string;
  userId?: number;
}) => {
  const { page = 1, limit = 20, status, search, userId } = params || {};
  const requiresClientFiltering = Boolean(search || userId);
  const requestParams = requiresClientFiltering
    ? { page: 1, limit: 1000, ...(status ? { status } : {}) }
    : { page, limit, ...(status ? { status } : {}) };

  return apiGet<KeyedListPayload<Order, 'orders'>>('/api/v1/admin/orders', requestParams).then((payload) => {
    const normalized = normalizeKeyedList(payload, 'orders', requestParams.page, requestParams.limit);
    if (!requiresClientFiltering) {
      return normalized;
    }

    let filteredOrders = normalized.data;
    if (typeof userId === 'number') {
      filteredOrders = filteredOrders.filter((order) => order.userId === userId);
    }
    if (search?.trim()) {
      const query = search.toLowerCase().trim();
      filteredOrders = filteredOrders.filter((order) =>
        order.orderNumber.toLowerCase().includes(query) ||
        order.customerName.toLowerCase().includes(query) ||
        order.customerPhone.toLowerCase().includes(query),
      );
    }

    return paginateClientSide(filteredOrders, page, limit);
  });
};

export const getOrder = (id: number) =>
  apiGet<Order>(`/api/v1/admin/orders/${id}`);

export const updateOrderStatus = (id: number, data: { status: string }) =>
  apiPatch<{ message: string; orderId: number; status: string }>(
    `/api/v1/admin/orders/${id}/status`,
    data
  );

export const getUsers = (params?: { page?: number; limit?: number }) =>
  apiGet<KeyedListPayload<User, 'users'>>('/api/v1/admin/users', params).then((payload) =>
    normalizeKeyedList(payload, 'users', params?.page, params?.limit),
  );

export const updateUserStatus = (id: number, data: { isActive: boolean }) =>
  apiPatch<{ message: string; userId: number; isActive: boolean }>(
    `/api/v1/admin/users/${id}/status`,
    data
  );

export const getMenuItems = (params?: {
  page?: number;
  limit?: number;
  categoryId?: number;
  query?: string;
}) => apiGet<MenuItem[]>('/api/v1/admin/menu/items', params);

export const getMenuItemByID = (id: number) =>
  apiGet<MenuItem>(`/api/v1/admin/menu/items/${id}`);

export const createMenuItem = (data: MenuItem) =>
  apiPost<MenuItem>('/api/v1/admin/menu/items', data);

export const updateMenuItem = (id: number, data: Partial<MenuItem>) =>
  apiPut<MenuItem>(`/api/v1/admin/menu/items/${id}`, data);

export const deleteMenuItem = (id: number) =>
  apiDelete<{ message: string; itemId: number }>(`/api/v1/admin/menu/items/${id}`);

export const getCategories = () =>
  apiGet<MenuCategory[]>('/api/v1/menu/categories');

export const createCategory = (data: MenuCategory) =>
  apiPost<MenuCategory>('/api/v1/admin/menu/categories', data);

export const updateCategory = (id: number, data: Partial<MenuCategory>) =>
  apiPut<MenuCategory>(`/api/v1/admin/menu/categories/${id}`, data);

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

export const getCateringRequests = (params?: {
  page?: number;
  limit?: number;
  status?: string;
  date?: string;
  search?: string;
  id?: number;
}) => {
  const { page = 1, limit = 20, status, date, search, id } = params || {};

  return apiGet<KeyedListPayload<CateringRequest, 'requests'>>('/api/v1/admin/catering/requests').then((payload) => {
    const normalized = normalizeKeyedList(payload, 'requests', 1, limit);
    let filteredRequests = normalized.data;

    if (typeof id === 'number') {
      filteredRequests = filteredRequests.filter((request) => request.id === id);
    }
    if (status) {
      filteredRequests = filteredRequests.filter((request) => request.status === status);
    }
    if (date) {
      filteredRequests = filteredRequests.filter((request) => request.eventDate.startsWith(date));
    }
    if (search?.trim()) {
      const query = search.toLowerCase().trim();
      filteredRequests = filteredRequests.filter((request) =>
        (request.eventName || '').toLowerCase().includes(query) ||
        request.contactName.toLowerCase().includes(query) ||
        request.contactPhone.toLowerCase().includes(query) ||
        (request.contactEmail || '').toLowerCase().includes(query),
      );
    }

    return paginateClientSide(filteredRequests, page, limit);
  });
};

export const updateCateringStatus = (id: number, data: { status: string }) =>
  apiPatch<{ message: string; requestId: number; status: string }>(
    `/api/v1/admin/catering/requests/${id}/status`,
    data
  );

export const getPromotions = (params?: {
  page?: number;
  limit?: number;
  status?: string;
  type?: string;
  search?: string;
}) => {
  const { page = 1, limit = 20, status, type, search } = params || {};
  const requiresClientFiltering = Boolean(type || search);
  const requestParams = requiresClientFiltering
    ? { page: 1, limit: 1000, ...(status ? { status } : {}) }
    : { page, limit, ...(status ? { status } : {}) };

  return apiGetResponse<Promotion[]>('/api/v1/admin/promotions', requestParams).then(({ data, meta }) => {
    const normalized = normalizeMetaPaginated(Array.isArray(data) ? data : [], meta, requestParams.page, requestParams.limit);
    if (!requiresClientFiltering) {
      return normalized;
    }

    let filteredPromotions = normalized.data;
    if (type) {
      filteredPromotions = filteredPromotions.filter((promotion) => promotion.promotionType === type);
    }
    if (search?.trim()) {
      const query = search.toLowerCase().trim();
      filteredPromotions = filteredPromotions.filter((promotion) =>
        promotion.title.toLowerCase().includes(query) ||
        promotion.description.toLowerCase().includes(query) ||
        promotion.code.toLowerCase().includes(query),
      );
    }

    return paginateClientSide(filteredPromotions, page, limit);
  });
};

export const getPromotionByID = (id: number) =>
  apiGet<Promotion>(`/api/v1/admin/promotions/${id}`);

export const createPromotion = (data: Promotion) =>
  apiPost<Promotion>('/api/v1/admin/promotions', data);

export const updatePromotion = (id: number, data: Partial<Promotion>) =>
  apiPut<Promotion>(`/api/v1/admin/promotions/${id}`, data);

export const deletePromotion = (id: number) =>
  apiDelete<{ message: string; id: number }>(`/api/v1/admin/promotions/${id}`);

export const getSettings = () =>
  apiGet<BusinessSettings>('/api/v1/admin/settings');

export const updateSettings = (data: Partial<BusinessSettings>) =>
  apiPut<BusinessSettings>('/api/v1/admin/settings', data);



