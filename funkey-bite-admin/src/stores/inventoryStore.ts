import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { InventoryItem, InventoryAlert, InventoryUpdate } from '../types';

interface InventoryState {
  // Inventory Items
  items: InventoryItem[];
  selectedItem: InventoryItem | null;
  
  // Alerts
  alerts: InventoryAlert[];
  unreadAlerts: number;
  
  // Low stock items
  lowStockItems: InventoryItem[];
  
  // Loading & Error
  isLoading: boolean;
  error: string | null;
  
  // Actions
  setItems: (items: InventoryItem[]) => void;
  setSelectedItem: (item: InventoryItem | null) => void;
  addItem: (item: InventoryItem) => void;
  updateItem: (id: number, updates: Partial<InventoryItem>) => void;
  updateStock: (menuItemId: number, operation: 'add' | 'subtract' | 'set', quantity: number) => void;
  deleteItem: (id: number) => void;
  
  // Alerts
  setAlerts: (alerts: InventoryAlert[]) => void;
  addAlert: (alert: InventoryAlert) => void;
  updateAlert: (id: number, updates: Partial<InventoryAlert>) => void;
  markAlertAsRead: (id: number) => void;
  markAllAlertsAsRead: () => void;
  deleteAlert: (id: number) => void;
  
  // Low stock
  updateLowStockItems: () => void;
  
  // Loading & Error
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
}

export const useInventoryStore = create<InventoryState>()(
  persist(
    (set, get) => ({
      // Initial state
      items: [],
      selectedItem: null,
      alerts: [],
      unreadAlerts: 0,
      lowStockItems: [],
      isLoading: false,
      error: null,
      
      // Actions
      setItems: (items) => {
        set({ items });
        get().updateLowStockItems();
      },
      
      setSelectedItem: (item) => set({ selectedItem: item }),
      
      addItem: (item) =>
        set((state) => {
          const newItems = [item, ...state.items];
          return { items: newItems };
        }),
      
      updateItem: (id, updates) =>
        set((state) => {
          const updatedItems = state.items.map(item =>
            item.id === id ? { ...item, ...updates } : item
          );
          
          // Check if update triggers low stock
          const updatedItem = updatedItems.find(item => item.id === id);
          const wasLowStock = state.lowStockItems.some(item => item.id === id);
          const isNowLowStock = updatedItem && 
            updatedItem.currentStock <= updatedItem.reorderPoint;
          
          let lowStockItems = [...state.lowStockItems];
          
          if (wasLowStock && !isNowLowStock) {
            lowStockItems = lowStockItems.filter(item => item.id !== id);
          } else if (!wasLowStock && isNowLowStock) {
            lowStockItems.push(updatedItem!);
          } else if (wasLowStock && isNowLowStock) {
            lowStockItems = lowStockItems.map(item =>
              item.id === id ? updatedItem! : item
            );
          }
          
          return {
            items: updatedItems,
            lowStockItems,
            selectedItem: state.selectedItem?.id === id 
              ? { ...state.selectedItem, ...updates } 
              : state.selectedItem,
          };
        }),
      
      updateStock: (menuItemId, operation, quantity) =>
        set((state) => {
          const item = state.items.find(item => item.menuItemId === menuItemId);
          if (!item) return state;
          
          let newStock = item.currentStock;
          
          switch (operation) {
            case 'add':
              newStock += quantity;
              break;
            case 'subtract':
              newStock = Math.max(0, newStock - quantity);
              break;
            case 'set':
              newStock = Math.max(0, quantity);
              break;
          }
          
          return get().updateItem(item.id, { 
            currentStock: newStock,
            lastRestocked: new Date().toISOString(),
          });
        }),
      
      deleteItem: (id) =>
        set((state) => ({
          items: state.items.filter(item => item.id !== id),
          lowStockItems: state.lowStockItems.filter(item => item.id !== id),
          selectedItem: state.selectedItem?.id === id ? null : state.selectedItem,
        })),
      
      // Alerts
      setAlerts: (alerts) => {
        const unreadCount = alerts.filter(alert => !alert.isResolved && !alert.readAt).length;
        set({ alerts, unreadAlerts: unreadCount });
      },
      
      addAlert: (alert) =>
        set((state) => {
          const newAlerts = [alert, ...state.alerts];
          const unreadCount = newAlerts.filter(a => !a.isResolved && !a.readAt).length;
          return { 
            alerts: newAlerts, 
            unreadAlerts: unreadCount 
          };
        }),
      
      updateAlert: (id, updates) =>
        set((state) => {
          const updatedAlerts = state.alerts.map(alert =>
            alert.id === id ? { ...alert, ...updates } : alert
          );
          const unreadCount = updatedAlerts.filter(a => !a.isResolved && !a.readAt).length;
          return { 
            alerts: updatedAlerts, 
            unreadAlerts: unreadCount 
          };
        }),
      
      markAlertAsRead: (id) => get().updateAlert(id, { readAt: new Date().toISOString() }),
      
      markAllAlertsAsRead: () =>
        set((state) => {
          const updatedAlerts = state.alerts.map(alert => 
            !alert.readAt ? { ...alert, readAt: new Date().toISOString() } : alert
          );
          return { 
            alerts: updatedAlerts, 
            unreadAlerts: 0 
          };
        }),
      
      deleteAlert: (id) =>
        set((state) => {
          const updatedAlerts = state.alerts.filter(alert => alert.id !== id);
          const unreadCount = updatedAlerts.filter(a => !a.isResolved && !a.readAt).length;
          return { 
            alerts: updatedAlerts, 
            unreadAlerts: unreadCount 
          };
        }),
      
      // Low stock
      updateLowStockItems: () => {
        const { items } = get();
        const lowStock = items.filter(item => 
          item.currentStock <= item.reorderPoint && item.isActive
        );
        set({ lowStockItems: lowStock });
      },
      
      // Loading & Error
      setLoading: (loading) => set({ isLoading: loading }),
      setError: (error) => set({ error }),
      clearError: () => set({ error: null }),
    }),
    {
      name: 'inventory-storage',
      partialize: (state) => ({
        items: state.items,
        alerts: state.alerts,
        lowStockItems: state.lowStockItems,
      }),
    }
  )
);


















