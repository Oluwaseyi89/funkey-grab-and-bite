import React, { useEffect, useState } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { ClipboardList, Home, LayoutGrid, PartyPopper, Soup } from 'lucide-react';
import Sheet from '../ui/Sheet';
import { NAV_GROUPS, activeHref } from './navigation';
import { badgeCounts, useDashboardOverview } from '../../hooks/useDashboardOverview';

const TABS = [
  { name: 'Home', href: '/dashboard', icon: Home },
  { name: 'Orders', href: '/orders', icon: ClipboardList, badge: 'newOrders' as const },
  { name: 'Menu', href: '/menu', icon: Soup },
  { name: 'Catering', href: '/catering', icon: PartyPopper, badge: 'cateringToQuote' as const },
];

// The phone's bottom bar: the four places staff go most, and "More" for the
// rest.
const BottomNav: React.FC = () => {
  const { pathname } = useLocation();
  const [moreOpen, setMoreOpen] = useState(false);
  const { data } = useDashboardOverview();
  const counts = badgeCounts(data);
  const current = activeHref(pathname);

  useEffect(() => setMoreOpen(false), [pathname]);

  const tabActive = (href: string) => pathname === href || pathname.startsWith(`${href}/`);
  const inMore = !TABS.some((tab) => tabActive(tab.href));
  const moreCount = counts.paymentsToReview + counts.lowStock;

  return (
    <>
      <nav className="fixed inset-x-0 bottom-0 z-30 border-t border-gray-200 bg-white/95 pb-safe backdrop-blur dark:border-gray-800 dark:bg-gray-900/95 lg:hidden" aria-label="Main">
        <ul className="grid h-16 grid-cols-5">
          {TABS.map((tab) => {
            const Icon = tab.icon;
            const active = tabActive(tab.href);
            const count = tab.badge ? counts[tab.badge] : 0;
            return (
              <li key={tab.href} className="flex">
                <NavLink
                  to={tab.href}
                  aria-current={active ? 'page' : undefined}
                  className={`relative flex flex-1 flex-col items-center justify-center gap-0.5 text-[11px] font-medium ${active ? 'text-primary-600 dark:text-primary-400' : 'text-gray-500 dark:text-gray-400'}`}
                >
                  <span className={`flex h-8 w-14 items-center justify-center rounded-full ${active ? 'bg-primary-50 dark:bg-primary-500/10' : ''}`}>
                    <Icon className="h-[22px] w-[22px]" strokeWidth={active ? 2.4 : 1.9} aria-hidden="true" />
                  </span>
                  {tab.name}
                  {count > 0 && (
                    <span className="absolute left-1/2 top-1 ml-2 inline-flex h-[18px] min-w-[18px] items-center justify-center rounded-full bg-primary-600 px-1 text-[10px] font-bold text-white">
                      {count}
                      <span className="sr-only"> waiting</span>
                    </span>
                  )}
                </NavLink>
              </li>
            );
          })}
          <li className="flex">
            <button
              type="button"
              aria-haspopup="dialog"
              onClick={() => setMoreOpen(true)}
              className={`relative flex flex-1 flex-col items-center justify-center gap-0.5 text-[11px] font-medium ${inMore || moreOpen ? 'text-primary-600 dark:text-primary-400' : 'text-gray-500 dark:text-gray-400'}`}
            >
              <span className={`flex h-8 w-14 items-center justify-center rounded-full ${inMore || moreOpen ? 'bg-primary-50 dark:bg-primary-500/10' : ''}`}>
                <LayoutGrid className="h-[22px] w-[22px]" aria-hidden="true" />
              </span>
              More
              {moreCount > 0 && (
                <span className="absolute left-1/2 top-1 ml-2 inline-flex h-[18px] min-w-[18px] items-center justify-center rounded-full bg-red-600 px-1 text-[10px] font-bold text-white">
                  {moreCount}
                  <span className="sr-only"> need attention</span>
                </span>
              )}
            </button>
          </li>
        </ul>
      </nav>

      <Sheet open={moreOpen} title="Everything" onClose={() => setMoreOpen(false)}>
        <div className="space-y-5">
          {NAV_GROUPS.map((group) => (
            <section key={group.label}>
              <h3 className="mb-2 text-xs font-semibold uppercase tracking-wider text-gray-400">{group.label}</h3>
              <ul className="grid grid-cols-2 gap-2">
                {group.items.map((item) => {
                  const Icon = item.icon;
                  const count = item.badge ? counts[item.badge] : 0;
                  const active = current === item.href;
                  return (
                    <li key={item.href}>
                      <NavLink
                        to={item.href}
                        end
                        className={`flex min-h-[56px] items-center gap-2.5 rounded-2xl border px-3 text-sm font-medium ${
                          active ? 'border-primary-200 bg-primary-50 text-primary-700 dark:border-primary-500/30 dark:bg-primary-500/10 dark:text-primary-300' : 'border-gray-200 text-gray-800 dark:border-gray-800 dark:text-gray-200'
                        }`}
                      >
                        <Icon className="h-5 w-5 shrink-0" aria-hidden="true" />
                        <span className="flex-1 leading-tight">{item.name}</span>
                        {count > 0 && <span className="rounded-full bg-primary-600 px-1.5 text-xs font-bold text-white">{count}</span>}
                      </NavLink>
                    </li>
                  );
                })}
              </ul>
            </section>
          ))}
        </div>
      </Sheet>
    </>
  );
};

export default BottomNav;
