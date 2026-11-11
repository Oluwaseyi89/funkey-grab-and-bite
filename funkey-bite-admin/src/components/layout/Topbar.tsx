import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { LogOut, Moon, Settings, Sun, User } from 'lucide-react';
import { useTheme } from '../../contexts/ThemeContext';
import { useAuthStore } from '../../stores/authStore';
import { useSocket } from '../../contexts/SocketContext';
import BrandLogo from '../ui/BrandLogo';
import Sheet from '../ui/Sheet';
import NotificationBell from './NotificationBell';

// The slim top bar: logo (phones), live status, theme, alerts and account.
const Topbar: React.FC = () => {
  const navigate = useNavigate();
  const { theme, toggleTheme } = useTheme();
  const { user, logout } = useAuthStore();
  const { isConnected, isRealtimeEnabled } = useSocket();
  const [accountOpen, setAccountOpen] = useState(false);

  const signOut = () => {
    setAccountOpen(false);
    logout();
    localStorage.removeItem('admin_token');
    navigate('/login', { replace: true });
  };

  const liveLabel = isConnected ? 'Live' : isRealtimeEnabled ? 'Reconnecting' : 'Live alerts off';

  return (
    <header className="fixed inset-x-0 top-0 z-30 border-b border-gray-200 bg-white/95 backdrop-blur dark:border-gray-800 dark:bg-gray-900/95 lg:left-64">
      <div className="flex h-14 items-center gap-2 px-3 md:h-16 md:px-6">
        <Link to="/dashboard" className="flex h-8 items-center lg:hidden" aria-label="Dashboard">
          <BrandLogo className="h-8" />
        </Link>

        <span
          className={`ml-1 inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-semibold ${
            isConnected ? 'bg-leaf-50 text-leaf-700 dark:bg-leaf-500/10 dark:text-leaf-400' : 'bg-amber-50 text-amber-800 dark:bg-amber-500/10 dark:text-amber-300'
          }`}
          role="status"
          title={isConnected ? 'New orders and payments arrive instantly' : 'Live alerts will resume when the connection is back'}
        >
          <span className="relative flex h-2 w-2" aria-hidden="true">
            {isConnected && <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-leaf-400 opacity-60 motion-reduce:hidden" />}
            <span className={`relative inline-flex h-2 w-2 rounded-full ${isConnected ? 'bg-leaf-500' : 'bg-amber-500'}`} />
          </span>
          <span className="hidden sm:inline">{liveLabel}</span>
          <span className="sr-only sm:hidden">{liveLabel}</span>
        </span>

        <div className="ml-auto flex items-center gap-1">
          <button type="button" className="icon-btn" onClick={toggleTheme} aria-label={theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}>
            {theme === 'light' ? <Moon className="h-5 w-5" aria-hidden="true" /> : <Sun className="h-5 w-5" aria-hidden="true" />}
          </button>
          <NotificationBell />
          <button
            type="button"
            className="ml-1 flex items-center gap-2 rounded-full p-1 pr-1 hover:bg-gray-100 dark:hover:bg-gray-800 md:pr-3"
            onClick={() => setAccountOpen(true)}
            aria-label="Account"
          >
            <span className="flex h-9 w-9 items-center justify-center rounded-full bg-primary-600 text-sm font-bold text-white" aria-hidden="true">
              {(user?.username?.[0] ?? 'A').toUpperCase()}
            </span>
            <span className="hidden text-left md:block">
              <span className="block text-sm font-medium leading-tight text-gray-900 dark:text-white">{user?.username ?? 'Admin'}</span>
              <span className="block text-xs capitalize leading-tight text-gray-500">{user?.role ?? ''}</span>
            </span>
          </button>
        </div>
      </div>

      <Sheet open={accountOpen} title="Account" onClose={() => setAccountOpen(false)}>
        <div className="flex items-center gap-3 rounded-2xl bg-gray-50 p-4 dark:bg-gray-800/60">
          <span className="flex h-12 w-12 items-center justify-center rounded-full bg-primary-600 text-lg font-bold text-white" aria-hidden="true">
            {(user?.username?.[0] ?? 'A').toUpperCase()}
          </span>
          <div className="min-w-0">
            <p className="truncate font-semibold text-gray-900 dark:text-white">{user?.username ?? 'Admin'}</p>
            <p className="truncate text-sm text-gray-500">{user?.email}</p>
          </div>
        </div>
        <div className="mt-4 grid gap-2">
          <Link to="/settings/profile" className="btn-secondary justify-start" onClick={() => setAccountOpen(false)}>
            <User className="h-4 w-4" aria-hidden="true" /> My profile and password
          </Link>
          <Link to="/settings" className="btn-secondary justify-start" onClick={() => setAccountOpen(false)}>
            <Settings className="h-4 w-4" aria-hidden="true" /> Business settings
          </Link>
          <button type="button" className="btn justify-start text-red-600 hover:bg-red-50 dark:hover:bg-red-500/10" onClick={signOut}>
            <LogOut className="h-4 w-4" aria-hidden="true" /> Sign out
          </button>
        </div>
      </Sheet>
    </header>
  );
};

export default Topbar;
