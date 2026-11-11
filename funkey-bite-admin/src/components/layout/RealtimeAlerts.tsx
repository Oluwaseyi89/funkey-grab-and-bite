import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { useSocket } from '../../contexts/SocketContext';
import { selectUnreadCount, useNotificationStore } from '../../stores/notificationStore';
import { playChime, showBrowserNotification, unlockAudio } from '../../utils/alerts';
import { naira } from '../../utils/format';
import type { Notification } from '../../types';

// Data that live events can change, refreshed when they arrive.
const ORDER_QUERIES = ['orders', 'recent-orders', 'order', 'dashboard-overview', 'online-payments', 'calendar-orders', 'customer-orders'];
const CATERING_QUERIES = ['catering-requests', 'catering-request', 'calendar-catering', 'dashboard-overview'];

type OrderEvent = {
  id?: number;
  orderNumber?: string;
  customerName?: string;
  totalAmount?: number;
  orderType?: string;
  status?: string;
  paymentStatus?: string;
};

// Turns live events into alerts on every page: a toast, a chime, a system
// notification when the tab is in the background, an entry in the bell, and
// fresh data on screen. Also keeps the unread count in the tab title.
const RealtimeAlerts: React.FC = () => {
  const { socket } = useSocket();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const unread = useNotificationStore(selectUnreadCount);

  // Sound needs one tap on the page first.
  useEffect(() => {
    const unlock = () => unlockAudio();
    window.addEventListener('pointerdown', unlock, { once: true });
    window.addEventListener('keydown', unlock, { once: true });
    return () => {
      window.removeEventListener('pointerdown', unlock);
      window.removeEventListener('keydown', unlock);
    };
  }, []);

  useEffect(() => {
    document.title = unread > 0 ? `(${unread}) Funkey Admin` : 'Funkey Admin';
  }, [unread]);

  useEffect(() => {
    if (!socket) return;

    const refresh = (keys: string[]) => keys.forEach((key) => queryClient.invalidateQueries({ queryKey: [key] }));

    const raise = (input: Omit<Notification, 'id' | 'createdAt' | 'isRead'>, sound: boolean) => {
      const { add, soundEnabled, browserEnabled } = useNotificationStore.getState();
      const notification = add(input);
      const open = () => {
        useNotificationStore.getState().markAsRead(notification.id);
        if (input.link) navigate(input.link);
      };
      if (sound && soundEnabled) playChime(input.urgent);
      if (browserEnabled) showBrowserNotification(input.title, input.body, open);
      toast(
        (t) => (
          <button
            type="button"
            className="flex w-full items-start gap-3 text-left"
            onClick={() => {
              toast.dismiss(t.id);
              open();
            }}
          >
            <span className={`mt-1.5 h-2.5 w-2.5 shrink-0 rounded-full ${input.urgent ? 'bg-red-500' : 'bg-leaf-500'}`} aria-hidden="true" />
            <span>
              <span className="block font-semibold">{input.title}</span>
              <span className="block text-sm opacity-80">{input.body}</span>
            </span>
          </button>
        ),
        { duration: input.urgent ? 12000 : 8000 },
      );
    };

    const onNewOrder = (order: OrderEvent) => {
      refresh(ORDER_QUERIES);
      const parts = [order.customerName, order.orderType, order.totalAmount != null ? naira(order.totalAmount) : ''].filter(Boolean);
      raise(
        { kind: 'order', title: `New order ${order.orderNumber ? `#${order.orderNumber}` : ''}`.trim(), body: parts.join(' · '), link: order.id ? `/orders/${order.id}` : '/orders' },
        true,
      );
    };

    // Status changes made by staff are already on their screen; only payments
    // arriving from customers are announced.
    const onOrderUpdated = (update: OrderEvent) => {
      refresh(ORDER_QUERIES);
      const label = update.orderNumber ? `#${update.orderNumber}` : update.id ? `order ${update.id}` : 'an order';
      const link = update.id ? `/orders/${update.id}` : '/payments';
      if (update.paymentStatus === 'paid' && update.orderNumber) {
        raise({ kind: 'payment', title: 'Payment received', body: `Paid online for ${label}`, link }, true);
      } else if (update.paymentStatus === 'needs_review') {
        raise({ kind: 'payment', title: 'Payment needs checking', body: `A payment for ${label} doesn't match. Review it in Payments.`, link: '/payments', urgent: true }, true);
      }
    };

    const onNewCatering = (request: { id?: number; reference?: string; contactName?: string; guestCount?: number }) => {
      refresh(CATERING_QUERIES);
      raise(
        {
          kind: 'catering',
          title: `New catering request ${request.reference ?? ''}`.trim(),
          body: [request.contactName, request.guestCount ? `${request.guestCount} guests` : ''].filter(Boolean).join(' · '),
          link: request.id ? `/catering/${request.id}` : '/catering',
        },
        true,
      );
    };

    const onCateringUpdated = (update: { id?: number; reference?: string; status?: string; paymentReference?: string }) => {
      refresh(CATERING_QUERIES);
      const link = update.id ? `/catering/${update.id}` : '/catering';
      if (update.status === 'accepted') {
        raise({ kind: 'catering', title: 'Catering quote accepted', body: `${update.reference ?? 'A booking'} was accepted by the customer`, link }, true);
      } else if (update.paymentReference) {
        raise({ kind: 'payment', title: 'Catering payment received', body: 'A customer paid online for a catering booking', link: '/catering' }, true);
      }
    };

    socket.on('new_order', onNewOrder);
    socket.on('order_updated', onOrderUpdated);
    socket.on('new_catering_request', onNewCatering);
    socket.on('catering_request_updated', onCateringUpdated);
    return () => {
      socket.off('new_order', onNewOrder);
      socket.off('order_updated', onOrderUpdated);
      socket.off('new_catering_request', onNewCatering);
      socket.off('catering_request_updated', onCateringUpdated);
    };
  }, [socket, queryClient, navigate]);

  return null;
};

export default RealtimeAlerts;
