import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { CateringRequest, CateringStatus } from '../types';

interface CateringState {
  // Catering Requests
  requests: CateringRequest[];
  selectedRequest: CateringRequest | null;
  isLoading: boolean;
  error: string | null;
  
  // Filters
  statusFilter: CateringStatus | 'all';
  dateFilter: string | null;
  searchQuery: string;
  
  // Calendar view
  calendarView: 'month' | 'week' | 'day';
  selectedDate: string;
  
  // Actions
  setRequests: (requests: CateringRequest[]) => void;
  setSelectedRequest: (request: CateringRequest | null) => void;
  addRequest: (request: CateringRequest) => void;
  updateRequest: (id: number, updates: Partial<CateringRequest>) => void;
  deleteRequest: (id: number) => void;
  updateRequestStatus: (id: number, status: CateringStatus) => void;
  
  // Filtering
  setStatusFilter: (status: CateringStatus | 'all') => void;
  setDateFilter: (date: string | null) => void;
  setSearchQuery: (query: string) => void;
  applyFilters: () => CateringRequest[];
  clearFilters: () => void;
  
  // Calendar
  setCalendarView: (view: 'month' | 'week' | 'day') => void;
  setSelectedDate: (date: string) => void;
  getEventsForDate: (date: string) => CateringRequest[];
  getUpcomingEvents: (limit?: number) => CateringRequest[];
  
  // Loading & Error
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useCateringStore = create<CateringState>()(
  persist(
    (set, get) => ({
      // Initial state
      requests: [],
      selectedRequest: null,
      isLoading: false,
      error: null,
      
      // Filters
      statusFilter: 'all',
      dateFilter: null,
      searchQuery: '',
      
      // Calendar
      calendarView: 'month',
      selectedDate: new Date().toISOString().split('T')[0],
      
      // Actions
      setRequests: (requests) => set({ requests }),
      
      setSelectedRequest: (request) => set({ selectedRequest: request }),
      
      addRequest: (request) =>
        set((state) => ({ requests: [request, ...state.requests] })),
      
      updateRequest: (id, updates) =>
        set((state) => {
          const updatedRequests = state.requests.map(req =>
            req.id === id ? { ...req, ...updates } : req
          );
          return {
            requests: updatedRequests,
            selectedRequest: state.selectedRequest?.id === id 
              ? { ...state.selectedRequest, ...updates } 
              : state.selectedRequest,
          };
        }),
      
      deleteRequest: (id) =>
        set((state) => ({
          requests: state.requests.filter(req => req.id !== id),
          selectedRequest: state.selectedRequest?.id === id ? null : state.selectedRequest,
        })),
      
      updateRequestStatus: (id, status) => {
        const { updateRequest } = get();
        updateRequest(id, { status });
      },
      
      // Filtering
      setStatusFilter: (status) => set({ statusFilter: status }),
      setDateFilter: (date) => set({ dateFilter: date }),
      setSearchQuery: (query) => set({ searchQuery: query }),
      
      applyFilters: () => {
        const { requests, statusFilter, dateFilter, searchQuery } = get();
        
        let filtered = [...requests];
        
        // Filter by status
        if (statusFilter !== 'all') {
          filtered = filtered.filter(request => request.status === statusFilter);
        }
        
        // Filter by date
        if (dateFilter) {
          const filterDate = new Date(dateFilter);
          filterDate.setHours(0, 0, 0, 0);
          const nextDay = new Date(filterDate);
          nextDay.setDate(nextDay.getDate() + 1);
          
          filtered = filtered.filter(request => {
            const requestDate = new Date(request.eventDate);
            return requestDate >= filterDate && requestDate < nextDay;
          });
        }
        
        // Filter by search query
        if (searchQuery.trim()) {
          const query = searchQuery.toLowerCase().trim();
          filtered = filtered.filter(request =>
            request.eventName?.toLowerCase().includes(query) ||
            request.contactName.toLowerCase().includes(query) ||
            request.contactPhone.includes(query) ||
            request.contactEmail?.toLowerCase().includes(query) ||
            request.eventType.toLowerCase().includes(query)
          );
        }
        
        return filtered;
      },
      
      clearFilters: () => set({
        statusFilter: 'all',
        dateFilter: null,
        searchQuery: '',
      }),
      
      // Calendar
      setCalendarView: (view) => set({ calendarView: view }),
      setSelectedDate: (date) => set({ selectedDate: date }),
      
      getEventsForDate: (date) => {
        const { requests } = get();
        const targetDate = new Date(date);
        targetDate.setHours(0, 0, 0, 0);
        const nextDay = new Date(targetDate);
        nextDay.setDate(nextDay.getDate() + 1);
        
        return requests.filter(request => {
          const requestDate = new Date(request.eventDate);
          return requestDate >= targetDate && requestDate < nextDay;
        });
      },
      
      getUpcomingEvents: (limit = 10) => {
        const { requests } = get();
        const now = new Date();
        
        return requests
          .filter(request => {
            const eventDate = new Date(request.eventDate);
            return eventDate >= now && request.status !== 'completed';
          })
          .sort((a, b) => new Date(a.eventDate).getTime() - new Date(b.eventDate).getTime())
          .slice(0, limit);
      },
      
      // Loading & Error
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),
      clearError: () => set({ error: null }),
    }),
    {
      name: 'catering-storage',
      partialize: (state) => ({
        requests: state.requests,
        filters: {
          statusFilter: state.statusFilter,
          dateFilter: state.dateFilter,
          searchQuery: state.searchQuery,
        },
        calendar: {
          view: state.calendarView,
          selectedDate: state.selectedDate,
        },
      }),
    }
  )
);