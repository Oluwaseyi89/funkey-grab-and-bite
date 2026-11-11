import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { Notification } from '../types';

const LIMIT = 50;

interface NotificationState {
  notifications: Notification[];
  // Per device: a chime for new orders and payments (on unless turned off),
  // and the browser's own notifications (off until staff turn them on).
  soundEnabled: boolean;
  browserEnabled: boolean;

  add: (notification: Omit<Notification, 'id' | 'createdAt' | 'isRead'>) => Notification;
  markAsRead: (id: string) => void;
  markAllAsRead: () => void;
  clearAll: () => void;
  setSoundEnabled: (enabled: boolean) => void;
  setBrowserEnabled: (enabled: boolean) => void;
}

const newId = () =>
  typeof crypto !== 'undefined' && 'randomUUID' in crypto
    ? crypto.randomUUID()
    : `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`;

// The bell's alerts, newest first, kept on this device so they survive a
// reload. They come from live events while the admin is open.
export const useNotificationStore = create<NotificationState>()(
  persist(
    (set) => ({
      notifications: [],
      soundEnabled: true,
      browserEnabled: false,

      add: (input) => {
        const notification: Notification = { ...input, id: newId(), createdAt: new Date().toISOString(), isRead: false };
        set((state) => ({ notifications: [notification, ...state.notifications].slice(0, LIMIT) }));
        return notification;
      },
      markAsRead: (id) =>
        set((state) => ({ notifications: state.notifications.map((n) => (n.id === id ? { ...n, isRead: true } : n)) })),
      markAllAsRead: () =>
        set((state) => ({ notifications: state.notifications.map((n) => (n.isRead ? n : { ...n, isRead: true })) })),
      clearAll: () => set({ notifications: [] }),
      setSoundEnabled: (soundEnabled) => set({ soundEnabled }),
      setBrowserEnabled: (browserEnabled) => set({ browserEnabled }),
    }),
    { name: 'admin-notifications' },
  ),
);

export const selectUnreadCount = (state: NotificationState) => state.notifications.filter((n) => !n.isRead).length;
