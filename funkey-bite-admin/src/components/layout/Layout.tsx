import React, { useEffect } from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import toast from 'react-hot-toast';
import { useQueryClient } from '@tanstack/react-query';
import Sidebar from './Sidebar';
import Topbar from './Topbar';
import Breadcrumbs from './Breadcrumbs';
import { useSocket } from '../../contexts/SocketContext';
import { Bell, Wifi, WifiOff } from 'lucide-react';
import { useBusinessSettings } from '../../hooks/useBusinessSettings';
import ErrorBoundary from '../common/ErrorBoundary';

const Layout: React.FC = () => {
  const { socket, isConnected, isRealtimeEnabled } = useSocket();
  const { settings } = useBusinessSettings();
  const queryClient = useQueryClient();
  const location = useLocation();

  // New catering requests are announced on every page, not just Catering.
  useEffect(() => {
    if (!socket) return;
    const handleNewCateringRequest = (request: { reference?: string; contactName?: string; guestCount?: number }) => {
      toast.success(`New catering request ${request.reference ?? ''} from ${request.contactName ?? 'a customer'} (${request.guestCount ?? '?'} guests)`, { duration: 8000 });
      queryClient.invalidateQueries({ queryKey: ['catering-requests'] });
    };
    // Quotes accepted and payments made by customers.
    const handleCateringUpdate = (update: { reference?: string; status?: string; paymentReference?: string }) => {
      if (update.status === 'accepted') {
        toast.success(`Catering quote ${update.reference ?? ''} accepted by the customer`, { duration: 8000 });
      } else if (update.paymentReference) {
        toast.success('Online catering payment received', { duration: 8000 });
      }
      queryClient.invalidateQueries({ queryKey: ['catering-requests'] });
      queryClient.invalidateQueries({ queryKey: ['catering-request'] });
    };
    socket.on('new_catering_request', handleNewCateringRequest);
    socket.on('catering_request_updated', handleCateringUpdate);
    return () => {
      socket.off('new_catering_request', handleNewCateringRequest);
      socket.off('catering_request_updated', handleCateringUpdate);
    };
  }, [socket, queryClient]);

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      
      <div className="print:hidden">
        <Topbar />
      </div>

      <div className="flex pt-16 print:pt-0">
        
        <div className="print:hidden">
          <Sidebar />
        </div>

        
        <main className="flex-1 p-4 md:p-6 lg:p-8 overflow-x-hidden">
          
          <div className="mb-6 flex flex-wrap items-center justify-between gap-4 print:hidden">
            
            <div className="flex items-center space-x-3">
              <div className="flex items-center space-x-2">
                {isConnected ? (
                  <Wifi className="h-4 w-4 text-green-500" />
                ) : !isRealtimeEnabled ? (
                  <WifiOff className="h-4 w-4 text-yellow-500" />
                ) : (
                  <WifiOff className="h-4 w-4 text-red-500" />
                )}
                <span className="text-sm text-gray-600 dark:text-gray-400">
                  {isConnected ? 'Connected' : isRealtimeEnabled ? 'Disconnected' : 'Realtime Disabled'}
                </span>
              </div>
              
              
              <button className="relative p-2 text-gray-600 dark:text-gray-400 hover:text-primary-500 transition-colors">
                <Bell className="h-5 w-5" />
                <span className="absolute -top-1 -right-1 h-5 w-5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center">
                  3
                </span>
              </button>
            </div>

            
            <Breadcrumbs />
          </div>

          
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-4 md:p-6">
            <ErrorBoundary resetKey={location.pathname}>
              <Outlet />
            </ErrorBoundary>
          </div>

          
          <footer className="mt-8 text-center text-sm text-gray-500 dark:text-gray-400 print:hidden">
            <p>{settings?.businessName ? `${settings.businessName} ` : ''}Admin Dashboard © {new Date().getFullYear()}</p>
            <p className="mt-1">Version 1.0.0</p>
          </footer>
        </main>
      </div>
    </div>
  );
};

export default Layout;



























              


