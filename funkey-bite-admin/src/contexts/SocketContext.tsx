// src/contexts/SocketContext.tsx
import React, { createContext, useContext, useEffect, useState } from 'react';
import { io, Socket } from 'socket.io-client';
import { useAuthStore } from '../stores/authStore';
import { useOrderStore } from '../stores/orderStore';
import { useNotificationStore } from '../stores/notificationStore';
import { useCateringStore } from '../stores/cateringStore';
import { useInventoryStore } from '../stores/inventoryStore';
import { useUserStore } from '../stores/userStore';
import type { Order, CateringRequest, InventoryAlert, User } from '../types';

interface SocketContextType {
  socket: Socket | null;
  isConnected: boolean;
}

const SocketContext = createContext<SocketContextType>({
  socket: null,
  isConnected: false,
});

export const useSocket = () => useContext(SocketContext);

export const SocketProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [socket, setSocket] = useState<Socket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const { token } = useAuthStore();

  useEffect(() => {
    if (!token) return;

    const socketUrl = import.meta.env.VITE_WS_URL || 'http://localhost:8080';
    const newSocket = io(socketUrl, {
      auth: { token },
      transports: ['websocket', 'polling'],
    });

    newSocket.on('connect', () => {
      console.log('Socket connected');
      setIsConnected(true);
    });

    newSocket.on('disconnect', () => {
      console.log('Socket disconnected');
      setIsConnected(false);
    });

    newSocket.on('connect_error', (error) => {
      console.error('Socket connection error:', error);
    });

    // Real-time event handlers
    // Order updates
    newSocket.on('new_order', (order: Order) => {
      useOrderStore.getState().addOrder(order);
      useNotificationStore.getState().addNotification({
        id: Date.now(),
        userId: 0,
        type: 'order',
        title: 'New Order Received',
        message: `Order #${order.orderNumber} from ${order.customerName}`,
        isRead: false,
        referenceId: order.id,
        referenceType: 'order',
        createdAt: new Date().toISOString(),
      });
    });

    newSocket.on('order_updated', (order: Order) => {
      useOrderStore.getState().updateOrder(order.id, order);
      useNotificationStore.getState().addNotification({
        id: Date.now(),
        userId: 0,
        type: 'order',
        title: 'Order Updated',
        message: `Order #${order.orderNumber} status changed to ${order.status}`,
        isRead: false,
        referenceId: order.id,
        referenceType: 'order',
        createdAt: new Date().toISOString(),
      });
    });

    // Catering updates
    newSocket.on('new_catering_request', (request: CateringRequest) => {
      useCateringStore.getState().addRequest(request);
      useNotificationStore.getState().addNotification({
        id: Date.now(),
        userId: 0,
        type: 'catering',
        title: 'New Catering Request',
        message: `${request.contactName} - ${request.eventName || 'Event'}`,
        isRead: false,
        referenceId: request.id,
        referenceType: 'catering',
        createdAt: new Date().toISOString(),
      });
    });

    newSocket.on('catering_request_updated', (request: CateringRequest) => {
      useCateringStore.getState().updateRequest(request.id, request);
      useNotificationStore.getState().addNotification({
        id: Date.now(),
        userId: 0,
        type: 'catering',
        title: 'Catering Request Updated',
        message: `${request.contactName}'s request status: ${request.status}`,
        isRead: false,
        referenceId: request.id,
        referenceType: 'catering',
        createdAt: new Date().toISOString(),
      });
    });

    // Inventory alerts
    newSocket.on('inventory_alert', (alert: InventoryAlert) => {
      useInventoryStore.getState().addAlert(alert);
      useNotificationStore.getState().addNotification({
        id: Date.now(),
        userId: 0,
        type: 'inventory',
        title: 'Inventory Alert',
        message: alert.message,
        isRead: false,
        referenceId: alert.id,
        referenceType: 'inventory',
        createdAt: new Date().toISOString(),
      });
    });

    newSocket.on('inventory_updated', (alert: InventoryAlert) => {
      useInventoryStore.getState().updateAlert(alert.id, alert);
    });

    // Customer updates
    newSocket.on('new_customer', (customer: User) => {
      useUserStore.getState().addCustomer(customer);
      useNotificationStore.getState().addNotification({
        id: Date.now(),
        userId: 0,
        type: 'customer',
        title: 'New Customer Registered',
        message: `New customer: ${customer.fullName}`,
        isRead: false,
        referenceId: customer.id,
        referenceType: 'customer',
        createdAt: new Date().toISOString(),
      });
    });

    // System notifications
    newSocket.on('system_notification', (notification: any) => {
      useNotificationStore.getState().addNotification({
        id: Date.now(),
        userId: 0,
        type: notification.type || 'system',
        title: notification.title || 'System Notification',
        message: notification.message || 'New system notification',
        isRead: false,
        referenceId: notification.referenceId,
        referenceType: notification.referenceType,
        createdAt: new Date().toISOString(),
      });
    });

    // Menu updates (optional)
    newSocket.on('menu_updated', () => {
      useNotificationStore.getState().addNotification({
        id: Date.now(),
        userId: 0,
        type: 'menu',
        title: 'Menu Updated',
        message: 'Menu items have been updated',
        isRead: false,
        createdAt: new Date().toISOString(),
      });
    });

    // Promotion updates (optional)
    newSocket.on('promotion_updated', () => {
      useNotificationStore.getState().addNotification({
        id: Date.now(),
        userId: 0,
        type: 'promotion',
        title: 'Promotion Updated',
        message: 'Promotions have been updated',
        isRead: false,
        createdAt: new Date().toISOString(),
      });
    });

    setSocket(newSocket);

    return () => {
      // Clean up all event listeners
      newSocket.off('connect');
      newSocket.off('disconnect');
      newSocket.off('connect_error');
      newSocket.off('new_order');
      newSocket.off('order_updated');
      newSocket.off('new_catering_request');
      newSocket.off('catering_request_updated');
      newSocket.off('inventory_alert');
      newSocket.off('inventory_updated');
      newSocket.off('new_customer');
      newSocket.off('system_notification');
      newSocket.off('menu_updated');
      newSocket.off('promotion_updated');
      newSocket.disconnect();
    };
  }, [token]);

  return (
    <SocketContext.Provider value={{ socket, isConnected }}>
      {children}
    </SocketContext.Provider>
  );
};
