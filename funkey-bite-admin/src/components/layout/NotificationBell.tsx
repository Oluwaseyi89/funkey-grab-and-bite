import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Bell, BellOff, ChefHat, CreditCard, PartyPopper, Volume2 } from 'lucide-react';
import Sheet from '../ui/Sheet';
import { selectUnreadCount, useNotificationStore } from '../../stores/notificationStore';
import { useSocket } from '../../contexts/SocketContext';
import { browserPermission, playChime, requestBrowserPermission, unlockAudio } from '../../utils/alerts';
import { timeAgo } from '../../utils/format';
import type { NotificationKind } from '../../types';

const KIND_ICON: Record<NotificationKind, React.ElementType> = {
  order: ChefHat,
  payment: CreditCard,
  catering: PartyPopper,
  system: Bell,
};

// The bell: unread count, the latest alerts on this device, and the alert
// settings (sound and system notifications) for this device.
const NotificationBell: React.FC = () => {
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);
  const [, tick] = useState(0);
  const { isConnected, isRealtimeEnabled } = useSocket();
  const unread = useNotificationStore(selectUnreadCount);
  const { notifications, soundEnabled, browserEnabled, markAsRead, markAllAsRead, clearAll, setSoundEnabled, setBrowserEnabled } =
    useNotificationStore();
  const [permission, setPermission] = useState(browserPermission());

  // Keep "5 min ago" current while the drawer is open.
  useEffect(() => {
    if (!open) return;
    const id = setInterval(() => tick((n) => n + 1), 30000);
    return () => clearInterval(id);
  }, [open]);

  const openItem = (id: string, link?: string) => {
    markAsRead(id);
    setOpen(false);
    if (link) navigate(link);
  };

  const toggleBrowser = async () => {
    if (browserEnabled) {
      setBrowserEnabled(false);
      return;
    }
    const result = await requestBrowserPermission();
    setPermission(result);
    setBrowserEnabled(result === 'granted');
  };

  return (
    <>
      <button
        type="button"
        className="icon-btn relative"
        aria-label={unread ? `Notifications, ${unread} unread` : 'Notifications'}
        onClick={() => setOpen(true)}
      >
        <Bell className="h-5 w-5" aria-hidden="true" />
        {unread > 0 && (
          <span key={unread} className="absolute right-1 top-1 inline-flex h-[18px] min-w-[18px] items-center justify-center rounded-full bg-primary-600 px-1 text-[10px] font-bold text-white animate-pop" aria-hidden="true">
            {unread > 99 ? '99+' : unread}
          </span>
        )}
      </button>

      <Sheet
        open={open}
        title="Notifications"
        onClose={() => setOpen(false)}
        footer={
          notifications.length > 0 ? (
            <div className="flex gap-3">
              <button type="button" className="btn-secondary flex-1" onClick={markAllAsRead} disabled={!unread}>Mark all read</button>
              <button type="button" className="btn-ghost" onClick={clearAll}>Clear</button>
            </div>
          ) : undefined
        }
      >
        <p className={`mb-3 flex items-center gap-2 rounded-xl px-3 py-2 text-sm ${isConnected ? 'bg-leaf-50 text-leaf-700 dark:bg-leaf-500/10 dark:text-leaf-400' : 'bg-amber-50 text-amber-800 dark:bg-amber-500/10 dark:text-amber-300'}`}>
          <span className={`h-2 w-2 rounded-full ${isConnected ? 'bg-leaf-500' : 'bg-amber-500'}`} aria-hidden="true" />
          {isConnected ? 'Live: new orders and payments appear here instantly.' : isRealtimeEnabled ? 'Reconnecting… alerts will resume shortly.' : 'Live alerts are switched off for this site.'}
        </p>

        {notifications.length === 0 ? (
          <div className="py-10 text-center">
            <BellOff className="mx-auto h-10 w-10 text-gray-300" aria-hidden="true" />
            <p className="mt-3 font-medium text-gray-900 dark:text-white">No alerts yet</p>
            <p className="mt-1 text-sm text-gray-500">New orders, payments and catering requests will show up here.</p>
          </div>
        ) : (
          <ul className="divide-y divide-gray-100 dark:divide-gray-800">
            {notifications.map((n) => {
              const Icon = KIND_ICON[n.kind] ?? Bell;
              return (
                <li key={n.id}>
                  <button type="button" className="flex w-full items-start gap-3 py-3 text-left" onClick={() => openItem(n.id, n.link)}>
                    <span className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-xl ${n.urgent ? 'bg-red-50 text-red-600 dark:bg-red-500/10' : 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300'}`}>
                      <Icon className="h-5 w-5" aria-hidden="true" />
                    </span>
                    <span className="min-w-0 flex-1">
                      <span className={`block text-sm ${n.isRead ? 'text-gray-700 dark:text-gray-300' : 'font-semibold text-gray-900 dark:text-white'}`}>{n.title}</span>
                      {n.body && <span className="block truncate text-sm text-gray-500 dark:text-gray-400">{n.body}</span>}
                      <span className="block text-xs text-gray-400">{timeAgo(n.createdAt)}</span>
                    </span>
                    {!n.isRead && <span className="mt-2 h-2.5 w-2.5 shrink-0 rounded-full bg-primary-600" aria-label="Unread" />}
                  </button>
                </li>
              );
            })}
          </ul>
        )}

        <section className="mt-5 space-y-3 rounded-2xl bg-gray-50 p-4 dark:bg-gray-800/60" aria-label="Alert settings on this device">
          <p className="text-sm font-semibold text-gray-900 dark:text-white">On this device</p>
          <label className="flex items-center justify-between gap-3 text-sm">
            <span>Chime for new orders and payments</span>
            <input type="checkbox" className="h-5 w-5 accent-primary-600" checked={soundEnabled} onChange={(e) => setSoundEnabled(e.target.checked)} />
          </label>
          <label className="flex items-center justify-between gap-3 text-sm">
            <span>
              Phone or computer notifications
              {permission === 'denied' && <span className="block text-xs text-red-600">Blocked in your browser settings.</span>}
              {permission === 'unsupported' && <span className="block text-xs text-gray-500">Not supported in this browser.</span>}
            </span>
            <input
              type="checkbox"
              className="h-5 w-5 accent-primary-600"
              checked={browserEnabled && permission === 'granted'}
              disabled={permission === 'denied' || permission === 'unsupported'}
              onChange={toggleBrowser}
            />
          </label>
          <button
            type="button"
            className="btn-ghost -ml-2 text-sm"
            onClick={() => {
              unlockAudio();
              playChime();
            }}
          >
            <Volume2 className="h-4 w-4" aria-hidden="true" /> Play the chime
          </button>
        </section>
      </Sheet>
    </>
  );
};

export default NotificationBell;