// import { create } from 'zustand';
// import { persist } from 'zustand/middleware';
// import type { InventoryItem, InventoryAlert } from '../types';

// interface InventoryState {
//   // Inventory Items
//   items: InventoryItem[];
//   selectedItem: InventoryItem | null;
  
//   // Alerts with read tracking (local only)
//   alerts: (InventoryAlert & { adminRead?: boolean })[];
//   unreadAlerts: number;
  
//   // Low stock items
//   lowStockItems: InventoryItem[];
  
//   // Loading & Error
//   isLoading: boolean;
//   error: string | null;
  
//   // Actions (same as before, but with adminRead field)
//   setItems: (items: InventoryItem[]) => void;
//   setSelectedItem: (item: InventoryItem | null) => void;
//   addItem: (item: InventoryItem) => void;
//   updateItem: (id: number, updates: Partial<InventoryItem>) => void;
//   updateStock: (menuItemId: number, operation: 'add' | 'subtract' | 'set', quantity: number) => void;
//   deleteItem: (id: number) => void;
  
//   // Alerts - Using adminRead for tracking
//   setAlerts: (alerts: InventoryAlert[]) => void;
//   addAlert: (alert: InventoryAlert) => void;
//   updateAlert: (id: number, updates: Partial<InventoryAlert>) => void;
//   markAlertAsRead: (id: number) => void;
//   markAllAlertsAsRead: () => void;
//   deleteAlert: (id: number) => void;
  
//   // Low stock
//   updateLowStockItems: () => void;
  
//   // Loading & Error
//   setLoading: (loading: boolean) => void;
//   setError: (error: string | null) => void;
//   clearError: () => void;
// }

// export const useInventoryStore = create<InventoryState>()(
//   persist(
//     (set, get) => ({
//       // Initial state
//       items: [],
//       selectedItem: null,
//       alerts: [],
//       unreadAlerts: 0,
//       lowStockItems: [],
//       isLoading: false,
//       error: null,
      
//       // Actions (same as before)
//       setItems: (items) => {
//         set({ items });
//         get().updateLowStockItems();
//       },
      
//       setSelectedItem: (item) => set({ selectedItem: item }),
      
//       addItem: (item) =>
//         set((state) => {
//           const newItems = [item, ...state.items];
//           return { items: newItems };
//         }),
      
//       updateItem: (id, updates) =>
//         set((state) => {
//           const updatedItems = state.items.map(item =>
//             item.id === id ? { ...item, ...updates } : item
//           );
          
//           const updatedItem = updatedItems.find(item => item.id === id);
//           const wasLowStock = state.lowStockItems.some(item => item.id === id);
//           const isNowLowStock = updatedItem && 
//             updatedItem.currentStock <= updatedItem.reorderPoint;
          
//           let lowStockItems = [...state.lowStockItems];
          
