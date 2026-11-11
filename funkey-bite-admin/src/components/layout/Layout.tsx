import React from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import Sidebar from './Sidebar';
import Topbar from './Topbar';
import BottomNav from './BottomNav';
import RealtimeAlerts from './RealtimeAlerts';
import ErrorBoundary from '../common/ErrorBoundary';

// The admin frame: top bar, the sidebar from laptop width, the bottom bar on
// phones, and live alerts on every page.
const Layout: React.FC = () => {
  const location = useLocation();

  return (
    <div className="min-h-dvh bg-gray-50 dark:bg-gray-950">
      <a
        href="#admin-main"
        className="sr-only focus:not-sr-only focus:fixed focus:left-3 focus:top-3 focus:z-[100] focus:rounded-xl focus:bg-primary-600 focus:px-4 focus:py-2 focus:text-white"
      >
        Skip to content
      </a>
      <div className="print:hidden">
        <Topbar />
        <Sidebar />
      </div>
      <RealtimeAlerts />

      <main
        id="admin-main"
        tabIndex={-1}
        className="min-w-0 px-3 pb-[calc(5.5rem+env(safe-area-inset-bottom))] pt-[4.25rem] outline-none sm:px-4 md:pt-20 lg:pb-10 lg:pl-[17rem] lg:pr-8 print:p-0"
      >
        <div className="mx-auto max-w-6xl">
          <ErrorBoundary resetKey={location.pathname}>
            <Outlet />
          </ErrorBoundary>
        </div>
      </main>

      <div className="print:hidden">
        <BottomNav />
      </div>
    </div>
  );
};

export default Layout;
