import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { MenuItem, MenuCategory } from '../types';

interface MenuState {
  menuItems: MenuItem[];
  selectedMenuItem: MenuItem | null;
  isLoading: boolean;
  error: string | null;
  
  categories: MenuCategory[];
  selectedCategory: MenuCategory | null;
  
  setMenuItems: (items: MenuItem[]) => void;
  setCategories: (categories: MenuCategory[]) => void;
  setSelectedMenuItem: (item: MenuItem | null) => void;
  setSelectedCategory: (category: MenuCategory | null) => void;
  addMenuItem: (item: MenuItem) => void;
  updateMenuItem: (id: number, updates: Partial<MenuItem>) => void;
  deleteMenuItem: (id: number) => void;
  addCategory: (category: MenuCategory) => void;
  updateCategory: (id: number, updates: Partial<MenuCategory>) => void;
  deleteCategory: (id: number) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useMenuStore = create<MenuState>()(
  persist(
    (set) => ({
      menuItems: [],
      selectedMenuItem: null,
      categories: [],
      selectedCategory: null,
      isLoading: false,
      error: null,
      
      setMenuItems: (items) => set({ menuItems: items }),
      setCategories: (categories) => set({ categories }),
      setSelectedMenuItem: (item) => set({ selectedMenuItem: item }),
      setSelectedCategory: (category) => set({ selectedCategory: category }),
      
      addMenuItem: (item) => 
        set((state) => ({ menuItems: [item, ...state.menuItems] })),
      
      updateMenuItem: (id, updates) =>
        set((state) => ({
          menuItems: state.menuItems.map(item =>
            item.id === id ? { ...item, ...updates } : item
          ),
        })),
      
      deleteMenuItem: (id) =>
        set((state) => ({
          menuItems: state.menuItems.filter(item => item.id !== id),
        })),
      
      addCategory: (category) =>
        set((state) => ({ categories: [category, ...state.categories] })),
      
      updateCategory: (id, updates) =>
        set((state) => ({
          categories: state.categories.map(cat =>
            cat.id === id ? { ...cat, ...updates } : cat
          ),
        })),
      
      deleteCategory: (id) =>
        set((state) => ({
          categories: state.categories.filter(cat => cat.id !== id),
        })),
      
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),
      clearError: () => set({ error: null }),
    }),
    {
      name: 'menu-storage',
      partialize: (state) => ({
        menuItems: state.menuItems,
        categories: state.categories,
      }),
    }
  )
);