//           if (wasLowStock && !isNowLowStock) {
//             lowStockItems = lowStockItems.filter(item => item.id !== id);
//           } else if (!wasLowStock && isNowLowStock) {
//             lowStockItems.push(updatedItem!);
//           } else if (wasLowStock && isNowLowStock) {
//             lowStockItems = lowStockItems.map(item =>
//               item.id === id ? updatedItem! : item
//             );
//           }
          
//           return {
//             items: updatedItems,
//             lowStockItems,
//             selectedItem: state.selectedItem?.id === id 
//               ? { ...state.selectedItem, ...updates } 
//               : state.selectedItem,
//           };
//         }),
      
//       updateStock: (menuItemId, operation, quantity) =>
//         set((state) => {
//           const item = state.items.find(item => item.menuItemId === menuItemId);
//           if (!item) return state;
          
//           let newStock = item.currentStock;
          
//           switch (operation) {
//             case 'add':
//               newStock += quantity;
//               break;
//             case 'subtract':
//               newStock = Math.max(0, newStock - quantity);
//               break;
//             case 'set':
//               newStock = Math.max(0, quantity);
//               break;
//           }
          
//           return get().updateItem(item.id, { 
//             currentStock: newStock,
//             lastRestocked: new Date().toISOString(),
//           });
//         }),
      
//       deleteItem: (id) =>
//         set((state) => ({
//           items: state.items.filter(item => item.id !== id),
//           lowStockItems: state.lowStockItems.filter(item => item.id !== id),
//           selectedItem: state.selectedItem?.id === id ? null : state.selectedItem,
//         })),
      
//       // Alerts - MODIFIED: Use adminRead instead of readAt
//       setAlerts: (alerts: InventoryAlert[]) => {
//         // Add adminRead field if not present
//         const alertsWithRead = alerts.map(alert => ({
//           ...alert,
//           adminRead: alert.isResolved // Mark resolved alerts as read by default
//         }));
//         const unreadCount = alertsWithRead.filter(alert => !alert.isResolved && !alert.adminRead).length;
//         set({ alerts: alertsWithRead, unreadAlerts: unreadCount });
//       },
      
//       addAlert: (alert: InventoryAlert) =>
//         set((state) => {
//           const alertWithRead = { ...alert, adminRead: false };
//           const newAlerts = [alertWithRead, ...state.alerts];
//           const unreadCount = newAlerts.filter(a => !a.isResolved && !a.adminRead).length;
//           return { 
//             alerts: newAlerts, 
//             unreadAlerts: unreadCount 
//           };
//         }),
      
//       updateAlert: (id, updates) =>
//         set((state) => {
//           const updatedAlerts = state.alerts.map(alert =>
//             alert.id === id ? { ...alert, ...updates } : alert
//           );
//           const unreadCount = updatedAlerts.filter(a => !a.isResolved && !a.adminRead).length;
//           return { 
//             alerts: updatedAlerts, 
//             unreadAlerts: unreadCount 
//           };
//         }),
      
//       markAlertAsRead: (id) => get().updateAlert(id, { 
//         adminRead: true 
//       }),
      
//       markAllAlertsAsRead: () =>
//         set((state) => {
//           const updatedAlerts = state.alerts.map(alert => 
//             !alert.adminRead ? { ...alert, adminRead: true } : alert
//           );
//           return { 
//             alerts: updatedAlerts, 
//             unreadAlerts: 0 
//           };
//         }),
      
//       deleteAlert: (id) =>
//         set((state) => {
//           const alert = state.alerts.find(a => a.id === id);
//           const updatedAlerts = state.alerts.filter(alert => alert.id !== id);
//           const unreadCount = alert?.adminRead 
//             ? state.unreadAlerts 
//             : Math.max(0, state.unreadAlerts - (alert?.isResolved ? 0 : 1));
//           return { 
//             alerts: updatedAlerts, 
//             unreadAlerts: unreadCount 
//           };
//         }),
      
//       // Low stock
//       updateLowStockItems: () => {
//         const { items } = get();
//         const lowStock = items.filter(item => 
//           item.currentStock <= item.reorderPoint && item.isActive
//         );
//         set({ lowStockItems: lowStock });
//       },
      
//       // Loading & Error
//       setLoading: (loading) => set({ isLoading: loading }),
//       setError: (error) => set({ error }),
//       clearError: () => set({ error: null }),
//     }),
//     {
//       name: 'inventory-storage',
//       partialize: (state) => ({
//         items: state.items,
//         alerts: state.alerts,
//         lowStockItems: state.lowStockItems,
//       }),
//     }
//   )
// );
