import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Promotion, PromotionType, PromotionStatus } from '../types';

interface PromotionState {
  promotions: Promotion[];
  selectedPromotion: Promotion | null;
  activePromotions: Promotion[];
  expiredPromotions: Promotion[];
  
  isLoading: boolean;
  error: string | null;
  
  statusFilter: PromotionStatus | 'all';
  typeFilter: PromotionType | 'all';
  searchQuery: string;
  
  setPromotions: (promotions: Promotion[]) => void;
  setSelectedPromotion: (promotion: Promotion | null) => void;
  addPromotion: (promotion: Promotion) => void;
  updatePromotion: (id: number, updates: Partial<Promotion>) => void;
  deletePromotion: (id: number) => void;
  togglePromotionStatus: (id: number) => void;
  
  setStatusFilter: (status: PromotionStatus | 'all') => void;
  setTypeFilter: (type: PromotionType | 'all') => void;
  setSearchQuery: (query: string) => void;
  applyFilters: () => Promotion[];
  clearFilters: () => void;
  
  validatePromotionCode: (code: string) => { isValid: boolean; promotion?: Promotion };
  
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const usePromotionStore = create<PromotionState>()(
  persist(
    (set, get) => ({
      promotions: [],
      selectedPromotion: null,
      activePromotions: [],
      expiredPromotions: [],
      isLoading: false,
      error: null,
      
      statusFilter: 'all',
      typeFilter: 'all',
      searchQuery: '',
      
      setPromotions: (promotions) => {
        const now = new Date();
        
        const active = promotions.filter(promo => {
          const validFrom = new Date(promo.validFrom);
          const validUntil = new Date(promo.validUntil);
          return promo.isActive && now >= validFrom && now <= validUntil;
        });
        
        const expired = promotions.filter(promo => {
          const validUntil = new Date(promo.validUntil);
          return !promo.isActive || now > validUntil;
        });
        
        set({ 
          promotions, 
          activePromotions: active, 
          expiredPromotions: expired 
        });
      },
      
      setSelectedPromotion: (promotion) => set({ selectedPromotion: promotion }),
      
      addPromotion: (promotion) =>
        set((state) => {
          const newPromotions = [promotion, ...state.promotions];
          get().setPromotions(newPromotions);
          return { promotions: newPromotions };
        }),
      
      updatePromotion: (id, updates) =>
        set((state) => {
          const updatedPromotions = state.promotions.map(promo =>
            promo.id === id ? { ...promo, ...updates } : promo
          );
          get().setPromotions(updatedPromotions);
          return {
            promotions: updatedPromotions,
            selectedPromotion: state.selectedPromotion?.id === id 
              ? { ...state.selectedPromotion, ...updates } 
              : state.selectedPromotion,
          };
        }),
      
      deletePromotion: (id) =>
        set((state) => {
          const updatedPromotions = state.promotions.filter(promo => promo.id !== id);
          get().setPromotions(updatedPromotions);
          return {
            promotions: updatedPromotions,
            selectedPromotion: state.selectedPromotion?.id === id ? null : state.selectedPromotion,
          };
        }),
      
      togglePromotionStatus: (id) => {
        const { promotions, updatePromotion } = get();
        const promotion = promotions.find(p => p.id === id);
        if (promotion) {
          updatePromotion(id, { isActive: !promotion.isActive });
        }
      },
      
      setStatusFilter: (status) => set({ statusFilter: status }),
      setTypeFilter: (type) => set({ typeFilter: type }),
      setSearchQuery: (query) => set({ searchQuery: query }),
      
      applyFilters: () => {
        const { promotions, statusFilter, typeFilter, searchQuery } = get();
        
        let filtered = [...promotions];
        
        if (statusFilter !== 'all') {
          const now = new Date();
          filtered = filtered.filter(promo => {
            const validFrom = new Date(promo.validFrom);
            const validUntil = new Date(promo.validUntil);
            
            if (statusFilter === 'active') {
              return promo.isActive && now >= validFrom && now <= validUntil;
            } else if (statusFilter === 'expired') {
              return !promo.isActive || now > validUntil;
            } else {
              return !promo.isActive;
            }
          });
        }
        
        if (typeFilter !== 'all') {
          filtered = filtered.filter(promo => promo.promotionType === typeFilter);
        }
        
        if (searchQuery.trim()) {
          const query = searchQuery.toLowerCase().trim();
          filtered = filtered.filter(promo =>
            promo.code.toLowerCase().includes(query) ||
            promo.title.toLowerCase().includes(query) ||
            promo.description?.toLowerCase().includes(query)
          );
        }
        
        return filtered;
      },
      
      clearFilters: () => set({
        statusFilter: 'all',
        typeFilter: 'all',
        searchQuery: '',
      }),
      
      validatePromotionCode: (code) => {
        const { activePromotions } = get();
        const now = new Date();
        
        const promotion = activePromotions.find(promo => {
          const codeMatch = promo.code.toLowerCase() === code.toLowerCase();
          const withinLimit = !promo.usageLimit || promo.usedCount < promo.usageLimit;
          return codeMatch && withinLimit;
        });
        
        return {
          isValid: !!promotion,
          promotion: promotion,
        };
      },
      
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),
      clearError: () => set({ error: null }),
    }),
    {
      name: 'promotion-storage',
      partialize: (state) => ({
        promotions: state.promotions,
        filters: {
          statusFilter: state.statusFilter,
          typeFilter: state.typeFilter,
          searchQuery: state.searchQuery,
        },
      }),
    }
  )
);