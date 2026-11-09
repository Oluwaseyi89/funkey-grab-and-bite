import React from 'react';
import { NavLink } from 'react-router-dom';

const TABS = [
  { to: '/settings', label: 'General', end: true },
  { to: '/settings/profile', label: 'Profile', end: false },
  { to: '/settings/admins', label: 'Admin users', end: false },
  { to: '/settings/activity', label: 'System activity', end: false },
];

// Links the settings pages to each other.
const SettingsTabs: React.FC = () => (
  <nav aria-label="Settings" className="flex flex-wrap gap-2 border-b border-gray-200 dark:border-gray-700 pb-3">
    {TABS.map((tab) => (
      <NavLink
        key={tab.to}
        to={tab.to}
        end={tab.end}
        className={({ isActive }) =>
          `px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
            isActive
              ? 'bg-primary-50 text-primary-600 dark:bg-primary-900/20 dark:text-primary-400'
              : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'
          }`
        }
      >
        {tab.label}
      </NavLink>
    ))}
  </nav>
);

export default SettingsTabs;
