import type { CateringStatus, ServiceStyle } from '../../types';

export const SERVICE_STYLE_LABELS: Record<ServiceStyle, string> = {
  drop_off: 'Drop-off (food only)',
  buffet: 'Buffet with staff',
  full_service: 'Full service',
};

export const EVENT_TYPE_LABELS: Record<string, string> = {
  wedding: 'Wedding',
  party: 'Party / Owambe',
  birthday: 'Birthday',
  naming: 'Naming ceremony',
  burial: 'Burial / Remembrance',
  church: 'Church event',
  corporate: 'Corporate / Office',
  graduation: 'Graduation',
  other: 'Other',
};

export const DIETARY_LABELS: Record<string, string> = {
  vegetarian: 'Vegetarian',
  vegan: 'Vegan',
  halal: 'Halal',
  no_pork: 'No pork',
  nut_free: 'Nut-free',
  gluten_free: 'Gluten-free',
  dairy_free: 'Dairy-free',
};

export const STATUS_STYLES: Record<CateringStatus, { label: string; className: string }> = {
  pending: { label: 'Pending', className: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/40 dark:text-yellow-200' },
  quoted: { label: 'Quoted', className: 'bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-200' },
  confirmed: { label: 'Confirmed', className: 'bg-green-100 text-green-800 dark:bg-green-900/40 dark:text-green-200' },
  declined: { label: 'Declined', className: 'bg-red-100 text-red-800 dark:bg-red-900/40 dark:text-red-200' },
  completed: { label: 'Completed', className: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200' },
  cancelled: { label: 'Cancelled', className: 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300' },
};

// The statuses a request can move to next.
export const NEXT_STATUSES: Record<CateringStatus, CateringStatus[]> = {
  pending: ['quoted', 'confirmed', 'declined', 'cancelled'],
  quoted: ['confirmed', 'declined', 'cancelled', 'pending'],
  confirmed: ['completed', 'cancelled'],
  declined: ['pending'],
  completed: [],
  cancelled: ['pending'],
};

export const eventTypeLabel = (eventType: string, other?: string) =>
  eventType === 'other' && other ? other : EVENT_TYPE_LABELS[eventType] ?? eventType;

export const naira = (amount?: number) =>
  amount == null ? '—' : `₦${amount.toLocaleString('en-NG', { minimumFractionDigits: 0, maximumFractionDigits: 2 })}`;

export const formatEventDate = (date: string, time?: string) => {
  const parsed = new Date(`${date}T${time ?? '00:00'}`);
  const day = parsed.toLocaleDateString('en-NG', { weekday: 'short', day: 'numeric', month: 'short', year: 'numeric' });
  return time ? `${day}, ${time}` : day;
};
