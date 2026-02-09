import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Notification } from '../types';

interface NotificationState {
  notifications: Notification[];
  unreadCount: number;
  isLoading: boolean;
  error: string | null;
  
  // Actions
  setNotifications: (notifications: Notification[]) => void;
  addNotification: (notification: Notification) => void;
  markAsRead: (id: number) => void;
  markAllAsRead: () => void;
  deleteNotification: (id: number) => void;
  clearAll: () => void;
  
  // Filtering
  getUnreadNotifications: () => Notification[];
  getNotificationsByType: (type: string) => Notification[];
  getRecentNotifications: (limit?: number) => Notification[];
  
  // Loading & Error
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useNotificationStore = create<NotificationState>()(
  persist(
    (set, get) => ({
      // Initial state
      notifications: [],
      unreadCount: 0,
      isLoading: false,
      error: null,
      
      // Actions
      setNotifications: (notifications) => {
        const unreadCount = notifications.filter(n => !n.isRead).length;
        set({ notifications, unreadCount });
      },
      
      addNotification: (notification) =>
        set((state) => {
          const newNotifications = [notification, ...state.notifications];
          const unreadCount = newNotifications.filter(n => !n.isRead).length;
          return { 
            notifications: newNotifications, 
            unreadCount 
          };
        }),
      
      markAsRead: (id) =>
        set((state) => {
          const updatedNotifications = state.notifications.map(notif =>
            notif.id === id ? { ...notif, isRead: true, readAt: new Date().toISOString() } : notif
          );
          const unreadCount = updatedNotifications.filter(n => !n.isRead).length;
          return { 
            notifications: updatedNotifications, 
            unreadCount 
          };
        }),
      
      markAllAsRead: () =>
        set((state) => {
          const updatedNotifications = state.notifications.map(notif => 
            !notif.isRead ? { ...notif, isRead: true, readAt: new Date().toISOString() } : notif
          );
          return { 
            notifications: updatedNotifications, 
            unreadCount: 0 
          };
        }),
      
      deleteNotification: (id) =>
        set((state) => {
          const notification = state.notifications.find(n => n.id === id);
          const updatedNotifications = state.notifications.filter(notif => notif.id !== id);
          const unreadCount = notification?.isRead 
            ? state.unreadCount 
            : Math.max(0, state.unreadCount - 1);
          return { 
            notifications: updatedNotifications, 
            unreadCount 
          };
        }),
      
      clearAll: () => set({ notifications: [], unreadCount: 0 }),
      
      // Filtering
      getUnreadNotifications: () => {
        const { notifications } = get();
        return notifications.filter(notif => !notif.isRead);
      },
      
      getNotificationsByType: (type) => {
        const { notifications } = get();
        return notifications.filter(notif => notif.type === type);
      },
      
      getRecentNotifications: (limit = 10) => {
        const { notifications } = get();
        return notifications
          .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
          .slice(0, limit);
      },
      
      // Loading & Error
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),
      clearError: () => set({ error: null }),
    }),
    {
      name: 'notification-storage',
      partialize: (state) => ({
        notifications: state.notifications,
        unreadCount: state.unreadCount,
      }),
    }
  )
);