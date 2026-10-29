import React, { useState } from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import {
  Home,
  Package,
  Users,
  ChefHat,
  ShoppingCart,
  BarChart3,
  Settings,
  Tag,
  BellRing,
  LogOut,
  Menu as MenuIcon,
  X,
  FolderOpen,
  PieChart,
  Package2,
  UserSquare,
} from 'lucide-react';
import { useAuthStore } from '../../stores/authStore';
import { useTheme } from '../../contexts/ThemeContext';
import { useBusinessSettings, businessInitials } from '../../hooks/useBusinessSettings';

const Sidebar: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const { logout, user } = useAuthStore();
  const { settings, isLoading: isSettingsLoading } = useBusinessSettings();
  const businessName = settings?.businessName ?? '';
  const { theme } = useTheme();
  const location = useLocation();

  const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: Home },
    { name: 'Orders', href: '/orders', icon: Package },
    { name: 'Menu', href: '/menu', icon: ChefHat },
    { name: 'Customers', href: '/customers', icon: Users },
    { name: 'Catering', href: '/catering', icon: ShoppingCart },
    { name: 'Inventory', href: '/inventory', icon: Package2 },
    { name: 'Promotions', href: '/promotions', icon: Tag },
    { name: 'Reports', href: '/reports', icon: BarChart3 },
    { name: 'Team', href: '/team', icon: UserSquare },
    { name: 'Settings', href: '/settings', icon: Settings },
  ];

  return (
    <>
      
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="lg:hidden fixed top-4 left-4 z-50 p-2 rounded-lg bg-primary-500 text-white shadow-lg"
      >
        {isOpen ? <X className="h-6 w-6" /> : <MenuIcon className="h-6 w-6" />}
      </button>

      
      {isOpen && (
        <div
          className="lg:hidden fixed inset-0 bg-black bg-opacity-50 z-40"
          onClick={() => setIsOpen(false)}
        />
      )}

      
      {/* Fixed at full height on every screen size: only the menu list scrolls,
          so the brand and the account/logout section always stay in view. */}
      <aside
        className={`${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        } lg:translate-x-0 fixed inset-y-0 left-0 z-40 w-64 h-screen bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 transform transition-transform duration-200 ease-in-out flex flex-col overflow-hidden`}
      >
        
        <div className="h-16 shrink-0 pl-20 pr-4 lg:pl-6 flex items-center border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3 min-w-0 w-full">
            <div className="h-10 w-10 shrink-0 rounded-lg bg-primary-500 flex items-center justify-center">
              <span className="text-white font-bold text-sm">{businessInitials(businessName)}</span>
            </div>
            <div className="min-w-0">
              {isSettingsLoading ? (
                <div className="h-6 w-36 rounded bg-gray-200 dark:bg-gray-700 animate-pulse" aria-label="Loading business name" />
              ) : (
                <h1 className="text-lg font-bold leading-tight text-gray-900 dark:text-white truncate" title={businessName}>
                  {businessName || 'Business name not set'}
                </h1>
              )}
              <p className="text-xs text-gray-500 dark:text-gray-400">Admin Panel</p>
            </div>
          </div>
        </div>

        
        <nav className="flex-1 min-h-0 p-3 space-y-1 overflow-y-auto">
          {navigation.map((item) => {
            const Icon = item.icon;
            const isActive = location.pathname === item.href || 
                           location.pathname.startsWith(item.href + '/');
            
            return (
              <NavLink
                key={item.name}
                to={item.href}
                onClick={() => setIsOpen(false)}
                className={`
                  flex items-center space-x-3 px-4 py-2 rounded-lg transition-colors
                  ${isActive
                    ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400'
                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                  }
                `}
              >
                <Icon className="h-5 w-5 flex-shrink-0" />
                <span className="font-medium">{item.name}</span>
                {isActive && (
                  <div className="ml-auto h-2 w-2 rounded-full bg-primary-500" />
                )}
              </NavLink>
            );
          })}
        </nav>

        
        <div className="shrink-0 p-3 border-t border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-3 px-1">
            <div className="h-9 w-9 shrink-0 rounded-full bg-primary-100 dark:bg-primary-900 flex items-center justify-center">
              <span className="text-primary-600 dark:text-primary-300 font-semibold">
                {(user?.username?.[0] ?? '?').toUpperCase()}
              </span>
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-gray-900 dark:text-white truncate">
                {user?.username ?? 'Admin'}
              </p>
              <p className="text-xs text-gray-500 dark:text-gray-400 truncate">
                {user?.email ?? ''}
              </p>
            </div>
            <button
              onClick={() => {
                logout();
                localStorage.removeItem('admin_token');
                window.location.href = '/login';
              }}
              title="Logout"
              aria-label="Logout"
              className="shrink-0 p-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
            >
              <LogOut className="h-5 w-5" />
            </button>
          </div>
        </div>
      </aside>
    </>
  );
};

export default Sidebar;


































