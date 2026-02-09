import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '../types';

interface UserState {
  // Customers
  customers: User[];
  selectedCustomer: User | null;
  isLoading: boolean;
  error: string | null;
  
  // Filters
  statusFilter: 'all' | 'active' | 'inactive';
  searchQuery: string;
  
  // Actions
  setCustomers: (customers: User[]) => void;
  setSelectedCustomer: (customer: User | null) => void;
  addCustomer: (customer: User) => void;
  updateCustomer: (id: number, updates: Partial<User>) => void;
  deleteCustomer: (id: number) => void;
  toggleCustomerStatus: (id: number) => void;
  
  // Filtering
  setStatusFilter: (status: 'all' | 'active' | 'inactive') => void;
  setSearchQuery: (query: string) => void;
  applyFilters: () => User[];
  clearFilters: () => void;
  
  // Analytics
  getActiveCustomers: () => User[];
  getNewCustomers: (days?: number) => User[];
  getCustomerStats: () => {
    total: number;
    active: number;
    newToday: number;
    newThisWeek: number;
  };
  
  // Loading & Error
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      // Initial state
      customers: [],
      selectedCustomer: null,
      isLoading: false,
      error: null,
      
      // Filters
      statusFilter: 'all',
      searchQuery: '',
      
      // Actions
      setCustomers: (customers) => set({ customers }),
      
      setSelectedCustomer: (customer) => set({ selectedCustomer: customer }),
      
      addCustomer: (customer) =>
        set((state) => ({ customers: [customer, ...state.customers] })),
      
      updateCustomer: (id, updates) =>
        set((state) => {
          const updatedCustomers = state.customers.map(customer =>
            customer.id === id ? { ...customer, ...updates } : customer
          );
          return {
            customers: updatedCustomers,
            selectedCustomer: state.selectedCustomer?.id === id 
              ? { ...state.selectedCustomer, ...updates } 
              : state.selectedCustomer,
          };
        }),
      
      deleteCustomer: (id) =>
        set((state) => ({
          customers: state.customers.filter(customer => customer.id !== id),
          selectedCustomer: state.selectedCustomer?.id === id ? null : state.selectedCustomer,
        })),
      
      toggleCustomerStatus: (id) => {
        const { customers, updateCustomer } = get();
        const customer = customers.find(c => c.id === id);
        if (customer) {
          updateCustomer(id, { isActive: !customer.isActive });
        }
      },
      
      // Filtering
      setStatusFilter: (status) => set({ statusFilter: status }),
      setSearchQuery: (query) => set({ searchQuery: query }),
      
      applyFilters: () => {
        const { customers, statusFilter, searchQuery } = get();
        
        let filtered = [...customers];
        
        // Filter by status
        if (statusFilter !== 'all') {
          const isActive = statusFilter === 'active';
          filtered = filtered.filter(customer => customer.isActive === isActive);
        }
        
        // Filter by search query
        if (searchQuery.trim()) {
          const query = searchQuery.toLowerCase().trim();
          filtered = filtered.filter(customer =>
            customer.fullName.toLowerCase().includes(query) ||
            customer.phone.includes(query) ||
            customer.email?.toLowerCase().includes(query)
          );
        }
        
        return filtered;
      },
      
      clearFilters: () => set({
        statusFilter: 'all',
        searchQuery: '',
      }),
      
      // Analytics
      getActiveCustomers: () => {
        const { customers } = get();
        return customers.filter(customer => customer.isActive);
      },
      
      getNewCustomers: (days = 7) => {
        const { customers } = get();
        const cutoffDate = new Date();
        cutoffDate.setDate(cutoffDate.getDate() - days);
        
        return customers.filter(customer => {
          const createdDate = new Date(customer.createdAt);
          return createdDate >= cutoffDate;
        });
      },
      
      getCustomerStats: () => {
        const { customers } = get();
        const now = new Date();
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        const weekAgo = new Date();
        weekAgo.setDate(weekAgo.getDate() - 7);
        
        const active = customers.filter(c => c.isActive).length;
        const newToday = customers.filter(c => {
          const created = new Date(c.createdAt);
          return created >= today;
        }).length;
        const newThisWeek = customers.filter(c => {
          const created = new Date(c.createdAt);
          return created >= weekAgo;
        }).length;
        
        return {
          total: customers.length,
          active,
          newToday,
          newThisWeek,
        };
      },
      
      // Loading & Error
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),
      clearError: () => set({ error: null }),
    }),
    {
      name: 'user-storage',
      partialize: (state) => ({
        customers: state.customers,
        filters: {
          statusFilter: state.statusFilter,
          searchQuery: state.searchQuery,
        },
      }),
    }
  )
);