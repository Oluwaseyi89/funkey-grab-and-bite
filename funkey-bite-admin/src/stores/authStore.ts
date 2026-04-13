import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { AdminUser } from '../types';

interface AuthState {
  token: string | null;
  user: AdminUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (token: string, user: AdminUser) => void;
  logout: () => void;
  setLoading: (loading: boolean) => void;
  updateUser: (user: Partial<AdminUser>) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      isAuthenticated: false,
      isLoading: false,
      login: (token, user) => 
        set({ 
          token, 
          user, 
          isAuthenticated: true,
          isLoading: false 
        }),
      logout: () => 
        set({ 
          token: null, 
          user: null, 
          isAuthenticated: false,
          isLoading: false 
        }),
      setLoading: (loading) => set({ isLoading: loading }),
      updateUser: (updatedUser) => 
        set((state) => ({ 
          user: state.user ? { ...state.user, ...updatedUser } : null 
        })),
    }),
    {
      name: 'admin-auth-storage',
      partialize: (state) => ({ 
        token: state.token, 
        user: state.user,
        isAuthenticated: state.isAuthenticated 
      }),
    }
  )
);




















