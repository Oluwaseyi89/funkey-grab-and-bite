import {
  BarChart3,
  BellRing,
  CalendarDays,
  ChefHat,
  ClipboardList,
  CreditCard,
  FolderTree,
  History,
  Home,
  Package2,
  PartyPopper,
  Settings,
  ShieldCheck,
  Soup,
  Tag,
  UserSquare,
  Users,
} from 'lucide-react';
import type { LucideIcon } from 'lucide-react';

export interface NavItem {
  name: string;
  href: string;
  icon: LucideIcon;
  // Shown with a count from the dashboard overview.
  badge?: 'newOrders' | 'paymentsToReview' | 'lowStock' | 'cateringToQuote';
  // Paths that also count as this section (sub-pages).
  match?: string[];
}

export interface NavGroup {
  label: string;
  items: NavItem[];
}

// Every section of the admin, grouped as staff think about the day.
export const NAV_GROUPS: NavGroup[] = [
  {
    label: 'Today',
    items: [
      { name: 'Dashboard', href: '/dashboard', icon: Home },
      { name: 'Orders', href: '/orders', icon: ClipboardList, badge: 'newOrders' },
      { name: 'Order calendar', href: '/orders/calendar', icon: CalendarDays },
      { name: 'Payments', href: '/payments', icon: CreditCard, badge: 'paymentsToReview' },
    ],
  },
  {
    label: 'Kitchen',
    items: [
      { name: 'Menu', href: '/menu', icon: Soup },
      { name: 'Categories', href: '/menu/categories', icon: FolderTree },
      { name: 'Inventory', href: '/inventory', icon: Package2 },
      { name: 'Stock alerts', href: '/inventory/alerts', icon: BellRing, badge: 'lowStock' },
    ],
  },
  {
    label: 'Customers',
    items: [
      { name: 'Customers', href: '/customers', icon: Users },
      { name: 'Catering', href: '/catering', icon: PartyPopper, badge: 'cateringToQuote' },
      { name: 'Catering packages', href: '/catering/packages', icon: ChefHat },
      { name: 'Promotions', href: '/promotions', icon: Tag },
    ],
  },
  {
    label: 'Business',
    items: [
      { name: 'Reports', href: '/reports', icon: BarChart3 },
      { name: 'Team', href: '/team', icon: UserSquare },
      { name: 'Settings', href: '/settings', icon: Settings },
      { name: 'Admin users', href: '/settings/admins', icon: ShieldCheck },
      { name: 'Activity log', href: '/settings/activity', icon: History },
    ],
  },
];

const ALL_ITEMS = NAV_GROUPS.flatMap((group) => group.items);

// The most specific section the path belongs to (so /menu/categories is
// Categories, not Menu).
export const activeHref = (pathname: string): string | undefined =>
  ALL_ITEMS.map((item) => item.href)
    .filter((href) => pathname === href || pathname.startsWith(`${href}/`))
    .sort((a, b) => b.length - a.length)[0];
