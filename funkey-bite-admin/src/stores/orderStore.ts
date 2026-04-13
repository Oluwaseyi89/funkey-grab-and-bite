import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Order, OrderStatus, OrderType } from '../types';

interface OrderState {
  orders: Order[];
  filteredOrders: Order[];
  selectedOrder: Order | null;
  isLoading: boolean;
  error: string | null;
  
  statusFilter: OrderStatus | 'all';
  typeFilter: OrderType | 'all';
  dateRange: { start: string | null; end: string | null };
  searchQuery: string;
  
  currentPage: number;
  itemsPerPage: number;
  totalItems: number;
  
  setOrders: (orders: Order[]) => void;
  setSelectedOrder: (order: Order | null) => void;
  addOrder: (order: Order) => void;
  updateOrder: (id: number, updates: Partial<Order>) => void;
  deleteOrder: (id: number) => void;
  updateOrderStatus: (id: number, status: OrderStatus) => void;
  
  setStatusFilter: (status: OrderStatus | 'all') => void;
  setTypeFilter: (type: OrderType | 'all') => void;
  setDateRange: (start: string | null, end: string | null) => void;
  setSearchQuery: (query: string) => void;
  applyFilters: () => void;
  clearFilters: () => void;
  
  setCurrentPage: (page: number) => void;
  setItemsPerPage: (items: number) => void;
  setTotalItems: (total: number) => void;
  
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useOrderStore = create<OrderState>()(
  persist(
    (set, get) => ({
      orders: [],
      filteredOrders: [],
      selectedOrder: null,
      isLoading: false,
      error: null,
      
      statusFilter: 'all',
      typeFilter: 'all',
      dateRange: { start: null, end: null },
      searchQuery: '',
      
      currentPage: 1,
      itemsPerPage: 20,
      totalItems: 0,
      
      setOrders: (orders) => {
        set({ orders, filteredOrders: orders });
      },
      
      setSelectedOrder: (order) => set({ selectedOrder: order }),
      
      addOrder: (order) =>
        set((state) => ({
          orders: [order, ...state.orders],
          filteredOrders: [order, ...state.filteredOrders],
        })),
      
      updateOrder: (id, updates) =>
        set((state) => {
          const updatedOrders = state.orders.map(order =>
            order.id === id ? { ...order, ...updates } : order
          );
          const updatedFiltered = state.filteredOrders.map(order =>
            order.id === id ? { ...order, ...updates } : order
          );
          return {
            orders: updatedOrders,
            filteredOrders: updatedFiltered,
            selectedOrder: state.selectedOrder?.id === id 
              ? { ...state.selectedOrder, ...updates } 
              : state.selectedOrder,
          };
        }),
      
      deleteOrder: (id) =>
        set((state) => ({
          orders: state.orders.filter(order => order.id !== id),
          filteredOrders: state.filteredOrders.filter(order => order.id !== id),
          selectedOrder: state.selectedOrder?.id === id ? null : state.selectedOrder,
        })),
      
      updateOrderStatus: (id, status) => {
        const { updateOrder } = get();
        updateOrder(id, { status });
      },
      
      setStatusFilter: (status) => set({ statusFilter: status }),
      setTypeFilter: (type) => set({ typeFilter: type }),
      setDateRange: (start, end) => set({ dateRange: { start, end } }),
      setSearchQuery: (query) => set({ searchQuery: query }),
      
      applyFilters: () => {
        const { orders, statusFilter, typeFilter, dateRange, searchQuery } = get();
        
        let filtered = [...orders];
        
        if (statusFilter !== 'all') {
          filtered = filtered.filter(order => order.status === statusFilter);
        }
        
        if (typeFilter !== 'all') {
          filtered = filtered.filter(order => order.orderType === typeFilter);
        }
        
        if (dateRange.start && dateRange.end) {
          const startDate = new Date(dateRange.start);
          const endDate = new Date(dateRange.end);
          endDate.setHours(23, 59, 59, 999);
          
          filtered = filtered.filter(order => {
            const orderDate = new Date(order.createdAt);
            return orderDate >= startDate && orderDate <= endDate;
          });
        }
        
        if (searchQuery.trim()) {
          const query = searchQuery.toLowerCase().trim();
          filtered = filtered.filter(order =>
            order.customerName.toLowerCase().includes(query) ||
            order.orderNumber.toLowerCase().includes(query) ||
            order.customerPhone.includes(query) ||
            order.customerEmail?.toLowerCase().includes(query)
          );
        }
        
        set({ filteredOrders: filtered });
      },
      
      clearFilters: () =>
        set({
          statusFilter: 'all',
          typeFilter: 'all',
          dateRange: { start: null, end: null },
          searchQuery: '',
          filteredOrders: get().orders,
        }),
      
      setCurrentPage: (page) => set({ currentPage: page }),
      setItemsPerPage: (items) => set({ itemsPerPage: items }),
      setTotalItems: (total) => set({ totalItems: total }),
      
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),
      clearError: () => set({ error: null }),
    }),
    {
      name: 'order-storage',
      partialize: (state) => ({
        orders: state.orders,
        filters: {
          statusFilter: state.statusFilter,
          typeFilter: state.typeFilter,
          dateRange: state.dateRange,
          searchQuery: state.searchQuery,
        },
        pagination: {
          currentPage: state.currentPage,
          itemsPerPage: state.itemsPerPage,
        },
      }),
    }
  )
);