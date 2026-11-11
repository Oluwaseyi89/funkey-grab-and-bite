// An alert shown in the bell, kept on this device.
export type NotificationKind = 'order' | 'payment' | 'catering' | 'system';

export interface Notification {
  id: string;
  kind: NotificationKind;
  title: string;
  body: string;
  // Where tapping it goes, e.g. /orders/12.
  link?: string;
  // Highlighted in red (e.g. a payment that needs checking).
  urgent?: boolean;
  createdAt: string;
  isRead: boolean;
}
