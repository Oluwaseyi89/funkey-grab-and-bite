// src/stores/settingsStore.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { BusinessSettings, OpeningHours } from '../types';

interface SettingsState {
  settings: BusinessSettings | null;
  isLoading: boolean;
  error: string | null;
  
  // Actions
  setSettings: (settings: BusinessSettings) => void;
  updateSettings: (updates: Partial<BusinessSettings>) => void;
  updateOpeningHours: (hours: OpeningHours[]) => void;
  toggleService: (service: 'delivery' | 'pickup', value: boolean) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  resetSettings: () => void;
}

export const useSettingsStore = create<SettingsState>()(
  persist(
    (set) => ({
      // Initial state
      settings: null,
      isLoading: false,
      error: null,
      
      // Actions
      setSettings: (settings) => set({ settings }),
      
      updateSettings: (updates) =>
        set((state) => ({
          settings: state.settings ? { ...state.settings, ...updates } : null,
        })),
      
      updateOpeningHours: (hours) =>
        set((state) => ({
          settings: state.settings 
            ? { ...state.settings, openingHours: hours } 
            : null,
        })),
      
      toggleService: (service, value) =>
        set((state) => ({
          settings: state.settings 
            ? { 
                ...state.settings, 
                ...(service === 'delivery' 
                  ? { isDeliveryOpen: value } 
                  : { isPickupOpen: value }) 
              } 
            : null,
        })),
      
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),
      clearError: () => set({ error: null }),
      resetSettings: () => set({ settings: null }),
    }),
    {
      name: 'settings-storage',
      partialize: (state) => ({
        settings: state.settings,
      }),
    }
  )
);