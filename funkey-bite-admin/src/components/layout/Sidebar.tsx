import React from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { NAV_GROUPS, activeHref } from './navigation';
import BrandLogo from '../ui/BrandLogo';
import { badgeCounts, useDashboardOverview } from '../../hooks/useDashboardOverview';

// The full navigation, shown on the left from laptop width up. On phones the
// same sections are in the bottom bar and its "More" sheet.
const Sidebar: React.FC = () => {
  const { pathname } = useLocation();
  const current = activeHref(pathname);
  const { data } = useDashboardOverview();
  const counts = badgeCounts(data);

  return (
    <aside className="fixed inset-y-0 left-0 z-40 hidden w-64 flex-col border-r border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900 lg:flex" aria-label="Admin sections">
      <div className="flex h-16 shrink-0 items-center border-b border-gray-200 px-5 dark:border-gray-800">
        <NavLink to="/dashboard" aria-label="Dashboard" className="flex items-center gap-2">
          <BrandLogo className="h-9" />
          <span className="rounded-md bg-gray-100 px-1.5 py-0.5 text-[10px] font-bold uppercase tracking-wider text-gray-500 dark:bg-gray-800 dark:text-gray-400">Admin</span>
        </NavLink>
      </div>
      <nav className="min-h-0 flex-1 space-y-5 overflow-y-auto px-3 py-4">
        {NAV_GROUPS.map((group) => (
          <div key={group.label}>
            <p className="px-3 pb-1 text-[11px] font-semibold uppercase tracking-wider text-gray-400">{group.label}</p>
            <ul className="space-y-0.5">
              {group.items.map((item) => {
                const Icon = item.icon;
                const active = current === item.href;
                const count = item.badge ? counts[item.badge] : 0;
                return (
                  <li key={item.href}>
                    <NavLink
                      to={item.href}
                      end
                      aria-current={active ? 'page' : undefined}
                      className={`flex min-h-[40px] items-center gap-3 rounded-xl px-3 text-sm font-medium transition ${
                        active
                          ? 'bg-primary-50 text-primary-700 dark:bg-primary-500/10 dark:text-primary-300'
                          : 'text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800'
                      }`}
                    >
                      <Icon className="h-[18px] w-[18px] shrink-0" aria-hidden="true" />
                      <span className="flex-1">{item.name}</span>
                      {count > 0 && (
                        <span className={`rounded-full px-2 py-0.5 text-xs font-bold ${item.badge === 'paymentsToReview' ? 'bg-red-600 text-white' : 'bg-gray-900 text-white dark:bg-white dark:text-gray-900'}`}>
                          {count}
                          <span className="sr-only"> waiting</span>
                        </span>
                      )}
                    </NavLink>
                  </li>
                );
              })}
            </ul>
          </div>
        ))}
      </nav>
    </aside>
  );
};

export default Sidebar;
